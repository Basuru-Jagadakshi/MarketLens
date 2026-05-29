import { describe, it, expect, beforeEach } from "vitest";
import { http, HttpResponse } from "msw";
import { server } from "@/vitest.setup";
import { GET as getVacancies } from "@/app/api/v1/vacancies/route";
import { GET as getFilterValues } from "@/app/api/v1/vacancies/filter-values/route";
import { GET as getDashboard } from "@/app/api/v1/analytics/dashboard/route";
import { GET as getCategoryAnalytics } from "@/app/api/v1/analytics/category/[category]/route";
import { GoJobPostResponse } from "@/types/models";
import { JobVacancyModel, FilterMetadataPayload, DashboardDataPayload, CategoryAnalyticsPayload, BffErrorResponse } from "@/types/payloads";

const BACKEND_JOBS_URL = "http://backend:8080/api/v1/jobs";

const mockJobsDataset: GoJobPostResponse[] = [
  {
    id: 101,
    employer: "WSO2",
    job_role: "Software Engineer",
    key_responsibilities: "Develop core engine components.",
    qualifications: "BSc in Computer Science",
    location: "Colombo",
    is_remote: true,
    job_type_id: 1,
    job_type: { id: 1, name: "Full Time" },
    created_at: "2026-03-15T00:00:00Z",
    meta_data: {
      id: 1,
      job_post_id: 101,
      standardized_category: "Software Engineering",
      source: "LinkedIn",
      seniority: "Senior",
      confidence_score: 98.9,
      ai_version: "chatgpt",
      error: false,
      geo_id: 1,
      geo: { province: "Western" },
    },
    skills: [{ id: 1, name: "Go" }, { id: 2, name: "Docker" }],
  },
  {
    id: 102,
    employer: "Dialog",
    job_role: "Data Analyst",
    key_responsibilities: "Build reporting metrics.",
    qualifications: "BSc in Statistics",
    location: "Trace Expert City",
    is_remote: false,
    job_type_id: 2,
    job_type: { id: 2, name: "Contract" },
    created_at: "2026-01-10T00:00:00Z",
    meta_data: {
      id: 2,
      job_post_id: 102,
      standardized_category: "Data Science",
      source: "LinkedIn",
      seniority: "Junior",
      confidence_score: 98.9,
      ai_version: "chatgpt",
      error: false,
      geo_id: 1,
      geo: { province: "Western" },
    },
    skills: [{ id: 3, name: "Python" }],
  },
  {
    id: 103,
    employer: "Sysco Labs",
    job_role: "QA Automation Engineer",
    key_responsibilities: "Develop testing automation tests.",
    qualifications: "BSc in Computer Science",
    location: "Kandy",
    is_remote: false,
    job_type_id: 1,
    job_type: { id: 1, name: "Full Time" },
    created_at: "2026-05-20T00:00:00Z",
    meta_data: {
      id: 3,
      job_post_id: 103,
      standardized_category: "Quality Assurance",
      source: "TopJobs",
      seniority: "Lead",
      confidence_score: 94.0,
      ai_version: "chatgpt",
      error: false,
      geo_id: 2,
      geo: { province: "Central" },
    },
    skills: [{ id: 2, name: "Docker" }, { id: 4, name: "Java" }],
  }
];

describe("BFF REST API Endpoints Integration Tests", () => {
  beforeEach(() => {
    
    server.use(
      http.get(BACKEND_JOBS_URL, () => {
        return HttpResponse.json(mockJobsDataset);
      })
    );
  });

  describe("GET /api/v1/vacancies", () => {
    it("should retrieve all vacancies successfully and format to standard JobVacancyModel schema", async () => {
      const request = new Request("http://localhost/api/v1/vacancies");
      const response = await getVacancies(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as JobVacancyModel[];
      
      expect(data).toHaveLength(3);
      
      const wso2Vacancy = data.find(j => j.id === 101);
      expect(wso2Vacancy).toBeDefined();
      expect(wso2Vacancy?.employer).toBe("WSO2");
      expect(wso2Vacancy?.job_type).toBe("Full Time");
      expect(wso2Vacancy?.is_remote).toBe(true);
      expect(wso2Vacancy?.skills).toEqual(["Go", "Docker"]);
      expect(wso2Vacancy?.meta_data.source).toBe("LinkedIn");
      expect(wso2Vacancy?.meta_data.standardized_category).toBe("Software Engineering");
      expect(wso2Vacancy?.meta_data.seniority).toBe("Senior");
      expect(wso2Vacancy?.meta_data.geo.province).toBe("Western");
      expect(wso2Vacancy?.meta_data.geo.lat).toBe(6.9271);
      expect(wso2Vacancy?.meta_data.geo.lng).toBe(79.8612);
    });

    it("should filter results correctly based on query search parameters", async () => {
      const request = new Request(
        "http://localhost/api/v1/vacancies?category=Software+Engineering&seniority=Senior&province=Western&contractType=Full+Time"
      );
      const response = await getVacancies(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as JobVacancyModel[];
      
      expect(data).toHaveLength(1);
      expect(data[0].id).toBe(101);
      expect(data[0].employer).toBe("WSO2");
    });

    it("should treat 'All' parameters as wildcards and skip exclusion filters", async () => {
      const request = new Request(
        "http://localhost/api/v1/vacancies?category=All&seniority=All&province=All&contractType=All"
      );
      const response = await getVacancies(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as JobVacancyModel[];
      expect(data).toHaveLength(3);
    });

    it("should return empty list if no record matches the filters", async () => {
      const request = new Request("http://localhost/api/v1/vacancies?category=Data+Science&seniority=Senior");
      const response = await getVacancies(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as JobVacancyModel[];
      expect(data).toHaveLength(0);
    });

    it("should handle downstream core server outages with structured BffErrorResponse", async () => {
      // Simulate outage
      server.use(
        http.get(BACKEND_JOBS_URL, () => {
          return new HttpResponse(null, { status: 502, statusText: "Bad Gateway" });
        })
      );

      const request = new Request("http://localhost/api/v1/vacancies");
      const response = await getVacancies(request);
      
      expect(response.status).toBe(502);
      const errBody = (await response.json()) as BffErrorResponse;
      expect(errBody.status).toBe(502);
      expect(errBody.error).toBe("Downstream core engine error");
      expect(errBody.path).toBe("/api/v1/vacancies");
      expect(errBody.message).toContain("Failed fetching datasets");
    });
  });

  describe("GET /api/v1/vacancies/filter-values", () => {
    it("should aggregate unique lookup values without duplicating categories, provinces, types or seniorities", async () => {
      const request = new Request("http://localhost/api/v1/vacancies/filter-values");
      const response = await getFilterValues(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as FilterMetadataPayload;
      
      // Expected deduplicated entries
      expect(data.categories).toHaveLength(3);
      expect(data.categories).toContainEqual({ value: "Software Engineering" });
      expect(data.categories).toContainEqual({ value: "Data Science" });
      expect(data.categories).toContainEqual({ value: "Quality Assurance" });

      expect(data.provinces).toHaveLength(2);
      expect(data.provinces.map(p => p.name)).toContain("Western");
      expect(data.provinces.map(p => p.name)).toContain("Central");

      expect(data.contractTypes).toHaveLength(2);
      expect(data.contractTypes).toContainEqual({ value: "Full Time" });
      expect(data.contractTypes).toContainEqual({ value: "Contract" });

      expect(data.seniorityLevels).toHaveLength(3);
      expect(data.seniorityLevels).toContainEqual({ value: "Senior" });
      expect(data.seniorityLevels).toContainEqual({ value: "Junior" });
      expect(data.seniorityLevels).toContainEqual({ value: "Lead" });
    });

    it("should bubble up downstream exceptions gracefully", async () => {
      server.use(
        http.get(BACKEND_JOBS_URL, () => {
          return new HttpResponse(null, { status: 500, statusText: "Internal Server Error" });
        })
      );

      const request = new Request("http://localhost/api/v1/vacancies/filter-values");
      const response = await getFilterValues(request);
      
      expect(response.status).toBe(500);
      const errBody = (await response.json()) as BffErrorResponse;
      expect(errBody.status).toBe(500);
      expect(errBody.error).toBe("Downstream core engine error");
      expect(errBody.path).toBe("/api/v1/vacancies/filter-values");
    });
  });

  describe("GET /api/v1/analytics/dashboard", () => {
    it("should correctly compute KPI metrics, chronological monthly trends, and structural configuration percentages", async () => {
      const request = new Request("http://localhost/api/v1/analytics/dashboard");
      const response = await getDashboard(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as DashboardDataPayload;
      
      expect(data.kpiSummary.totalVacancies).toBe(3);
      expect(data.kpiSummary.sectorsTracked).toBe(3);
      expect(data.kpiSummary.skillsIdentified).toBe(4); 
      
      expect(data.monthlyTrends).toHaveLength(3);
      expect(data.monthlyTrends[0].month).toBe("Jan 2026");
      expect(data.monthlyTrends[0].vacancies).toBe(1);
      expect(data.monthlyTrends[1].month).toBe("Mar 2026");
      expect(data.monthlyTrends[1].vacancies).toBe(1);
      expect(data.monthlyTrends[2].month).toBe("May 2026");
      expect(data.monthlyTrends[2].vacancies).toBe(1);
      
      expect(data.distributionTracks.remoteConfiguration).toContainEqual({
        name: "Remote Available",
        share: 33,
      });
      expect(data.distributionTracks.remoteConfiguration).toContainEqual({
        name: "Office Based",
        share: 67, 
      });
      
      expect(data.districtGeoData).toHaveLength(2);
      const westernGeo = data.districtGeoData.find(g => g.province === "Western");
      const centralGeo = data.districtGeoData.find(g => g.province === "Central");
      expect(westernGeo?.jobs).toBe(2);
      expect(westernGeo?.nationalShare).toBe(66.7); 
      expect(centralGeo?.jobs).toBe(1);
      expect(centralGeo?.nationalShare).toBe(33.3); 
    });

    it("should handle fallbacks safely if metadata structures are empty or missing in records", async () => {
      const incompleteJob: Partial<GoJobPostResponse> = {
        id: 104,
        employer: "Incomplete Corp",
        job_role: "Software Developer",
        is_remote: true,
        created_at: "2026-05-20T00:00:00Z",
        skills: [],
      };

      server.use(
        http.get(BACKEND_JOBS_URL, () => {
          return HttpResponse.json([incompleteJob as GoJobPostResponse]);
        })
      );

      const request = new Request("http://localhost/api/v1/analytics/dashboard");
      const response = await getDashboard(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as DashboardDataPayload;
      
      expect(data.categoryData[0].category).toBe("Unclassified Operations");
      expect(data.ingestionSources[0].name).toBe("Direct Scrape");
      expect(data.districtGeoData[0].province).toBe("Unknown Region");
      expect(data.distributionTracks.seniority[0].name).toBe("Mid-Level");
    });

    it("should gracefully handle empty dataset state without crash", async () => {
      server.use(
        http.get(BACKEND_JOBS_URL, () => {
          return HttpResponse.json([]);
        })
      );

      const request = new Request("http://localhost/api/v1/analytics/dashboard");
      const response = await getDashboard(request);
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as DashboardDataPayload;
      
      expect(data.kpiSummary.totalVacancies).toBe(0);
      expect(data.categoryData).toHaveLength(0);
      expect(data.monthlyTrends).toHaveLength(0);
      expect(data.distributionTracks.seniority).toHaveLength(0);
      expect(data.distributionTracks.remoteConfiguration).toContainEqual({
        name: "Office Based",
        share: 100,
      });
    });

    it("should catch and propagate backend errors", async () => {
      server.use(
        http.get(BACKEND_JOBS_URL, () => {
          return HttpResponse.error(); 
        })
      );

      const request = new Request("http://localhost/api/v1/analytics/dashboard");
      const response = await getDashboard(request);
      
      expect(response.status).toBe(502);
      const errBody = (await response.json()) as BffErrorResponse;
      expect(errBody.status).toBe(502);
      expect(errBody.error).toBe("Bad Gateway Trigger State");
    });
  });

  describe("GET /api/v1/analytics/category/[category]", () => {
    it("should fetch deep dive analytics for a valid category and order skills and employers by demand", async () => {
      const request = new Request("http://localhost/api/v1/analytics/category/Software%20Engineering");
      const response = await getCategoryAnalytics(request, {
        params: Promise.resolve({ category: "Software%20Engineering" }),
      });
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as CategoryAnalyticsPayload;
      
      expect(data.provinces).toHaveLength(1);
      expect(data.provinces[0].name).toBe("Western");
      expect(data.provinces[0].vacancies).toBe(1);

      expect(data.skills).toHaveLength(2);
      expect(data.skills.map(s => s.skill)).toContain("Go");
      expect(data.skills.map(s => s.skill)).toContain("Docker");
      expect(data.skills.every(s => s.category === "Software Engineering")).toBe(true);

      expect(data.employers).toHaveLength(1);
      expect(data.employers[0].name).toBe("WSO2");
      expect(data.employers[0].openRoles).toBe(1);
    });

    it("should bypass filter constraints when wildcard target 'All Categories' is specified", async () => {
      const request = new Request("http://localhost/api/v1/analytics/category/All%20Categories");
      const response = await getCategoryAnalytics(request, {
        params: Promise.resolve({ category: "All%20Categories" }),
      });
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as CategoryAnalyticsPayload;
      
      expect(data.skills).toHaveLength(4);
      
      expect(data.skills[0].skill).toBe("Docker");
      expect(data.skills[0].demand).toBe(2);

      expect(data.employers).toHaveLength(3);
      expect(data.provinces).toHaveLength(2);
    });

    it("should return empty metrics gracefully for a non-existent category", async () => {
      const request = new Request("http://localhost/api/v1/analytics/category/Marketing");
      const response = await getCategoryAnalytics(request, {
        params: Promise.resolve({ category: "Marketing" }),
      });
      
      expect(response.status).toBe(200);
      const data = (await response.json()) as CategoryAnalyticsPayload;
      
      expect(data.skills).toHaveLength(0);
      expect(data.provinces).toHaveLength(0);
      expect(data.employers).toHaveLength(0);
    });

    it("should bubble up downstream exceptions gracefully", async () => {
      server.use(
        http.get(BACKEND_JOBS_URL, () => {
          return new HttpResponse(null, { status: 500, statusText: "Internal Server Error" });
        })
      );

      const request = new Request("http://localhost/api/v1/analytics/category/Software%20Engineering");
      const response = await getCategoryAnalytics(request, {
        params: Promise.resolve({ category: "Software%20Engineering" }),
      });
      
      expect(response.status).toBe(500);
      const errBody = (await response.json()) as BffErrorResponse;
      expect(errBody.status).toBe(500);
      expect(errBody.error).toBe("Downstream core engine error");
      expect(errBody.path).toBe("/api/v1/analytics/category/Software%20Engineering");
    });
  });
});
