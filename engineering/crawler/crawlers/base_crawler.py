from abc import ABC, abstractmethod
import logging
from typing import List

import httpx

from config import BACKEND_BASE_URL
from models.raw_job import RawJobInput

logger = logging.getLogger(__name__)


class BaseJobCrawler(ABC):

    @abstractmethod
    async def crawl_jobs(
        self,
        crawler_run_id: int,
        async_client: httpx.AsyncClient,
    ) -> None:
        pass

    async def _flush_batch(
        self,
        async_client: httpx.AsyncClient,
        auth_headers: dict,
        job_batch: List[RawJobInput],
    ) -> None:
        try:
            response = await async_client.post(
                f"{BACKEND_BASE_URL}/jobs/batch-save",
                json=[job.model_dump() for job in job_batch],
                headers=auth_headers,
            )
            response.raise_for_status()
        except httpx.TimeoutException as e:
            logger.error(
                f"Batch POST timed out after waiting for a response ({len(job_batch)} "
                f"jobs pending) — keeping batch buffered to retry on the next flush: {e}"
            )
            return
        except httpx.RequestError as e:
            logger.error(
                f"Batch POST failed due to a network/connection error ({len(job_batch)} "
                f"jobs pending) — keeping batch buffered to retry on the next flush: {e}"
            )
            return
        except httpx.HTTPStatusError as e:
            status = e.response.status_code
            if status >= 500:
                logger.error(
                    f"Batch POST rejected by backend with server error {status} "
                    f"({len(job_batch)} jobs pending) — keeping batch buffered to "
                    f"retry on the next flush: {e.response.text}"
                )
            else:
                logger.error(
                    f"Batch POST rejected by backend with client error {status} — "
                    f"dropping this batch of {len(job_batch)} jobs, NOT retrying "
                    f"(same payload would fail again): {e.response.text}"
                )
                job_batch.clear()
            return

        job_batch.clear()