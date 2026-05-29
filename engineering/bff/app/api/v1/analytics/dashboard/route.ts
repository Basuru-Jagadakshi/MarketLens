// app/api/analytics/dashboard/route.ts
import { NextResponse } from "next/server";
import { getAggregatedDashboardState } from "@/core/services/dashboard.service";
import { handleBffError } from "@/core/errors/bff-error";

export async function GET(request: Request) {
  const path = "/api/v1/analytics/dashboard";
  try {
    const dashboardData = await getAggregatedDashboardState();
    return NextResponse.json(dashboardData, { status: 200 });
  } catch (error) {
    return handleBffError(error, path);
  }
}