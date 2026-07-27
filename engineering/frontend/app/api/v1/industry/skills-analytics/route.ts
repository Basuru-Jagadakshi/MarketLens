import { NextRequest, NextResponse } from "next/server";

const GO_API = process.env.GO_BACKEND_URL;

export async function GET(request: NextRequest) {
  const { searchParams } = new URL(request.url);
  const industryId = searchParams.get("industry_id");

  if (!industryId) {
    return NextResponse.json(
      { error: "industry_id query parameter is required" },
      { status: 400 }
    );
  }

  const param = `industry_id=${industryId}`;

  try {
    const [
      skillsCount,
      topDemand,
      top15Skills,
      allSkills,
      topEmployers,
    ] = await Promise.all([
      fetch(`${GO_API}/industries/skills/count?${param}`).then((r) => {
        if (!r.ok) throw new Error(`skills/count failed: ${r.status}`);
        return r.json();
      }),
      fetch(`${GO_API}/industries/skills/top-demand?${param}`).then((r) => {
        if (!r.ok) throw new Error(`skills/top-demand failed: ${r.status}`);
        return r.json();
      }),
      fetch(`${GO_API}/industries/skills/top15?${param}`).then((r) => {
        if (!r.ok) throw new Error(`skills/top15 failed: ${r.status}`);
        return r.json();
      }),
      fetch(`${GO_API}/industries/skills?${param}`).then((r) => {
        if (!r.ok) throw new Error(`skills failed: ${r.status}`);
        return r.json();
      }),
      fetch(`${GO_API}/industries/employers?${param}`).then((r) => {
        if (!r.ok) throw new Error(`employers failed: ${r.status}`);
        return r.json();
      }),
    ]);

    return NextResponse.json({
      industry_id:        Number(industryId),
      unique_skills_count: skillsCount.unique_skills_count ?? 0,
      most_in_demand_skill: topDemand.most_in_demand_skill ?? null,
      top15_skills:        top15Skills.skills              ?? [],
      all_skills:          allSkills.skills                ?? [],
      top_employers:       topEmployers.employers          ?? [],
    });
  } catch (error) {
    console.error("[industry/skills-analytics] error:", error);
    return NextResponse.json(
      {
        error:   "Failed to fetch industry skills analytics",
        details: error instanceof Error ? error.message : String(error),
      },
      { status: 500 }
    );
  }
}