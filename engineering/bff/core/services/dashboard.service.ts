import { fetchAllJobsFromCore } from "../clients/go-backend.client";
import { 
  DashboardDataPayload, 
  CategoryAnalyticsPayload,
  SkillCompetencyNode,
  ProvinceDataNode,
  EmployerMetricsNode
} from "@/types/payloads";


export async function getAggregatedDashboardState(): Promise<DashboardDataPayload> {
  const rawJobs = await fetchAllJobsFromCore();
  
  const totalVacancies = rawJobs.length;
  const uniqueSectors = new Set(rawJobs.map(j => j.meta_data?.standardized_category).filter(Boolean));
  
  const skillSet = new Set<string>();
  rawJobs.forEach(j => j.skills?.forEach(s => skillSet.add(s.name)));

  const categoryCount: Record<string, number> = {};
  const platformCount: Record<string, number> = {};
  const employerMap: Record<string, { count: number; sector: string; loc: string }> = {};
  const provinceCount: Record<string, number> = {};
  const seniorityCount: Record<string, number> = {};
  const contractCount: Record<string, number> = {};
  const trendCount: Record<string, number> = {};
  const remoteCount = { Office: 0, Remote: 0 };

  rawJobs.forEach(job => {
    const cat = job.meta_data?.standardized_category || "Unclassified Operations";
    categoryCount[cat] = (categoryCount[cat] || 0) + 1;

    const src = job.meta_data?.source || "Direct Scrape";
    platformCount[src] = (platformCount[src] || 0) + 1;

    const prov = job.meta_data?.geo?.province || "Unknown Region";
    provinceCount[prov] = (provinceCount[prov] || 0) + 1;

    const senior = job.meta_data?.seniority || "Mid-Level";
    seniorityCount[senior] = (seniorityCount[senior] || 0) + 1;

    const cType = job.job_type?.name || "Full Time";
    contractCount[cType] = (contractCount[cType] || 0) + 1;

    if (job.is_remote) remoteCount.Remote++; else remoteCount.Office++;

    if (job.employer) {
      if (!employerMap[job.employer]) {
        employerMap[job.employer] = { count: 0, sector: cat, loc: job.location || "Sri Lanka" };
      }
      employerMap[job.employer].count++;
    }

    const dateRecord = job.created_at ? new Date(job.created_at) : new Date();

    const monthYearKey = dateRecord.toLocaleDateString("en-US", {
      month: "short",
      year: "numeric",
    });

    trendCount[monthYearKey] = (trendCount[monthYearKey] || 0) + 1;
  });

  const categoryData = Object.entries(categoryCount).map(([category, vacancies]) => ({ category, vacancies }));
  const ingestionSources = Object.entries(platformCount).map(([name, vacancies]) => ({ name, vacancies }));
  
  const leadingEmployers: EmployerMetricsNode[] = Object.entries(employerMap)
    .map(([name, meta]) => ({ name, location: meta.loc, sector: meta.sector, openRoles: meta.count }))
    .sort((a, b) => b.openRoles - a.openRoles)
    .slice(0, 5);

  const districtGeoData = Object.entries(provinceCount).map(([province, jobs], idx) => ({
    id: String(idx + 1),
    province,
    jobs,
    nationalShare: totalVacancies > 0 ? parseFloat(((jobs / totalVacancies) * 100).toFixed(1)) : 0
  }));

  const monthlyTrends = Object.entries(trendCount)
    .map(([month, vacancies]) => ({ month, vacancies }))
    .sort((a, b) => new Date(a.month).getTime() - new Date(b.month).getTime());

  const totalSeniority = Object.values(seniorityCount).reduce((a, b) => a + b, 0) || 1;
  const totalContract = Object.values(contractCount).reduce((a, b) => a + b, 0) || 1;

  return {
    kpiSummary: {
      totalVacancies,
      vacancyGrowthPct: 12.4, 
      sectorsTracked: uniqueSectors.size,
      skillsIdentified: skillSet.size
    },
    categoryData,
    monthlyTrends,
    ingestionSources,
    leadingEmployers,
    distributionTracks: {
      seniority: Object.entries(seniorityCount).map(([name, val]) => ({ name, share: Math.round((val / totalSeniority) * 100) })),
      contractTypes: Object.entries(contractCount).map(([name, val]) => ({ name, share: Math.round((val / totalContract) * 100) })),
      remoteConfiguration: [
        { name: "Office Based", share: totalVacancies > 0 ? Math.round((remoteCount.Office / totalVacancies) * 100) : 100 },
        { name: "Remote Available", share: totalVacancies > 0 ? Math.round((remoteCount.Remote / totalVacancies) * 100) : 0 }
      ]
    },
    districtGeoData
  };
}


export async function getCategoryDeepDiveAnalysis(targetCategory: string): Promise<CategoryAnalyticsPayload> {
  const rawJobs = await fetchAllJobsFromCore();
  
  const filtered = targetCategory === "All" || targetCategory === "All Categories" 
    ? rawJobs 
    : rawJobs.filter(j => j.meta_data?.standardized_category === targetCategory);

  const skillMetrics: Record<string, number> = {};
  const provinceMetrics: Record<string, number> = {};
  const employerMetrics: Record<string, number> = {};

  filtered.forEach(job => {
    job.skills?.forEach(s => {
      skillMetrics[s.name] = (skillMetrics[s.name] || 0) + 1;
    });
    if (job.meta_data?.geo?.province) {
      provinceMetrics[job.meta_data.geo.province] = (provinceMetrics[job.meta_data.geo.province] || 0) + 1;
    }
    if (job.employer) {
      employerMetrics[job.employer] = (employerMetrics[job.employer] || 0) + 1;
    }
  });

  const skills: SkillCompetencyNode[] = Object.entries(skillMetrics)
    .map(([skill, demand]) => ({ skill, demand, category: targetCategory }))
    .sort((a, b) => b.demand - a.demand).slice(0, 10);

  const provinces: ProvinceDataNode[] = Object.entries(provinceMetrics)
    .map(([name, vacancies], id) => ({ id: id + 1, name, vacancies }));

  const employers: EmployerMetricsNode[] = Object.entries(employerMetrics)
    .map(([name, openRoles]) => ({ name, openRoles }))
    .sort((a, b) => b.openRoles - a.openRoles).slice(0, 5);

  return { skills, provinces, employers };
}