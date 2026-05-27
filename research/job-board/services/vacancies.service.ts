import { apiClient } from "./api-client";
import { MOCK_METADATA, MOCK_VACANCIES, MOCK_ANALYTICS_DATA, DashboardMetricsResponse, getDynamicDashboardMetrics } from "./mock-data";
import { FilterMetadataResponse, JobVacancy, CategoryAnalyticsResponse } from "@/types/job";

export interface QueryParams {
  category?: string;
  province?: string;
  contractType?: string;
  seniority?: string;
}

export const vacanciesService = {
  async getFilterOptions(): Promise<FilterMetadataResponse> {
    try {
      const { data } = await apiClient.get<FilterMetadataResponse>("/vacancies/meta-options");
      return data;
    } catch {
      // Clean fallback orchestration for testing environments
      return MOCK_METADATA;
    }
  },

  async getAllVacancies(filters: QueryParams): Promise<JobVacancy[]> {
    try {
      const { data } = await apiClient.get<JobVacancy[]>("/vacancies", { params: filters });
      return data;
    } catch {
      // Fallback matching server-side param operations locally
      return MOCK_VACANCIES.filter((v) => {
        if (filters.category && filters.category !== "All" && v.meta_data.standardized_category !== filters.category) return false;
        if (filters.province && filters.province !== "All" && v.meta_data.geo.province !== filters.province) return false;
        if (filters.contractType && filters.contractType !== "All" && v.job_type.toLowerCase() !== filters.contractType.toLowerCase()) return false;
        if (filters.seniority && filters.seniority !== "All" && v.meta_data.seniority.toLowerCase() !== filters.seniority.toLowerCase()) return false;
        return true;
      });
    }
  },

  async getAnalyticsByCategory(category: string): Promise<CategoryAnalyticsResponse> {
    try {
      const normalizedCategory = category === "All Categories" ? "All" : category;
      const { data } = await apiClient.get<CategoryAnalyticsResponse>(`/analytics/category`, {
        params: { category: normalizedCategory },
      });
      return data;
    } catch {
      // Local development environment circuit breaker fallback
      return MOCK_ANALYTICS_DATA[category] || MOCK_ANALYTICS_DATA["All Categories"];
    }
  },

  async getSummaryMetrics(): Promise<DashboardMetricsResponse> {
    try {
      const { data } = await apiClient.get<DashboardMetricsResponse>("/dashboard/summary");
      return data;
    } catch {
      // Local circuit breaker fallback executing calculation over latest raw vacancies list
      return getDynamicDashboardMetrics(MOCK_VACANCIES);
    }
  },
};