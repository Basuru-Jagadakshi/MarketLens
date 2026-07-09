import asyncio
import os
import time

from datetime import datetime, timedelta
from zoneinfo import ZoneInfo
from pathlib import Path
from dotenv import load_dotenv
from schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION
from orchestrator.ikman_orchestrator import run_pipeline_orchestrator

env_path = Path(__file__).resolve().parent / '.env'

load_dotenv(dotenv_path=env_path)

print(f"API Key Loaded: {os.getenv('DEEPSEEK_API_KEY') is not None}")

async def crawl_job():
    """Encapsulates the actual crawling logic execution using the new pipeline."""
    print(f"\n--- Execution Started at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')} ---")
    
    try:
        # ✔ SWAP OUT THE OLD LOGIC FOR THE NEW DEDUPLICATING ORCHESTRATOR
        # Set max_pages to whatever scanning depth you want to process per loop execution
        await run_pipeline_orchestrator(max_pages=2)
        
    except Exception as e:
        print(f"CRITICAL ERROR encountered during execution lifecycle: {e}")
        
    print("--- Execution Complete ---")

async def main():

    #This part is for the server
    print("=== Automated 10-Day Interval Scheduler Started ===")

    sl_tz = ZoneInfo("Asia/Colombo")
    
    while True:
        now = datetime.now(sl_tz)
        
        target_time = now.replace(hour=8, minute=16, second=0, microsecond=0)
        
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