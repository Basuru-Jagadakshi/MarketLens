export interface CrawlerKPIs {
  last_crawl_job_count: number;
  last_crawled_at:      string;
  gap_seconds:          number;
  gap_human:            string;
}

export interface CrawlerSource {
  id:             number;
  source:         string;
  open_job_count: number;
}

export interface CrawlerRun {
  id:          number;
  started_at:  string;
  finished_at: string | null;
  status:      "COMPLETED" | "RUNNING" | "FAILED";
}

export interface CrawlerOverviewResponse {
  kpis:         CrawlerKPIs;
  sources:      CrawlerSource[];
  crawler_runs: CrawlerRun[];
}