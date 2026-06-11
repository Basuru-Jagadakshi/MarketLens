# import json
# import asyncio
# import re
# import math
# import gc
# from typing import List
# from crawl4ai import AsyncWebCrawler, CrawlerRunConfig, BrowserConfig, LLMConfig, LLMExtractionStrategy, MemoryAdaptiveDispatcher, CrawlerMonitor
# from crawlers.schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION


# async def get_last_page_from_text():
#    async with AsyncWebCrawler() as crawler:
      
#        result = await crawler.arun(url="https://ikman.lk/en/ads/sri-lanka/jobs")
      
#        html_content = result.html
      
#        match = re.search(r'of ([\d,]+) ads', html_content)
      
#        if match:
          
#            total_ads = int(match.group(1).replace(',', ''))
#            ads_per_page = 25
#            last_page = math.ceil(total_ads / ads_per_page)
          
#            print(f"Total Ads: {total_ads}")
#            print(f"Calculated Last Page: {last_page}")
#            return last_page
#        else:
#            print("Could not find the total ad count text.")
#            return 1

# async def ikman_jobs_extraction(max_pages: int = 50):
   
#    max_pages = await get_last_page_from_text()

#    dispatcher = MemoryAdaptiveDispatcher(
#        memory_threshold_percent=80.0,
#        check_interval=1.0,     
#        max_session_permit=10
#    )


#    extracted_jobs_list = []
  
#    detail_extraction_strategy = LLMExtractionStrategy(
#        llm_config=LLMConfig(
#            provider="deepseek/deepseek-chat",
#            api_token="env:DEEPSEEK_API_KEY",
#        ),
#        instruction=BASE_JOB_INSTRUCTION,
#        schema=json.dumps(JOB_EXTRACTION_SCHEMA),
#        extra_args={"base_url": "https://api.deepseek.com", "temperature": 0.0},
#    )


#    browser_config = BrowserConfig(
#        headless=True,
#        verbose=True,
#        extra_args=[
#            "--disable-gpu",
#            "--disable-dev-shm-usage",
#            "--no-sandbox",
#            "--js-flags=--max-old-space-size=512" # Limits JS memory
#        ]
#    )

#    async with AsyncWebCrawler(config=browser_config) as crawler:
#        all_job_urls = []
      
#        for page in range(1, max_pages + 1):
          
#            url = f"https://ikman.lk/en/ads/sri-lanka/jobs?page={page}"
#            print(f"Crawling listing page {page}: {url}")
          
#            listing_result = await crawler.arun(
#                url=url,
#                config=CrawlerRunConfig(cache_mode="BYPASS")
#            )


#            if listing_result.success:
              
#                current_page_links = [
#                    link['href'] for link in listing_result.links.get("internal", [])
#                    if "/en/ad/" in link['href']
#                ]
#                all_job_urls.extend(current_page_links)
#                print(f"Found {len(current_page_links)} links on page {page}.")
#            else:
#                print(f"Failed to crawl page {page}")


#        unique_urls = list(set(all_job_urls))
#        full_urls = [f"https://ikman.lk{url}" if url.startswith('/') else url for url in unique_urls]


#        if not full_urls:
#            print("No job URLs found across all pages!")
#            return


#        #limit_urls = full_urls[:]
#        # ... (previous setup code)


#        print(f"Total jobs found: {len(full_urls)}.")
      
#        detail_config = CrawlerRunConfig(
#            extraction_strategy=detail_extraction_strategy,
#            cache_mode="BYPASS",
#            stream=True,        
#            page_timeout=60000  
#        )

#        results_generator = await crawler.arun_many(
#            urls=full_urls,
#            config=detail_config,
#            dispatcher=dispatcher
#        )

#        async for result in results_generator:
           
#            if not result.success:
              
#                print(f"SKIPPING: {result.url} | Error: {result.error_message}")
#                continue


#            if not result.extracted_content:
#                print(f"WARNING: No content extracted from {result.url}")
#                continue


#            try:
#                data = json.loads(result.extracted_content)
              
#                if isinstance(data, list):
#                    extracted_jobs_list.extend(data)
#                else:
#                    extracted_jobs_list.append(data)
              
#                print(f"SUCCESS: Extracted data from {result.url}")


#            except json.JSONDecodeError:
#                print(f"ERROR: LLM returned invalid JSON for {result.url}")
#            except Exception as e:
#                print(f"ERROR: Unexpected error processing {result.url}: {e}")


#        return extracted_jobs_list















































# crawlers/ikman_crawler.py

import json
import asyncio
import re
import math
import tiktoken

from crawl4ai import (
    AsyncWebCrawler,
    BrowserConfig,
    CrawlerRunConfig,
    CacheMode,
    LLMConfig,
)
from crawl4ai import LLMExtractionStrategy
from crawl4ai.async_dispatcher import MemoryAdaptiveDispatcher, RateLimiter
from crawl4ai import CrawlerMonitor

#from crawlers.schemas.job_schema import LLMJobFields, LLM_ONLY_INSTRUCTION
from crawlers.parsers.ikman_parser import parse_rule_based_fields, merge_llm_fields

def calculate_projected_cost(avg_input_tokens: int, total_ads: int) -> dict:
    """
    Calculates the estimated DeepSeek V4-Flash API cost based on 
    average tokens per page and total site volume.
    """
    total_input_tokens = avg_input_tokens * total_ads
    
    total_output_tokens = 800 * total_ads 
    cache_hits = int(total_input_tokens * 0.90)
    cache_misses = total_input_tokens - cache_hits
    
    cost_miss = (cache_misses / 1_000_000) * 0.14
    cost_hit = (cache_hits / 1_000_000) * 0.0028
    cost_output = (total_output_tokens / 1_000_000) * 0.28
    
    total_cost_usd = cost_miss + cost_hit + cost_output
    
    print("=" * 50)
    print("           DEEPSEEK API COST ESTIMATE          ")
    print("=" * 50)
    print(f"Total Target Postings    : {total_ads:,} pages")
    print(f"Avg. Input Tokens / Page : {avg_input_tokens:,}")
    print(f"Avg. Output Tokens / Page: 800")
    print("-" * 50)
    print(f"Projected Input (Misses) : {cache_misses:,} tokens  -> ${cost_miss:.4f}")
    print(f"Projected Input (Hits)   : {cache_hits:,} tokens  -> ${cost_hit:.4f}")
    print(f"Projected Output         : {total_output_tokens:,} tokens  -> ${cost_output:.4f}")
    print("-" * 50)
    print(f"TOTAL ESTIMATED BUDGET   : ${total_cost_usd:.4f} USD")
    print("=" * 50)
    
    return {
        "total_cost_usd": total_cost_usd,
        "total_tokens": total_input_tokens + total_output_tokens
    }

async def get_last_page_from_text(crawler: AsyncWebCrawler) -> int:
    result = await crawler.arun(
        url="https://ikman.lk/en/ads/sri-lanka/jobs",
        config=CrawlerRunConfig(cache_mode=CacheMode.BYPASS)
    )
    match = re.search(r'of ([\d,]+) ads', result.html or "")
    if match:
        total = int(match.group(1).replace(',', ''))
        last_page = math.ceil(total / 25)
        print(f"Total ads: {total} → {last_page} pages")
        return last_page
    return 1


async def ikman_jobs_extraction(max_pages: int = 1) -> list:

    totalTokenCount = 0

    browser_config = BrowserConfig(
        headless=True,
        verbose=False,
        extra_args=[
            "--disable-gpu",
            "--disable-dev-shm-usage",
            "--no-sandbox",
            "--js-flags=--max-old-space-size=512",
        ]
    )

    # LLM extracts ONLY 4 fields — minimal token cost
    # llm_strategy = LLMExtractionStrategy(
    #     llm_config=LLMConfig(
    #         provider="deepseek/deepseek-chat",
    #         api_token="env:DEEPSEEK_API_KEY",
    #         base_url="https://api.deepseek.com",
    #     ),
    #     schema=LLMJobFields.model_json_schema(),
    #     extraction_type="schema",
    #     instruction=LLM_ONLY_INSTRUCTION,
    #     apply_chunking=False,          # single job page — no chunking needed
    #     input_format="fit_markdown",   # filtered markdown = fewer tokens
    #     extra_args={"temperature": 0.0, "max_tokens": 500},
    #     verbose=False,
    # )

    dispatcher = MemoryAdaptiveDispatcher(
        memory_threshold_percent=80.0,
        check_interval=1.0,
        max_session_permit=10,
        rate_limiter=RateLimiter(
            base_delay=(1.0, 2.0),
            max_delay=30.0,
            max_retries=2,
        ),
        monitor=CrawlerMonitor(
            
        )
    )

    extracted_jobs = []

    async with AsyncWebCrawler(config=browser_config) as crawler:

        # --- Step 1: Collect job URLs from listing pages ---
        all_job_urls = []
        for page in range(1, max_pages + 1):
            url = f"https://ikman.lk/en/ads/sri-lanka/jobs?page={page}"
            print(f"Scraping listing page {page}: {url}")

            listing_result = await crawler.arun(
                url=url,
                config=CrawlerRunConfig(cache_mode=CacheMode.BYPASS)
            )

            if listing_result.success:
                links = [
                    link["href"]
                    for link in listing_result.links.get("internal", [])
                    if "/en/ad/" in link["href"]
                ]
                all_job_urls.extend(links)
                print(f"  Found {len(links)} job links on page {page}")
            else:
                print(f"  Failed: {listing_result.error_message}")

        # Deduplicate and build full URLs
        full_urls = list({
            f"https://ikman.lk{u}" if u.startswith("/") else u
            for u in all_job_urls
        })

        if not full_urls:
            print("No job URLs found.")
            return []

        print(f"\nTotal unique jobs to crawl: {len(full_urls)}")

        # --- Step 2: Crawl detail pages ---
        # detail_config = CrawlerRunConfig(
        #     extraction_strategy=llm_strategy,
        #     cache_mode=CacheMode.BYPASS,
        #     stream=True,
        #     page_timeout=60000,
        # )

        results = await crawler.arun_many(
            urls=full_urls,
            #config=detail_config,
            dispatcher=dispatcher,
        )

        for result in results:
            if not result.success:
                print(f"  SKIP (failed): {result.url} | {result.error_message}")
                continue

            # Rule-based fields from markdown — free
            rule_fields = parse_rule_based_fields(
                result.markdown or "", source="Ikman"
            )

            tokenizer = tiktoken.get_encoding("cl100k_base")
            token_ids = tokenizer.encode(result.markdown)
            totalTokenCount+=len(token_ids)

            # LLM fields (only 4)
            # llm_fields = {}
            # if result.extracted_content:
            #     try:
            #         parsed = json.loads(result.extracted_content)
            #         llm_fields = parsed[0] if isinstance(parsed, list) else parsed
            #     except (json.JSONDecodeError, IndexError, TypeError):
            #         print(f"  LLM parse error: {result.url}")
            #         rule_fields["meta_data"]["error"] = True

            # final_job = merge_llm_fields(rule_fields, llm_fields)
            # final_job["url"] = result.url

            extracted_jobs.append(rule_fields)
            #print(f"  OK: {final_job['job_role']} @ {final_job['employer']}")

        # Show token usage summary at the end
        #llm_strategy.show_usage()
    
    averageTokenCountPerPage = totalTokenCount / len(full_urls)
    metrics = calculate_projected_cost(averageTokenCountPerPage, total_ads=8000)


    return extracted_jobs

