import os
import requests

def save_jobs_to_api(jobs: list):
    """
    Streams scraped job data directly to the Go backend API endpoint.
     Bypasses the database completely and relies on the Go container
    to process and save the data.
    """
    if not jobs:
        print("No jobs found to send. Skipping API update.")
        return

    api_url = os.getenv("API_URL", "http://localhost:8080/api/v1/jobs")
    
    success_count = 0
    fail_count = 0
    total_jobs = len(jobs)

    print(f"Starting to stream {len(jobs)} jobs to the backend API at: {api_url}", flush=True)

    for index, job in enumerate(jobs, start=1):
        try:
            response = requests.post(
                api_url, 
                json=job, 
                headers={"Content-Type": "application/json"},
                timeout=10
            )
            
            if response.status_code in [200, 201]:
                success_count += 1
            else:
                fail_count += 1
                print(f"[{index}] Failed to process job. Status: {response.status_code}, Server Error: {response.text}")

        except requests.exceptions.Timeout:
            fail_count += 1
            print(f"[{index}/{total_jobs}] Timeout error: Backend took too long to reply.")
                
        except requests.exceptions.RequestException as e:
            fail_count += 1
            print(f"[{index}] Network error connecting to backend: {e}")

        if index % 10 == 0 or index == total_jobs:
            print(f"Progress: Processed {index}/{total_jobs} jobs...")

    print(f"\n--- Sync Pipeline Summary ---")
    print(f"Successfully streamed to Go API: {success_count} jobs")
    if fail_count > 0:
        print(f"Failed pipeline connections  : {fail_count} jobs")