import asyncio
import httpx
import logging

from typing import List
from bs4 import BeautifulSoup
from pydantic import ValidationError
from crawlers.base_crawler import BaseJobCrawler
from utils.thunder_id_client import ThunderIDClient
from parsers.xpressjobs_parser import XpressJobsParser
from models.raw_job import RawJobInput
from config import BACKEND_BASE_URL, BATCH_SIZE

logger = logging.getLogger(__name__)

class XpressJobsCrawler(BaseJobCrawler):

    def __init__(self):
        self._parser = XpressJobsParser()
        self.async_client = httpx.AsyncClient(timeout=30.0)
        self._thunder_client = ThunderIDClient() 

    def _clean_html(self, html_content):
        if not html_content:
            return ""
        soup = BeautifulSoup(html_content, "html.parser")
        return soup.get_text(separator=" ").strip()

    async def _fetch_job_details(self, job_id):
        try:
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
        except Exception as e:
            logger.warning(f"Failed to fetch details for job {job_id}: {e}")
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
                logger.error(f"Error fetching page {page}: {e}")
                break
                
            # Break the loop if the list is empty
            if not jobs_list:
                logger.info("No more jobs found. Finishing.")
                break
            
            # Process each job on the current page
            for job_summary in jobs_list:
                job_id = job_summary['jobId']
                logger.info(f"Processing job {job_id}: {job_summary['jobTitle']}")
                
                details = await self._fetch_job_details(job_id)
                if details:
                    final_data.append(details)
                
                await asyncio.sleep(10)
                
            # Move to next page
            page += 1
            
            await asyncio.sleep(10)
                
        return final_data

    #This funtion starts the crawler and save or update the job after checking whether job already exists or not
    async def crawl_jobs(
        self,
        crawler_run_id: int,
        async_client: httpx.AsyncClient,
    ) -> None:
 
        logger.info("Xpress jobs crawl started.")
 
        try:
            token = await self._thunder_client.get_access_token()
        except Exception as e:
            logger.error(f"Failed to obtain ThunderID access token: {e}")
            raise
        auth_headers = {"Authorization": f"Bearer {token}"} 
        
        job_data_list = await self._process_all_jobs()
 
        job_batch: List[RawJobInput] = []
 
        for result in job_data_list:
            try:
                job_input = self._parser.parse_rule_based_fields(result, crawler_run_id)
            except ValidationError as e:
                logger.warning(f"Skipping malformed job: {e}")
                continue
 
            job_batch.append(job_input)
 
            if len(job_batch) >= BATCH_SIZE:
                logger.info(f"Flushing full batch of {len(job_batch)} job records to backend.")
                await self._flush_batch(async_client, auth_headers, job_batch)
 
        if job_batch:
            logger.info(f"Flushing remaining {len(job_batch)} job records to backend.")
            await self._flush_batch(async_client, auth_headers, job_batch)
 
        logger.info("Xpress jobs crawl pass concluded.")