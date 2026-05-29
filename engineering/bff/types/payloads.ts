export interface KpiSummaryBlock {
  totalVacancies: number;
  vacancyGrowthPct: number;
  sectorsTracked: number;
  skillsIdentified: number;
}

export interface CategoryWeightNode {
  category: string;
  vacancies: number;
}

export interface MonthlyTrendNode {
  month: string;
  vacancies: number;
}

export interface IngestionSourceNode {
  name: string;
  vacancies: number;
}

export interface ShareMetricNode {
  name: string;
  share: number;
}

export interface DistributionTracksBlock {
  seniority: ShareMetricNode[];
  contractTypes: ShareMetricNode[];
  remoteConfiguration: ShareMetricNode[];
}

export interface DistrictGeoNode {
  id: string;
  province: string;
  jobs: number;
  nationalShare?: number;
}

export interface DashboardDataPayload {
  kpiSummary: KpiSummaryBlock;
  categoryData: CategoryWeightNode[];
  monthlyTrends: MonthlyTrendNode[];
  ingestionSources: IngestionSourceNode[];
  leadingEmployers: EmployerMetricsNode[];
  distributionTracks: DistributionTracksBlock;
  districtGeoData: DistrictGeoNode[];
}

export interface LookupOptionNode {
  value: string;
}

export interface ProvinceDataNode {
  id: number;
  name: string;
  vacancies?: number;
}

export interface EmployerMetricsNode {
  name: string;
  location?: string;
  sector?: string;
  openRoles: number;
}

export interface GeoCoordinatesBlock {
  lat: number;
  lng: number;
  province: string;
}

export interface VacancyMetaBlock {
  posted_at: string;
  source: string;
  standardized_category: string;
  seniority: string;
  geo: GeoCoordinatesBlock;
  confidence_score: number;
  ai_version: string;
  error: boolean;
}

export interface JobVacancyModel {
  id: number;
  employer: string;
  job_role: string;
  job_type: string;
  key_responsibilities: string;
  qualifications: string;
  location: string;
  offers: string;
  is_remote: boolean;
  skills: string[];
  meta_data: VacancyMetaBlock;
}

export interface FilterMetadataPayload {
  categories: LookupOptionNode[];
  provinces: ProvinceDataNode[];
  contractTypes: LookupOptionNode[];
  seniorityLevels: LookupOptionNode[];
}

export interface SkillCompetencyNode {
  skill: string;
  demand: number;
  category: string;
}

export interface CategoryAnalyticsPayload {
  skills: SkillCompetencyNode[];
  provinces: ProvinceDataNode[];
  employers: EmployerMetricsNode[];
}

export interface BffErrorResponse {
  timestamp: string;
  status: number;
  error: string;
  message: string;
  path: string;
}