import { useQuery } from "@tanstack/react-query";
import { ENV } from "@/config/env";
import { CrawlerOverviewResponse } from "@/types/crawler";

const BASE = ENV.NEXT_PUBLIC_API_BASE_URL;

async function fetchCrawlerOverview(): Promise<CrawlerOverviewResponse> {
  const res = await fetch(`${BASE}/crawler/overview`);
  if (!res.ok) throw new Error("Failed to fetch crawler overview");
  return res.json();
}

export function useCrawlerOverview() {
  return useQuery<CrawlerOverviewResponse>({
    queryKey: ["crawler-overview"],
    queryFn: fetchCrawlerOverview,
    staleTime: 1000 * 60 * 1, // Refresh every minute
    refetchInterval: 1000 * 30, // Poll every 30 seconds for live updates
  });
}