import { describe, it, expect, vi, beforeEach } from "vitest";
import { getAggregatedDashboardState, getCategoryDeepDiveAnalysis} from "./dashboard.service";
import { fetchAllJobsFromCore } from "../clients/go-backend.client";
import { GoJobPostResponse, GoJobMetaData } from "@/types/models";


vi.mock("../clients/go-backend.client", () => ({
  fetchAllJobsFromCore: vi.fn(),
}));

describe("Dashboard Service Unit Tests", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should correctly aggregate KPI data, metrics, and sort monthly trends chronologically", async () => {
    
    const mockRawJobs: GoJobPostResponse[] = [
      {
        id: 1,
        employer: "WSO2",
        job_role: "Software Engineer",
        key_responsibilities: "Write code",
        qualifications: "BSc",
        location: "Colombo",
        is_remote: true,
        job_type_id: 1,
        job_type: { id: 1, name: "Full Time" },
        created_at: "2026-03-15T00:00:00Z", 
        meta_data: {
          id:1,
          job_post_id: 1,
          standardized_category: "Software Engineering",
          source: "LinkedIn",
          seniority: "Junior",
          confidence_score: 98.9,
          ai_version: "chatgpt",
          error: false,
          geo_id: 1,
          geo: { province: "Western" },
        },
        skills: [{ id: 1, name: "Go" }, { id: 2, name: "Java" }],
      },
      {
        id: 2,
        employer: "Sysco Labs",
        job_role: "Senior QA Engineer",
        key_responsibilities: "Test code",
        qualifications: "BSc",
        location: "Colombo",
        is_remote: false,
        job_type_id: 1,
        job_type: { id: 1, name: "Full Time" },
        created_at: "2026-01-10T00:00:00Z", 
        meta_data: {
          id:2,
          job_post_id: 2,
          standardized_category: "Quality Assurance",
          source: "TopJobs",
          seniority: "Senior",
          confidence_score: 98.9,
          ai_version: "chatgpt",
          error: false,
          geo_id: 2,
          geo: { province: "Western" },
        },
        skills: [{ id: 2, name: "Java" }, { id: 3, name: "Selenium" }],
      },
    ];

    vi.mocked(fetchAllJobsFromCore).mockResolvedValue(mockRawJobs);

    const result = await getAggregatedDashboardState();
    
    // Test KPI Summary Metrics
    expect(result.kpiSummary.totalVacancies).toBe(2);
    expect(result.kpiSummary.sectorsTracked).toBe(2); 
    expect(result.kpiSummary.skillsIdentified).toBe(3); 

    // Test Chronological Monthly Trend Sorting
    expect(result.monthlyTrends).toHaveLength(2);
    expect(result.monthlyTrends[0].month).toBe("Jan 2026"); 
    expect(result.monthlyTrends[0].vacancies).toBe(1);
    expect(result.monthlyTrends[1].month).toBe("Mar 2026");
    expect(result.monthlyTrends[1].vacancies).toBe(1);

    // Test Geographical Aggregations
    expect(result.districtGeoData).toHaveLength(1);
    expect(result.districtGeoData[0].province).toBe("Western");
    expect(result.districtGeoData[0].jobs).toBe(2);
    expect(result.districtGeoData[0].nationalShare).toBe(100.0);

    // Test Structural Distributions (Percentages)
    expect(result.distributionTracks.remoteConfiguration).toContainEqual({
      name: "Remote Available",
      share: 50,
    });
    expect(result.distributionTracks.remoteConfiguration).toContainEqual({
      name: "Office Based",
      share: 50,
    });
  });

  it("should handle fallbacks gracefully when metadata attributes are missing", async () => {
    const mockMinimalJob: Partial<GoJobPostResponse> = {
      id: 3,
      employer: "", 
      is_remote: false,
      job_type: { id: 1, name: "Contract" },
      created_at: "2026-05-20T00:00:00Z",
      meta_data: {} as GoJobMetaData, 
      skills: [],
    };

    vi.mocked(fetchAllJobsFromCore).mockResolvedValue([mockMinimalJob as GoJobPostResponse]);

    const result = await getAggregatedDashboardState();

    expect(result.categoryData[0].category).toBe("Unclassified Operations");
    expect(result.ingestionSources[0].name).toBe("Direct Scrape");
    expect(result.districtGeoData[0].province).toBe("Unknown Region");
    expect(result.leadingEmployers).toHaveLength(0); 
  });

  it("should handle empty database state gracefully without throwing divide-by-zero errors", async () => {

    vi.mocked(fetchAllJobsFromCore).mockResolvedValue([]);

    const result = await getAggregatedDashboardState();

    expect(result.kpiSummary.totalVacancies).toBe(0);
    expect(result.categoryData).toHaveLength(0);
    expect(result.monthlyTrends).toHaveLength(0);
    expect(result.districtGeoData).toHaveLength(0);
    
    expect(result.distributionTracks.seniority).toHaveLength(0);
    expect(result.distributionTracks.remoteConfiguration).toContainEqual({
      name: "Office Based",
      share: 100,
    });
  });
});


describe("getCategoryDeepDiveAnalysis Unit Tests", () => {
    
    const mockDeepDiveJobs: GoJobPostResponse[] = [
      {
        id: 1,
        employer: "WSO2",
        job_role: "Go Developer",
        key_responsibilities: "BFF dev",
        qualifications: "BSc",
        location: "Colombo",
        is_remote: true,
        job_type_id: 1,
        job_type: { id: 1, name: "Full Time" },
        created_at: "2026-05-28T00:00:00Z",
        meta_data: {
          id:1,
          job_post_id: 1,
          standardized_category: "Software Engineering",
          source: "TopJobs",
          seniority: "Senior",
          confidence_score: 98.9,
          ai_version: "chatgpt",
          error: false,
          geo_id: 2,
          geo: { province: "Western" },
        },
        skills: [{ id: 1, name: "Go" }, { id: 2, name: "Docker" }],
      },
      {
        id: 2,
        employer: "Sysco Labs",
        job_role: "Java Architect",
        key_responsibilities: "Spring microservices",
        qualifications: "BSc",
        location: "Colombo",
        is_remote: false,
        job_type_id: 1,
        job_type: { id: 1, name: "Full Time" },
        created_at: "2026-05-28T00:00:00Z",
        meta_data: {
          id:2,
          job_post_id: 2,
          standardized_category: "Software Engineering",
          source: "TopJobs",
          seniority: "Senior",
          confidence_score: 98.9,
          ai_version: "chatgpt",
          error: false,
          geo_id: 2,
          geo: { province: "Western" },
        },
        skills: [{ id: 2, name: "Docker" }, { id: 3, name: "Java" }],
      },
      {
        id: 3,
        employer: "Dialog",
        job_role: "Data Analyst",
        key_responsibilities: "Build dashboards",
        qualifications: "BSc",
        location: "Colombo",
        is_remote: false,
        job_type_id: 1,
        job_type: { id: 1, name: "Full Time" },
        created_at: "2026-05-28T00:00:00Z",
        meta_data: {
          id:3,
          job_post_id: 3,
          standardized_category: "Data Science",
          source: "TopJobs",
          seniority: "Senior",
          confidence_score: 98.9,
          ai_version: "chatgpt",
          error: false,
          geo_id: 2,
          geo: { province: "Central" },
        },
        skills: [{ id: 4, name: "Python" }],
      }
    ];

    it("should filter datasets accurately and extract top sorted skills/employers for a target category", async () => {
      vi.mocked(fetchAllJobsFromCore).mockResolvedValue(mockDeepDiveJobs);

      const result = await getCategoryDeepDiveAnalysis("Software Engineering");

      expect(result.provinces).toHaveLength(1);
      expect(result.provinces[0].name).toBe("Western");
      expect(result.provinces[0].vacancies).toBe(2);

      expect(result.skills).toHaveLength(3);
      expect(result.skills[0].skill).toBe("Docker"); 
      expect(result.skills[0].demand).toBe(2);
      expect(result.skills[0].category).toBe("Software Engineering");

      expect(result.employers).toHaveLength(2);
      expect(result.employers.map(e => e.name)).toEqual(["WSO2", "Sysco Labs"]);
    });

    it("should bypass filtering conditions when wildcard categories are specified", async () => {
      vi.mocked(fetchAllJobsFromCore).mockResolvedValue(mockDeepDiveJobs);

      const result = await getCategoryDeepDiveAnalysis("All Categories");

      const totalProvincesVacancies = result.provinces.reduce((sum, p) => sum + (p.vacancies || 0), 0);
      expect(totalProvincesVacancies).toBe(3); 
      
      expect(result.skills).toHaveLength(4); 
      expect(result.employers).toHaveLength(3); 
    });

    it("should return empty arrays safely if the target category does not match any records", async () => {
      vi.mocked(fetchAllJobsFromCore).mockResolvedValue(mockDeepDiveJobs);

      const result = await getCategoryDeepDiveAnalysis("Marketing");
      
      expect(result.skills).toHaveLength(0);
      expect(result.provinces).toHaveLength(0);
      expect(result.employers).toHaveLength(0);
    });
  });