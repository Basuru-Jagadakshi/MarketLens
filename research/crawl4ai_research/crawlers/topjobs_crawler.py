import json
import asyncio
import requests
import math
import time
from crawl4ai import LLMExtractionStrategy, LLMConfig
from crawlers.schemas.job_schema import JOB_EXTRACTION_SCHEMA, BASE_JOB_INSTRUCTION
from crawlers.parsers.topjobs_parser import run_crawler

async def topjobs_jobs_extraction():

    job_data_list = await run_crawler()

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
                print(f"SUCCESS: Normalized job ID {job.get('row_id')}")
            else:
                print(f"WARNING: No extraction result for job {job.get('row_id')} — ocr_text may be empty/unreadable")
            print(f"SUCCESS: Normalized job ID {job.get('row_id')}")

        except Exception as e:
            print(f"ERROR: Failed to process job {job.get('row_id')}: {e}")

    return extracted_jobs_list