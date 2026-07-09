import logging
import asyncio
import httpx
from datetime import datetime, timezone
from typing import Callable, Awaitable, Dict, List, Optional

from config import BACKEND_BASE_URL

from crawlers.ikman_crawler import run_ikman_pipeline as run_ikman

CRAWLER_REGISTRY: Dict[str, Callable[..., Awaitable[None]]] = {
    "ikman": run_ikman
}

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)


async def _start_run(client: httpx.AsyncClient) -> int:
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


async def _finalize_run(client: httpx.AsyncClient, crawler_run_id: int) -> None:
    try:
        logger.info("Executing pipeline reconciliation. Retiring dead listings from active pool.")
        await client.post(f"{BACKEND_BASE_URL}/jobs/reconcile", json={"crawler_run_id": crawler_run_id})
        await client.post(
            f"{BACKEND_BASE_URL}/runs/{crawler_run_id}/complete",
            json={"id": crawler_run_id, "status": "COMPLETED"},
        )
    except Exception as e:
        logger.error(f"Failed to finalize crawler run {crawler_run_id}: {e}")


async def _run_one_crawler(
    name: str,
    orchestrator: Callable[..., Awaitable[None]],
    crawler_run_id: int,
    client: httpx.AsyncClient,
    **kwargs,
) -> None:
    try:
        logger.info(f"--- Starting crawler: {name} ---")
        await orchestrator(crawler_run_id=crawler_run_id, async_client=client, **kwargs)
        logger.info(f"--- Finished crawler: {name} ---")
    except Exception as e:
        logger.error(f"Crawler '{name}' failed and was skipped: {e}", exc_info=True)


async def run_all_crawlers(
    crawler_names: Optional[List[str]] = None,
    concurrent: bool = False,
    **crawler_kwargs,
) -> None:
    names = crawler_names or list(CRAWLER_REGISTRY.keys())
    unknown = [n for n in names if n not in CRAWLER_REGISTRY]
    if unknown:
        logger.warning(f"Unknown crawler name(s) requested and skipped: {unknown}")

    async_client = httpx.AsyncClient(timeout=30.0)
    try:
        crawler_run_id = await _start_run(async_client)

        runnable = [n for n in names if n in CRAWLER_REGISTRY]
        tasks = [
            _run_one_crawler(name, CRAWLER_REGISTRY[name], crawler_run_id, async_client, **crawler_kwargs)
            for name in runnable
        ]

        if concurrent:
            await asyncio.gather(*tasks, return_exceptions=True)
        else:
            for task in tasks:
                await task

        await _finalize_run(async_client, crawler_run_id)
    finally:
        await async_client.aclose()
        logger.info("Scraper execution pipeline concluded.")