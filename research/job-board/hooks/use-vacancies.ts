import { useQuery } from "@tanstack/react-query";
import { vacanciesService, QueryParams } from "@/services/vacancies.service";

export const vacancyKeys = {
  all: ["vacancies"] as const,
  meta: () => [...vacancyKeys.all, "metadata"] as const,
  lists: (filters: QueryParams) => [...vacancyKeys.all, "list", filters] as const,
  analytics: (category: string) => [...vacancyKeys.all, "analytics", category] as const,
  dashboard: () => [...vacancyKeys.all, "dashboard-summary"] as const,
};

export function useVacancyMetadata() {
  return useQuery({
    queryKey: vacancyKeys.meta(),
    queryFn: vacanciesService.getFilterOptions,
    staleTime: 1000 * 60 * 30, // Metadata drops infrequently; cache for 30 mins
  });
}

export function useVacanciesList(filters: QueryParams) {
  return useQuery({
    queryKey: vacancyKeys.lists(filters),
    queryFn: () => vacanciesService.getAllVacancies(filters),
    staleTime: 1000 * 60 * 2, // Re-verify freshness every 2 mins
  });
}

export function useCategoryAnalytics(category: string) {
  return useQuery({
    queryKey: vacancyKeys.analytics(category),
    queryFn: () => vacanciesService.getAnalyticsByCategory(category),
    staleTime: 1000 * 60 * 5,
  });
}

export function useDashboardSummary() {
  return useQuery({
    queryKey: vacancyKeys.dashboard(),
    queryFn: vacanciesService.getSummaryMetrics,
    staleTime: 1000 * 60 * 5,
    refetchOnWindowFocus: false,
  });
}