from unittest.mock import AsyncMock, MagicMock

import pytest

from crawlers.base_crawler import BaseJobCrawler
from utils.thunder_id_client import ThunderIDClient


class _ConcreteCrawler(BaseJobCrawler):
    """BaseJobCrawler is abstract; _flush_batch is what we're testing and
    doesn't need a real crawl_jobs implementation behind it."""

    async def crawl_jobs(self, crawler_run_id, async_client, thunder_client):
        raise NotImplementedError


@pytest.fixture
def crawler():
    return _ConcreteCrawler()


@pytest.fixture
def auth_headers():
    return {"Authorization": "Bearer test-token"}


@pytest.fixture
def thunder_client():
    """A stand-in for ThunderIDClient used by _flush_batch and
    crawler_run_manager — these tests are about batch/retry/finalize
    behavior, not token caching, so a mock that always hands back the same
    token is enough. get_access_token is async like the real client;
    invalidate_token is left as a plain (sync) MagicMock call, since
    ThunderIDClient.invalidate_token() is never awaited by callers."""
    mock = MagicMock()
    mock.get_access_token = AsyncMock(return_value="test-token")
    return mock


@pytest.fixture
def client():
    return ThunderIDClient()


@pytest.fixture
def another_client():
    """A second, independent ThunderIDClient — for tests that must confirm
    two instances don't share cache state."""
    return ThunderIDClient()
