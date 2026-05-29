import { GoJobPostResponse } from "@/types/models";
import { BffHttpException } from "../errors/bff-error";

const GO_BACKEND_URL = process.env.GO_BACKEND_URL || "http://backend:8080/api/v1";

export async function fetchAllJobsFromCore(): Promise<GoJobPostResponse[]> {
  try {
    const response = await fetch(`${GO_BACKEND_URL}/jobs`, {
      method: "GET",
      headers: {
        "Accept": "application/json",
        "Content-Type": "application/json",
      },
      next: { revalidate: 10 }, // Cache resolution layer for 10 seconds to safeguard database throughput
    });

    if (!response.ok) {
      throw new BffHttpException(
        response.status,
        "Downstream core engine error",
        `Failed fetching datasets from relational core storage. Status: ${response.statusText}`
      );
    }

    return await response.json();
  } catch (error) {
    if (error instanceof BffHttpException) throw error;
    
    throw new BffHttpException(
      502,
      "Bad Gateway Trigger State",
      `The Core Go service appears offline or unreachable. Context: ${error instanceof Error ? error.message : String(error)}`
    );
  }
}