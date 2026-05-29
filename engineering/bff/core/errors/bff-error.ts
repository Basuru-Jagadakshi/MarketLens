import { NextResponse } from "next/server";
import { BffErrorResponse } from "@/types/payloads";

export class BffHttpException extends Error {
  constructor(
    public status: number,
    public errorKey: string,
    message: string
  ) {
    super(message);
    this.name = "BffHttpException";
  }
}

export function handleBffError(error: unknown, path: string): NextResponse<BffErrorResponse> {
  console.error(`[BFF Error Engine Execution Fallback Context] Route: ${path}`, error);

  const timestamp = new Date().toISOString();

  if (error instanceof BffHttpException) {
    return NextResponse.json(
      {
        timestamp,
        status: error.status,
        error: error.errorKey,
        message: error.message,
        path,
      },
      { status: error.status }
    );
  }

  // Fallback for unexpected failures (e.g., Go Engine Offline)
  return NextResponse.json(
    {
      timestamp,
      status: 500,
      error: "Internal Server Error",
      message: error instanceof Error ? error.message : "Downstream microservice communications failure.",
      path,
    },
    { status: 500 }
  );
}