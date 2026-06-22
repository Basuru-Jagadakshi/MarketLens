export interface ActiveJobStats {
  active_job_count: number;
  last_month_count: number;
  change_percent: number;
  trend: "up" | "down" | "stable";
}

export interface OccupationStat {
  id: number;
  name: string;
  open_job_count: number;
}

export interface IndustryStat {
  id: number;
  name: string;
  open_job_count: number;
}

export interface ExperienceStat {
  id: number;
  name: string;
  open_job_count: number;
}

export interface EducationStat {
  id: number;
  level: string;
  open_job_count: number;
}

export interface RemoteOnSiteStat {
  remote_count: number;
  on_site_count: number;
}

export interface JobTypeStat {
  id: number;
  type: string;
  open_job_count: number;
}

export interface DashboardOverview {
  active_jobs:     ActiveJobStats;
  by_occupation:   { count: number; occupations: OccupationStat[] };
  by_industry:     { count: number; industries: IndustryStat[] };
  by_experience:   { count: number; experiences: ExperienceStat[] };
  by_education:    { count: number; education_levels: EducationStat[] };
  remote_vs_onsite: RemoteOnSiteStat;
  by_job_type:     { count: number; job_types: JobTypeStat[] };
}

export interface YearlyTrend {
  year: number;
  open_job_count: number;
}

export interface TopJobRole {
  job_role: string;
  open_job_count: number;
}

export interface OccupationAnalytics {
  occupation_id: number;
  yearly_trend:  { occupation_id: number; count: number; yearly_trend: YearlyTrend[] };
  top_job_roles: { occupation_id: number; top_job_roles: TopJobRole[] };
}

export interface IndustryYearlyTrend {
  year: number;
  open_job_count: number;
}

export interface IndustryExperience {
  id: number;
  name: string;
  open_job_count: number;
}

export interface IndustryProvince {
  id: number;
  province: string;
  lat: number;
  lng: number;
  open_job_count: number;
}

export interface IndustryEducation {
  id: number;
  level: string;
  open_job_count: number;
}

export interface IndustryEmployer {
  id: number;
  name: string;
  open_job_count: number;
}

export interface IndustryAnalytics {
  industry_id:  number;
  year:         number;
  yearly_trend: { count: number; yearly_trend: IndustryYearlyTrend[] };
  by_experience: { count: number; experiences: IndustryExperience[] };
  by_province:   { count: number; provinces: IndustryProvince[] };
  by_education:  { count: number; education_levels: IndustryEducation[] };
  top_employers: { count: number; employers: IndustryEmployer[] };
}