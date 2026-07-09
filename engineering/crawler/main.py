import asyncio
import os

from datetime import datetime, timedelta
from zoneinfo import ZoneInfo
from pathlib import Path
from dotenv import load_dotenv
from utils.crawler_run_manager import run_all_crawlers

env_path = Path(__file__).resolve().parent / '.env'
load_dotenv(dotenv_path=env_path)

print(f"API Key Loaded: {os.getenv('DEEPSEEK_API_KEY') is not None}")

SL_TZ = ZoneInfo("Asia/Colombo")

RUN_HOUR = 8
RUN_MINUTE = 16
RUN_INTERVAL_DAYS = 10
MAX_PAGES_PER_CRAWLER = 2


async def crawl_job():
    print(f"\n--- Execution Started at {datetime.now(SL_TZ).strftime('%Y-%m-%d %H:%M:%S')} ---")

    try:
        await run_all_crawlers(max_pages=MAX_PAGES_PER_CRAWLER)
    except Exception as e:
        print(f"CRITICAL ERROR encountered during execution lifecycle: {e}")

    print("--- Execution Complete ---")


def _next_run_time(now: datetime) -> datetime:
    target = now.replace(hour=RUN_HOUR, minute=RUN_MINUTE, second=0, microsecond=0)
    if now >= target:
        target += timedelta(days=1)
    return target


async def main():
    print("=== Automated Scheduler Started ===")
    print(f"Runs every {RUN_INTERVAL_DAYS} days at {RUN_HOUR:02d}:{RUN_MINUTE:02d} (Asia/Colombo)")

    now = datetime.now(SL_TZ)
    target_time = _next_run_time(now)

    while True:
        now = datetime.now(SL_TZ)
        seconds_to_wait = (target_time - now).total_seconds()

        if seconds_to_wait > 0:
            print(f"Next run scheduled for: {target_time.strftime('%Y-%m-%d %H:%M:%S %Z')}")
            print(f"Sleeping for {round(seconds_to_wait / 3600, 2)} hours...")
            await asyncio.sleep(seconds_to_wait)

        await crawl_job()

        target_time = target_time + timedelta(days=RUN_INTERVAL_DAYS)
        print(f"Cycle complete. Next run in {RUN_INTERVAL_DAYS} days on: {target_time.strftime('%Y-%m-%d %H:%M:%S %Z')}")


if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\nCrawling stopped by user.")