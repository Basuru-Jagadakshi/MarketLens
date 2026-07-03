import asyncio
import requests
from bs4 import BeautifulSoup

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
        
        # Extract and clean data
        return {
            "job_title": data.get("jobTitle"),
            "employer": data.get("jobItem", {}).get("organizationName"),
            "location": data.get("jobItem", {}).get("locations"),
            "description": clean_html(data.get("jobInfo", ""))
        }
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