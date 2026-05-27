// ============================================================================
// BFF API — OpenAPI 3.0.3 Schema Types
// Labour Market Intelligence BFF API
// ============================================================================

// --- KPI Summary ---
export interface KpiSummaryBlock {
  totalVacancies: number;
  vacancyGrowthPct: number;
  sectorsTracked: number;
  skillsIdentified: number;
}

// --- Category Weight (Pie Chart) ---
export interface CategoryWeightNode {
  category: string;
  vacancies: number;
}

// --- Monthly Trend (Area Chart) ---
export interface MonthlyTrendNode {
  month: string;
  vacancies: number;
}

// --- Ingestion Source (Bar Chart) ---
export interface IngestionSourceNode {
  name: string;
  vacancies: number;
}

// --- Share Metric (Progress Bars) ---
export interface ShareMetricNode {
  name: string;
  share: number;
}

// --- Distribution Tracks Block ---
export interface DistributionTracksBlock {
  seniority: ShareMetricNode[];
  contractTypes: ShareMetricNode[];
  remoteConfiguration: ShareMetricNode[];
}

// --- District Geo Node (Map) ---
export interface DistrictGeoNode {
  id: string;
  province: string;
  jobs: number;
  nationalShare?: number;
}

// --- Dashboard Data Payload ---
export interface DashboardDataPayload {
  kpiSummary: KpiSummaryBlock;
  categoryData: CategoryWeightNode[];
  monthlyTrends: MonthlyTrendNode[];
  ingestionSources: IngestionSourceNode[];
  leadingEmployers: EmployerMetricsNode[];
  distributionTracks: DistributionTracksBlock;
  districtGeoData: DistrictGeoNode[];
}

// --- Filter Metadata ---
export interface LookupOptionNode {
  value: string;
}

export interface ProvinceDataNode {
  id: number;
  name: string;
  vacancies?: number;
}

export interface FilterMetadataPayload {
  categories: LookupOptionNode[];
  provinces: ProvinceDataNode[];
  contractTypes: LookupOptionNode[];
  seniorityLevels: LookupOptionNode[];
}

// --- Employer Metrics ---
export interface EmployerMetricsNode {
  name: string;
  openRoles: number;
  location?: string;
  sector?: string;
}

// --- Geo Coordinates ---
export interface GeoCoordinatesBlock {
  lat: number;
  lng: number;
  province: string;
}

// --- Vacancy Meta Block ---
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

// --- Job Vacancy Model ---
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

// --- Skill Competency Node ---
export interface SkillCompetencyNode {
  skill: string;
  demand: number;
  category: string;
}

// --- Category Analytics Payload ---
export interface CategoryAnalyticsPayload {
  skills: SkillCompetencyNode[];
  provinces: ProvinceDataNode[];
  employers: EmployerMetricsNode[];
}

// --- Error Response ---
export interface ErrorResponse {
  timestamp: string;
  status: number;
  error: string;
  message: string;
  path: string;
}

// ============================================================================
// Backward-Compatible Aliases (used by mock-data helpers & legacy hooks)
// ============================================================================

/** @deprecated Use JobVacancyModel. Kept for mock-data compatibility. */
export interface JobVacancy {
  id: string;
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

/** @deprecated Use SkillCompetencyNode */
export type SkillDemand = SkillCompetencyNode;

/** @deprecated Use ProvinceDataNode */
export interface ProvinceVacancy {
  province: string;
  vacancies: number;
}

/** @deprecated Use EmployerMetricsNode */
export interface HiringEmployer {
  name: string;
  openRoles: number;
  location: string;
}

/** @deprecated Use CategoryAnalyticsPayload */
export interface CategoryAnalyticsResponse {
  skills: SkillCompetencyNode[];
  provinces: ProvinceVacancy[];
  employers: HiringEmployer[];
}

/** @deprecated Use FilterMetadataPayload */
export interface FilterMetadataResponse {
  categories: string[];
  provinces: string[];
  contractTypes: string[];
  seniorityLevels: string[];
}