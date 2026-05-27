export interface Job {
  employer: string;
  job_role: string;
  job_type: string;
  key_responsibilities: string;
  qualifications: string;
  location: string;
  offers: string;
}

export interface SkillDemand {
  skill: string;
  demand: number;
  category: string;
}

export interface ProvinceVacancy {
  province: string;
  vacancies: number;
}

export interface HiringEmployer {
  name: string;
  openRoles: number;
  location: string;
}

export interface CategoryAnalyticsResponse {
  skills: SkillDemand[];
  provinces: ProvinceVacancy[];
  employers: HiringEmployer[];
}

export interface FilterMetadataResponse {
  categories: string[];
  provinces: string[];
  contractTypes: string[];
  seniorityLevels: string[];
}

export interface JobVacancy {
  id: string; // Added ID for key tracking loops
  employer: string;
  job_role: string;
  job_type: string;
  key_responsibilities: string;
  qualifications: string;
  location: string;
  offers: string;
  is_remote: boolean;
  skills: string[];
  meta_data: {
    posted_at: string;
    source: string;
    standardized_category: string;
    seniority: string;
    geo: {
      lat: number;
      lng: number;
      province: string;
    };
    confidence_score: number;
    ai_version: string;
    error: boolean;
  };
}

export interface KpiSummary {
  totalVacancies: number;
  vacancyGrowthPct: number;
  sectorsTracked: number;
  skillsIdentified: number;
}

export interface MonthlyTrend {
  month: string;
  domestic: number;
  overseas: number;
}

export interface SectorShare {
  sector: string;
  vacancies: number;
}

export interface IngestionSource {
  name: string;
  vacancies: number;
}

export interface LeadingEmployer {
  name: string;
  activePosts: number;
  sector: string;
}

export interface ShareDistribution {
  level?: string;
  type?: string;
  share: number;
}

export interface DistrictVacancy {
  id: string;
  name: string;
  province: string;
  jobs: number;
  path: string;
}

export interface ContractType {
  id: string;
  name: string;
}

export interface DashboardMetricsResponse {
  kpiSummary: KpiSummary;
  monthlyTrends: MonthlyTrend[];
  sectorDistribution: SectorShare[];
  ingestionSources: IngestionSource[];
  leadingEmployers: LeadingEmployer[];
  seniorityData: ShareDistribution[];
  contractTypes: ShareDistribution[];
  remoteConfiguration: ShareDistribution[];
  regionalVacancies: DistrictVacancy[];
}