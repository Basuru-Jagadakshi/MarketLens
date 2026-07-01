import requests
import math
import time

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
        
        # Parse and append to list
        for job in response['body']['data']:
            all_jobs.append(parse_job(job))
            
        time.sleep(1) # Polite crawler delay
        
    return all_jobs

def parse_job(job):
    """Extracts specific fields from the raw API data."""
    return {
        "employer": job.get("company_name"),
        "job_role": job.get("title"),
        "location": job.get("location"),
        "description": job.get("description")
    }

