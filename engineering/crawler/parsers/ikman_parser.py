# crawlers/parsers/ikman_parser.py

import re
from datetime import datetime, timezone

PROVINCE_MAP = {
    # Western Province
    "colombo": "Western",
    "gampaha": "Western",
    "kalutara": "Western",

    # Central Province
    "kandy": "Central",
    "matale": "Central",
    "nuwara eliya": "Central",

    # Southern Province
    "galle": "Southern",
    "matara": "Southern",
    "hambantota": "Southern",

    # Northern Province
    "jaffna": "Northern",
    "kilinochchi": "Northern",
    "mannar": "Northern",
    "mullaitivu": "Northern",
    "vavuniya": "Northern",

    # Eastern Province
    "batticaloa": "Eastern",
    "ampara": "Eastern",
    "trincomalee": "Eastern",

    # North Western Province
    "kurunegala": "North Western",
    "puttalam": "North Western",

    # North Central Province
    "anuradhapura": "North Central",
    "polonnaruwa": "North Central",

    # Uva Province
    "badulla": "Uva",
    "monaragala": "Uva",

    # Sabaragamuwa Province
    "ratnapura": "Sabaragamuwa",
    "kegalle": "Sabaragamuwa",
}

# Noise markers where job description ends and footer begins
NOISE_MARKERS = [
    "Show more",
    "Call employer",
    "WhatsApp",
    "Boost this ad",
    "Report this ad",
    "Stay Alert",
    "Popular Jobs",
    "Trending Jobs",
    "Jobs by District",
    "Jobs by City",
    "More from ikman",
    "More ads from",
    "Help & Support",
    "About ikman",
    "Blog & Guides",
    "Download our app",
    "© 20",
]

# Mapping ikman categories → your standardized categories
IKMAN_CATEGORY_MAP = {
    # Technology & Telecom
    "IT and Network Industry":          "Technology & Telecom",
    "Digital Marketing":                "Technology & Telecom",
    "Data Entry":                       "Technology & Telecom",
    "Content Writers":                  "Technology & Telecom",

    # Finance & Banking
    "Accounting & Finance":             "Finance & Banking",
    "Banking and Financial services":   "Finance & Banking",
    "Collection and Recoveries":        "Finance & Banking",

    # Sales, Marketing & Retail
    "Sales and Marketing":              "Sales, Marketing & Retail",
    "Retail and Showroom":              "Sales, Marketing & Retail",
    "Merchandiser":                     "Sales, Marketing & Retail",
    "Business Development":             "Sales, Marketing & Retail",

    # Corporate, HR & Legal
    "HR and Recruitments":              "Corporate, HR & Legal",
    "Legal":                            "Corporate, HR & Legal",
    "Office Administration and Operations": "Corporate, HR & Legal",
    "Business Operations":              "Corporate, HR & Legal",
    "Consultant":                       "Corporate, HR & Legal",
    "Analyst":                          "Corporate, HR & Legal",
    "Clerk":                            "Corporate, HR & Legal",
    "Government Jobs":                  "Corporate, HR & Legal",

    # Customer Support & Service
    "Customer Service":                 "Customer Support & Service",
    "Call Centre":                      "Customer Support & Service",

    # Logistics, Manufacturing & Trade
    "Shipping and Logistics":           "Logistics, Manufacturing & Trade",
    "Stores and Warehouse":             "Logistics, Manufacturing & Trade",
    "Factory and Manufacturing":        "Logistics, Manufacturing & Trade",
    "Production and Operations":        "Logistics, Manufacturing & Trade",
    "Procurement and Purchasing":       "Logistics, Manufacturing & Trade",
    "Apparel and Garment":              "Logistics, Manufacturing & Trade",
    "Operator":                         "Logistics, Manufacturing & Trade",
    "Packing Officer":                  "Logistics, Manufacturing & Trade",

    # Hospitality, Food & Tourism
    "Hotel and Hospitality":            "Hospitality, Food & Tourism",
    "Restaurant and Cafe Staff":        "Hospitality, Food & Tourism",
    "Chef":                             "Hospitality, Food & Tourism",
    "Cook":                             "Hospitality, Food & Tourism",
    "Travel and Ticketing":             "Hospitality, Food & Tourism",
    "Event Planner & Coordinators":     "Hospitality, Food & Tourism",

    # Healthcare & Education
    "Hospital and Healthcare":          "Healthcare & Education",
    "Education and Training":           "Healthcare & Education",
    "Counsellor":                       "Healthcare & Education",
    "Sports and Wellness":              "Healthcare & Education",

    # Services, Driving & Blue-Collar
    "Driver":                           "Services, Driving & Blue-Collar",
    "Delivery Rider":                   "Services, Driving & Blue-Collar",
    "Mechanic":                         "Services, Driving & Blue-Collar",
    "Technician":                       "Services, Driving & Blue-Collar",
    "Helper":                           "Services, Driving & Blue-Collar",
    "Cleaner":                          "Services, Driving & Blue-Collar",
    "Security and Military":            "Services, Driving & Blue-Collar",
    "Domestic Jobs":                    "Services, Driving & Blue-Collar",
    "Building and Construction":        "Services, Driving & Blue-Collar",
    "Agriculture":                      "Services, Driving & Blue-Collar",
    "Crew member":                      "Services, Driving & Blue-Collar",
    "Cashier":                          "Services, Driving & Blue-Collar",
    "Supervisor":                       "Services, Driving & Blue-Collar",
    "Fire Fighter":                     "Services, Driving & Blue-Collar",

    # Creative, Media
    "Designer":                         "Creative, Media & Other",
    "Photography and Videography":      "Creative, Media & Other",
    "Editor":                           "Creative, Media & Other",
    "Architecture":                     "Creative, Media & Other",
    "Draftsman":                        "Creative, Media & Other",
    "Florist":                          "Creative, Media & Other",
    "Beauty & Hairdressing":            "Creative, Media & Other",
    "Freelance and NGO Worker":         "Creative, Media & Other",
    "Translator":                       "Creative, Media & Other",
    "Internships & Trainees":           "Creative, Media & Other",
    "Other":                            "Creative, Media & Other",
}


def extract_ikman_category(markdown: str) -> str:
    """Extract category from breadcrumb item 5 in numbered list."""
    
    # Match exactly "  5. [Category Name](url)"
    match = re.search(
        r'^\s*5\.\s*\[([^\]]+)\]\([^\)]+\)',
        markdown,
        re.MULTILINE
    )
    
    if match:
        raw_category = match.group(1).strip()
        # Remove " Jobs" suffix e.g. "Beauty & Hairdressing Jobs" → "Beauty & Hairdressing"
        raw_category = re.sub(r'\s+Jobs$', '', raw_category, flags=re.IGNORECASE).strip()
        return raw_category
    
    return ""


def map_to_standardized_category(ikman_category: str) -> str:
    """Map ikman category to your standardized category."""
    if not ikman_category:
        return "xxx"

    # Exact match first
    if ikman_category in IKMAN_CATEGORY_MAP:
        return IKMAN_CATEGORY_MAP[ikman_category]

    # Partial match fallback — in case of slight wording differences
    for key, value in IKMAN_CATEGORY_MAP.items():
        if key.lower() in ikman_category.lower() or ikman_category.lower() in key.lower():
            return value

    return "xxx"


def infer_province(location: str) -> str:
    loc = location.lower()
    for city, province in PROVINCE_MAP.items():
        if city in loc:
            return province
    return "Western"


def clean_markdown_links(text: str) -> str:
    """Convert [text](url) → text, and remove broken ](url) fragments."""
    text = re.sub(r'\[([^\]]*)\]\([^\)]+\)', r'\1', text)
    text = re.sub(r'\]\([^\)]+\)', '', text)
    return text.strip()


def clean_noise(text: str) -> str:
    """Cut off footer/noise content from description."""
    for marker in NOISE_MARKERS:
        idx = text.find(marker)
        if idx != -1:
            text = text[:idx]
    text = re.sub(r'\n+', ' ', text)
    text = re.sub(r'\s{2,}', ' ', text)
    return text.strip()


def extract_label(markdown: str, label: str) -> str:
    match = re.search(rf"{label}\s*[:\-]\s*([^\n]+)", markdown, re.IGNORECASE)
    return match.group(1).strip() if match else ""

def extract_location_from_posted_line(markdown: str) -> str:
    """Extract district (second span) from 'Posted on ... City, District' line."""
    match = re.search(r"Posted on [^,]+,\s*([^,\n]+),\s*([^,\n]+)", markdown, re.IGNORECASE)
    if match:
        # group(2) is the second span — the parent district
        location = match.group(2).strip().rstrip(".,")
        location = clean_markdown_links(location)
        return location
    return ""


def extract_description(markdown: str) -> str:
    """Grab job body content — saved to key_responsibilities."""
    patterns = [
        r"About the role\s*\n+([\s\S]+?)(?=\n#{1,3}\s|\Z)",
        r"Key Responsibilities\s*\n+([\s\S]+?)(?=\n#{1,3}\s|\Z)",
        r"Job Description\s*\n+([\s\S]+?)(?=\n#{1,3}\s|\Z)",
    ]
    for pattern in patterns:
        match = re.search(pattern, markdown, re.IGNORECASE)
        if match:
            return match.group(1).strip()

    # Fallback: everything after the meta label-value block
    meta_end = re.search(
        r"(Application deadline|Required work experience)[^\n]*\n",
        markdown, re.IGNORECASE
    )
    if meta_end:
        return markdown[meta_end.end():].strip()

    return ""


def parse_rule_based_fields(markdown: str, job_link: str = "", source: str = "Ikman") -> dict:
    """Build job dict from markdown using regex — zero LLM cost."""

    employer   = extract_label(markdown, "Employer")
    job_role   = extract_label(markdown, "Role")
    job_type   = extract_label(markdown, "Job type")

    ikman_category       = extract_ikman_category(markdown)
    standardized_category = map_to_standardized_category(ikman_category)

    # Clean markdown link syntax from job_role e.g. [Other Driver](url) → Other Driver
    job_role = clean_markdown_links(job_role)

    # Normalize private poster
    if employer.lower() in ("private poster", "private", ""):
        employer = ""

    # Fallback: H1 as job_role
    if not job_role:
        h1 = re.search(r"^#\s+(.+)$", markdown, re.MULTILINE)
        job_role = clean_markdown_links(h1.group(1).strip()) if h1 else ""

    location = extract_location_from_posted_line(markdown)
    if not location:
        location = extract_label(markdown, "Location")
    if not location and job_role:
        dash = re.search(r"-\s+([A-Za-z\s]+)$", job_role)
        if dash:
            location = dash.group(1).strip()

    location = location or "Sri Lanka"

    location      = location or "Sri Lanka"
    job_type_name = "Part Time" if "part" in job_type.lower() else "Full Time"
    is_remote     = bool(re.search(r"\bremote\b|\bwork from home\b|\bwfh\b", markdown, re.IGNORECASE))
    province      = infer_province(location)

    # Extract and clean description
    description = extract_description(markdown)
    description = clean_noise(description)

    return {
        "employer":             employer,
        "job_role":             job_role,
        "job_type":             {"name": job_type_name},
        "location":             location,
        "is_remote":            is_remote,
        "key_responsibilities": description,
        "qualifications":       "",
        "offers":               "",
        "meta_data": {
            "geo":                   {"province": province},
            "posted_at":             datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            "source":                source,
            "job_link":              job_link,
            "ai_version":            "v1.2.0",
            "error":                 False,
            # filled by LLM:
            "standardized_category": standardized_category,
            "seniority":             "Not Specified",
            "confidence_score":      0.0,
        },
        "skills": []
    }


def merge_llm_fields(rule_fields: dict, llm_fields: dict) -> dict:
    rule_fields["skills"] = llm_fields.get("skills", [])
    #rule_fields["meta_data"]["standardized_category"] = llm_fields.get("standardized_category", "Other")
    rule_fields["meta_data"]["seniority"]             = llm_fields.get("seniority", "Not Specified")
    rule_fields["meta_data"]["confidence_score"]      = llm_fields.get("confidence_score", 0.0)
    return rule_fields