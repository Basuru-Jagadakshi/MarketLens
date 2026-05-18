import json
import asyncio
import re
import math
import gc
from typing import List
from crawl4ai import AsyncWebCrawler, CrawlerRunConfig, BrowserConfig, LLMConfig, LLMExtractionStrategy, MemoryAdaptiveDispatcher, CrawlerMonitor
from crawlers.schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION




BATCH_SIZE = 15


async def get_last_page_from_text():
   async with AsyncWebCrawler() as crawler:
      
       result = await crawler.arun(url="https://ikman.lk/en/ads/sri-lanka/jobs")
      
       html_content = result.html
      
       match = re.search(r'of ([\d,]+) ads', html_content)
      
       if match:
          
           total_ads = int(match.group(1).replace(',', ''))
           ads_per_page = 25
           last_page = math.ceil(total_ads / ads_per_page)
          
           print(f"Total Ads: {total_ads}")
           print(f"Calculated Last Page: {last_page}")
           return last_page
       else:
           print("Could not find the total ad count text.")
           return 1
      






async def ikman_jobs_extraction(max_pages: int = 3):


   dispatcher = MemoryAdaptiveDispatcher(
       memory_threshold_percent=70.0,
       check_interval=1.0,     
       max_session_permit=1
   )


   extracted_jobs_list = []
  
   detail_extraction_strategy = LLMExtractionStrategy(
       llm_config=LLMConfig(
           provider="deepseek/deepseek-chat",
           api_token="env:DEEPSEEK_API_KEY",
       ),
       instruction=BASE_JOB_INSTRUCTION,
       schema=json.dumps(JOB_EXTRACTION_SCHEMA),
       extra_args={"base_url": "https://api.deepseek.com", "temperature": 0.0},
   )


   browser_config = BrowserConfig(
       headless=True,
       extra_args=[
           "--disable-gpu",
           "--disable-dev-shm-usage",
           "--no-sandbox",
           "--js-flags=--max-old-space-size=512" # Limits JS memory
       ]
   )




   async with AsyncWebCrawler(config=browser_config) as crawler:


       # ── 1. Collect all job URLs from every listing page ───────────
       max_pages    = await get_last_page_from_text()
       all_job_urls = []


       for page in range(1, max_pages + 1):
           url = f"https://ikman.lk/en/ads/sri-lanka/jobs?page={page}"
           print(f"Listing {page}/{max_pages}: {url}")


           listing_result = await crawler.arun(
               url=url,
               config=CrawlerRunConfig(cache_mode="ENABLED"),
           )


           if listing_result.success:
               links = [
                   link["href"]
                   for link in listing_result.links.get("internal", [])
                   if "/en/ad/" in link["href"]
               ]
               all_job_urls.extend(links)
               print(f"  └─ {len(links)} links")
           else:
               print(f"  └─ FAILED to crawl listing page {page}")


       unique_urls = list(set(all_job_urls))
       full_urls   = [
           f"https://ikman.lk{u}" if u.startswith("/") else u
           for u in unique_urls
       ]


       if not full_urls:
           print("No job URLs found.")
           return []


       print(f"\n  Total URLs to crawl: {len(full_urls)}\n")


       # ── 2. Crawl detail pages in batches ──────────────────────────
       detail_config = CrawlerRunConfig(
           extraction_strategy=detail_extraction_strategy,
           cache_mode="BYPASS",
           stream=True,
           page_timeout=60000,
       )


       total_batches = math.ceil(len(full_urls) / BATCH_SIZE)


       for batch_num, i in enumerate(range(0, len(full_urls), BATCH_SIZE), 1):
           batch = full_urls[i : i + BATCH_SIZE]
           print(f"\n─── Batch {batch_num}/{total_batches}  ({len(batch)} URLs) ───")


           results_generator = await crawler.arun_many(
               urls=batch,
               config=detail_config,
               dispatcher=dispatcher,
           )


           async for result in results_generator:
               if not result.success:
                   print(f"  SKIP  : {result.url} | {result.error_message}")
                   continue


               if not result.extracted_content:
                   print(f"  WARN  : No content from {result.url}")
                   continue


               try:
                   data = json.loads(result.extracted_content)
                   jobs = data if isinstance(data, list) else [data]


                   #crawled_at = datetime.now(timezone.utc).isoformat()
                   # for job in jobs:
                   #     job["crawled_at"] = crawled_at
                   #     job["run_date"]   = run_date
                   #     job["source_url"] = result.url


                   extracted_jobs_list.extend(jobs)
                   print(f"  OK    : {result.url}")


               except json.JSONDecodeError:
                   print(f"  ERROR : Invalid JSON from {result.url}")
               except Exception as e:
                   print(f"  ERROR : {result.url} → {e}")


           gc.collect()
           await asyncio.sleep(5)
   print(f"\n  Total jobs extracted: {len(extracted_jobs_list)}")
   return extracted_jobs_list


   # async with AsyncWebCrawler(config=browser_config) as crawler:
   #     all_job_urls = []
      
   #     for page in range(1, max_pages + 1):
          
   #         url = f"https://ikman.lk/en/ads/sri-lanka/jobs?page={page}"
   #         print(f"Crawling listing page {page}: {url}")
          
   #         listing_result = await crawler.arun(
   #             url=url,
   #             config=CrawlerRunConfig(cache_mode="BYPASS")
   #         )


   #         if listing_result.success:
              
   #             current_page_links = [
   #                 link['href'] for link in listing_result.links.get("internal", [])
   #                 if "/en/ad/" in link['href']
   #             ]
   #             all_job_urls.extend(current_page_links)
   #             print(f"Found {len(current_page_links)} links on page {page}.")
   #         else:
   #             print(f"Failed to crawl page {page}")


   #     unique_urls = list(set(all_job_urls))
   #     full_urls = [f"https://ikman.lk{url}" if url.startswith('/') else url for url in unique_urls]


   #     if not full_urls:
   #         print("No job URLs found across all pages!")
   #         return


   #     #limit_urls = full_urls[:]
   #     # ... (previous setup code)


   #     print(f"Total jobs found: {len(full_urls)}.")
      
   #     detail_config = CrawlerRunConfig(
   #         extraction_strategy=detail_extraction_strategy,
   #         cache_mode="BYPASS",
   #         stream=True,        
   #         page_timeout=60000  
   #     )


   #     # Start the multi-page crawl
   #     results_generator = await crawler.arun_many(
   #         urls=full_urls,
   #         config=detail_config,
   #         dispatcher=dispatcher
   #     )


   #     # Use async for to process each result as it finishes
   #     async for result in results_generator:
   #         # 1. HANDLE TIMEOUTS / ERRORS
   #         if not result.success:
   #             # This logs the error and moves to the NEXT result automatically
   #             print(f"SKIPPING: {result.url} | Error: {result.error_message}")
   #             continue


   #         # 2. HANDLE SUCCESSFUL FETCH BUT EMPTY CONTENT
   #         if not result.extracted_content:
   #             print(f"WARNING: No content extracted from {result.url}")
   #             continue


   #         # 3. ATTEMPT DATA PARSING
   #         try:
   #             data = json.loads(result.extracted_content)
              
   #             # Success! Add to your list
   #             if isinstance(data, list):
   #                 extracted_jobs_list.extend(data)
   #             else:
   #                 extracted_jobs_list.append(data)
              
   #             print(f"SUCCESS: Extracted data from {result.url}")


   #         except json.JSONDecodeError:
   #             print(f"ERROR: LLM returned invalid JSON for {result.url}")
   #         except Exception as e:
   #             print(f"ERROR: Unexpected error processing {result.url}: {e}")


       # return extracted_jobs_list
