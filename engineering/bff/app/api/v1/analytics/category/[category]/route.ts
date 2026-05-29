import { NextResponse } from "next/server";
import { getCategoryDeepDiveAnalysis } from "@/core/services/dashboard.service";
import { handleBffError } from "@/core/errors/bff-error";

type RouteParams = {
  params: Promise<{ category: string }>;
};

export async function GET(request: Request, { params }: RouteParams) {
  const { category } = await params;
  const path = `/api/v1/analytics/category/${category}`;

  try {
    const decodedCategory = decodeURIComponent(category);
    const analysisPayload = await getCategoryDeepDiveAnalysis(decodedCategory);
    return NextResponse.json(analysisPayload, { status: 200 });
  } catch (error) {
    return handleBffError(error, path);
  }
}