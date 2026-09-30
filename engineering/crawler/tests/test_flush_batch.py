"""
Tests for BaseJobCrawler._flush_batch.

_flush_batch is the single point where crawled jobs leave the crawler and
reach the backend, and it owns the retry/drop decision for anything that
fails. A silent regression here (e.g. clearing jobs that should be retried,
or retrying jobs that should be dropped) would lose data quietly in
production with no visible error, so it's covered in isolation from the
rest of the crawl loop.

Run with:
    pytest tests/test_flush_batch.py -v
"""
import json
from unittest.mock import AsyncMock, MagicMock

import httpx
import pytest

from crawlers.base_crawler import BaseJobCrawler, MAX_RETRIES
from models.raw_job import RawJobInput


class _ConcreteCrawler(BaseJobCrawler):
    """BaseJobCrawler is abstract; _flush_batch is what we're testing and
    doesn't need a real crawl_jobs implementation behind it."""

    async def crawl_jobs(self, crawler_run_id, async_client):
        raise NotImplementedError


@pytest.fixture
def crawler():
    return _ConcreteCrawler()


@pytest.fixture
def auth_headers():
    return {"Authorization": "Bearer test-token"}


def make_job(job_id: str = "job-1", crawler_run_id: int = 1) -> RawJobInput:
    return RawJobInput(
        job_id=job_id,
        employer="Acme",
        job_role="Engineer",
        location="Colombo",
        description="desc",
        crawler_run_id=crawler_run_id,
        source="ikman",
    )


def make_response(status_code: int, json_body: dict | None = None, text: str = ""):
    """Builds a mock httpx.Response-like object whose raise_for_status()
    behaves like the real thing, so the except blocks in _flush_batch are
    exercised exactly as they would be against a live response."""
    response = MagicMock(spec=httpx.Response)
    response.status_code = status_code
    response.text = text if not json_body else json.dumps(json_body)

    if json_body is not None:
        response.json.return_value = json_body
    else:
        response.json.side_effect = json.JSONDecodeError("bad json", "", 0)

    if status_code >= 400:
        request = httpx.Request("POST", "https://api.example.com/jobs/batch-save")
        response.raise_for_status.side_effect = httpx.HTTPStatusError(
            f"{status_code} error", request=request, response=response
        )
    else:
        response.raise_for_status.return_value = None

    return response


# ---------------------------------------------------------------------------
# All jobs saved successfully
# ---------------------------------------------------------------------------

class TestAllSucceeded:

    @pytest.mark.asyncio
    async def test_empty_failed_jobs_list_clears_batch(self, crawler, auth_headers):
        job_batch = [make_job("a"), make_job("b")]
        client = AsyncMock()
        client.post.return_value = make_response(200, {"failed_jobs": []})

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == []

    @pytest.mark.asyncio
    async def test_missing_failed_jobs_key_clears_batch(self, crawler, auth_headers):
        """The backend contract says absence of the key means no failures,
        same as an explicit empty list — .get() must treat them the same."""
        job_batch = [make_job("a")]
        client = AsyncMock()
        client.post.return_value = make_response(200, {"message": "ok"})

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == []

    @pytest.mark.asyncio
    async def test_sends_correct_payload_shape(self, crawler, auth_headers):
        """retry_count must never be sent to the backend, and the batch
        must be posted to the batch-save endpoint with auth headers."""
        job = make_job("a")
        job.retry_count = 2
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(200, {"failed_jobs": []})

        await crawler._flush_batch(client, auth_headers, job_batch)

        client.post.assert_awaited_once()
        _, kwargs = client.post.call_args
        assert kwargs["headers"] == auth_headers
        sent_payload = kwargs["json"][0]
        assert "retry_count" not in sent_payload
        assert sent_payload["job_id"] == "a"


# ---------------------------------------------------------------------------
# Partial failures: retryable vs non-retryable
# ---------------------------------------------------------------------------

class TestPartialFailures:

    @pytest.mark.asyncio
    async def test_retryable_failure_kept_in_batch_with_incremented_retry_count(
        self, crawler, auth_headers
    ):
        job = make_job("a")
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"job_id": "a", "status_code": 503, "error": "db down"}]},
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == [job]
        assert job.retry_count == 1

    @pytest.mark.asyncio
    async def test_non_retryable_failure_dropped_from_batch(self, crawler, auth_headers):
        job = make_job("a")
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"job_id": "a", "status_code": 422, "error": "bad province"}]},
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == []

    @pytest.mark.asyncio
    async def test_mixed_batch_only_keeps_the_retryable_failure(self, crawler, auth_headers):
        saved = make_job("saved")
        retryable = make_job("retryable")
        permanent = make_job("permanent")
        job_batch = [saved, retryable, permanent]

        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {
                "failed_jobs": [
                    {"job_id": "retryable", "status_code": 503, "error": "timeout"},
                    {"job_id": "permanent", "status_code": 422, "error": "bad data"},
                ]
            },
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == [retryable]
        assert retryable.retry_count == 1

    @pytest.mark.asyncio
    async def test_matching_is_by_job_id_not_position(self, crawler, auth_headers):
        """Guards against a regression back to positional zip() matching:
        put the failure for the *last* job first in the response list and
        confirm the right job is the one kept."""
        first = make_job("first")
        second = make_job("second")
        job_batch = [first, second]

        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"job_id": "second", "status_code": 503, "error": "x"}]},
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == [second]
        assert second.retry_count == 1
        assert first.retry_count == 0

    @pytest.mark.asyncio
    async def test_unknown_status_code_treated_as_non_retryable(self, crawler, auth_headers):
        """Any status code outside RETRYABLE_STATUS_CODES should be dropped,
        not just 422 specifically — this is the safer default for a code
        the crawler doesn't recognize."""
        job = make_job("a")
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"job_id": "a", "status_code": 418, "error": "teapot"}]},
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == []


# ---------------------------------------------------------------------------
# Retry exhaustion (poison-pill protection)
# ---------------------------------------------------------------------------

class TestRetryExhaustion:

    @pytest.mark.asyncio
    async def test_job_dropped_once_max_retries_reached(self, crawler, auth_headers):
        job = make_job("a")
        job.retry_count = MAX_RETRIES - 1  # one more failure should exhaust it
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"job_id": "a", "status_code": 503, "error": "still down"}]},
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == []

    @pytest.mark.asyncio
    async def test_job_survives_up_to_but_not_including_max_retries(
        self, crawler, auth_headers
    ):
        job = make_job("a")
        job.retry_count = MAX_RETRIES - 2
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"job_id": "a", "status_code": 503, "error": "still down"}]},
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == [job]
        assert job.retry_count == MAX_RETRIES - 1

    @pytest.mark.asyncio
    async def test_repeated_flushes_eventually_drop_a_permanently_failing_job(
        self, crawler, auth_headers
    ):
        """Simulates the real failure mode this design fixes: a job that
        always 503s should be retried a bounded number of times, then
        dropped — never resent forever, never silently kept growing."""
        job = make_job("a")
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"job_id": "a", "status_code": 503, "error": "still down"}]},
        )

        for _ in range(MAX_RETRIES):
            await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == []
        assert client.post.await_count == MAX_RETRIES


# ---------------------------------------------------------------------------
# Transport-level failures (request never got a usable response)
# ---------------------------------------------------------------------------

class TestTransportFailures:

    @pytest.mark.asyncio
    async def test_timeout_keeps_entire_batch_buffered(self, crawler, auth_headers):
        job_batch = [make_job("a"), make_job("b")]
        client = AsyncMock()
        client.post.side_effect = httpx.TimeoutException("timed out")

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert len(job_batch) == 2

    @pytest.mark.asyncio
    async def test_connection_error_keeps_entire_batch_buffered(self, crawler, auth_headers):
        job_batch = [make_job("a"), make_job("b")]
        client = AsyncMock()
        client.post.side_effect = httpx.ConnectError("connection refused")

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert len(job_batch) == 2

    @pytest.mark.asyncio
    async def test_server_error_status_keeps_entire_batch_buffered(self, crawler, auth_headers):
        job_batch = [make_job("a"), make_job("b")]
        client = AsyncMock()
        client.post.return_value = make_response(503, text="upstream down")

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert len(job_batch) == 2

    @pytest.mark.asyncio
    async def test_client_error_status_drops_entire_batch(self, crawler, auth_headers):
        """A 4xx at the whole-request level (bad auth, malformed array) is
        not a per-job failure — the same request would fail again, so the
        whole batch is dropped rather than retried forever."""
        job_batch = [make_job("a"), make_job("b")]
        client = AsyncMock()
        client.post.return_value = make_response(401, text="unauthorized")

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert job_batch == []

    @pytest.mark.asyncio
    async def test_malformed_json_body_on_200_keeps_batch_buffered(self, crawler, auth_headers):
        """A 200 with a body that isn't valid JSON (e.g. a proxy error page)
        must not crash the crawl loop, and must not be treated as success."""
        job_batch = [make_job("a")]
        client = AsyncMock()
        client.post.return_value = make_response(200, json_body=None, text="<html>502</html>")

        await crawler._flush_batch(client, auth_headers, job_batch)

        assert len(job_batch) == 1

    @pytest.mark.asyncio
    async def test_failed_job_entry_missing_job_id_is_ignored_not_crashed(
        self, crawler, auth_headers
    ):
        """Defensive contract check: a malformed failed_jobs entry from the
        backend (missing job_id) should be skipped, not raise KeyError."""
        job = make_job("a")
        job_batch = [job]
        client = AsyncMock()
        client.post.return_value = make_response(
            200,
            {"failed_jobs": [{"status_code": 503, "error": "missing job_id field"}]},
        )

        await crawler._flush_batch(client, auth_headers, job_batch)

        # the malformed entry can't be matched to any job, so from the
        # crawler's point of view no known job failed -> batch clears
        assert job_batch == []