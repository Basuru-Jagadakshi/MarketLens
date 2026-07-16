import logging
import asyncio
import httpx
from datetime import datetime, timezone
from typing import Callable, Awaitable, Dict, List, Optional, Type

from config import BACKEND_BASE_URL

from crawlers.base_crawler import BaseJobCrawler
from crawlers.ikman_crawler import IkmanCrawler

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

class CrawlerManager:

    def __init__(self):
        self._registry: Dict[str, Type[BaseJobCrawler]] = {
            "ikman": IkmanCrawler
        }

    #This function created the new crawler session and return the new crawler run id
    async def _start_run(self, client: httpx.AsyncClient) -> int:
        current_time_iso = datetime.now(timezone.utc).astimezone().isoformat()
        start_payload = {
            "started_at": current_time_iso,
            "finished_at": None,
            "status": "RUNNING",
        }
        try:
            init_res = await client.post(f"{BACKEND_BASE_URL}/runs", json=start_payload)
            crawler_run_id = init_res.json().get("id", 1)
            logger.info(f"Initialized Tracking Session Run ID: {crawler_run_id}")
            return crawler_run_id
        except Exception as e:
            logger.warning(f"Could not connect to tracking backend. Defaulting fallback to run sequence ID 1: {e}")
            return 1

    #This function sets the status of the current crawling session to "COMPLETED" 
    #and sets the end date of the jobs that are not equal to current crawler run id
    async def _finalize_run(self, client: httpx.AsyncClient, crawler_run_id: int) -> None:
        try:
            logger.info("Executing pipeline reconciliation. Retiring dead listings from active pool.")
            await client.post(f"{BACKEND_BASE_URL}/jobs/reconcile", json={"crawler_run_id": crawler_run_id})
            await client.post(
                f"{BACKEND_BASE_URL}/runs/{crawler_run_id}/complete",
                json={"id": crawler_run_id, "status": "COMPLETED"},
            )
        except Exception as e:
            logger.error(f"Failed to finalize crawler run {crawler_run_id}: {e}")

    #This function calls the relevant crawlers
    async def _run_crawler(
        self,
        name: str,
        crawler_run_id: int,
        client: httpx.AsyncClient
    ) -> None:
        crawler_class = self._registry.get(name)
        if not crawler_class:
            logger.warning(f"Crawler '{name}' not found in registry.")
            return

        try:
            logger.info(f"--- Starting crawler: {name} ---")
            crawler_instance = crawler_class()
            await crawler_instance.crawl_jobs(crawler_run_id=crawler_run_id, async_client=client)
            logger.info(f"--- Finished crawler: {name} ---")
        except Exception as e:
            logger.error(f"Crawler '{name}' failed: {e}", exc_info=True)

    #This function calls the all crawlers one by one and send it to _run_crawler function
    async def run_all_crawlers(
        self,
        crawler_names: Optional[List[str]] = None,
        concurrent: bool = False
    ) -> None:
        names = crawler_names or list(self._registry.keys())

        async with httpx.AsyncClient(timeout=30.0) as client:
            crawler_run_id = await self._start_run(client)
            
            tasks = [self._run_crawler(name, crawler_run_id, client) for name in names if name in self._registry]
            
            if concurrent:
                await asyncio.gather(*tasks, return_exceptions=True)
            else:
                for task in tasks:
                    await task
            
            await self._finalize_run(client, crawler_run_id)
            logger.info("Scraper execution pipeline concluded.")