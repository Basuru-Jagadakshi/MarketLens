import asyncio
import requests
import httpx
import typing

from bs4 import BeautifulSoup
from crawlers.base_crawler import BaseJobCrawler
from utils.dedup_utils import JobDuplicationCheck
from parsers.xpressjobs_parser import XpressJobsParser
from schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION
from config import BACKEND_BASE_URL, BATCH_SIZE

logging.basicConfig(level=logging.INFO, format="%(asctime)s [%(levelname)s] %(message)s")
logger = logging.getLogger("xpressjobs_crawler")

class XpreeJobsCrawler(BaseJobCrawler):

    def __init__(self):
        self._parser = XpressJobsParser()
        self.duplication_checker = JobDuplicationCheck()
        self.async_client = httpx.AsyncClient()

    def _clean_html(self, html_content):
        if not html_content:
            return ""
        soup = BeautifulSoup(html_content, "html.parser")
        return soup.get_text(separator=" ").strip()

    async def _fetch_job_details(self, job_id):
        url = f"https://xpress.jobs/api/jobs/publishedJob?jobId={job_id}"
        response = await self.async_client.get(url)
        if response.status_code == 200:
            data = response.json()
            
            return {
                "job_title": data.get("jobTitle"),
                "employer": data.get("jobItem", {}).get("organizationName"),
                "location": data.get("jobItem", {}).get("locations"),
                "description": self._clean_html(data.get("jobInfo", ""))
            }
        return None

    async def _process_all_jobs(self):
        final_data = []
        page = 1
        
        while True:
            logger.info(f"--- Fetching page {page} ---")
            
            # Build the URL with the current page
            list_url = f"https://xpress.jobs/api/jobs/searchJobs?page={page}&pageSize=20&keyword=&locations=&sectors=&jobTypes=&careerLevels=&sortBy=SortedCreateDate+DESC&byCVLess=false&byWalkIn=false"
            
            try:
                response = await self.async_client.get(list_url, timeout=10)
                jobs_list = response.json()
            except Exception as e:
                print(f"Error fetching page {page}: {e}")
                break
                
            # Break the loop if the list is empty
            if not jobs_list:
                print("No more jobs found. Finishing.")
                break
            
            # Process each job on the current page
            for job_summary in jobs_list:
                job_id = job_summary['jobId']
                print(f"Processing job {job_id}: {job_summary['jobTitle']}")
                
                details = await self._fetch_job_details(job_id)
                if details:
                    final_data.append(details)
                
                await asyncio.sleep(10)
                
            # Move to next page
            page += 1

            if page >= 3:
                break
            
            await asyncio.sleep(10)
                
        return final_data

    #This funtion starts the crawler and save or update the job after checking whether job already exists or not
    async def crawl_jobs(
        self,
        crawler_run_id: int,
        async_client: httpx.AsyncClient
    ) -> None:

        logger.info("Xpress jobs crawl started.")
        
        job_data_list = await self._process_all_jobs()

        new_jobs_buffer: List[Dict[str, Any]] = []
        lsh_index_buffer: List[Dict[str, Any]] = []
        updated_jobs_buffer: List[Dict[str, Any]] = []

        detail_extraction_strategy = LLMExtractionStrategy(
            llm_config=LLMConfig(
                provider="deepseek/deepseek-chat",
                api_token=os.getenv("DEEPSEEK_API_KEY"),
            ),
            instruction=BASE_JOB_INSTRUCTION,
            schema=json.dumps(JOB_EXTRACTION_SCHEMA),
            extraction_type="schema", 
            apply_chunking=False,          
            extra_args={"base_url": "https://api.deepseek.com", "temperature": 0.0},
        )

        for result in job_data_list:

            temp_payload = self._parser.parse_rule_based_fields(result)

            raw_text = json.dumps(result)

            minhash_sig, lsh_indexes = self.duplication_checker.generate_production_minhash_and_lsh(temp_payload)
            is_duplicate, matched_id = await self.duplication_checker.check_duplicate_via_backend(
                async_client,
                BACKEND_BASE_URL,
                lsh_indexes,
                minhash_sig,
                incoming_location=temp_payload.get("location", ""),
            )

            if is_duplicate:
                logger.info(f"Duplicate Match Found: Routing job reference {matched_id} to keep-alive updates.")
                updated_jobs_buffer.append({
                    "job_post_id": matched_id,
                    "crawler_run_id": crawler_run_id,
                })

                if len(updated_jobs_buffer) >= BATCH_SIZE:
                    await async_client.post(f"{BACKEND_BASE_URL}/jobs/batch-update", json={"duplicates": updated_jobs_buffer})
                    updated_jobs_buffer.clear()
            else:
                logger.info(f"Unique entry found. Calling LLM to parse entire schema")
                llm_res = await asyncio.to_thread(
                    detail_extraction_strategy.run, url="", sections=[raw_text]
                )

                if llm_res and isinstance(llm_res, list) and len(llm_res) > 0:
                    try:
                        extracted_job = llm_res[0]

                        extracted_job["meta_data"]["crawler_run_id"] = crawler_run_id
                        extracted_job["meta_data"]["minhash_signature"] = minhash_sig

                        new_jobs_buffer.append(extracted_job)

                        for idx_item in lsh_indexes:
                            lsh_index_buffer.append(idx_item)

                    except Exception as e:
                        logger.error(f"Failed to unmarshal LLM response into schema format: {e}")

                if len(new_jobs_buffer) >= BATCH_SIZE:
                    payload = {"new_jobs": new_jobs_buffer, "lsh_indexes": lsh_index_buffer}
                    await async_client.post(f"{BACKEND_BASE_URL}/jobs/batch-save", json=payload)
                    new_jobs_buffer.clear()
                    lsh_index_buffer.clear()

        if updated_jobs_buffer:
            logger.info(f"Flushing remaining {len(updated_jobs_buffer)} update records to backend.")
            await async_client.post(f"{BACKEND_BASE_URL}/jobs/batch-update", json={"duplicates": updated_jobs_buffer})
            updated_jobs_buffer.clear()

        if new_jobs_buffer:
            logger.info(f"Flushing remaining {len(new_jobs_buffer)} insertion records to backend.")
            payload = {"new_jobs": new_jobs_buffer, "lsh_indexes": lsh_index_buffer}
            await async_client.post(f"{BACKEND_BASE_URL}/jobs/batch-save", json=payload)
            new_jobs_buffer.clear()
            lsh_index_buffer.clear()

        logger.info("Xpress jobs crawl pass concluded.") 