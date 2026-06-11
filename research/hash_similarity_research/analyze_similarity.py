#python3 analyze_similarity.py

import json
import re
import csv

def prepare_text(title, employer, description):
    clean_title = re.sub(r'[^a-zA-Z0-9]', '', title.lower())
    clean_employer = re.sub(r'[^a-zA-Z0-9]', '', employer.lower())
    clean_desc = re.sub(r'[^a-z0-9\s]', '', description.lower())
    return f"{clean_title} {clean_employer} {clean_desc}"

def generate_3_word_shingles(text):
    words = text.split()
    shingles = set()
    for i in range(len(words) - 2):
        shingle = f"{words[i]} {words[i+1]} {words[i+2]}"
        shingles.add(shingle)
    return shingles

def calculate_jaccard(set_a, set_b):
    if not set_a or not set_b:
        return 0.0
    return len(set_a.intersection(set_b)) / len(set_a.union(set_b))

with open("scraped_job_market.json", "r", encoding="utf-8") as f:
    job_matrix = json.load(f)

research_results = []

for base_group in job_matrix:
    base_flat = prepare_text(base_group["base_title"], base_group["base_employer"], base_group["base_description"])
    base_shingles = generate_3_word_shingles(base_flat)
    
    for perm in base_group["permutations"]:
        comp_flat = prepare_text(perm["title"], perm["employer"], perm["description"])
        comp_shingles = generate_3_word_shingles(comp_flat)
        
        score = calculate_jaccard(base_shingles, comp_shingles)
            
        research_results.append({
            "Base Job ID": base_group["base_id"],
            "Base Job Title": base_group["base_title"],
            "Tested Case ID": perm["case_id"],
            "Permutation Profile": perm["name"],
            "Calculated Jaccard Score": round(score, 4)
        })

csv_filename = "lsh_similarity_research.csv"
fields = ["Base Job ID", "Base Job Title", "Tested Case ID", "Permutation Profile", "Calculated Jaccard Score", "System Action Decision"]

with open(csv_filename, mode="w", newline="", encoding="utf-8") as f:
    writer = csv.DictWriter(f, fieldnames=fields)
    writer.writeheader()
    writer.writerows(research_results)