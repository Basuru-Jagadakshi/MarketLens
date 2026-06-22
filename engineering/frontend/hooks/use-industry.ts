import { useQuery } from "@tanstack/react-query";
import { ENV } from "@/config/env";
import { IndustryListResponse, IndustrySkillsAnalytics } from "@/types/industry";

const BASE = ENV.NEXT_PUBLIC_API_BASE_URL;

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