import { useQuery } from "@tanstack/react-query";
import { ENV } from "@/config/env";
import { IndustryListResponse, 
  IndustrySkillsAnalytics,
  IndustryDivisionChildrenResponse,
  IndustryGroupChildrenResponse,
  IndustryClassChildrenResponse,
  IndustrySubclassChildrenResponse,
  IndustryAnalysisResponse,
} from "@/types/industry";

const BASE = ENV.NEXT_PUBLIC_API_BASE_URL;

async function fetchJson<T>(url: string): Promise<T> {
  const res = await fetch(url);
  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new Error(body?.error || body?.details || `Request failed: ${url}`);
  }
  return res.json();
}

export function useIndustrySectors(fromDate: string, toDate: string) {
  return useQuery<IndustryListResponse>({
    queryKey: ["industry", "sectors", fromDate, toDate],
    queryFn: () => {
      const params = new URLSearchParams({ "from-date": fromDate, "to-date": toDate });
      return fetchJson(`${BASE}/industry/industry-sectors?${params.toString()}`);
    },
    enabled: Boolean(fromDate) && Boolean(toDate),
    staleTime: 1000 * 60 * 10,
  });
}

export function useIndustryDivisions(industrySectorId: number | null, fromDate: string, toDate: string) {
  return useQuery<IndustryDivisionChildrenResponse>({
    queryKey: ["industry", "industry-divisions", industrySectorId, fromDate, toDate],
    queryFn: () => {
      const params = new URLSearchParams({ "from-date": fromDate, "to-date": toDate });
      return fetchJson(`${BASE}/industry/industry-sectors/${industrySectorId}/industry-divisions?${params.toString()}`);
    },
    enabled: industrySectorId !== null && Boolean(fromDate) && Boolean(toDate),
    staleTime: 1000 * 60 * 10,
  });
}

export function useIndustryGroups(industryDivisionId: number | null, fromDate: string, toDate: string) {
  return useQuery<IndustryGroupChildrenResponse>({
    queryKey: ["industry", "industry-groups", industryDivisionId, fromDate, toDate],
    queryFn: () => {
      const params = new URLSearchParams({ "from-date": fromDate, "to-date": toDate });
      return fetchJson(`${BASE}/industry/industry-divisions/${industryDivisionId}/industry-groups?${params.toString()}`);
    },
    enabled: industryDivisionId !== null && Boolean(fromDate) && Boolean(toDate),
    staleTime: 1000 * 60 * 10,
  });
}

export function useIndustryClasses(industryGroupId: number | null, fromDate: string, toDate: string) {
  return useQuery<IndustryClassChildrenResponse>({
    queryKey: ["industry", "industry-classes", industryGroupId, fromDate, toDate],
    queryFn: () => {
      const params = new URLSearchParams({ "from-date": fromDate, "to-date": toDate });
      return fetchJson(`${BASE}/industry/industry-groups/${industryGroupId}/industry-classes?${params.toString()}`);
    },
    enabled: industryGroupId !== null && Boolean(fromDate) && Boolean(toDate),
    staleTime: 1000 * 60 * 10,
  });
}

export function useIndustrySubclasses(industryClassId: number | null, fromDate: string, toDate: string) {
  return useQuery<IndustrySubclassChildrenResponse>({
    queryKey: ["industry", "industry-subclasses", industryClassId, fromDate, toDate],
    queryFn: () => {
      const params = new URLSearchParams({ "from-date": fromDate, "to-date": toDate });
      return fetchJson(`${BASE}/industry/industry-classes/${industryClassId}/industry-subclasses?${params.toString()}`);
    },
    enabled: industryClassId !== null && Boolean(fromDate) && Boolean(toDate),
    staleTime: 1000 * 60 * 10,
  });
}

async function fetchIndustryAnalysis(
  level: string,
  id: number,
  fromDate: string,
  toDate: string,
): Promise<IndustryAnalysisResponse> {
  const params = new URLSearchParams({ "from-date": fromDate, "to-date": toDate });
  const res = await fetch(`${BASE}/analysis/industry/${level}/${id}?${params.toString()}`);
  if (!res.ok) {
    const body = await res.json().catch(() => null);
    throw new Error(body?.error || body?.details || "Failed to fetch industry analysis");
  }
  return res.json();
}

export function useIndustryAnalysis(
  level: string | null,
  id: number | null,
  fromDate: string,
  toDate: string,
) {
  return useQuery<IndustryAnalysisResponse>({
    queryKey: ["industry-analysis", level, id, fromDate, toDate],
    queryFn: () => fetchIndustryAnalysis(level as string, id as number, fromDate, toDate),
    enabled: level !== null && id !== null && Boolean(fromDate) && Boolean(toDate) && fromDate <= toDate,
    staleTime: 1000 * 60 * 5,
  });
}

async function fetchIndustries(): Promise<IndustryListResponse> {
  const res = await fetch(`${BASE}/industry`);
  if (!res.ok) throw new Error("Failed to fetch industries");
  return res.json();
}

export function useIndustries() {
  return useQuery<IndustryListResponse>({
    queryKey: ["industries"],
    queryFn:  fetchIndustries,
    staleTime: 1000 * 60 * 10,
  });
}

async function fetchIndustrySkillsAnalytics(industryId: number): Promise<IndustrySkillsAnalytics> {
  const res = await fetch(`${BASE}/industry/skills-analytics?industry_id=${industryId}`);
  if (!res.ok) throw new Error("Failed to fetch industry skills analytics");
  return res.json();
}

export function useIndustrySkillsAnalytics(industryId: number | null) {
  return useQuery<IndustrySkillsAnalytics>({
    queryKey: ["industry", "skills-analytics", industryId],
    queryFn:  () => fetchIndustrySkillsAnalytics(industryId!),
    enabled:  industryId !== null,
    staleTime: 1000 * 60 * 5,
  });
}