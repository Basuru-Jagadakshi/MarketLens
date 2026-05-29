import { describe, it, expect, vi, beforeEach } from "vitest";
import { 
  mapToGoVacancyModel, 
  getFilteredVacancies, 
  getFilterLookupMetadata 
} from "./vacancy.service"; 
import { fetchAllJobsFromCore } from "../clients/go-backend.client";
import { GoJobPostResponse, GoJobMetaData } from "@/types/models";

vi.mock("../clients/go-backend.client", () => ({
  fetchAllJobsFromCore: vi.fn(),
}));

describe("Vacancy Service Unit Tests", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const mockDataset: GoJobPostResponse[] = [
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
      created_at: "2026-05-29T00:00:00Z",
      meta_data: {
        id: 1,
        job_post_id: 1,
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
      employer: "Dialog Axiata",
      job_role: "Data Analyst",
      key_responsibilities: "Build reporting metrics.",
      qualifications: "BSc in Statistics",
      location: "Trace Expert City",
      is_remote: false,
      job_type_id: 2,
      job_type: { id: 2, name: "Contract" },
      created_at: "2026-05-29T00:00:00Z",
      meta_data: {
        id: 2,
        job_post_id: 2,
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
  ];

  describe("mapToGoVacancyModel", () => {
    it("should mapping structure flawlessly and preserve coordinate fallback constants", () => {
      const targetJob = mockDataset[0];
      const mappedOutput = mapToGoVacancyModel(targetJob);

      expect(mappedOutput.id).toBe(101);
      expect(mappedOutput.employer).toBe("WSO2");
      expect(mappedOutput.job_type).toBe("Full Time");
      expect(mappedOutput.offers).toBe("Salary Negotiable."); 
      expect(mappedOutput.skills).toEqual(["Go", "Docker"]);
      expect(mappedOutput.meta_data.source).toBe("LinkedIn");
      expect(mappedOutput.meta_data.geo.province).toBe("Western");
      expect(mappedOutput.meta_data.geo.lat).toBe(6.9271); 
      expect(mappedOutput.meta_data.geo.lng).toBe(79.8612);
      expect(isNaN(Date.parse(mappedOutput.meta_data.posted_at))).toBe(false);
    });

    it("should handle missing properties safely using fallback string design conventions", () => {
      const brokenJob: Partial<GoJobPostResponse> = {
        id: 500,
        employer: "Unknown Corp",
        meta_data: {} as GoJobMetaData, 
        skills: undefined,    
      };

      const mappedOutput = mapToGoVacancyModel(brokenJob as GoJobPostResponse);

      expect(mappedOutput.job_type).toBe("Full Time");
      expect(mappedOutput.skills).toEqual([]);
      expect(mappedOutput.meta_data.standardized_category).toBe("Unclassified");
      expect(mappedOutput.meta_data.seniority).toBe("Associate");
      expect(mappedOutput.meta_data.geo.province).toBe("Western");
    });
  });

  describe("getFilteredVacancies", () => {
    it("should correctly isolate records based on multiple filter parameters simultaneously", async () => {
      vi.mocked(fetchAllJobsFromCore).mockResolvedValue(mockDataset);

      const activeFilters = {
        category: "Software Engineering",
        seniority: "Senior", 
      };

      const results = await getFilteredVacancies(activeFilters);

      expect(results).toHaveLength(1);
      expect(results[0].id).toBe(101); 
      expect(results[0].employer).toBe("WSO2");
    });

    it("should treat 'All' parameters as wildcards and skip exclusion filters", async () => {
      vi.mocked(fetchAllJobsFromCore).mockResolvedValue(mockDataset);

      const activeFilters = {
        category: "All",
        province: "Western",
        contractType: "All",
        seniority: "All",
      };

      const results = await getFilteredVacancies(activeFilters);
      expect(results).toHaveLength(2); 
    });

    it("should return an empty array if no records meet the specified filtering matrix", async () => {
      vi.mocked(fetchAllJobsFromCore).mockResolvedValue(mockDataset);

      const impossibleFilters = {
        category: "Data Science",
        seniority: "Senior", 
      };

      const results = await getFilteredVacancies(impossibleFilters);
      expect(results).toHaveLength(0);
    });
  });

  describe("getFilterLookupMetadata", () => {
    it("should extract unique values and format clean selection option matrices", async () => {
      const duplicateProneDataset: GoJobPostResponse[] = [
        ...mockDataset,
        {
          id: 103,
          employer: "Dialog Axiata",
          job_role: "Data Analyst",
          key_responsibilities: "Build reporting metrics.",
          qualifications: "BSc in Statistics",
          location: "Trace Expert City",
          is_remote: false,
          job_type_id: 2,
          job_type: { id: 2, name: "Contract" },
          created_at: "2026-05-29T00:00:00Z",
          meta_data: {
            id: 2,
            job_post_id: 2,
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
        } as unknown as GoJobPostResponse,
      ];

      vi.mocked(fetchAllJobsFromCore).mockResolvedValue(duplicateProneDataset);

      const metadata = await getFilterLookupMetadata();

      expect(metadata.categories).toHaveLength(2);
      expect(metadata.categories).toContainEqual({ value: "Software Engineering" });
      expect(metadata.categories).toContainEqual({ value: "Data Science" }); 

      expect(metadata.provinces).toHaveLength(1);
      expect(metadata.provinces[0]).toEqual({ id: 1, name: "Western" });

      expect(metadata.contractTypes).toHaveLength(2);
      expect(metadata.contractTypes).toContainEqual({ value: "Full Time" });
      expect(metadata.contractTypes).toContainEqual({ value: "Contract" });

      expect(metadata.seniorityLevels).toHaveLength(2);
      expect(metadata.seniorityLevels).toContainEqual({ value: "Senior" });
      expect(metadata.seniorityLevels).toContainEqual({ value: "Junior" });
    });
  });
});