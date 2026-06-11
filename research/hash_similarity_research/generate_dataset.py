#python3 generate_dataset.py

import json

base_jobs = [
    {
        "id": "JOB_01", "title": "Senior Java Developer", "employer": "WSO2",
        "description": "We are looking for a Senior Java Developer with Spring Boot and microservices experience to join our cloud platform team. You will design scalable APIs and optimize database queries."
    },
    {
        "id": "JOB_02", "title": "Data Engineer", "employer": "Sysco Labs",
        "description": "Seeking a Data Engineer expert in Apache Spark, Python, and AWS data pipelines. You will build scalable ETL pipelines, manage data warehousing solutions, and optimize complex SQL analytical queries."
    },
    {
        "id": "JOB_03", "title": "Golang Backend Engineer", "employer": "Axiata Digital",
        "description": "Join our core platform division building high-throughput microservices in Go. Experience with the Gin framework, gRPC communication protocols, Docker containerization, and MongoDB is strictly required."
    },
    {
        "id": "JOB_04", "title": "Frontend React Developer", "employer": "Mitra Innovation",
        "description": "Looking for a Frontend UI specialist skilled in React, TypeScript, and Tailwind CSS. You will transform complex user flows into interactive web dashboards, optimize bundle sizes, and manage global application state."
    },
    {
        "id": "JOB_05", "title": "DevOps Cloud Engineer", "employer": "Virtusa",
        "description": "Expand our infrastructure automation via Terraform, AWS, and Kubernetes. You will design resilient CI/CD pipelines using GitHub Actions, manage cluster deployments, and enforce security policies."
    },
    {
        "id": "JOB_06", "title": "QA Automation Engineer", "employer": "99x",
        "description": "Build comprehensive automated test coverage suites using Selenium, Cypress, and Postman. Implement quality assurance gate metrics into delivery cycles and track regression defects systematically."
    },
    {
        "id": "JOB_07", "title": "Machine Learning Engineer", "employer": "Codegen",
        "description": "Deploy advanced NLP and predictive intelligence models into production frameworks. Strong experience with Python, PyTorch, Scikit-Learn, data preprocessing, and model optimization is necessary."
    },
    {
        "id": "JOB_08", "title": "Mobile iOS Developer", "employer": "Calcey Technologies",
        "description": "Craft high-performance native iOS mobile applications using Swift and SwiftUI. Implement clean architectural patterns like MVVM, integrate core RESTful web sockets, and handle local CoreData caching systems."
    },
    {
        "id": "JOB_09", "title": "Cybersecurity Analyst", "employer": "LSEG",
        "description": "Monitor corporate network architecture security boundaries for vulnerabilities. Experience with penetration testing, SIEM logging tools, identity management protocols, and executing incident response plans is required."
    },
    {
        "id": "JOB_10", "title": "Scrum Product Owner", "employer": "Pearson",
        "description": "Responsible for translating stakeholder business requirements into clear technical product backlogs. You will manage agile sprint planning schedules, coordinate daily standups, and direct product feature delivery roadmap targets."
    }
]

alt_titles = ["UI/UX Designer", "Product Owner", "QA Lead", "System Admin", "Data Scientist", "Support Engineer", "Scrum Master", "Android Developer", "Network Engineer", "Technical Writer"]
alt_employers = ["Dialog Axiata", "IFS", "Aeturnum", "Fortude", "Arimac", "Cambio", "SimCentric", "Zone24x7", "Eyepax", "Synopsys"]
alt_descriptions = [
    "Focus on responsive layout frameworks, user persona map journeys, low-fidelity wireframes, and managing component design libraries using Figma.",
    "Draft enterprise feature workflows, calculate resource budget projections, prioritize client tickets, and maintain long-term corporate product strategy sheets.",
    "Oversee release cycle compliance strategies, audit load metrics testing configurations, run black-box exploratory assessments, and direct engineer assignments.",
    "Configure Linux internal network architectures, manage virtualization boundaries, patch server infrastructure setups, and execute data storage failovers.",
    "Analyze predictive business KPIs using Jupyter, R scripts, Tableau visual dashboards, and extract insights from unstructured enterprise tables.",
    "Triage customer bug escalations, monitor active application error tracing systems, issue software patches, and draft release documentation summaries.",
    "Facilitate team velocity optimizations, remove cross-departmental project blockers, run post-sprint retrospective logs, and verify agile framework adoptions.",
    "Build clean Android applications using Kotlin and Jetpack Compose. Implement background worker tasks, manage threading pools, and optimize local SQLite storage layers.",
    "Configure enterprise Cisco firewall routers, maintain active virtual private network channels, audit hardware rack allocations, and debug packet routing switches.",
    "Author API markdown integration blueprints, draft internal system architecture wiki reference material, and maintain user onboarding instruction documentation sets."
]

all_scraped_records = []

for idx, base in enumerate(base_jobs):
    cases = [
        {"case_id": "Case 1", "name": "Completely same jobs (Exact word match)", "title": base["title"], "employer": base["employer"], "description": base["description"]},
        {"case_id": "Case 2", "name": "Job title same but different employer and description", "title": base["title"], "employer": alt_employers[idx], "description": alt_descriptions[idx]},
        {"case_id": "Case 3", "name": "Job employer same but different title and description", "title": alt_titles[idx], "employer": base["employer"], "description": alt_descriptions[idx]},
        {"case_id": "Case 4", "name": "Same description but different employer and title", "title": alt_titles[idx], "employer": alt_employers[idx], "description": base["description"]},
        {"case_id": "Case 5", "name": "Job title and employer same but different description", "title": base["title"], "employer": base["employer"], "description": alt_descriptions[idx]},
        {"case_id": "Case 6", "name": "Employer and description same but different title", "title": alt_titles[idx], "employer": base["employer"], "description": base["description"]},
        {"case_id": "Case 7", "name": "Title and description same but different employer", "title": base["title"], "employer": alt_employers[idx], "description": base["description"]},
        {"case_id": "Case 8", "name": "All attributes completely different", "title": alt_titles[idx], "employer": alt_employers[idx], "description": alt_descriptions[idx]},
        {"case_id": "Case 9", "name": "Title & Employer same, Description highly similar (Near-Duplicate)", "title": base["title"], "employer": base["employer"], "description": base["description"].replace("We are looking for", "Our growing company is actively seeking") + " Splitting text slightly changes some words here."}
    ]
    
    all_scraped_records.append({
        "base_id": base["id"], "base_title": base["title"], "base_employer": base["employer"], "base_description": base["description"],
        "permutations": cases
    })

with open("scraped_job_market.json", "w", encoding="utf-8") as f:
    json.dump(all_scraped_records, f, indent=4)

print("Generated updated 'scraped_job_market.json' containing 10 base fields and 90 research cases.")