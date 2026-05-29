export interface GoJobType {
  id?: number;
  name: string;
}

export interface GoSkill {
  id?: number;
  name: string;
}

export interface GoGeoData {
  id?: number;
  province: string;
}

export interface GoJobMetaData {
  id: number;
  job_post_id: number;
  source: string;
  standardized_category: string;
  seniority: string;
  confidence_score: number;
  ai_version: string;
  error: boolean;
  geo_id: number;
  geo: GoGeoData;
}

export interface GoJobPostResponse {
  id: number;
  employer: string;
  job_role: string;
  key_responsibilities: string;
  qualifications: string;
  location: string;
  is_remote: boolean;
  created_at: string;
  job_type_id: number;
  job_type: GoJobType;
  meta_data: GoJobMetaData;
  skills: GoSkill[];
}