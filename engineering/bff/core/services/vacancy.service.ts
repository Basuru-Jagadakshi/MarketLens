import { fetchAllJobsFromCore } from "../clients/go-backend.client";
import { GoJobPostResponse } from "@/types/models";
import { JobVacancyModel, FilterMetadataPayload } from "@/types/payloads";


export function mapToGoVacancyModel(job: GoJobPostResponse): JobVacancyModel {
  return {
    id: job.id,
    employer: job.employer,
    job_role: job.job_role,
    job_type: job.job_type?.name || "Full Time",
    key_responsibilities: job.key_responsibilities,
    qualifications: job.qualifications,
    location: job.location,
    offers: "Salary Negotiable.", 
    is_remote: job.is_remote,
    skills: job.skills?.map(s => s.name) || [],
    meta_data: {
      posted_at: new Date().toISOString(),
      source: job.meta_data?.source || "System Scraper Engine",
      standardized_category: job.meta_data?.standardized_category || "Unclassified",
      seniority: job.meta_data?.seniority || "Associate",
      confidence_score: job.meta_data?.confidence_score || 1.0,
      ai_version: job.meta_data?.ai_version || "Production",
      error: job.meta_data?.error || false,
      geo: {
        lat: 6.9271, // Colombo base coordinate fallback
        lng: 79.8612,
        province: job.meta_data?.geo?.province || "Western"
      }
    }
  };
}


export async function getFilteredVacancies(filters: {
  category?: string;
  province?: string;
  contractType?: string;
  seniority?: string;
}): Promise<JobVacancyModel[]> {
  const rawJobs = await fetchAllJobsFromCore();

  return rawJobs
    .filter(job => {
      if (filters.category && filters.category !== "All" && job.meta_data?.standardized_category !== filters.category) return false;
      if (filters.province && filters.province !== "All" && job.meta_data?.geo?.province !== filters.province) return false;
      if (filters.contractType && filters.contractType !== "All" && job.job_type?.name !== filters.contractType) return false;
      if (filters.seniority && filters.seniority !== "All" && job.meta_data?.seniority !== filters.seniority) return false;
      return true;
    })
    .map(mapToGoVacancyModel);
}


export async function getFilterLookupMetadata(): Promise<FilterMetadataPayload> {
  const rawJobs = await fetchAllJobsFromCore();

  const categories = new Set<string>();
  const provinces = new Set<string>();
  const contractTypes = new Set<string>();
  const seniorityLevels = new Set<string>();

  rawJobs.forEach(j => {
    if (j.meta_data?.standardized_category) categories.add(j.meta_data.standardized_category);
    if (j.meta_data?.geo?.province) provinces.add(j.meta_data.geo.province);
    if (j.job_type?.name) contractTypes.add(j.job_type.name);
    if (j.meta_data?.seniority) seniorityLevels.add(j.meta_data.seniority);
  });

  return {
    categories: Array.from(categories).map(value => ({ value })),
    provinces: Array.from(provinces).map((name, id) => ({ id: id + 1, name })),
    contractTypes: Array.from(contractTypes).map(value => ({ value })),
    seniorityLevels: Array.from(seniorityLevels).map(value => ({ value }))
  };
}