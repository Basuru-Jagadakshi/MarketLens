import asyncio
import httpx
import logging
import math
from typing import List
from pydantic import ValidationError
from utils.thunder_id_client import ThunderIDClient  
from config import BACKEND_BASE_URL, BATCH_SIZE

from crawlers.base_crawler import BaseJobCrawler
from parsers.rooster_parser import RoosterParser
from models.raw_job import RawJobInput

logger = logging.getLogger(__name__)

class RoosterCrawler(BaseJobCrawler):

    def __init__(self):
        self._parser = RoosterParser()
        self._thunder_client = ThunderIDClient() 

    async def _fetch_all_jobs(self, async_client: httpx.AsyncClient):
        base_url = "https://api.rooster.jobs/jobSearch/jobs/search"
        limit = 20
        all_jobs = []
        
        # Initial call to get total count
        payload = {"query": [], "limit": limit, "page": 1, "filters": {"country": "Sri Lanka"}}
        response = (await async_client.post(base_url, json=payload)).json()
        total_jobs = response['body']['count']
        total_pages = math.ceil(total_jobs / limit)
        
        logger.info(f"Total jobs to fetch: {total_jobs} over {total_pages} pages.")

        for page in range(1, total_pages + 1):
            payload['page'] = page
            response = (await async_client.post(base_url, json=payload)).json()

            for job in response['body']['data']:
                all_jobs.append(job)
                
            await asyncio.sleep(1)
            
        return all_jobs

    #This funtion starts the crawler and save or update the job after checking whether job already exists or not
    async def crawl_jobs(
        self,
        crawler_run_id: int,
        async_client: httpx.AsyncClient,
    ) -> None:
 
        logger.info("Rooster crawl started.")
 
        try:
            token = await self._thunder_client.get_access_token()
        except Exception as e:
            logger.error(f"Failed to obtain ThunderID access token: {e}")
            raise
        auth_headers = {"Authorization": f"Bearer {token}"}
 
        job_data_list = await self._fetch_all_jobs(async_client)
 
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
 
        logger.info("Rooster crawl pass concluded.")