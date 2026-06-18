import os
import json
import httpx
import asyncio
import logging
from typing import List, Dict, Any

logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s - %(message)s")
logger = logging.getLogger(__name__)

# Route directly to your backend Go application container on the bridge network
BACKEND_BASE_URL = os.getenv("BACKEND_URL", "http://backend:8080/api/v1/crawler")

from orchestrator.ikman_orchestrator import (
    generate_production_minhash_and_lsh,
    check_duplicate_via_backend
)

# =====================================================================
# SEED CONFIGURATION: 3 HIGHLY DIVERSE TARGET BASELINE RECORDS
# =====================================================================
SEED_BASE_JOBS = [
    {
        "url_ref": "https://ikman.lk/en/ad/wso2-backend-engineer-500",
        "payload": {
            "employer": "WSO2",
            "job_role": "Software Engineer - Customer Success",
            "job_type": {"name": "Full Time"},
            "key_responsibilities": """Job Summary

                Do you get a thrill out of watching software unfold before your eyes? Do you dream about code every night? If so, we’d love to talk to you about your career with us. We’re looking for top-notch Software Engineers who would see technical glitches as an enjoyable challenge and are willing to walk extra mile to see a project through to completion.


                Responsibilities & Duties

                Ideal candidates should have working competency in the technology domains, programming languages,  and OOAD.
                Ability to design and implement solutions adhering to overall architecture and system design goals including performance, security, scalability, quality of code, etc.
                Think of all possible scenarios and the ‘big picture’ when implementing some functionality.
                Ability to estimate effort on functional areas worked on. On time delivery.
                Ability to identify user stories for the product and document them accordingly following the process.
                Proactively own the functional areas of the product you work on. Own all other aspects of the product including the below.
                Research on functional and technical improvements.
                Ability to communicate clearly, articulate both on written and verbal communication.
                Ability to conduct product demos, trainings, and presentations.
                Ability to successfully contribute to technical and non-technical discussions on email and in-person.

                Qualifications & Skills

                Fresh graduates with BSc in Computer Science/Engineering or Equivalent or with a minimum of 1-2 years industry experience
                Strong development skills and proficiency in at least one programming language. Having experience in Java, C# or C/C++ will be an added advantage.
                Strong analytical skills
                Experience and knowledge on Distributed Systems is an added advantage

                In Addition to a Competitive Compensation Package, WSO2 Offers:

                A flexible vacation policy, as long as you have demonstrated your commitment to shared goals.
                You are entitled to Health Benefits, as per the Company Healthcare plan provided by WSO2.
                A Laptop as per standard company specifications.
                Additional equipment on request (Mouse, Keyboard, Monitor)""",
            "qualifications": "BSc in Computer Science or equivalent engineering field. Foundational knowledge of Go (Golang) or Java with Spring Boot. Understanding of JWT authentication systems is a plus.",
            "location": "Colombo 03, Sri Lanka",
            "offers": "International exposure, competitive salary package, health insurance, and hybrid work flexibility.",
            "is_remote": False,
            "meta_data": {
                "geo": {"province": "Western"},
                "posted_at": "2026-06-12T11:45:00Z",
                "source": "Direct Careers Page",
                "standardized_category": "Software Engineering & Technology",
                "seniority": "Entry-Level",
                "confidence_score": 0.97,
                "ai_version": "v1.2.0",
                "error": False
            },
            "skills": [{"name": "Go (Golang)"}, {"name": "REST APIs"}, {"name": "Microservices"}]
        }
    },
    {
        "url_ref": "https://ikman.lk/en/ad/cinnamon-hr-executive-600",
        "payload": {
            "employer": "Cinnamon Hotels & Resorts",
            "job_role": "HR Recruiting Executive",
            "job_type": {"name": "Full Time"},
            "key_responsibilities": "Manage undergraduate engineering internship tracks, coordinate cross-department technical interviews, and oversee employee onboarding schedules.",
            "qualifications": "Degree or Diploma in Human Resources Management, strong organizational capability, and exceptional communication skills.",
            "location": "Colombo, Sri Lanka",
            "offers": "Competitive base package, free corporate meals, dynamic hospitality environment exposure.",
            "is_remote": False,
            "meta_data": {
                "geo": {"province": "Western"},
                "posted_at": "2026-06-11T09:00:00Z",
                "source": "TopJobs",
                "standardized_category": "Human Resources",
                "seniority": "Associate",
                "confidence_score": 0.94,
                "ai_version": "v1.2.0",
                "error": False
            },
            "skills": [{"name": "Talent Acquisition"}, {"name": "Sourcing"}, {"name": "Onboarding"}]
        }
    },
    {
        "url_ref": "https://ikman.lk/en/ad/dialog-ui-designer-700",
        "payload": {
            "employer": "Dialog Axiata",
            "job_role": "UI FrontEnd Engineer",
            "job_type": {"name": "Contract"},
            "key_responsibilities": "Develop responsive modern web applications using React, NextJS, and Tailwind CSS. Optimize frontend dashboard components for maximum rendering speed and cross-browser consistency.",
            "qualifications": "2+ years web engineering experience, high proficiency in modern JavaScript/TypeScript framework development ecosystems.",
            "location": "Technical Hub, Colombo",
            "offers": "Premium tech equipment, medical allowance insurance access, structured milestone bonuses.",
            "is_remote": True,
            "meta_data": {
                "geo": {"province": "Western"},
                "posted_at": "2026-06-13T15:30:00Z",
                "source": "Ikman",
                "standardized_category": "Software Engineering & Technology",
                "seniority": "Mid-Level",
                "confidence_score": 0.98,
                "ai_version": "v1.2.0",
                "error": False
            },
            "skills": [{"name": "React"}, {"name": "TypeScript"}, {"name": "Tailwind CSS"}]
        }
    }
]

# =====================================================================
# EVALUATION MATRIX: SCENARIO PROFILES TESTING RADAR AND JACCARD CAPABILITIES
# =====================================================================
TEST_SCENARIOS = [
    {
        "name": "Strict Near-Duplicate Check (Expected: DUPLICATE of Backend Job)",
        "payload": {
            "employer": "WSO2",
            "job_role": "Software Engineer - Customer Success",
            "location": "Colombo 03, Sri Lanka",
            "key_responsibilities": """Job Summary

                Do you get a thrill out of watching software unfold before your eyes? Do you dream about code every night? If so, we’d love to talk to you about your career with us. We’re looking for top-notch Software Engineers who would see technical glitches as an enjoyable challenge and are willing to walk extra mile to see a project through to completion.


                Responsibilities & Duties

                Ideal candidates should have working competency in the technology domains, programming languages,  and OOAD.
                Ability to design and implement solutions adhering to overall architecture and system design goals including performance, security, scalability, quality of code, etc.
                Think of all possible scenarios and the ‘big picture’ when implementing some functionality.
                Ability to estimate effort on functional areas worked on. On time delivery.
                Ability to identify user stories for the product and document them accordingly following the process.
                Proactively own the functional areas of the product you work on. Own all other aspects of the product including the below.
                Research on functional and technical improvements.
                Ability to communicate clearly, articulate both on written and verbal communication.
                Ability to conduct product demos, trainings, and presentations.
                Ability to successfully contribute to technical and non-technical discussions on email and in-person.

                Qualifications & Skills

                Fresh graduates with BSc in Computer Science/Engineering or Equivalent or with a minimum of 1-2 years industry experience
                Strong development skills and proficiency in at least one programming language. Having experience in Java, C# or C/C++ will be an added advantage.
                Strong analytical skills
                Experience and knowledge on Distributed Systems is an added advantage

                In Addition to a Competitive Compensation Package, WSO2 Offers:
                 A flexible vacation policy, as long as you have demonstrated your commitment to shared goals.
                You are entitled to Health Benefits, as per the Company Healthcare plan provided by WSO2.
                A Laptop as per standard company specifications.
""",
        },
        "expected_duplicate": True
    },
     {
        "name": "Different locations",
        "payload": {
            "employer": "WSO2",
            "job_role": "Software Engineer - Customer Success",
            "location": "Jaffna",
            "key_responsibilities": """Job Summary

                Do you get a thrill out of watching software unfold before your eyes? Do you dream about code every night? If so, we’d love to talk to you about your career with us. We’re looking for top-notch Software Engineers who would see technical glitches as an enjoyable challenge and are willing to walk extra mile to see a project through to completion.


                Responsibilities & Duties

                Ideal candidates should have working competency in the technology domains, programming languages,  and OOAD.
                Ability to design and implement solutions adhering to overall architecture and system design goals including performance, security, scalability, quality of code, etc.
                Think of all possible scenarios and the ‘big picture’ when implementing some functionality.
                Ability to estimate effort on functional areas worked on. On time delivery.
                Ability to identify user stories for the product and document them accordingly following the process.
                Proactively own the functional areas of the product you work on. Own all other aspects of the product including the below.
                Research on functional and technical improvements.
                Ability to communicate clearly, articulate both on written and verbal communication.
                Ability to conduct product demos, trainings, and presentations.
                Ability to successfully contribute to technical and non-technical discussions on email and in-person.

                Qualifications & Skills

                Fresh graduates with BSc in Computer Science/Engineering or Equivalent or with a minimum of 1-2 years industry experience
                Strong development skills and proficiency in at least one programming language. Having experience in Java, C# or C/C++ will be an added advantage.
                Strong analytical skills
                Experience and knowledge on Distributed Systems is an added advantage

                In Addition to a Competitive Compensation Package, WSO2 Offers:

                A flexible vacation policy, as long as you have demonstrated your commitment to shared goals.
                You are entitled to Health Benefits, as per the Company Healthcare plan provided by WSO2.
                A Laptop as per standard company specifications.
                Additional equipment on request (Mouse, Keyboard, Monitor)""",
        },
        "expected_duplicate": False
    },
    {
        "name": "Same Company, Completely Different Field (Expected: UNIQUE - Prevents False Positives)",
        "payload": {
            "employer": "WSO2",
            "job_role": "Senior Software Engineer - APIM",
            "location": "Colombo 03, Sri Lanka",
            "key_responsibilities": """Job Summary:
WSO2 is looking for a highly motivated individual who wants to pursue a career in a fast-paced environment.


The ideal candidate must have the ability to prioritize well, communicate clearly, possess a consistent track record of delivery, and excellent software engineering skills are mandatory. Creativity bundled with high quality and a customer focus are very important and you must be able to work across multiple facets of the project and juggle multiple responsibilities.


Strong analytic capability and the ability to create innovative solutions are prerequisites for this position. Sounds interesting? Then we’d like to hear from you.


Responsibilities and Duties:
Working competency in the technology domains, programming languages, and OOAD.
Ability to design and implement solutions adhering to overall architecture and system design goals including performance, security, scalability, quality of code, etc.
Think of all possible scenarios and the ‘big picture’ when implementing some functionality, ability to estimate effort on functional areas worked on, and deliver on time.
Ability to identify user stories for the product and document them accordingly following the processes. Follow the WSO2 development process end-to-end when developing components / features for products i.e. coding best practices, patterns, unit testing, automated testing, documentation, etc.
Proactively own the functional areas of the product you work on and other aspects such as:
Marketing (blogs, social media, helping out with marketing campaigns)
Pre-sales (product demos)
Sales (anticipate future customer requirements and account expansion insights)
Documentation
Community engagement (answer questions on Stack Overflow)
Delivery and support (monitor and help with support issues, patches, etc).
Research on functional and technical improvements; introduce new ideas on how to improve the product and overall technical designs.
Ability to give technical leadership to a small team. Delegate and follow-up on work assigned. Mentor junior members in engineering best practices and processes.
Keep customers informed in a timely manner. Practice communication with empathy and be proactive in communicating.
Ability to communicate clearly, articulate both on written and verbal communication, conduct product demos, training, and presentations.
Ability to successfully contribute to technical and non-technical discussions on email and in-person.
Be an advocate for communication practices in open source development and successfully practice it.
Engage with the extended teams such as documentation, pre-sales, marketing, and sales on product and customer related activities.

Requirements:
BSc in Computer Science with Minimum 3-5 years of experience.
Strong development skills and proficiency in at least one programming language. Having experience in Java, C# or C/C++ will be an added advantage.
Excellent technical design skills.
Distributed computing skills will be an advantage.
The passion to learn and excel in the software engineering landscape.
Knowledge in project estimation techniques, design patterns, and performance engineering.
Other benefits:
We want you to enjoy being part of our team and feel happy and cared for while at WSO2; hence, in addition to the performance-driven commission, WSO2 offers you a host of wonderful benefits and facilities.
A flexible vacation policy, as long as you have demonstrated your commitment to shared goals.
Healthcare.
You are entitled to Health Benefits, as per the Company Healthcare plan provided by WSO2.
Additional Support.
A Laptop as per standard company specifications.
Additional equipment on request (Mouse, Keyboard, Monitor)""",
        },
        "expected_duplicate": False
    },
    {
        "name": "Similar Role, Different Company Context (Expected: UNIQUE - Prevents Cross-Matching)",
        "payload": {
            "employer": "Dialog Axiata",
            "job_role": "Junior Backend Developer",
            "location": "Technical Hub, Colombo",
            "key_responsibilities": "Design and implement robust microservices patterns, maintain API gateways, write comprehensive unit tests, and collaborate with open-source project documentation streams.",
        },
        "expected_duplicate": False
    },
    {
        "name": "Minor Content Shift Check (Expected: DUPLICATE of FrontEnd Job)",
        "payload": {
            "employer": "Dialog Axiata",
            "job_role": "UI FrontEnd Engineer",
            "location": "Technical Hub, Colombo",
            "key_responsibilities": "Develop responsive modern web applications using React, NextJS, Tailwind CSS. Optimize frontend dashboard components for maximum rendering speed, cross-browser consistency."
        },
        "expected_duplicate": True
    },
    {
        "name": "Completely Unseen Industry Post (Expected: UNIQUE)",
        "payload": {
            "employer": "Keells Super",
            "job_role": "Inventory Stock Supervisor",
            "location": "Kandy, Sri Lanka",
            "key_responsibilities": "Oversee grocery supply logistics lines, coordinate warehouse delivery trucks, and manage shelf layout plans.",
        },
        "expected_duplicate": False
    }
]


async def run_matrix_verification():
    logger.info("⚡ Booting up Closed-Network Matrix Sandbox Integration Test Suite...")
    async_client = httpx.AsyncClient(timeout=15.0)
    
    try:
        new_jobs_buffer = []
        lsh_index_buffer = []
        
        # 1. COMPILE STRUCUTRED MATRIX SEEDS & ENFORCE THE DATA RELATIONSHIPS
        for item in SEED_BASE_JOBS:
            job_node = item["payload"]
            target_url = item["url_ref"]
            
            # Generate local algorithmic values
            sig, lsh = generate_production_minhash_and_lsh(job_node)
            
            # Form structural metadata layout according to Go JSON expectations
            job_node["meta_data"]["crawler_run_id"] = 1
            job_node["meta_data"]["minhash_signature"] = sig
            new_jobs_buffer.append(job_node)
            
            # Attach structural parent linkage identifiers into sub-index entities
            for lsh_entry in lsh:
                lsh_index_buffer.append({
                    "band_no": lsh_entry["band_no"],
                    "bucket_key": lsh_entry["bucket_key"]
                })
                
        # Fire database seed tracking initialization execution
        logger.info(f"Seeding {len(new_jobs_buffer)} valid baseline configurations into DB endpoints...")
        seed_payload = {
            "new_jobs": new_jobs_buffer,
            "lsh_indexes": lsh_index_buffer
        }
        
        seed_res = await async_client.post(f"{BACKEND_BASE_URL}/jobs/batch-save", json=seed_payload)
        if seed_res.status_code not in (200, 201):
            logger.error(f"❌ Migration seeding script aborted by backend API status: {seed_res.status_code}")
            logger.error(f"Response data log context: {seed_res.text}")
            return
            
        logger.info("✅ Mock matrix database state baseline loaded successfully.")

        # 2. EVALUATE MATRIX TEST CRITERIA
        failed_assertions = 0
        
        for index, scenario in enumerate(TEST_SCENARIOS, start=1):
            logger.info(f"\n--- Running Evaluation Scenario {index}: {scenario['name']} ---")
            
            # Process lookup parameters
            target_sig, target_lsh = generate_production_minhash_and_lsh(scenario["payload"])
            
            # Challenge the Go backend radar service
            is_duplicate, matched_id = await check_duplicate_via_backend(async_client, target_lsh, target_sig, incoming_location=scenario["payload"].get("location", ""))
            
            logger.info(f"   -> API Evaluation Result: Is Duplicate = {is_duplicate} | Matched DB ID = {matched_id}")
            
            # Validate output matches expectations
            if is_duplicate != scenario["expected_duplicate"]:
                logger.error(f"   ❌ CRITICAL ASSERTION MISMATH: Expected outcome '{scenario['expected_duplicate']}', received '{is_duplicate}' instead.")
                failed_assertions += 1
            else:
                logger.info("   ✅ Success: System accurately resolved algorithm classification rules.")

        # 3. OUTPUT FINAL VERDICT LOGS
        print("\n==========================================================================")
        if failed_assertions == 0:
            print("🎉 INTEGRATION TEST SUITE PASSED SUCCESSFULLY!")
            print("Your MinHash calculations, LSH relations, and Go struct contracts match 100%.")
            print("You are clear to share this output execution layout with Vibhatha.")
        else:
            print(f"⚠️ TEST SUITE REJECTED: {failed_assertions} validation checks failed processing rules.")
        print("==========================================================================")

    except Exception as e:
        logger.error(f"Matrix run suite terminated due to connection/system framework failure: {e}")
    finally:
        await async_client.aclose()


if __name__ == "__main__":
    asyncio.run(run_matrix_verification())




# docker compose run --rm \
#   -e BACKEND_URL="http://backend:8080/api/v1/crawler" \
#   --entrypoint "python" \
#   crawler test_matrix_integration.py
