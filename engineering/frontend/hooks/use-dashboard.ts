import { useQuery } from "@tanstack/react-query";
import { ENV } from "@/config/env";
import {
  DashboardOverview,
  OccupationAnalytics,
  IndustryAnalytics,
} from "@/types/dashboard";

const BASE = ENV.NEXT_PUBLIC_API_BASE_URL;

// ── Overview (on page load) ───────────────────────────────────────────────────
async function fetchDashboardOverview(): Promise<DashboardOverview> {
  const res = await fetch(`${BASE}/dashboard/overview`);
  if (!res.ok) throw new Error("Failed to fetch dashboard overview");
  return res.json();
}

export function useDashboardOverview() {
  return useQuery<DashboardOverview>({
    queryKey: ["dashboard", "overview"],
    queryFn:  fetchDashboardOverview,
    staleTime: 1000 * 60 * 5, // 5 minutes
  });
}

// ── Occupation analytics (on "See Analytics" click) ───────────────────────────
async function fetchOccupationAnalytics(occupationId: number): Promise<OccupationAnalytics> {
  const res = await fetch(`${BASE}/dashboard/occupation-analytics?occupation_id=${occupationId}`);
  if (!res.ok) throw new Error("Failed to fetch occupation analytics");
  return res.json();
}

export function useOccupationAnalytics(occupationId: number | null) {
  return useQuery<OccupationAnalytics>({
    queryKey: ["dashboard", "occupation-analytics", occupationId],
    queryFn:  () => fetchOccupationAnalytics(occupationId!),
    enabled:  occupationId !== null,
    staleTime: 1000 * 60 * 5,
  });
}

// ── Industry analytics (on "See Analytics" click + year change) ───────────────
async function fetchIndustryAnalytics(
  industryId: number,
  year: number
): Promise<IndustryAnalytics> {
  const res = await fetch(
    `${BASE}/dashboard/industry-analytics?industry_id=${industryId}&year=${year}`
  );
  if (!res.ok) throw new Error("Failed to fetch industry analytics");
  return res.json();
}

export function useIndustryAnalytics(industryId: number | null, year: number) {
  return useQuery<IndustryAnalytics>({
    queryKey: ["dashboard", "industry-analytics", industryId, year],
    queryFn:  () => fetchIndustryAnalytics(industryId!, year),
    enabled:  industryId !== null,
    staleTime: 1000 * 60 * 5,
  });
}