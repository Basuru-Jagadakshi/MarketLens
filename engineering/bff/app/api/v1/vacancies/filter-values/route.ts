import { NextResponse } from "next/server";
import { getFilterLookupMetadata } from "@/core/services/vacancy.service";
import { handleBffError } from "@/core/errors/bff-error";

export async function GET(request: Request) {
  const path = "/api/v1/vacancies/filter-values";
  try {
    const filters = await getFilterLookupMetadata();
    return NextResponse.json(filters, { status: 200 });
  } catch (error) {
    return handleBffError(error, path);
  }
}