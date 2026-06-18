import os
import json
import pytest
import httpx
from datetime import datetime

# Import your core pipeline deduplication engine features
from orchestrator.ikman_orchestrator import (
    generate_production_minhash_and_lsh,
    check_duplicate_via_backend
)

BACKEND_BASE_URL = os.getenv("BACKEND_URL", "http://backend:8080/api/v1/crawler")

# =====================================================================
# CORE MOCK DATA OBJECTS
# =====================================================================
JOB_A_REFERENCE = {
    "employer": "LDF",
    "job_role": "Software Engineer",
    "job_type": {"name": "Full Time"},
    "key_responsibilities": "Design and implement Goverment job analysis product called MarketLens. Should be able to learn new technologies",
    "qualifications": "BSc in Computer Science or equivalent engineering field. Foundational knowledge of Go (Golang) or Java with Spring Boot. Understanding of JWT authentication systems is a plus.",
    "location": "Colombo 03, Sri Lanka",
    "offers": "International exposure, competitive salary package, health insurance, and hybrid work flexibility.",
    "is_remote": False,
    "skills": [{"name": "Go (Golang)"}, {"name": "REST APIs"}, {"name": "Microservices"}]
}

JOB_B_UNIQUE = {
    "employer": "Sysco LABS",
    "job_role": "QA Automation Engineer",
    "job_type": {"name": "Full Time"},
    "key_responsibilities": "Develop automated testing suites for cloud native storage applications, write integration tests in Python, and maintain CI/CD automation pipelines.",
    "qualifications": "BSc in Software Engineering or equivalent experience. Experience with Selenium, PyTest, Docker, and Linux environments.",
    "location": "Colombo 07, Sri Lanka",
    "offers": "Competitive base, performance bonuses, wellness benefits, and cross-team training.",
    "is_remote": True,
    "skills": [{"name": "Python"}, {"name": "Selenium"}, {"name": "Docker"}]
}

JOB_C_UNIQUE = {
    "employer": "Dialog Axiata",
    "job_role": "Data Scientist",
    "job_type": {"name": "Full Time"},
    "key_responsibilities": "Build predictive models on customer telecom patterns, optimize data pipeline streaming infrastructure, and deploy production machine learning classifiers.",
    "qualifications": "MSc or BSc in Statistics, Data Science, or Mathematics. Proficient with pandas, SQL, Apache Spark, and cloud streaming components.",
    "location": "Colombo 02, Sri Lanka",
    "offers": "Telecom allowance, performance-linked pay scales, premium family health coverage.",
    "is_remote": False,
    "skills": [{"name": "SQL"}, {"name": "Apache Spark"}, {"name": "Machine Learning"}]
}


@pytest.mark.asyncio
async def test_pipeline_loop_deduplication_lifecycle():
    """
    Simulates full iterative production scraper loops across multiple sessions.
    Validates dynamic runtime branching between batch-saves and batch-updates.
    """
    async with httpx.AsyncClient(timeout=15.0) as client:

        # =====================================================================
        # SESSION 1: SCRAPING INITIAL 3 JOBS
        # =====================================================================
        print("\n🚀 [SESSION 1] Initializing first crawler orchestration run...")
        
        session_1_payload = {
            "started_at": "2026-06-14T14:35:10.124578+05:30",
            "finished_at": None,
            "status": "RUNNING"
        }
        init_res_1 = await client.post(f"{BACKEND_BASE_URL}/runs", json=session_1_payload)
        assert init_res_1.status_code in (200, 201)
        session_1_id = int(init_res_1.json().get("id"))
        print(f"   -> Session 1 Assigned ID: {session_1_id}")

        # List 1 containing 3 distinct job listings
        session_1_job_list = [
            JOB_A_REFERENCE.copy(),
            JOB_B_UNIQUE.copy(),
            JOB_C_UNIQUE.copy()
        ]

        print(f"   -> Iterating over {len(session_1_job_list)} collected job listings...")
        for incoming_job in session_1_job_list:
            # Generate local minhash signatures and LSH indexes
            sig, lsh = generate_production_minhash_and_lsh(incoming_job)
            
            # Enforce production meta data structure alignment
            incoming_job["meta_data"] = {
                "geo": {"province": "Western"},
                "posted_at": datetime.utcnow().isoformat() + "Z",
                "source": "Direct Careers Page",
                "standardized_category": "Software Engineering & Technology",
                "seniority": "Mid-Senior",
                "confidence_score": 0.98,
                "ai_version": "v1.2.0",
                "error": False,
                "crawler_run_id": session_1_id,
                "minhash_signature": sig
            }
            
            formatted_lsh = [{"band_no": item["band_no"], "bucket_key": item["bucket_key"]} for item in lsh]

            # Request backend duplicate classification check
            is_dup, matched_id = await check_duplicate_via_backend(client, formatted_lsh, sig, incoming_location=incoming_job.get("location", ""))
            
            if not is_dup:
                print(f"      🔹 New Unique [ {incoming_job['employer']} - {incoming_job['job_role']} ]. Executing Batch-Save...")
                save_res = await client.post(f"{BACKEND_BASE_URL}/jobs/batch-save", json={
                    "new_jobs": [incoming_job],
                    "lsh_indexes": formatted_lsh
                })
                assert save_res.status_code in (200, 201)
            else:
                pytest.fail(f"Unexpected collision encountered during clean session setup run for id: {matched_id}, {incoming_job}")

        # Reconcile Session 1
        print("   -> Running data reconciliation sync for Session 1...")
        reconcile_res_1 = await client.post(f"{BACKEND_BASE_URL}/jobs/reconcile", json={"crawler_run_id": session_1_id})
        assert reconcile_res_1.status_code == 200

        # Complete Session 1
        close_res_1 = await client.post(f"{BACKEND_BASE_URL}/runs/{session_1_id}/complete", json={"status": "COMPLETED"})
        assert close_res_1.status_code in (200, 201)
        print("🏁 [SESSION 1] Concluded successfully.")


        # =====================================================================
        # SESSION 2: NEXT CRON TRIGGER (1 NEAR-DUPLICATE + 2 BRAND NEW)
        # =====================================================================
        print("\n🚀 [SESSION 2] Initializing second sequential crawler orchestration run...")
        
        session_2_payload = {
            "started_at": "2026-06-15T10:15:00.000000+05:30",
            "finished_at": None,
            "status": "RUNNING"
        }
        init_res_2 = await client.post(f"{BACKEND_BASE_URL}/runs", json=session_2_payload)
        assert init_res_2.status_code in (200, 201)
        session_2_id = int(init_res_2.json().get("id"))
        print(f"   -> Session 2 Assigned ID: {session_2_id}")

        # Simulate Job A being re-scraped with minor modifications
        job_a_reposted_duplicate = JOB_A_REFERENCE.copy()
        job_a_reposted_duplicate["key_responsibilities"] = (
            "Design and implement Goverment job analysis product called MarketLens. Should be able to learn new technologies"
        )

        # 2 Brand new distinct jobs to append
        job_d_new = {
            "employer": "WSO2",
            "job_role": "Finance Officer",
            "job_type": {"name": "Full Time"},
            "key_responsibilities": "Develop core APIs using Go and Spring Boot, help maintain local Docker configurations, and learn cloud orchestration tools.",
            "qualifications": "Undergraduate student focusing on backend development. Strong interest in microservices architectures and technical blogging.",
            "location": "Colombo 03, Sri Lanka",
            "offers": "Mentorship, great engineering culture, hybrid flexibility.",
            "is_remote": False,
            "skills": [{"name": "Java"}, {"name": "Go (Golang)"}]
        }

        job_e_new = {
            "employer": "Virtusa",
            "job_role": "Cloud DevOps Engineer",
            "job_type": {"name": "Full Time"},
            "key_responsibilities": "Manage production AWS cloud architecture instances, maintain Terraform scripts, and monitor application health dashboards.",
            "qualifications": "BSc in Computer Science or IT. Experience with AWS cloud infrastructure components, bash scripting, and Kubernetes.",
            "location": "Colombo 05, Sri Lanka",
            "offers": "Global mobility options, certified training budgets.",
            "is_remote": False,
            "skills": [{"name": "AWS"}, {"name": "Terraform"}, {"name": "Kubernetes"}]
        }

        job_f_new = {
            "employer": "LDF",
            "job_role": "Software Engineer",
            "job_type": {"name": "Full Time"},
            "key_responsibilities": "Design and implement Goverment job analysis product called MarketLens. Should be able to learn new technologies",
            "qualifications": "BSc in Computer Science or equivalent engineering field. Foundational knowledge of Go (Golang) or Java with Spring Boot. Understanding of JWT authentication systems is a plus.",
            "location": "Galle, Sri Lanka",
            "offers": "International exposure, competitive salary package, health insurance, and hybrid work flexibility.",
            "is_remote": False,
            "skills": [{"name": "Go (Golang)"}, {"name": "REST APIs"}, {"name": "Microservices"}]
        }

        # List 2 contains: 1 duplicate item from Session 1, and 2 completely fresh job listings
        session_2_job_list = [
            job_a_reposted_duplicate,
            job_d_new,
            job_e_new,
            job_f_new
        ]

        print(f"   -> Iterating over {len(session_2_job_list)} collected job listings...")
        for incoming_job in session_2_job_list:
            sig, lsh = generate_production_minhash_and_lsh(incoming_job)
            
            incoming_job["meta_data"] = {
                "geo": {"province": "Western"},
                "posted_at": datetime.utcnow().isoformat() + "Z",
                "source": "Direct Careers Page",
                "standardized_category": "Software Engineering & Technology",
                "seniority": "Entry-Level",
                "confidence_score": 0.96,
                "ai_version": "v1.2.0",
                "error": False,
                "crawler_run_id": session_2_id,
                "minhash_signature": sig
            }
            
            formatted_lsh = [{"band_no": item["band_no"], "bucket_key": item["bucket_key"]} for item in lsh]

            # Request duplicate check evaluation from Go
            is_dup, matched_id = await check_duplicate_via_backend(client, formatted_lsh, sig, incoming_location=incoming_job.get("location", ""))
            
            if is_dup:
                # 🎯 DYNAMIC PATH A: Jaccard value threshold hit! Execute batch-update
                print(f"      ⚠️  Duplicate Spotted! [ {incoming_job['employer']} - {incoming_job['job_role']} ]. Linked to Job ID {matched_id}. Executing Batch-Update...")
                assert matched_id is not None
                
                update_payload = {
                    "duplicates": [
                        {
                            "job_post_id": int(matched_id),
                            "crawler_run_id": session_2_id
                        }
                    ]
                }
                update_res = await client.post(f"{BACKEND_BASE_URL}/jobs/batch-update", json=update_payload)
                assert update_res.status_code in (200, 204)
            else:
                # 🎯 DYNAMIC PATH B: Fresh item entry. Execute batch-save
                print(f"      🔹 New Unique [ {incoming_job['employer']} - {incoming_job['job_role']} ]. Executing Batch-Save...")
                save_res = await client.post(f"{BACKEND_BASE_URL}/jobs/batch-save", json={
                    "new_jobs": [incoming_job],
                    "lsh_indexes": formatted_lsh
                })
                assert save_res.status_code in (200, 201)

        # Reconcile Session 2
        print("   -> Running data reconciliation sync for Session 2...")
        reconcile_res_2 = await client.post(f"{BACKEND_BASE_URL}/jobs/reconcile", json={"crawler_run_id": session_2_id})
        assert reconcile_res_2.status_code == 200

        # Complete Session 2
        close_res_2 = await client.post(f"{BACKEND_BASE_URL}/runs/{session_2_id}/complete", json={"status": "COMPLETED"})
        assert close_res_2.status_code in (200, 201)
        print("🏁 [SESSION 2] Concluded successfully. Pipeline looping assertions verified.")






# docker compose run --rm \
#   -e BACKEND_URL="http://backend:8080/api/v1/crawler" \
#   --entrypoint "pytest test_cross_session_deduplication.py -v -s" \