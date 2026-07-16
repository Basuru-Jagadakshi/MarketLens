from abc import ABC, abstractmethod
import httpx

class BaseJobCrawler(ABC):

    @abstractmethod
    async def crawl_jobs(self, crawler_run_id: int, async_client: httpx.AsyncClient) -> None:
        pass