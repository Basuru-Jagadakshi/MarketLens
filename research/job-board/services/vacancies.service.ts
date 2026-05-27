import { apiClient } from "./api-client";
import {
  MOCK_FILTER_METADATA,
  MOCK_VACANCIES,
  MOCK_ANALYTICS_DATA,
  MOCK_DASHBOARD_DATA,
} from "./mock-data";
import {
  FilterMetadataPayload,
  JobVacancy,
  CategoryAnalyticsPayload,
  DashboardDataPayload,
} from "@/types/job";

export interface QueryParams {
  category?: string;
  province?: number;
  contractType?: string;
  seniority?: string;
}

export const vacanciesService = {
  // GET /vacancies/filter-values
  async getFilterOptions(): Promise<FilterMetadataPayload> {
    try {
      const { data } = await apiClient.get<FilterMetadataPayload>(
        "/vacancies/filter-values"
      );
      return data;
    } catch {
      return MOCK_FILTER_METADATA;
    }
  },

  // GET /vacancies
  async getAllVacancies(filters: QueryParams): Promise<JobVacancy[]> {
    try {
      const params: Record<string, string | number | undefined> = {};
      if (filters.category) params.category = filters.category;
      if (filters.province !== undefined) params.province = filters.province;
      if (filters.contractType) params.contractType = filters.contractType;
      if (filters.seniority) params.seniority = filters.seniority;

      const { data } = await apiClient.get<JobVacancy[]>("/vacancies", {
        params,
      });
      return data;
    } catch {
      // Local fallback: mirror server-side filter logic
      return MOCK_VACANCIES.filter((v) => {
        if (
          filters.category &&
          filters.category !== "All" &&
          v.meta_data.standardized_category !== filters.category
        )
          return false;
        if (
          filters.province !== undefined &&
          filters.province !== 0
        ) {
          const provinceObj = MOCK_FILTER_METADATA.provinces.find(
            (p) => p.id === filters.province
          );
          if (
            !provinceObj ||
            v.meta_data.geo.province.toLowerCase() !== provinceObj.name.toLowerCase()
          ) {
            return false;
          }
        }
        if (
          filters.contractType &&
          filters.contractType !== "All" &&
          v.job_type.toLowerCase() !== filters.contractType.toLowerCase()
        )
          return false;
        if (
          filters.seniority &&
          filters.seniority !== "All" &&
          v.meta_data.seniority.toLowerCase() !== filters.seniority.toLowerCase()
        )
          return false;
        return true;
      });
    }
  },

  // GET /analytics/category/{category}
  async getAnalyticsByCategory(
    category: string
  ): Promise<CategoryAnalyticsPayload> {
    try {
      const normalized =
        category === "All Categories" ? "All" : category;
      const { data } = await apiClient.get<CategoryAnalyticsPayload>(
        `/analytics/category/${encodeURIComponent(normalized)}`
      );
      return data;
    } catch {
      return (
        MOCK_ANALYTICS_DATA[category] ||
        MOCK_ANALYTICS_DATA["All Categories"]
      );
    }
  },

  // GET /analytics/dashboard
  async getSummaryMetrics(): Promise<DashboardDataPayload> {
    try {
      const { data } = await apiClient.get<DashboardDataPayload>(
        "/analytics/dashboard"
      );
      return data;
    } catch {
      return MOCK_DASHBOARD_DATA;
    }
  },
};