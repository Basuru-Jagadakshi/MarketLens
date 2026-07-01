import json
import asyncio
import requests
import math
import time
from crawl4ai import LLMExtractionStrategy, LLMConfig
from crawlers.schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION

def fetch_all_jobs():
    base_url = "https://api.rooster.jobs/jobSearch/jobs/search"
    limit = 20
    all_jobs = []
    
    # Initial call to get total count
    payload = {"query": [], "limit": limit, "page": 1, "filters": {"country": "Sri Lanka"}}
    response = requests.post(base_url, json=payload).json()
    total_jobs = response['body']['count']
    total_pages = math.ceil(total_jobs / limit)
    
    print(f"Total jobs to fetch: {total_jobs} over {total_pages} pages.")

    for page in range(1, total_pages + 1):
        print(f"Fetching page {page}...")
        payload['page'] = page
        response = requests.post(base_url, json=payload).json()

        for job in response['body']['data']:
            all_jobs.append(job)
            
        time.sleep(1) # Polite crawler delay
        
    return all_jobs

async def rooster_jobs_extraction():

    job_data_list = fetch_all_jobs()

    detail_extraction_strategy = LLMExtractionStrategy(
        llm_config=LLMConfig(
            provider="deepseek/deepseek-chat",
            api_token="env:DEEPSEEK_API_KEY",
        ),
        instruction=BASE_JOB_INSTRUCTION,
        schema=json.dumps(JOB_EXTRACTION_SCHEMA),
        extraction_type="schema", 
        apply_chunking=False,          
        extra_args={"base_url": "https://api.deepseek.com", "temperature": 0.0},
    )

    extracted_jobs_list = []
    print(f"Starting LLM normalization for {len(job_data_list)} jobs...")

    for job in job_data_list:
        raw_text = json.dumps(job)

        try:
            # run() is sync and takes `sections`, a list of strings
            result = detail_extraction_strategy.run(url="", sections=[raw_text])
            print(result)
            # result is already a list of parsed dicts — no json.loads needed
            extracted_jobs_list.append(result)
            print(f"SUCCESS: Normalized job ID {job.get('id')}")

        except Exception as e:
            print(f"ERROR: Failed to process job {job.get('id')}: {e}")

    return extracted_jobs_list