import os
import json
import httpx
import logging
import asyncio
import hashlib
from datetime import datetime, timezone
from typing import List, Dict, Any
from datasketch import MinHash

# Import only what we need for deduplication from your parser
from parsers.ikman_parser import parse_rule_based_fields

# Import your precise schema configurations and prompt rules
from schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION

from crawl4ai import (
    AsyncWebCrawler, 
    CrawlerRunConfig, 
    BrowserConfig, 
    LLMConfig, 
    LLMExtractionStrategy, 
    MemoryAdaptiveDispatcher
)

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

BACKEND_BASE_URL = os.getenv("BACKEND_URL", "http://backend:8080/api/v1/crawler")

async def get_last_page_from_text():
   async with AsyncWebCrawler() as crawler:
      
       result = await crawler.arun(url="https://ikman.lk/en/ads/sri-lanka/jobs")
      
       html_content = result.html
      
       match = re.search(r'of ([\d,]+) ads', html_content)
      
       if match:
          
           total_ads = int(match.group(1).replace(',', ''))
           ads_per_page = 25
           last_page = math.ceil(total_ads / ads_per_page)
          
           print(f"Total Ads: {total_ads}")
           print(f"Calculated Last Page: {last_page}")
           return last_page
       else:
           print("Could not find the total ad count text.")
           return 1


def generate_production_minhash_and_lsh(job_data: dict, num_perm: int = 128, num_bands: int = 8) -> tuple:
    """Uses text shingling to generate stable MinHash and 8 LSH bucket band keys."""
    text_content = f"{job_data['employer']} {job_data['job_role']} {job_data['location']} {job_data['key_responsibilities']}"
    text_content = " ".join(text_content.lower().split())
    
    k = 3
    shingles = set(text_content[i:i+k] for i in range(len(text_content) - k + 1))
    
    m = MinHash(num_perm=num_perm)
    for shingle in shingles:
        m.update(shingle.encode('utf-8'))
        
    raw_signature = [int(h) for h in m.hashvalues]
    lsh_indexes = []
    hashes_per_band = num_perm // num_bands
    
    for band_no in range(num_bands):
        start_idx = band_no * hashes_per_band
        end_idx = start_idx + hashes_per_band
        band_hashes = m.hashvalues[start_idx:end_idx]
        
        hasher = hashlib.md5()
        hasher.update(band_hashes.tobytes())
        bucket_key = f"b{band_no}_{hasher.hexdigest()[:16]}"
        
        lsh_indexes.append({
            "band_no": band_no,
            "bucket_key": bucket_key
        })
        
    return raw_signature, lsh_indexes


async def check_duplicate_via_backend(client: httpx.AsyncClient, lsh_indexes: List[Dict[str, Any]], current_sig: List[int], incoming_location: str = "") -> tuple:
    """
    Queries your Go radar routing registry to evaluate bucket collisions via Jaccard metrics.
    Handles the flat JSON array response matching full DB job schemas.
    """
    try:
        # Flatten your lsh_indexes structured array into a raw array of strings
        flat_bucket_keys = [item["bucket_key"] for item in lsh_indexes]
        
        request_payload = {
            "bucket_keys": flat_bucket_keys
        }
        
        # Hit your updated routing lookup endpoint
        response = await client.post(f"{BACKEND_BASE_URL}/radar/lookup", json=request_payload)
        if response.status_code != 200:
            logger.warning(f"Lookup server returned unexpected status code: {response.status_code}")
            return False, None
            
        # ✔ FIX 1: Parse the response directly as a flat list, defaulting to an empty list if empty
        similar_jobs = response.json()
        if not isinstance(similar_jobs, list) or not similar_jobs:
            return False, None

        # Compare signatures for all matching database candidates
        for job in similar_jobs:

            db_location = (job.get("location") or "").strip().lower()
            inc_location = incoming_location.strip().lower()
            
            if not _locations_compatible(db_location, inc_location):
                logger.info(f"Skipping Job ID {job.get('id')}: location mismatch ('{inc_location}' vs '{db_location}')")
                continue

            # ✔ FIX 2: Pull the signature from the nested 'meta_data' block
            meta = job.get("meta_data") or {}
            db_sig = meta.get("minhash_signature")
            
            # If signature is null or length doesn't match, skip to the next candidate
            if not db_sig or len(db_sig) != len(current_sig):
                continue
                
            # Compute Jaccard Intersection over Union
            intersection = sum(1 for i, j in zip(current_sig, db_sig) if i == j)
            union = len(current_sig) + len(db_sig) - intersection
            jaccard_value = intersection / union if union > 0 else 0.0
            
            # Match confirmed if similarity value passes 50% threshold
            if jaccard_value >= 0.65:
                # ✔ FIX 3: Pull the primary key ID from the root level field 'id'
                matched_job_post_id = job.get("id")
                logger.info(f"Collision detected! Jaccard score: {round(jaccard_value, 2)} against Job ID: {matched_job_post_id}")
                return True, matched_job_post_id
                
    except Exception as e:
        logger.error(f"Backend radar lookup collision check failed: {e}")
        
    return False, None

def _locations_compatible(loc_a: str, loc_b: str) -> bool:
    if not loc_a or not loc_b:
        return True  # can't determine, don't block

    if loc_a == loc_b:
        return True

    if loc_a in loc_b or loc_b in loc_a:
        return True

    # Share at least one significant word
    stopwords = {"the", "of", "in", "at", "and", "a"}
    words_a = set(loc_a.split()) - stopwords
    words_b = set(loc_b.split()) - stopwords
    return bool(words_a & words_b)

async def run_pipeline_orchestrator(max_pages: int = 5):

    #max_pages = await get_last_page_from_text()

    new_jobs_buffer: List[Dict[str, Any]] = []
    lsh_index_buffer: List[Dict[str, Any]] = []
    updated_jobs_buffer: List[Dict[str, Any]] = []
    
    async_client = httpx.AsyncClient(timeout=30.0)
    
    try:
        current_time_iso = datetime.now(timezone.utc).astimezone().isoformat()
        
        start_payload = {
            "started_at": current_time_iso,
            "finished_at": None,
            "status": "RUNNING"
        }
        
        # Hit the correct endpoint: /runs
        init_res = await async_client.post(f"{BACKEND_BASE_URL}/runs", json=start_payload)
        
        # Extract the sequence tracker ID returned by Go GORM/SQL insert
        crawler_run_id = init_res.json().get("id", 1)
        logger.info(f"Initialized Tracking Session Run ID: {crawler_run_id}")
    except Exception as e:
        logger.warning(f"Could not connect to tracking backend. Defaulting fallback to run sequence ID 1: {e}")
        crawler_run_id = 1

    browser_config = BrowserConfig(headless=True, extra_args=["--disable-gpu", "--no-sandbox"])
    dispatcher = MemoryAdaptiveDispatcher(memory_threshold_percent=80.0, max_session_permit=10)
    
    # Configure LLM strategy using your cleanly imported structural schema parameters
    llm_extraction_strategy = LLMExtractionStrategy(
        llm_config=LLMConfig(provider="deepseek/deepseek-chat", api_token=os.getenv("DEEPSEEK_API_KEY")),
        instruction=BASE_JOB_INSTRUCTION,
        schema=json.dumps(JOB_EXTRACTION_SCHEMA),
        extra_args={"base_url": "https://api.deepseek.com", "temperature": 0.0},
    )

    async with AsyncWebCrawler(config=browser_config) as crawler:
        all_detail_urls = []
        
        # Scrape index tracking paths
        for page in range(1, max_pages + 1):
            url = f"https://ikman.lk/en/ads/sri-lanka/jobs?page={page}"
            logger.info(f"Scanning Listing Page Index: {page}")
            res = await crawler.arun(url=url, config=CrawlerRunConfig(cache_mode="BYPASS"))
            if res.success:
                links = [f"https://ikman.lk{l['href']}" if l['href'].startswith('/') else l['href'] 
                         for l in res.links.get("internal", []) if "/en/ad/" in l['href']]
                all_detail_urls.extend(links)
                
        unique_urls = list(set(all_detail_urls))
        logger.info(f"Processing structural extraction queue for {len(unique_urls)} links.")

        # Stream pages sequentially to extract lightweight fields for deduplication check
        detail_config = CrawlerRunConfig(cache_mode="BYPASS", stream=True)
        results_generator = await crawler.arun_many(urls=unique_urls, config=detail_config, dispatcher=dispatcher)

        async for result in results_generator:
            if not result.success or not result.markdown:
                continue
            
            # Fast rule-based execution for zero-cost deduplication checking
            temp_payload = parse_rule_based_fields(
                markdown=result.markdown.raw_markdown, 
                job_link=result.url,
                source="Ikman"
            )
            
            minhash_sig, lsh_indexes = generate_production_minhash_and_lsh(temp_payload)
            is_duplicate, matched_id = await check_duplicate_via_backend(async_client, lsh_indexes, minhash_sig, incoming_location=temp_payload.get("location", ""))
            
            if is_duplicate:
                logger.info(f"Duplicate Match Found: Routing job reference {matched_id} to keep-alive updates.")
                updated_jobs_buffer.append({
                    "job_post_id": matched_id,
                    "crawler_run_id": crawler_run_id
                })
                
                if len(updated_jobs_buffer) >= 10:
                    await async_client.post(f"{BACKEND_BASE_URL}/jobs/batch-update", json={"duplicates": updated_jobs_buffer})
                    updated_jobs_buffer.clear()
            else:
                logger.info(f"Unique entry found. Calling LLM to parse entire schema: {result.url}")
                llm_res = await crawler.arun(
                    url=result.url, 
                    config=CrawlerRunConfig(extraction_strategy=llm_extraction_strategy, cache_mode="BYPASS")
                )
                
                if llm_res.success and llm_res.extracted_content:
                    try:
                        extracted_job = json.loads(llm_res.extracted_content)
                        if isinstance(extracted_job, list) and len(extracted_job) > 0:
                            extracted_job = extracted_job[0]
                        
                        # Inject system tracking meta blocks back into the structured object
                        extracted_job["meta_data"]["crawler_run_id"] = crawler_run_id
                        extracted_job["meta_data"]["minhash_signature"] = minhash_sig
                        
                        new_jobs_buffer.append(extracted_job)
                        
                        for idx_item in lsh_indexes:
                            lsh_index_buffer.append(idx_item)
                            
                    except Exception as e:
                        logger.error(f"Failed to unmarshal LLM response into schema format: {e}")
                        
                if len(new_jobs_buffer) >= 10:
                    payload = {"new_jobs": new_jobs_buffer, "lsh_indexes": lsh_index_buffer}
                    await async_client.post(f"{BACKEND_BASE_URL}/jobs/batch-save", json=payload)
                    new_jobs_buffer.clear()
                    lsh_index_buffer.clear()

        # --- Trailing record cleanup and synchronization ---
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

        # Run reconciliation cleanup
        logger.info("Executing pipeline reconciliation. Retiring dead listings from active pool.")
        await async_client.post(f"{BACKEND_BASE_URL}/jobs/reconcile", json={"crawler_run_id": crawler_run_id})
        await async_client.post(f"{BACKEND_BASE_URL}/runs/{crawler_run_id}/complete", json={"id": crawler_run_id, "status": "COMPLETED"})
        
    await async_client.aclose()
    logger.info("Scraper execution pipeline concluded successfully.")
