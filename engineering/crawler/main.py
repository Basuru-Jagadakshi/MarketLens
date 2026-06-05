import asyncio
import os
import time

from datetime import datetime, timedelta
from zoneinfo import ZoneInfo
from pathlib import Path
from dotenv import load_dotenv
from crawlers.ikman_jobs_crawler.ikman_jobs_crawler import ikman_jobs_extraction, get_last_page_from_text
from storage_handlers.handler import save_jobs_to_api
from crawlers.schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION

env_path = Path(__file__).resolve().parent / '.env'

load_dotenv(dotenv_path=env_path)

print(f"API Key Loaded: {os.getenv('DEEPSEEK_API_KEY') is not None}")

async def crawl_job():
    """Encapsulates the actual crawling logic execution."""
    print(f"\n--- Execution Started at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')} ---")
    ikman_task = ikman_jobs_extraction()
    ikman_jobs = await ikman_task
    save_jobs_to_api(ikman_jobs)
    print("--- Execution Complete ---")

async def main():
    #This part is used to testing
    # print("=== Starting Multi-Site Crawl... ===")



    # ikman_task = ikman_jobs_extraction()

    # ikman_jobs = await ikman_task

    # all_combined_jobs = ikman_jobs

    # print(all_combined_jobs)

    # save_jobs_to_api(all_combined_jobs)



    # print("\n=== Crawling Complete ===")





    #This part is for the server
    print("=== Automated 10-Day Interval Scheduler Started ===")

    sl_tz = ZoneInfo("Asia/Colombo")
    
    while True:
        now = datetime.now(sl_tz)
        
        target_time = now.replace(hour=17, minute=40, second=0, microsecond=0)
        
        if now >= target_time:
            target_time += timedelta(days=1)
            
        seconds_to_wait = (target_time - now).total_seconds()
        
        print(f"Next routine run scheduled for: {target_time.strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"Sleeping for {round(seconds_to_wait / 3600, 2)} hours until 10 PM...")
        
        await asyncio.sleep(seconds_to_wait)
        
        await crawl_job()
        
        next_cycle_target = datetime.now(sl_tz) + timedelta(days=10)
        print(f"Cycle complete. Next run in 10 days on: {next_cycle_target.strftime('%Y-%m-%d %H:%M:%S')}")
        
        await asyncio.sleep(10 * 24 * 60 * 60)

if __name__ == "__main__":
    try:
        asyncio.run(main())
    except KeyboardInterrupt:
        print("\nCrawling stopped by user.")