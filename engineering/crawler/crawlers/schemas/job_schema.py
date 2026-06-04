JOB_EXTRACTION_SCHEMA = {
    "type": "object",
    "properties": {
        "employer": {"type": "string"},
        "job_role": {"type": "string"},
        "job_type": {
            "type": "object",
            "properties": {
                "name": {"type": "string"}
            },
            "required": ["name"]
        },
        "key_responsibilities": {"type": "string"},
        "qualifications": {"type": "string"},
        "location": {"type": "string"},
        "offers": {"type": "string"},
        "is_remote": {"type": "boolean"},
        "meta_data": {
            "type": "object",
            "properties": {
                "geo": {
                    "type": "object",
                    "properties": {
                        "province": {"type": "string"}
                    },
                    "required": ["province"]
                },
                "posted_at": {
                    "type": "string",
                    "pattern": r"^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$"
                },
                "source": {"type": "string"},
                "standardized_category": {"type": "string"},
                "seniority": {"type": "string"},
                "confidence_score": {"type": "number"},
                "ai_version": {"type": "string"},
                "error": {"type": "boolean"}
            },
            "required": ["geo", "source", "standardized_category", "error"]
        },
        "skills": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "name": {"type": "string"}
                },
                "required": ["name"]
            }
        }
    },
    "required": ["employer", "job_role", "location", "job_type", "meta_data", "skills"]
}



BASE_JOB_INSTRUCTION = """
Extract structural data from the raw job text into the specified JSON format.
CRITICAL: If a field is missing, use the strict default provided. Do not hallucinate.

Fields & Defaults:
1. 'employer': Company name. Default: "".
2. 'job_role': Vacancy title. Default: "".
3. 'job_type': {"name": "Full Time"} or {"name": "Part Time"}. Default: {"name": "Full Time"}.
4. 'key_responsibilities': Brief summary text of tasks. Default: "".
5. 'qualifications': Background/skills required text. Default: "".
6. 'location': City/country (e.g. "Colombo, Sri Lanka"). Default: "Sri Lanka".
7. 'offers': Perks/salary text. Default: "".
8. 'is_remote': true ONLY if explicit remote wording exists, else false.
9. 'meta_data': Object containing:
   - 'geo': {"province": "Western"} or other inferred Sri Lankan province. Default: {"province": "Western"}.
   - 'posted_at': Current ISO timestamp.
   - 'source': Crawled source (e.g., "TopJobs", "Ikman").
   - 'standardized_category': Match exactly ONE value from this list, else default to "Other": ["Technology & Telecom","Sales, Marketing & Retail","Finance & Banking","Corporate, HR & Legal","Customer Support & Service","Logistics, Manufacturing & Trade","Hospitality, Food & Tourism","Healthcare & Education","Services, Driving & Blue-Collar","Creative, Media & Other"]
   - 'seniority': Match exactly ONE: ["Internship","Trainee","Entry-Level","Mid-Level","Senior","Not Specified"]. Default: "Not Specified".
   - 'confidence_score': Float between 0.00 and 1.00.
   - 'ai_version': "v1.2.0".
   - 'error': false unless text is corrupted/unparseable.
10. 'skills': Array of objects like [{"name": "Skill"}]. Default: [].
"""