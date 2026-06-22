export interface Industry {
  id:   number;
  name: string;
}

export interface IndustryListResponse {
  count:      number;
  industries: Industry[];
}

export interface SkillDemand {
  id:             number;
  skill:          string;
  open_job_count: number;
}

export interface EmployerDemand {
  id:             number;
  name:           string;
  open_job_count: number;
}

export interface MostInDemandSkill {
  id:             number;
  skill:          string;
  open_job_count: number;
}

export interface IndustrySkillsAnalytics {
  industry_id:          number;
  unique_skills_count:  number;
  most_in_demand_skill: MostInDemandSkill | null;
  top15_skills:         SkillDemand[];
  all_skills:           SkillDemand[];
  top_employers:        EmployerDemand[];
}