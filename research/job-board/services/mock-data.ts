import {SkillDemand, ProvinceVacancy, HiringEmployer, CategoryAnalyticsResponse, FilterMetadataResponse, JobVacancy} from "@/types/job";


//Dashboard Page
export interface KpiSummary {
  totalVacancies: number;
  vacancyGrowthPct: number;
  sectorsTracked: number;
  skillsIdentified: number;
}

export interface MonthlyTrend {
  month: string;
  domestic: number;
  overseas: number;
}

export interface SectorShare {
  sector: string;
  vacancies: number;
}

export interface IngestionSource {
  name: string;
  vacancies: number;
}

export interface LeadingEmployer {
  name: string;
  activePosts: number;
  sector: string;
}

export interface ShareDistribution {
  level?: string;
  type?: string;
  share: number;
}

export interface DistrictVacancy {
  id: string;
  name: string;
  province: string;
  jobs: number;
  path: string;
}

export interface DashboardMetricsResponse {
  kpiSummary: KpiSummary;
  monthlyTrends: MonthlyTrend[];
  sectorDistribution: SectorShare[];
  ingestionSources: IngestionSource[];
  leadingEmployers: LeadingEmployer[];
  seniorityData: ShareDistribution[];
  contractTypes: ShareDistribution[];
  remoteConfiguration: ShareDistribution[];
  regionalVacancies: DistrictVacancy[];
}

export function getDynamicDashboardMetrics(vacanciesList: JobVacancy[]): DashboardMetricsResponse {
  const totalVacancies = vacanciesList.length;
  
  const sectors = new Set<string>();
  const skills = new Set<string>();
  
  const sectorCounts: Record<string, number> = {};
  const sourceCounts: Record<string, number> = {};
  const employerMap: Record<string, { posts: number; sector: string }> = {};
  const seniorityCounts: Record<string, number> = {};
  const contractCounts: Record<string, number> = {};
  
  let remoteCount = 0;
  const regionalMap: Record<string, number> = {};

  vacanciesList.forEach((vac) => {
    const category = vac.meta_data.standardized_category || "Unassigned";
    const source = vac.meta_data.source || "Unknown";
    const seniority = vac.meta_data.seniority || "Not Specified";
    const contract = vac.job_type || "Full Time";
    const districtName = vac.location.split(",").pop()?.trim() || "Colombo";

    sectors.add(category);
    vac.skills.forEach(s => skills.add(s));

    sectorCounts[category] = (sectorCounts[category] || 0) + 1;
    sourceCounts[source] = (sourceCounts[source] || 0) + 1;
    seniorityCounts[seniority] = (seniorityCounts[seniority] || 0) + 1;
    contractCounts[contract] = (contractCounts[contract] || 0) + 1;
    regionalMap[districtName.toLowerCase()] = (regionalMap[districtName.toLowerCase()] || 0) + 1;

    if (vac.is_remote) remoteCount++;

    if (!employerMap[vac.employer]) {
      employerMap[vac.employer] = { posts: 0, sector: category };
    }
    employerMap[vac.employer].posts++;
  });

  const sectorDistribution = Object.entries(sectorCounts)
    .map(([sector, vacancies]) => ({ sector, vacancies }))
    .sort((a, b) => b.vacancies - a.vacancies);

  const ingestionSources = Object.entries(sourceCounts)
    .map(([name, vacancies]) => ({ name, vacancies }))
    .sort((a, b) => b.vacancies - a.vacancies);

  const leadingEmployers = Object.entries(employerMap)
    .map(([name, data]) => ({ name, activePosts: data.posts, sector: data.sector }))
    .sort((a, b) => b.activePosts - a.activePosts)
    .slice(0, 5);

  const calcShare = (count: number) => totalVacancies > 0 ? Math.round((count / totalVacancies) * 100) : 0;

  const seniorityData = Object.entries(seniorityCounts).map(([level, count]) => ({
    level,
    share: calcShare(count)
  }));

  const contractTypes = Object.entries(contractCounts).map(([type, count]) => ({
    type,
    share: calcShare(count)
  }));

  const remoteAvailableShare = calcShare(remoteCount);
  const remoteConfiguration = [
    { type: "Office Based", share: 100 - remoteAvailableShare },
    { type: "Remote Available", share: remoteAvailableShare }
  ];

  const baseDistricts = [
    { id: "LK-11", name: "Colombo", province: "western", path: "M 115,220 L 125,230..." },
    { id: "LK-12", name: "Gampaha", province: "western", path: "M 100,190 L 115,210..." },
    { id: "LK-13", name: "Kalutara", province: "western", path: "M 120,260 L 130,280..." },
    { id: "LK-21", name: "Kandy", province: "central", path: "M 160,210 L 175,230..." },
    { id: "LK-31", name: "Galle", province: "southern", path: "M 140,320 L 150,340..." },
    { id: "LK-41", name: "Kurunegala", province: "northwestern", path: "M 90,150 L 110,170..." },
    { id: "LK-42", name: "Puttalam", province: "northwestern", path: "M 70,120 L 85,140..." }
  ];

  const regionalVacancies: DistrictVacancy[] = baseDistricts.map(d => ({
    ...d,
    jobs: regionalMap[d.name.toLowerCase()] || 0
  }));

  return {
    kpiSummary: {
      totalVacancies,
      vacancyGrowthPct: 14.2,
      sectorsTracked: sectors.size,
      skillsIdentified: skills.size
    },
    monthlyTrends: [
      { month: "Jan", domestic: Math.round(totalVacancies * 0.8), overseas: Math.round(totalVacancies * 0.2) },
      { month: "Feb", domestic: Math.round(totalVacancies * 0.85), overseas: Math.round(totalVacancies * 0.15) },
      { month: "Mar", domestic: totalVacancies, overseas: Math.round(totalVacancies * 0.22) }
    ],
    sectorDistribution,
    ingestionSources,
    leadingEmployers,
    seniorityData,
    contractTypes,
    remoteConfiguration,
    regionalVacancies
  };
}





//Categories Page
export const MOCK_ANALYTICS_DATA: Record<string, CategoryAnalyticsResponse> = {
  "All Categories": {
    skills: [
      { skill: "TypeScript / React", demand: 4250, category: "IT" },
      { skill: "Critical Care Nursing", demand: 3900, category: "Medicine" },
      { skill: "Cloud Architecture (AWS/Azure)", demand: 3800, category: "IT" },
      { skill: "BIM Modelling (Revit)", demand: 2450, category: "Construction" },
    ],
    provinces: [
      { province: "Western", vacancies: 14850 },
      { province: "Central", vacancies: 5900 },
      { province: "Southern", vacancies: 4200 },
    ],
    employers: [
      { name: "Virtusa Labs", openRoles: 320, location: "Western Province" },
      { name: "Asiri Health Group", openRoles: 240, location: "Multi-Region" },
    ],
  },
  "IT": {
    skills: [
      { skill: "TypeScript / React", demand: 4250, category: "IT" },
      { skill: "Cloud Architecture (AWS/Azure)", demand: 3800, category: "IT" },
      { skill: "Python & Machine Learning", demand: 3100, category: "IT" },
    ],
    provinces: [
      { province: "Western", vacancies: 11800 },
      { province: "Central", vacancies: 1400 },
    ],
    employers: [
      { name: "Virtusa Labs", openRoles: 320, location: "Western Province" },
      { name: "Sysco LABS", openRoles: 180, location: "Western Province" },
    ],
  },
  "Medicine": {
    skills: [
      { skill: "Critical Care Nursing", demand: 3900, category: "Medicine" },
      { skill: "Health Informatics & EHR", demand: 2600, category: "Medicine" },
    ],
    provinces: [
      { province: "Western", vacancies: 4100 },
      { province: "Central", vacancies: 2600 },
    ],
    employers: [
      { name: "Asiri Health Group", openRoles: 240, location: "Multi-Region" },
      { name: "Nawaloka Hospitals", openRoles: 140, location: "Western Province" },
    ],
  },
  "Construction": {
    skills: [
      { skill: "BIM Modelling (Revit)", demand: 2450, category: "Construction" },
      { skill: "Structural Engineering CAD", demand: 2200, category: "Construction" },
    ],
    provinces: [
      { province: "Western", vacancies: 5200 },
      { province: "Central", vacancies: 1900 },
    ],
    employers: [
      { name: "Access Engineering", openRoles: 195, location: "Western Province" },
      { name: "MAGA Engineering", openRoles: 160, location: "Western Province" },
    ],
  },
};




//Vcancies Page
export const MOCK_METADATA: FilterMetadataResponse = {
  categories: ["IT", "Medicine", "Construction", "Furniture Design"],
  provinces: ["Western", "Central", "Southern", "North Western", "Northern", "Eastern", "Sabaragamuwa", "Uva", "North Central"],
  contractTypes: ["Full Time", "Contract", "Part Time", "Freelance", "Internship"],
  seniorityLevels: ["Entry-level", "Junior", "Mid Level", "Senior"],
};

export const MOCK_VACANCIES: JobVacancy[] = [
  {
    id: "vac_001",
    employer: "Rimaco Furniture Pvt Ltd.",
    job_role: "3D Sketchup Designer",
    job_type: "Full Time",
    key_responsibilities: "Creating 3D designs using Sketchup software for MDF furniture manufacturing.",
    qualifications: "Ordinary Level (OL) education; ability to type in English and Sinhala; age above 18; no prior experience required.",
    location: "Boralesgamuwa, Colombo",
    offers: "Salary negotiable; priority given to applicants from Maharagama and Boralesgamuwa areas.",
    is_remote: false,
    skills: ["Sketchup", "3D Design", "MDF Manufacturing"],
    meta_data: {
      posted_at: "2026-05-14T08:30:00Z",
      source: "TopJobs.lk",
      standardized_category: "IT",
      seniority: "Entry-level",
      geo: { lat: 6.8488, lng: 79.9101, province: "Western" },
      confidence_score: 0.98,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_002",
    employer: "Softlogic Holdings PLC",
    job_role: "Junior Software Engineer",
    job_type: "Full Time",
    key_responsibilities: "Developing and maintaining web applications using React and Node.js.",
    qualifications: "Diploma or degree in IT/Computer Science; basic knowledge of JavaScript frameworks.",
    location: "Colombo 03",
    offers: "Health insurance, annual bonuses, and hybrid work opportunities.",
    is_remote: true,
    skills: ["React", "Node.js", "JavaScript", "REST APIs"],
    meta_data: {
      posted_at: "2026-05-12T09:00:00Z",
      source: "LinkedIn",
      standardized_category: "IT",
      seniority: "Junior",
      geo: { lat: 6.9271, lng: 79.8612, province: "Western" },
      confidence_score: 0.97,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_003",
    employer: "Dialog Axiata PLC",
    job_role: "Customer Support Executive",
    job_type: "Full Time",
    key_responsibilities: "Handling customer inquiries and resolving service-related issues via phone and email.",
    qualifications: "Excellent communication skills in Sinhala and English; prior customer service experience preferred.",
    location: "Kandy",
    offers: "EPF/ETF, performance incentives, and career growth opportunities.",
    is_remote: false,
    skills: ["Communication", "Customer Service", "CRM"],
    meta_data: {
      posted_at: "2026-05-10T11:15:00Z",
      source: "TopJobs.lk",
      standardized_category: "Customer Service",
      seniority: "Mid-level",
      geo: { lat: 7.2906, lng: 80.6337, province: "Central" },
      confidence_score: 0.96,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_004",
    employer: "MAS Holdings",
    job_role: "HR Assistant",
    job_type: "Full Time",
    key_responsibilities: "Assisting recruitment processes and maintaining employee records.",
    qualifications: "Degree or diploma in Human Resource Management; proficiency in MS Office.",
    location: "Katunayake",
    offers: "Transport facilities, meals, and training programs.",
    is_remote: false,
    skills: ["HR Management", "Recruitment", "MS Office"],
    meta_data: {
      posted_at: "2026-05-08T07:45:00Z",
      source: "XpressJobs",
      standardized_category: "Human Resources",
      seniority: "Entry-level",
      geo: { lat: 7.1699, lng: 79.8841, province: "Western" },
      confidence_score: 0.95,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_005",
    employer: "Commercial Bank of Ceylon",
    job_role: "Banking Assistant",
    job_type: "Full Time",
    key_responsibilities: "Managing customer transactions and supporting daily banking operations.",
    qualifications: "Passed GCE A/L; strong numerical and interpersonal skills.",
    location: "Galle",
    offers: "Medical insurance and annual performance bonuses.",
    is_remote: false,
    skills: ["Banking", "Cash Handling", "Customer Relations"],
    meta_data: {
      posted_at: "2026-05-15T10:20:00Z",
      source: "Observer Jobs",
      standardized_category: "Banking & Finance",
      seniority: "Entry-level",
      geo: { lat: 6.0535, lng: 80.2210, province: "Southern" },
      confidence_score: 0.94,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_006",
    employer: "Virtusa Sri Lanka",
    job_role: "QA Engineer",
    job_type: "Full Time",
    key_responsibilities: "Performing manual and automated software testing for enterprise applications.",
    qualifications: "Degree in Computer Science or related field; knowledge of Selenium is an advantage.",
    location: "Colombo 01",
    offers: "Flexible working hours and learning opportunities.",
    is_remote: true,
    skills: ["QA Testing", "Selenium", "Automation"],
    meta_data: {
      posted_at: "2026-05-11T14:40:00Z",
      source: "LinkedIn",
      standardized_category: "Quality Assurance",
      seniority: "Associate",
      geo: { lat: 6.9344, lng: 79.8428, province: "Western" },
      confidence_score: 0.98,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_007",
    employer: "Cargills Food City",
    job_role: "Cashier",
    job_type: "Part Time",
    key_responsibilities: "Handling billing operations and assisting customers at checkout counters.",
    qualifications: "GCE O/L qualification; friendly personality and basic computer literacy.",
    location: "Negombo",
    offers: "Flexible shifts and overtime payments.",
    is_remote: false,
    skills: ["POS Systems", "Customer Handling", "Cash Management"],
    meta_data: {
      posted_at: "2026-05-13T16:00:00Z",
      source: "TopJobs.lk",
      standardized_category: "Retail",
      seniority: "Entry-level",
      geo: { lat: 7.2084, lng: 79.8358, province: "Western" },
      confidence_score: 0.93,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_008",
    employer: "PickMe",
    job_role: "Digital Marketing Executive",
    job_type: "Full Time",
    key_responsibilities: "Planning and executing social media campaigns and digital advertisements.",
    qualifications: "Degree or diploma in Marketing; experience with Meta Ads and Google Ads preferred.",
    location: "Colombo 05",
    offers: "Performance bonuses and hybrid work model.",
    is_remote: true,
    skills: ["Digital Marketing", "SEO", "Google Ads", "Social Media"],
    meta_data: {
      posted_at: "2026-05-09T12:30:00Z",
      source: "LinkedIn",
      standardized_category: "Marketing",
      seniority: "Mid-level",
      geo: { lat: 6.8916, lng: 79.8547, province: "Western" },
      confidence_score: 0.97,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_009",
    employer: "Lanka Hospitals",
    job_role: "Receptionist",
    job_type: "Full Time",
    key_responsibilities: "Managing front desk operations and coordinating patient appointments.",
    qualifications: "Good communication skills and prior front-office experience preferred.",
    location: "Colombo 06",
    offers: "Uniform allowance and medical benefits.",
    is_remote: false,
    skills: ["Reception", "Communication", "Scheduling"],
    meta_data: {
      posted_at: "2026-05-16T08:10:00Z",
      source: "XpressJobs",
      standardized_category: "Healthcare Administration",
      seniority: "Entry-level",
      geo: { lat: 6.8741, lng: 79.8590, province: "Western" },
      confidence_score: 0.95,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },

  {
    id: "vac_010",
    employer: "Hayleys PLC",
    job_role: "Logistics Coordinator",
    job_type: "Contract",
    key_responsibilities: "Coordinating shipments and maintaining inventory and logistics records.",
    qualifications: "Diploma in Supply Chain Management; knowledge of ERP systems preferred.",
    location: "Kurunegala",
    offers: "Travel allowance and career advancement opportunities.",
    is_remote: false,
    skills: ["Logistics", "Inventory Management", "ERP"],
    meta_data: {
      posted_at: "2026-05-17T13:25:00Z",
      source: "Observer Jobs",
      standardized_category: "Supply Chain",
      seniority: "Mid-level",
      geo: { lat: 7.4863, lng: 80.3623, province: "North Western" },
      confidence_score: 0.96,
      ai_version: "gemini-2.0-flash-v1",
      error: false,
    },
  },
];