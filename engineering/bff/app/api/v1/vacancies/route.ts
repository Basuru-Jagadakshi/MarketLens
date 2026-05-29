import { NextResponse } from "next/server";
import { getFilteredVacancies } from "@/core/services/vacancy.service";
import { handleBffError } from "@/core/errors/bff-error";

export async function GET(request: Request) {
  const path = "/api/v1/vacancies";
  try {
    const { searchParams } = new URL(request.url);
    
    const filters = {
      category: searchParams.get("category") || undefined,
      province: searchParams.get("province") || undefined,
      contractType: searchParams.get("contractType") || undefined,
      seniority: searchParams.get("seniority") || undefined,
    };

    const records = await getFilteredVacancies(filters);
    return NextResponse.json(records, { status: 200 });
  } catch (error) {
    return handleBffError(error, path);
  }
}