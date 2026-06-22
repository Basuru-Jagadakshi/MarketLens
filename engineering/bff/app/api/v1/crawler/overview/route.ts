import { NextResponse } from "next/server";

const GO_API = process.env.GO_BACKEND_URL;

export async function GET() {
  try {
    const [lastJobCount, timeGap, crawlerRuns, sources] = await Promise.all([
      fetch(`${GO_API}/crawler/last-job-count`).then((r) => {
        if (!r.ok) throw new Error(`crawler/last-job-count failed: ${r.status}`);
        return r.json();
      }),
      fetch(`${GO_API}/crawler/time-gap`).then((r) => {
        if (!r.ok) throw new Error(`crawler/time-gap failed: ${r.status}`);
        return r.json();
      }),
      fetch(`${GO_API}/crawler/runs`).then((r) => {
        if (!r.ok) throw new Error(`crawler/runs failed: ${r.status}`);
        return r.json();
      }),
      fetch(`${GO_API}/sources`).then((r) => {
        if (!r.ok) throw new Error(`sources failed: ${r.status}`);
        return r.json();
      }),
    ]);

    return NextResponse.json({
      kpis: {
        last_crawl_job_count: lastJobCount.last_crawl_job_count ?? 0,
        last_crawled_at:      timeGap.last_crawled_at          ?? null,
        gap_seconds:          timeGap.gap_seconds               ?? 0,
        gap_human:            timeGap.gap_human                 ?? "N/A",
      },
      sources:      sources.sources       ?? [],
      crawler_runs: crawlerRuns.runs      ?? [],
    });
  } catch (error) {
    console.error("[crawler/overview] error:", error);
    return NextResponse.json(
      {
        error:   "Failed to fetch crawler overview data",
        details: error instanceof Error ? error.message : String(error),
      },
      { status: 500 }
    );
  }
}