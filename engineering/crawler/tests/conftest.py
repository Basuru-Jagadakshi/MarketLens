from unittest.mock import AsyncMock

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
    """A stand-in for ThunderIDClient used by _flush_batch — these tests are
    about batch/retry behavior, not token caching, so a mock that always
    hands back the same token is enough."""
    mock = AsyncMock()
    mock.get_access_token.return_value = "test-token"
    return mock


@pytest.fixture
def client():
    return ThunderIDClient()


@pytest.fixture
def another_client():
    """A second, independent ThunderIDClient — for tests that must confirm
    two instances don't share cache state."""
    return ThunderIDClient()
