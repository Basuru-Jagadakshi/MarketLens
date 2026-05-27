// ============================================================================
// 1. Core Map Types & Regional Configurations
// ============================================================================

export interface ProvinceConfig {
  label: string;
  color: string;
}

export interface DistrictGeoNode {
  id: string;
  name: string;
  province: string;
  jobs: number;
  path: string; // SVG path instruction coordinates
}

export const PROVINCES: Record<string, ProvinceConfig> = {
  western: { label: "Western", color: "#3b82f6" },
  central: { label: "Central", color: "#10b981" },
  southern: { label: "Southern", color: "#f59e0b" },
  northern: { label: "Northern", color: "#ef4444" },
  eastern: { label: "Eastern", color: "#8b5cf6" },
  northWestern: { label: "North Western", color: "#ec4899" },
  northCentral: { label: "North Central", color: "#06b6d4" },
  uva: { label: "Uva", color: "#f97316" },
  sabaragamuwa: { label: "Sabaragamuwa", color: "#14b8a6" }
};

// ============================================================================
// 2. Main KPI Statistics Summaries
// ============================================================================

export const kpiSummary = {
  totalVacancies: 13560,
  vacancyGrowthPct: 12.4,
  sectorsTracked: 42,
  skillsIdentified: 846
};

// ============================================================================
// 3. Sectional Graphs & Analytical Charts Data
// ============================================================================

export const sectorData1 = [
  { sector: "Software Engineering", vacancies: 5400 },
  { sector: "Data Science & AI", vacancies: 2800 },
  { sector: "DevOps & Cloud Systems", vacancies: 2100 },
  { sector: "Product Management", vacancies: 1800 },
  { sector: "Quality Assurance Automation", vacancies: 1460 }
];

export const monthlyTrends = [
  { month: "Jun 2025", domestic: 4100, overseas: 1200 },
  { month: "Jul 2025", domestic: 4400, overseas: 1350 },
  { month: "Aug 2025", domestic: 4800, overseas: 1400 },
  { month: "Sep 2025", domestic: 5100, overseas: 1600 },
  { month: "Oct 2025", domestic: 5300, overseas: 1550 },
  { month: "Nov 2025", domestic: 5200, overseas: 1700 },
  { month: "Dec 2025", domestic: 5600, overseas: 1900 },
  { month: "Jan 2026", domestic: 5900, overseas: 1850 },
  { month: "Feb 2026", domestic: 6200, overseas: 2100 },
  { month: "Mar 2026", domestic: 6400, overseas: 2300 },
  { month: "Apr 2026", domestic: 6800, overseas: 2250 },
  { month: "May 2026", domestic: 7160, overseas: 2400 }
];

export const seniorityData = [
  { level: "Junior / Associate", share: 25 },
  { level: "Mid-Level Professional", share: 45 },
  { level: "Senior Staff", share: 22 },
  { level: "Lead / Principal Architect", share: 8 }
];

export const contractTypes = [
  { type: "Full-Time Corporate", share: 70 },
  { type: "Independent Contract", share: 15 },
  { type: "Part-Time Engagement", share: 10 },
  { type: "Graduate Internship", share: 5 }
];

// Unused structural variables included for import safety compliance
export const sectorData = sectorData1; 

// ============================================================================
// 4. Detailed Sri Lanka District Geo-Coordinates (Minimal Structural Paths)
// ============================================================================

export const SRI_LANKA_DISTRICTS: DistrictGeoNode[] = [
  {
    id: "LK-11",
    name: "Colombo",
    province: "western",
    jobs: 4250,
    path: "M110,420 L130,420 L125,450 L105,445 Z"
  },
  {
    id: "LK-12",
    name: "Gampaha",
    province: "western",
    jobs: 1820,
    path: "M105,380 L135,385 L130,420 L110,420 Z"
  },
  {
    id: "LK-13",
    name: "Kalutara",
    province: "western",
    jobs: 940,
    path: "M105,445 L125,450 L120,490 L100,480 Z"
  },
  {
    id: "LK-21",
    name: "Kandy",
    province: "central",
    jobs: 1120,
    path: "M140,330 L180,340 L175,380 L135,370 Z"
  },
  {
    id: "LK-22",
    name: "Matale",
    province: "central",
    jobs: 480,
    path: "M145,270 L175,280 L180,340 L140,330 Z"
  },
  {
    id: "LK-23",
    name: "Nuwara Eliya",
    province: "central",
    jobs: 310,
    path: "M135,370 L175,380 L165,420 L130,410 Z"
  },
  {
    id: "LK-31",
    name: "Galle",
    province: "southern",
    jobs: 680,
    path: "M100,480 L120,490 L130,540 L105,530 Z"
  },
  {
    id: "LK-32",
    name: "Matara",
    province: "southern",
    jobs: 420,
    path: "M130,540 L155,545 L160,520 L120,490 Z"
  },
  {
    id: "LK-33",
    name: "Hambantota",
    province: "southern",
    jobs: 390,
    path: "M155,545 L200,510 L185,480 L160,520 Z"
  },
  {
    id: "LK-41",
    name: "Jaffna",
    province: "northern",
    jobs: 290,
    path: "M60,20 L110,30 L100,60 L65,50 Z"
  },
  {
    id: "LK-51",
    name: "Batticaloa",
    province: "eastern",
    jobs: 340,
    path: "M210,240 L240,280 L220,350 L200,300 Z"
  },
  {
    id: "LK-61",
    name: "Kurunegala",
    province: "northWestern",
    jobs: 710,
    path: "M90,290 L145,270 L140,330 L105,350 Z"
  },
  {
    id: "LK-71",
    name: "Anuradhapura",
    province: "northCentral",
    jobs: 540,
    path: "M110,150 L170,160 L165,240 L100,220 Z"
  },
  {
    id: "LK-81",
    name: "Badulla",
    province: "uva",
    jobs: 460,
    path: "M175,380 L215,390 L205,460 L165,420 Z"
  },
  {
    id: "LK-91",
    name: "Ratnapura",
    province: "sabaragamuwa",
    jobs: 400,
    path: "M125,450 L165,420 L185,480 L120,490 Z"
  }
];