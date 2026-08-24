export interface HierarchyNode {
  id: number;
  name: string;
  code: string;
}

export interface IndustryDivisionChildrenResponse {
  industry_sector_id: number;
  count: number;
  industry_divisions: HierarchyNode[];
}
export interface IndustryGroupChildrenResponse {
  industry_division_id: number;
  count: number;
  industry_groups: HierarchyNode[];
}
export interface IndustryClassChildrenResponse {
  industry_group_id: number;
  count: number;
  industry_classes: HierarchyNode[];
}
export interface IndustrySubclassChildrenResponse {
  industry_class_id: number;
  count: number;
  industry_subclasses: HierarchyNode[];
}

export interface TotalJobCountResult {
  standard: string;
  level: string;
  id: number;
  total_job_count: number;
}

export interface ChildJobCount {
  id: number;
  name: string;
  code: string;
  open_job_count: number;
}

export interface ChildrenResult {
  standard: string;
  level: string;
  id: number;
  child_level: string;
  count: number;
  children: ChildJobCount[];
}

export interface EmploymentSectorStat { id: number; sector: string; open_job_count: number; }
export interface ExperienceStat { id: number; name: string; open_job_count: number; }
export interface ProvinceStat { id: number; province: string; lat: number; lng: number; open_job_count: number; }
export interface EducationLevelStat { id: number; level: string; open_job_count: number; }
export interface FormalityStat { id: number; formality_type: string; open_job_count: number; }
export interface GenderStat { id: number; gender_type: string; open_job_count: number; }
export interface VocationalEducationStat { id: number; level: string; open_job_count: number; }
export interface JobTypeStat { id: number; type: string; open_job_count: number; }
export interface RemoteOnSiteStat { remote_count: number; on_site_count: number; }

export interface IndustryAnalysisResponse {
  standard: string;
  level: string;
  id: string;
  from_date: string | null;
  to_date: string | null;
  total_job_count: TotalJobCountResult;
  children: ChildrenResult | null;
  employment_sector: { employment_sectors: EmploymentSectorStat[] };
  experience: { experiences: ExperienceStat[] };
  province: { provinces: ProvinceStat[] };
  education: { education_levels: EducationLevelStat[] };
  formality: { formalities: FormalityStat[] };
  gender: { genders: GenderStat[] };
  vocational_education: { vocational_educations: VocationalEducationStat[] };
  remote_onsite: { remote_vs_onsite: RemoteOnSiteStat };
  job_type: { job_types: JobTypeStat[] };
}

export interface Industry {
  id:   number;
  name: string;
}

export interface IndustryListResponse {
  count:      number;
  industry_sectors: Industry[];
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