import asyncio
import logging
import time

import httpx

from config import (
    THUNDER_BASE_URL,
    THUNDER_CLIENT_ID,
    THUNDER_CLIENT_SECRET,
    THUNDER_RESOURCE,
    THUNDER_VERIFY_TLS,
)

logger = logging.getLogger(__name__)

TOKEN_TTL_SECONDS = 3600
TOKEN_REFRESH_BUFFER_SECONDS = 60


class ThunderIDClient:
    def __init__(self):
        self._cached_token: str | None = None
        self._cached_token_expiry: float = 0
        self._lock = asyncio.Lock()

    async def get_access_token(self):
        if self._is_cached_token_valid():
            return self._cached_token

        async with self._lock:
            if self._is_cached_token_valid():
                return self._cached_token

            token = await self._fetch_access_token()
            self._cached_token = token
            self._cached_token_expiry = (
                time.monotonic() + TOKEN_TTL_SECONDS - TOKEN_REFRESH_BUFFER_SECONDS
            )
            return token

    def invalidate_token(self):
        """Forces the next get_access_token() call to fetch a fresh token,
        e.g. after the backend rejects the cached one with a 401 — our local
        TTL has no way of knowing the token died earlier than expected."""
        self._cached_token = None
        self._cached_token_expiry = 0

    def _is_cached_token_valid(self):
        return self._cached_token is not None and time.monotonic() < self._cached_token_expiry

    async def _fetch_access_token(self):
        async with httpx.AsyncClient(verify=THUNDER_VERIFY_TLS) as client:
            response = await client.post(
                f"{THUNDER_BASE_URL}/oauth2/token",
                auth=(THUNDER_CLIENT_ID, THUNDER_CLIENT_SECRET),
                headers={"Content-Type": "application/x-www-form-urlencoded"},
                data={
                    "grant_type": "client_credentials",
                    "scope": "crawler:runs crawler:complete crawler:lookup crawler:batch-save crawler:batch-update crawler:reconcile",
                    "resource": THUNDER_RESOURCE,
                },
            )
            response.raise_for_status()
            body = response.json()
            token = body.get("access_token")
            if not token:
                raise ValueError(f"ThunderID token response missing 'access_token': {body}")

            logger.info("Fetched new ThunderID access token")
            return token


async def get_auth_headers(thunder_client: ThunderIDClient) -> dict:
    """Shared by every call site that needs a Bearer header: fetches the
    (possibly cached) token and logs clearly if that fails, instead of each
    caller repeating the same try/except around get_access_token()."""
    try:
        token = await thunder_client.get_access_token()
    except Exception as e:
        logger.error(f"Failed to obtain ThunderID access token: {e}")
        raise
    return {"Authorization": f"Bearer {token}"}
