import json
import asyncio
import requests
from bs4 import BeautifulSoup
import math
import time
from crawl4ai import LLMExtractionStrategy, LLMConfig
from crawlers.schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION

def clean_html(html_content):
    """Removes HTML tags and returns clean text."""
    if not html_content:
        return ""
    soup = BeautifulSoup(html_content, "html.parser")
    return soup.get_text(separator=" ").strip()

def fetch_job_details(job_id):
    """Fetches full details for a single job."""
    url = f"https://xpress.jobs/api/jobs/publishedJob?jobId={job_id}"
    response = requests.get(url)
    if response.status_code == 200:
        data = response.json()
        data['jobInfo'] = clean_html(data.get("jobInfo", ""))
        
        # Extract and clean data
        return data
    return None

async def process_all_jobs():
    final_data = []
    page = 1
    
    while True:
        print(f"--- Fetching page {page} ---")
        
        # Build the URL with the current page
        list_url = f"https://xpress.jobs/api/jobs/searchJobs?page={page}&pageSize=20&keyword=&locations=&sectors=&jobTypes=&careerLevels=&sortBy=SortedCreateDate+DESC&byCVLess=false&byWalkIn=false"
        
        try:
            response = requests.get(list_url, timeout=10)
            jobs_list = response.json()
        except Exception as e:
            print(f"Error fetching page {page}: {e}")
            break
            
        # Break the loop if the list is empty
        if not jobs_list:
            print("No more jobs found. Finishing.")
            break
        
        # Process each job on the current page
        for job_summary in jobs_list:
            job_id = job_summary['jobId']
            print(f"Processing job {job_id}: {job_summary['jobTitle']}")
            
            details = fetch_job_details(job_id)
            if details:
                final_data.append(details)
            
            await asyncio.sleep(10)
            
        # Move to next page
        page += 1

        if page >= 3:
            break
        
        await asyncio.sleep(10)
            
    return final_data


async def rooster_jobs_extraction():

    job_data_list = await process_all_jobs()

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
            if result:
                extracted_jobs_list.append(result[0])
                print(f"SUCCESS: Normalized job ID {job.get('jobId')}")
            else:
                print(f"WARNING: No extraction")

        except Exception as e:
            print(f"ERROR: Failed to process job {job.get('jobId')}: {e}")

    return extracted_jobs_list