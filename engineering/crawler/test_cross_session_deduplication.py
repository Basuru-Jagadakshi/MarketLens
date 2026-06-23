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

@pytest.mark.asyncio
async def test_cross_session_deduplication_lifecycle():
    """
    Validates state persistence across distinct execution chronologies using the
    exact production job JSON schema and invoking the reconciliation layer.
    """
    async with httpx.AsyncClient(timeout=15.0) as client:
        
        ad_url_session_1 = "https://ikman.lk/en/ad/wso2-java-backend-run-alpha"
        ad_url_session_2 = "https://ikman.lk/en/ad/wso2-java-backend-run-beta"

        # =====================================================================
        # CRON TASK RUN 1: SEEDING THE ARCHIVE STATE
        # =====================================================================
        print("\n🚀 [SESSION 1] Starting first isolated crawler execution run...")
        
        # 1. Initialize First Session Tracker
        session_1_payload = {
            "started_at": "2026-06-14T14:35:10.124578+05:30",
            "finished_at": None,
            "status": "RUNNING"
        }
        init_res_1 = await client.post(f"{BACKEND_BASE_URL}/runs", json=session_1_payload)
        assert init_res_1.status_code in (200, 201)
        
        # Ensure ID parsing matches your Go integer allocation format
        session_1_id = int(init_res_1.json().get("id"))
        print(f"   -> Session 1 Active ID: {session_1_id}")

        # 2. Build exact production JSON structure for Job 1
        job_data_1 = {
            "employer": "Global Logistics Lanka Ltd",
            "job_role": "Operations Coordinator",
            "job_type": {
                "name": "Full Time"
            },
            "key_responsibilities": "We are looking for a detail-oriented Operations Coordinator to manage supply chain documentation and warehouse schedules. Responsibilities include tracking shipments, communicating with logistics partners, and maintaining inventory records. Candidates must possess strong organizational skills and familiarity with warehouse management software.",
            "qualifications": "BSc in Computer Science or equivalent engineering field. Foundational knowledge of Go (Golang) or Java with Spring Boot. Understanding of JWT authentication systems is a plus.",
            "location": "Gampaha, Sri Lanka",
            "offers": "International exposure, competitive salary package, health insurance, and hybrid work flexibility.",
            "is_remote": False,
            "skills": [
                { "name": "Go (Golang)" },
                { "name": "REST APIs" },
                { "name": "Microservices" }
            ]
        }
        
        sig_1, lsh_1 = generate_production_minhash_and_lsh(job_data_1)
        
        # Enforce exact production metadata block placement
        job_data_1["meta_data"] = {
            "geo": {
                "province": "Western"
            },
            "posted_at": datetime.utcnow().isoformat() + "Z",
            "source": "Direct Careers Page",
            "standardized_category": "Software Engineering & Technology",
            "seniority": "Entry-Level",
            "confidence_score": 0.97,
            "ai_version": "v1.2.0",
            "error": False,
            "crawler_run_id": session_1_id,
            "minhash_signature": sig_1  # Production format integer array [5541, 1029, ...]
        }
        
        formatted_lsh_1 = [
            {"band_no": item["band_no"], "bucket_key": item["bucket_key"]}
            for item in lsh_1
        ]

        job_data_1_saving_json = {
        "employer": {
            "name": "Global Logistics Lanka Ltd"
        },
        "job_role": "Operations Coordinator",
        "job_type": {
            "type": "Full Time"
        },
        "job_description": "We are looking for a detail-oriented Operations Coordinator to manage supply chain documentation and warehouse schedules. Responsibilities include tracking shipments, communicating with logistics partners, and maintaining inventory records. Candidates must possess strong organizational skills and familiarity with warehouse management software.",
        "location": "Gampaha, Sri Lanka",
        "is_remote": False,
        "meta_data": {
            "geo_data": {
            "province": "Western"
            },
            "industry": {
            "name": "Transportation and Storage"
            },
            "occupation": {
            "name": "Clerks"
            },
            "education_level": {
            "level": "GCE A/L"
            },
            "experience": {
            "name": "Experience Required"
            },
            "source": {
            "source": "Ikman"
            },
            "crawler_run_id": session_1_id,
            "ai_version": {
            "version": "deepseek-v1"
            },
            "posted_at": "2026-06-22T09:30:00Z",
            "confidence_score": 0.92,
            "minhash_signature": sig_1
        },
        "skills": [
            { "skill": "Logistics Management" },
            { "skill": "Supply Chain Operations" },
            { "skill": "Inventory Control" },
            { "skill": "Microsoft Excel" }
        ]
        }
        
        # Persist base job data
        save_res_1 = await client.post(f"{BACKEND_BASE_URL}/jobs/batch-save", json={
            "new_jobs": [job_data_1_saving_json], 
            "lsh_indexes": formatted_lsh_1
        })
        assert save_res_1.status_code in (200, 201), "Session 1 baseline persistence failing."
        print("   ✅ Baseline job cached securely using production JSON schema format.")

        # 3. Reconcile Session 1 before shutdown
        print("   -> Invoking Session 1 Data Reconciliation sync...")
        reconcile_res_1 = await client.post(f"{BACKEND_BASE_URL}/jobs/reconcile", json={"crawler_run_id": session_1_id})
        assert reconcile_res_1.status_code == 200, "Reconciliation failed for Session 1"

        # 4. Cleanly close down session 1
        close_res_1 = await client.post(f"{BACKEND_BASE_URL}/runs/{session_1_id}/complete", json={
            "status": "COMPLETED"
        })
        assert close_res_1.status_code in (200, 201)
        print("🏁 [SESSION 1] Terminated cleanly. State successfully locked on disk.")


        # =====================================================================
        # CRON TASK RUN 2: THE CROSS-SESSION COLLISION DETECTION
        # =====================================================================
        print("\n🚀 [SESSION 2] Starting next sequential crawler execution run...")
        
        # 1. Initialize Second Session Tracker
        session_2_payload = {
            "started_at": "2026-06-15T10:15:00.654321+05:30",
            "finished_at": None,
            "status": "RUNNING"
        }
        init_res_2 = await client.post(f"{BACKEND_BASE_URL}/runs", json=session_2_payload)
        assert init_res_2.status_code in (200, 201)
        session_2_id = int(init_res_2.json().get("id"))
        print(f"   -> Session 2 Active ID: {session_2_id}")

        # 2. Build a near-duplicate job payload with minor text additions to hit the LSH window
        job_data_2 = job_data_1.copy()
        job_data_2["key_responsibilities"] = (
            "We are looking for a detail oriented Operations Coordinator to manage supply chain documentation, warehouse schedules. Responsibilities include tracking shipments, communicating with logistics partners and maintaining inventory records. Candidates must possess strong organizational skills and familiarity with warehouse management software."
        )
        
        sig_2, lsh_2 = generate_production_minhash_and_lsh(job_data_2)
        
        # Assign meta_data for the new session run
        job_data_2["meta_data"] = {
            "geo": {"province": "Western"},
            "posted_at": datetime.utcnow().isoformat() + "Z",
            "source": "Direct Careers Page",
            "standardized_category": "Software Engineering & Technology",
            "seniority": "Entry-Level",
            "confidence_score": 0.95,
            "ai_version": "v1.2.0",
            "error": False,
            "crawler_run_id": session_2_id,
            "minhash_signature": sig_2
        }

        # 3. Request Radar lookup check from Go backend
        print("   -> Evaluating duplicate check against historical logs...")
        is_duplicate, matched_parent_id = await check_duplicate_via_backend(client, lsh_2, sig_2, job_data_2.get("location", ""))
        
        assert is_duplicate is True, "Cross-session deduplication failed! System missed the signature collision."
        assert matched_parent_id is not None, "Failed to retrieve parent relationship pointer mapping cross sessions."
        print(f"   ✅ Success: Historical duplicate caught! Linked back to Root Job ID: {matched_parent_id}")

        # 4. Trigger the UPDATE endpoint payload logic to capture the duplication relationship
        update_payload = {
            "duplicates": [
                {
                    "job_post_id": int(matched_parent_id),
                    "crawler_run_id": session_2_id
                }
            ]
        }
        
        update_res = await client.post(f"{BACKEND_BASE_URL}/jobs/batch-update", json=update_payload)
        assert update_res.status_code in (200, 204), f"Failed to record duplicate connection mapping link: {update_res.text}"
        print(f"   ✅ Update lifecycle state record executed cleanly on active DB mapping keys.")

        # 5. Reconcile Session 2 tracking states
        print("   -> Invoking Session 2 Data Reconciliation sync...")
        reconcile_res_2 = await client.post(f"{BACKEND_BASE_URL}/jobs/reconcile", json={"crawler_run_id": session_2_id})
        assert reconcile_res_2.status_code == 200, "Reconciliation failed for Session 2"

        # 6. Cleanly close down session 2
        close_res_2 = await client.post(f"{BACKEND_BASE_URL}/runs/{session_2_id}/complete", json={
            "status": "COMPLETED", 
        })
        assert close_res_2.status_code in (200, 201)
        print("🏁 [SESSION 2] Lifecycle terminated cleanly. Multi-session pipeline validation finalized.")






# docker compose run --rm \ -e BACKEND_URL="http://backend:8080/api/v1/crawler" \ --entrypoint "pytest test_cross_session_deduplication.py -v -s" \
