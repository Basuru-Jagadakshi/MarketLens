"use client";

import { useState } from "react";
import {
  XAxis, YAxis, Tooltip, ResponsiveContainer,
  PieChart, Pie, Cell, Legend,
  BarChart, Bar, LineChart, Line,
} from "recharts";
import {
  useDashboardOverview,
  useOccupationAnalytics,
  useIndustryAnalytics,
} from "@/hooks/use-dashboard";
import { OccupationStat, IndustryStat } from "@/types/dashboard";

const currentYearNum = new Date().getFullYear();
const DYNAMIC_YEARS = [
  String(currentYearNum),
  String(currentYearNum - 1),
  String(currentYearNum - 2),
];

const CHART_COLORS = ["#2563eb", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#6366f1", "#ec4899"];

const TRANSLATIONS = {
  en: {
    title: "Labour Market Demand Dashboard",
    subtitle: "National Strategic Overview Driven by SLSO & SLSIC Registries",
    vacancies: "Current Vacancies",
    occFramework: "Occupations Framework",
    indFramework: "Industries Framework",
    slsoBased: "Based on SLSO",
    slsicBased: "Based on SLSIC",
    collapse: "Click to collapse",
    expand: "Click to view full breakdown",
    matrixTitle: "Registered Framework Classifications Matrix",
    open: "open",
    occChartTitle: "Current Job Distribution by Occupation (SLSO)",
    occChartSub: "Horizontal mapping showing all 10 standard occupation bands",
    indChartTitle: "Current Job Distribution by Industry (SLSIC)",
    indChartSub: "Vertical bar chart projection featuring rotated X-axis headers for all 21 divisions",
    analyticsBtn: "See Analytics",
    expChartTitle: "Current Job Distribution by Experience",
    eduChartTitle: "Current Job Distribution by Education Level",
    remoteTitle: "Remote / On-Site Configuration",
    contractTitle: "Contract Type Share",
    occModalTitle: "Occupation Analytics Breakdown",
    occModalSub: "Macro trend patterns & top jobs mapped by tier",
    indModalTitle: "Industry Segment Workspace Matrix",
    indModalSub: "Past year performance indices categorized by specific market sectors",
    selectOccLabel: "Select Occupation",
    trendChartHeader: "Variation Over Past Years (Historical Demand Trend)",
    demandingJobsHeader: "Current Demanding Jobs for",
    sectorTrendHeader: "Sector Variant Level Across Years",
    expAllocHeader: "Experience Allocation Distribution",
    regionalShareHeader: "Regional Province Share Allocation",
    eduThresholdHeader: "Minimum Educational Level Threshold",
    topEnterpriseHeader: "Top Sector Enterprise Hiring Groups",
    selectIndLabel: "Select Industry",
  },
  si: {
    title: "ශ්‍රම වෙළඳපල ඉල්ලුම උපකරණ පුවරුව",
    subtitle: "SLSO සහ SLSIC ලේඛන මගින් මෙහෙයවන ජාතික උපායමාර්ගික දළ විශ්ලේෂණය",
    vacancies: "වත්මන් පුරප්පාඩු",
    occFramework: "වෘත්තීය රාමුව",
    indFramework: "කර්මාන්ත රාමුව",
    slsoBased: "SLSO මත පදනම්ව",
    slsicBased: "SLSIC මත පදනම්ව",
    collapse: "හකුලන්න ක්ලික් කරන්න",
    expand: "සම්පූර්ණ විස්තරය බැලීමට ක්ලික් කරන්න",
    matrixTitle: "ලියාපදිංචි රාමු වර්ගීකරණ අනුකෘතිය",
    open: "විවෘතයි",
    occChartTitle: "වෘත්තිය අනුව වත්මන් රැකියා ව්‍යාප්තිය (SLSO)",
    occChartSub: "ප්‍රමිතිගත වෘත්තීය කාණ්ඩ 10ම පෙන්වන තිරස් සිතියම්කරණය",
    indChartTitle: "කර්මාන්තය අනුව වත්මන් රැකියා ව්‍යාප්තිය (SLSIC)",
    indChartSub: "අංශ 21 සඳහාම භ්‍රමණය වූ X-අක්ෂ ශීර්ෂයන් සහිත සිරස් තීරු ප්‍රස්තාරය",
    analyticsBtn: "විශ්ලේෂණ බලන්න",
    expChartTitle: "අත්දැකීම් අනුව වත්මන් රැකියා ව්‍යාප්තිය",
    eduChartTitle: "අධ්‍යාපන මට්ටම අනුව වත්මන් රැකියා ව්‍යාප්තිය",
    remoteTitle: "දුරස්ථ / සේවා ස්ථානගත වින්‍යාසය",
    contractTitle: "කොන්ත්‍රාත්තු වර්ගයේ කොටස",
    occModalTitle: "වෘත්තීය විශ්ලේෂණ බිඳවැටීම",
    occModalSub: "ස්ථර අනුව සිතියම්ගත කරන ලද සාර්ව ප්‍රවණතා රටා සහ ඉහළම රැකියා",
    indModalTitle: "කර්මාන්ත අංශ වැඩ අවකාශ අනුකෘතිය",
    indModalSub: "නිශ්චිත වෙළඳපල අංශ අනුව වර්ගීකරණය කරන ලද පසුගිය වසරේ කාර්ය සාධන දර්ශක",
    selectOccLabel: "වෘත්තිය තෝරන්න",
    trendChartHeader: "පසුගිය වසරවල විචලනය (ඓතිහාසික ඉල්ලුමේ ප්‍රවණතාවය)",
    demandingJobsHeader: "සඳහා ඉල්ලුමක් ඇති වත්මන් රැකියා",
    sectorTrendHeader: "වසර පුරා අංශ විචල්‍ය මට්ටම",
    expAllocHeader: "අත්දැකීම් වෙන් කිරීමේ ව්‍යාප්තිය",
    regionalShareHeader: "ප්‍රාදේශීය පළාත් කොටස් වෙන් කිරීම",
    eduThresholdHeader: "අවම අධ්‍යාපන මට්ටමේ සීමාව",
    topEnterpriseHeader: "ඉහළම අංශයේ ව්‍යවසාය බඳවා ගැනීමේ කණ්ඩායම්",
    selectIndLabel: "කර්මාන්තය තෝරන්න",
  },
  ta: {
    title: "தொழில் சந்தை தேவை தகவல் பலகை",
    subtitle: "SLSO & SLSIC பதிவேடுகளால் இயக்கப்படும் தேசிய மூலோபாய கண்ணோட்டம்",
    vacancies: "தற்போதைய காலியிடங்கள்",
    occFramework: "தொழில் கட்டமைப்பு",
    indFramework: "தொழில்துறை கட்டமைப்பு",
    slsoBased: "SLSO இன் அடிப்படையில்",
    slsicBased: "SLSIC இன் அடிப்படையில்",
    collapse: "சுருக்க கிளிக் செய்யவும்",
    expand: "முழு விபரங்களையும் பார்க்க கிளிக் செய்யவும்",
    matrixTitle: "பதிவுசெய்யப்பட்ட கட்டமைப்பு வகைப்பாடு அணி",
    open: "காலியிடம்",
    occChartTitle: "தொழில் வாரியான தற்போதைய வேலை விநியோகம் (SLSO)",
    occChartSub: "அனைத்து 10 நிலையான தொழில் குழுக்களையும் காட்டும் கிடைமட்ட வரைபடம்",
    indChartTitle: "தொழில்துறை வாரியான தற்போதைய வேலை விநியோகம் (SLSIC)",
    indChartSub: "அனைத்து 21 பிரிவுகளுக்கான சுழற்றப்பட்ட X-அச்சு தலைப்புகளைக் கொண்ட செங்குத்து பட்டை வரைபடம்",
    analyticsBtn: "பகுப்பாய்வைக் காண்க",
    expChartTitle: "அனுபவ வாரியான தற்போதைய வேலை விநியோகம்",
    eduChartTitle: "கல்வித் தகுதி வாரியான தற்போதைய வேலை விநியோகம்",
    remoteTitle: "தொலைதூர / தள வேலை கட்டமைப்பு",
    contractTitle: "ஒப்பந்த வகை பங்கீடு",
    occModalTitle: "தொழில் பகுப்பாய்வு முறிவு",
    occModalSub: "மேக்ரோ போக்கு வடிவங்கள் மற்றும் அடுக்கு வாரியாக வரைபடமாக்கப்பட்ட சிறந்த வேலைகள்",
    indModalTitle: "தொழில்துறை பிரிவு பணிவெளி அணி",
    indModalSub: "குறிப்பிட்ட சந்தைத் துறைகளால் வகைப்படுத்தப்பட்ட கடந்த ஆண்டு செயல்திறன் குறியீடுகள்",
    selectOccLabel: "தொழிலைத் தேர்ந்தெடுக்கவும்",
    trendChartHeader: "கடந்த ஆண்டுகளின் மாறுபாடு (வரலாற்று தேவை போக்கு)",
    demandingJobsHeader: "விருப்பமுள்ள வேலைகள்",
    sectorTrendHeader: "ஆண்டுகள் முழுவதும் துறை மாறுபாட்டின் அளவு",
    expAllocHeader: "அனுபவ ஒதுக்கீடு விநியோகம்",
    regionalShareHeader: "பிராந்திய மாகாணப் பங்கு ஒதுக்கீடு",
    eduThresholdHeader: "கையெழுத்து கல்வித் தகுதி வரம்பு",
    topEnterpriseHeader: "முன்னணி துறை நிறுவன வேலைவாய்ப்பு குழுக்கள்",
    selectIndLabel: "தொழில்துறையைத் தேர்ந்தெடுக்கவும்",
  },
};

export default function DashboardPage() {
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");
  const [activeKpiRow, setActiveKpiRow] = useState<"SLSO" | "SLSIC" | null>(null);
  const [activeModal, setActiveModal] = useState<"OCCUPATION" | "INDUSTRY" | null>(null);

  // IDs for analytics panels — null until user opens the panel
  const [selectedOccId, setSelectedOccId] = useState<number | null>(null);
  const [selectedOccName, setSelectedOccName] = useState<string>("");
  const [selectedIndId, setSelectedIndId] = useState<number | null>(null);
  const [selectedIndName, setSelectedIndName] = useState<string>("");
  const [analyticsYear, setAnalyticsYear] = useState(Number(DYNAMIC_YEARS[0]));

  const d = TRANSLATIONS[currentLang];

  // ── Data fetching ─────────────────────────────────────────────────────────
  const { data: overview, isLoading, isError } = useDashboardOverview();
  const { data: occAnalytics, isLoading: occLoading } = useOccupationAnalytics(selectedOccId);
  const { data: indAnalytics, isLoading: indLoading } = useIndustryAnalytics(selectedIndId, analyticsYear);

  const formattedDate = new Date().toLocaleDateString(
    currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
    { year: "numeric", month: "short", day: "numeric" }
  );

  // ── Derived chart data from API ───────────────────────────────────────────
  const occupationChartData = overview?.by_occupation.occupations.map((o) => ({
    name: o.name,
    count: o.open_job_count,
    id: o.id,
  })) ?? [];

  const industryChartData = overview?.by_industry.industries.map((i) => ({
    name: i.name,
    count: i.open_job_count,
    id: i.id,
  })) ?? [];

  const experienceChartData = overview?.by_experience.experiences.map((e) => ({
    name: e.name,
    value: e.open_job_count,
  })) ?? [];

  const educationChartData = overview?.by_education.education_levels.map((e) => ({
    name: e.level,
    value: e.open_job_count,
  })) ?? [];

  const remoteCount   = overview?.remote_vs_onsite.remote_count  ?? 0;
  const onsiteCount   = overview?.remote_vs_onsite.on_site_count ?? 0;
  const totalRemote   = remoteCount + onsiteCount;
  const remotePct     = totalRemote > 0 ? Math.round((remoteCount  / totalRemote) * 100) : 0;
  const onsitePct     = totalRemote > 0 ? Math.round((onsiteCount  / totalRemote) * 100) : 0;

  const jobTypeData = overview?.by_job_type.job_types ?? [];
  const totalJobTypes = jobTypeData.reduce((sum, jt) => sum + jt.open_job_count, 0);

  // ── Occupation panel helpers ──────────────────────────────────────────────
  const handleOpenOccModal = () => {
    // Pre-select the first occupation if none selected
    if (!selectedOccId && occupationChartData.length > 0) {
      setSelectedOccId(occupationChartData[0].id);
      setSelectedOccName(occupationChartData[0].name);
    }
    setActiveModal("OCCUPATION");
  };

  const handleOpenIndModal = () => {
    if (!selectedIndId && industryChartData.length > 0) {
      setSelectedIndId(industryChartData[0].id);
      setSelectedIndName(industryChartData[0].name);
    }
    setActiveModal("INDUSTRY");
  };

  // ── Occupation analytics derived data ────────────────────────────────────
  const occTrendData = occAnalytics?.yearly_trend?.yearly_trend?.map((t) => ({
    year: String(t.year),
    vacancies: t.open_job_count,
  })) ?? [];

  const topJobRoles = occAnalytics?.top_job_roles.top_job_roles ?? [];

  // ── Industry analytics derived data ──────────────────────────────────────
  const indTrendData = indAnalytics?.yearly_trend?.yearly_trend?.map((t) => ({
    year: String(t.year),
    vacancies: t.open_job_count,
  })) ?? [];

  const indExpData = indAnalytics?.by_experience.experiences.map((e) => ({
    label: e.name,
    value: e.open_job_count,
  })) ?? [];

  const indProvinceData = indAnalytics?.by_province.provinces.map((p) => ({
    name: p.province,
    value: p.open_job_count,
  })) ?? [];

  const indEduData = indAnalytics?.by_education.education_levels.map((e) => ({
    label: e.level,
    value: e.open_job_count,
  })) ?? [];

  const indEmployers = indAnalytics?.top_employers.employers ?? [];

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center space-y-3">
          <div className="w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto" />
          <p className="text-sm text-gray-400 font-medium">Loading dashboard data...</p>
        </div>
      </div>
    );
  }

  if (isError) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center text-red-500 text-sm font-medium">
          Failed to load dashboard. Please try again.
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800 font-sans pb-20">

      {/* HEADER */}
      <header className="bg-white border-b border-gray-100 px-8 py-5 flex flex-col md:flex-row justify-between items-start md:items-center gap-4 sticky top-0 z-40">
        <div>
          <h1 className="text-xl font-black text-gray-900 tracking-tight">{d.title}</h1>
          <p className="text-xs text-gray-400 mt-1 font-medium">{d.subtitle}</p>
        </div>
        <div className="flex flex-wrap items-center gap-4 ml-auto md:ml-0 w-full md:w-auto justify-end">
          <div className="flex items-center gap-2 bg-gray-50 px-3 py-1.5 rounded-xl text-xs text-gray-500 font-medium shadow-inner">
            <svg className="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
              <path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 002-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            <span>{formattedDate}</span>
          </div>
          <div className="flex bg-gray-100 p-1 rounded-xl shadow-sm">
            {(["en", "si", "ta"] as const).map((l) => (
              <button key={l} onClick={() => setCurrentLang(l)} className={`px-3 py-1 text-xs font-bold rounded-lg transition-all ${currentLang === l ? "bg-white text-blue-600 shadow-sm" : "text-gray-500 hover:text-gray-900"}`}>
                {l === "en" ? "English" : l === "si" ? "සිංහල" : "தமிழ்"}
              </button>
            ))}
          </div>
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white text-xs font-black shadow-md border border-white">BJ</div>
        </div>
      </header>

      <div className="p-8 space-y-8 max-w-[1650px] mx-auto">

        {/* TIER 1: KPI CARDS */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="bg-white p-6 rounded-xl shadow-sm flex flex-col justify-between">
            <div>
              <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">{d.vacancies}</p>
              <h3 className="text-3xl font-black text-gray-900 mt-2">
                {overview?.active_jobs.active_job_count.toLocaleString()}
              </h3>
            </div>
            <div className={`mt-4 flex items-center text-xs font-bold ${overview?.active_jobs.trend === "up" ? "text-emerald-600" : overview?.active_jobs.trend === "down" ? "text-red-500" : "text-gray-400"}`}>
              {overview?.active_jobs.trend === "up" ? "▲" : overview?.active_jobs.trend === "down" ? "▼" : "─"}{" "}
              {Math.abs(overview?.active_jobs.change_percent ?? 0).toFixed(1)}% vs last month
            </div>
          </div>

          <div onClick={() => setActiveKpiRow(activeKpiRow === "SLSO" ? null : "SLSO")} className={`bg-white p-6 rounded-xl transition-all cursor-pointer shadow-sm flex flex-col justify-between ${activeKpiRow === "SLSO" ? "ring-2 ring-blue-500" : ""}`}>
            <div>
              <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">{d.occFramework}</p>
              <h3 className="text-3xl font-black text-gray-900 mt-2">{overview?.by_occupation.count}</h3>
              <p className="text-[11px] font-semibold text-gray-400 mt-1">{d.slsoBased}</p>
            </div>
            <p className="mt-4 text-[11px] text-blue-600 font-medium">{activeKpiRow === "SLSO" ? d.collapse : d.expand}</p>
          </div>

          <div onClick={() => setActiveKpiRow(activeKpiRow === "SLSIC" ? null : "SLSIC")} className={`bg-white p-6 rounded-xl transition-all cursor-pointer shadow-sm flex flex-col justify-between ${activeKpiRow === "SLSIC" ? "ring-2 ring-emerald-500" : ""}`}>
            <div>
              <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">{d.indFramework}</p>
              <h3 className="text-3xl font-black text-gray-900 mt-2">{overview?.by_industry.count}</h3>
              <p className="text-[11px] font-semibold text-gray-400 mt-1">{d.slsicBased}</p>
            </div>
            <p className="mt-4 text-[11px] text-emerald-600 font-medium">{activeKpiRow === "SLSIC" ? d.collapse : d.expand}</p>
          </div>
        </div>

        {/* EXPANDABLE MATRIX */}
        {activeKpiRow && (
          <div className="bg-white p-6 rounded-xl shadow-inner">
            <h4 className="text-xs font-bold text-gray-500 uppercase tracking-wider mb-4 font-mono">{d.matrixTitle} ({activeKpiRow})</h4>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
              {(activeKpiRow === "SLSO" ? occupationChartData : industryChartData).map((item, idx) => (
                <div key={idx} className="flex justify-between items-center p-3 rounded-lg bg-gray-50 text-xs">
                  <span className="font-semibold text-gray-700 truncate mr-2">{item.name}</span>
                  <span className="font-mono font-bold bg-white px-2.5 py-1 rounded text-gray-600">{item.count} {d.open}</span>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* TIER 2: PRIMARY CHARTS */}
        <div className="space-y-6">
          <div className="bg-white p-5 rounded-xl shadow-sm h-[520px] flex flex-col">
            <div className="flex justify-between items-center border-b border-gray-50 pb-2">
              <div>
                <h4 className="text-sm font-bold text-gray-800">{d.occChartTitle}</h4>
                <p className="text-[11px] text-gray-400">{d.occChartSub}</p>
              </div>
              <button onClick={handleOpenOccModal} className="text-xs font-bold text-blue-600 bg-blue-50 hover:bg-blue-100 px-3 py-1.5 rounded-lg transition-all">{d.analyticsBtn}</button>
            </div>
            <div className="flex-1 mt-4">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={occupationChartData} layout="vertical" margin={{ left: 20, right: 20 }}>
                  <XAxis type="number" tick={{ fontSize: 10 }} />
                  <YAxis dataKey="name" type="category" tick={{ fontSize: 10 }} width={160} />
                  <Tooltip />
                  <Bar dataKey="count" fill="#2563eb" radius={[0, 4, 4, 0]} barSize={16} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          <div className="bg-white p-5 rounded-xl shadow-sm h-[620px] flex flex-col">
            <div className="flex justify-between items-center border-b border-gray-50 pb-2">
              <div>
                <h4 className="text-sm font-bold text-gray-800">{d.indChartTitle}</h4>
                <p className="text-[11px] text-gray-400">{d.indChartSub}</p>
              </div>
              <button onClick={handleOpenIndModal} className="text-xs font-bold text-emerald-600 bg-emerald-50 hover:bg-emerald-100 px-3 py-1.5 rounded-lg transition-all">{d.analyticsBtn}</button>
            </div>
            <div className="flex-1 mt-4 pb-16">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={industryChartData}>
                  <XAxis dataKey="name" tick={{ fontSize: 9 }} interval={0} angle={-45} textAnchor="end" height={100} tickFormatter={(value) => value.length > 20 ? `${value.substring(0, 20)}...` : value} />
                  <YAxis type="number" tick={{ fontSize: 10 }} />
                  <Tooltip />
                  <Bar dataKey="count" fill="#10b981" radius={[4, 4, 0, 0]} barSize={22} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>

        {/* TIER 3: EXPERIENCE & EDUCATION */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div className="bg-white p-5 rounded-xl shadow-sm h-[400px] flex flex-col">
            <h4 className="text-sm font-bold text-gray-800 border-b border-gray-50 pb-2">{d.expChartTitle}</h4>
            <div className="flex-1 mt-4">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={experienceChartData} margin={{ top: 10, right: 10, left: -20, bottom: 5 }}>
                  <XAxis dataKey="name" tick={{ fontSize: 10 }} />
                  <YAxis tick={{ fontSize: 10 }} />
                  <Tooltip />
                  <Bar dataKey="value" fill="#f59e0b" radius={[4, 4, 0, 0]} barSize={30} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          </div>

          <div className="bg-white p-5 rounded-xl shadow-sm h-[400px] flex flex-col">
            <h4 className="text-sm font-bold text-gray-800 border-b border-gray-50 pb-2">{d.eduChartTitle}</h4>
            <div className="flex-1 flex items-center justify-center">
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie data={educationChartData} cx="50%" cy="45%" innerRadius={65} outerRadius={90} paddingAngle={3} dataKey="value">
                    {educationChartData.map((_, i) => <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                  </Pie>
                  <Legend verticalAlign="bottom" iconSize={8} iconType="circle" wrapperStyle={{ fontSize: 11 }} />
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </div>
          </div>
        </div>

        {/* TIER 4: REMOTE & JOB TYPE PROGRESS BARS */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 pt-4">
          <div className="bg-white p-6 rounded-xl shadow-sm">
            <h3 className="text-sm font-bold text-gray-800 mb-4 flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-indigo-600" />{d.remoteTitle}
            </h3>
            <div className="space-y-4">
              {[
                { label: "On-Site", share: onsitePct },
                { label: "Remote",  share: remotePct },
              ].map((item, i) => (
                <div key={i}>
                  <div className="flex justify-between text-xs mb-1 font-medium text-gray-600">
                    <span>{item.label}</span>
                    <span className="font-mono">{item.share}%</span>
                  </div>
                  <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                    <div className="h-full bg-indigo-600 transition-all duration-1000" style={{ width: `${item.share}%` }} />
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="bg-white p-6 rounded-xl shadow-sm">
            <h3 className="text-sm font-bold text-gray-800 mb-4 flex items-center gap-2">
              <span className="w-2 h-2 rounded-full bg-emerald-600" />{d.contractTitle}
            </h3>
            <div className="space-y-4">
              {jobTypeData.map((jt, i) => {
                const pct = totalJobTypes > 0 ? Math.round((jt.open_job_count / totalJobTypes) * 100) : 0;
                return (
                  <div key={i}>
                    <div className="flex justify-between text-xs mb-1 font-medium text-gray-600">
                      <span>{jt.type}</span>
                      <span className="font-mono">{pct}%</span>
                    </div>
                    <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                      <div className="h-full bg-emerald-600 transition-all duration-1000" style={{ width: `${pct}%` }} />
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      </div>

      {/* SLIDING PANEL */}
      <div className={`fixed inset-y-0 right-0 z-50 w-full max-w-4xl bg-white shadow-2xl border-l transform transition-transform duration-300 ease-in-out ${activeModal ? "translate-x-0" : "translate-x-full"}`}>
        <div className="h-full flex flex-col">
          <div className="p-6 border-b flex justify-between items-center bg-white sticky top-0 z-10">
            <div>
              <h2 className="text-lg font-black text-gray-900">
                {activeModal === "OCCUPATION" ? d.occModalTitle : d.indModalTitle}
              </h2>
              <p className="text-xs text-gray-400 mt-1">
                {activeModal === "OCCUPATION" ? d.occModalSub : d.indModalSub}
              </p>
            </div>
            <button onClick={() => setActiveModal(null)} className="p-2 hover:bg-gray-100 rounded-lg text-gray-500 font-bold text-lg">✕</button>
          </div>

          <div className="flex-1 overflow-y-auto p-6">

            {/* OCCUPATION PANEL */}
            {activeModal === "OCCUPATION" && (
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <div className="bg-gray-50 p-4 rounded-xl space-y-2 h-fit">
                  <label className="text-xs font-bold text-gray-400 uppercase block mb-2">{d.selectOccLabel}</label>
                  {occupationChartData.map((occ) => (
                    <button
                      key={occ.id}
                      onClick={() => { setSelectedOccId(occ.id); setSelectedOccName(occ.name); }}
                      className={`w-full text-left p-3 text-xs rounded-xl transition-all ${selectedOccId === occ.id ? "bg-blue-600 text-white font-bold shadow-sm" : "bg-white hover:bg-gray-50"}`}
                    >
                      {occ.name}
                    </button>
                  ))}
                </div>

                <div className="md:col-span-2 space-y-6">
                  {occLoading ? (
                    <div className="flex items-center justify-center h-40">
                      <div className="w-6 h-6 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
                    </div>
                  ) : (
                    <>
                      <div className="bg-white p-4 rounded-xl border border-gray-100 shadow-sm">
                        <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">{d.trendChartHeader}</h4>
                        <div className="h-44">
                          <ResponsiveContainer>
                            <LineChart data={occTrendData}>
                              <XAxis dataKey="year" />
                              <YAxis />
                              <Tooltip />
                              <Line type="monotone" dataKey="vacancies" stroke="#2563eb" strokeWidth={3} />
                            </LineChart>
                          </ResponsiveContainer>
                        </div>
                      </div>

                      <div className="bg-gray-50 p-4 rounded-xl">
                        <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">{d.demandingJobsHeader} {selectedOccName}</h4>
                        <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                          {topJobRoles.map((role, i) => (
                            <div key={i} className="bg-white p-3 rounded-lg text-center shadow-sm border border-gray-100">
                              <p className="text-xs font-bold text-gray-800">{role.job_role}</p>
                              <p className="text-[10px] text-emerald-600 font-semibold mt-1">{role.open_job_count} open</p>
                            </div>
                          ))}
                        </div>
                      </div>
                    </>
                  )}
                </div>
              </div>
            )}

            {/* INDUSTRY PANEL */}
            {activeModal === "INDUSTRY" && (
              <div className="space-y-6">
                <div className="flex flex-wrap items-center gap-3 border-b pb-4">
                  <select
                    value={selectedIndId ?? ""}
                    onChange={(e) => {
                      const found = industryChartData.find((i) => i.id === Number(e.target.value));
                      if (found) { setSelectedIndId(found.id); setSelectedIndName(found.name); }
                    }}
                    className="bg-gray-50 rounded-xl p-2 text-xs font-bold border border-gray-200"
                  >
                    {industryChartData.map((ind) => (
                      <option key={ind.id} value={ind.id}>{ind.name}</option>
                    ))}
                  </select>

                  <div className="flex bg-gray-100 p-1 rounded-xl">
                    {DYNAMIC_YEARS.map((year) => (
                      <button
                        key={year}
                        onClick={() => setAnalyticsYear(Number(year))}
                        className={`px-3 py-1 text-[11px] font-bold rounded-lg ${analyticsYear === Number(year) ? "bg-white text-emerald-600 shadow-sm" : "text-gray-500"}`}
                      >
                        {year}
                      </button>
                    ))}
                  </div>
                </div>

                {indLoading ? (
                  <div className="flex items-center justify-center h-40">
                    <div className="w-6 h-6 border-2 border-emerald-600 border-t-transparent rounded-full animate-spin" />
                  </div>
                ) : (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    <div className="bg-white p-4 rounded-xl border border-gray-100">
                      <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">{d.sectorTrendHeader}</h4>
                      <div className="h-40">
                        <ResponsiveContainer>
                          <LineChart data={indTrendData}>
                            <XAxis dataKey="year" />
                            <YAxis />
                            <Tooltip />
                            <Line type="monotone" dataKey="vacancies" stroke="#10b981" strokeWidth={3} />
                          </LineChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div className="bg-white p-4 rounded-xl border border-gray-100">
                      <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">{d.expAllocHeader}</h4>
                      <div className="h-40">
                        <ResponsiveContainer>
                          <BarChart data={indExpData}>
                            <XAxis dataKey="label" tick={{ fontSize: 10 }} />
                            <YAxis />
                            <Tooltip />
                            <Bar dataKey="value" fill="#3b82f6" radius={[4, 4, 0, 0]} />
                          </BarChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div className="bg-white p-4 rounded-xl border border-gray-100">
                      <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">{d.regionalShareHeader}</h4>
                      <div className="h-40">
                        <ResponsiveContainer>
                          <PieChart>
                            <Pie data={indProvinceData} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={45}>
                              {indProvinceData.map((_, i) => <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                            </Pie>
                            <Tooltip />
                          </PieChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div className="bg-white p-4 rounded-xl border border-gray-100">
                      <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">{d.eduThresholdHeader}</h4>
                      <div className="h-40">
                        <ResponsiveContainer>
                          <BarChart data={indEduData}>
                            <XAxis dataKey="label" tick={{ fontSize: 10 }} />
                            <YAxis />
                            <Tooltip />
                            <Bar dataKey="value" fill="#f59e0b" radius={[4, 4, 0, 0]} />
                          </BarChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div className="bg-gray-50 p-4 rounded-xl md:col-span-2 shadow-sm">
                      <h4 className="text-xs font-bold text-gray-500 uppercase mb-3">{d.topEnterpriseHeader}</h4>
                      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                        {indEmployers.map((emp, i) => (
                          <div key={i} className="bg-white p-3 rounded-lg flex justify-between items-center text-xs shadow-sm border border-gray-100">
                            <span className="font-semibold text-gray-700">{emp.name}</span>
                            <span className="font-mono bg-emerald-50 text-emerald-600 font-bold px-2 py-0.5 rounded-md">{emp.open_job_count}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>

      {activeModal && <div onClick={() => setActiveModal(null)} className="fixed inset-0 bg-black/10 backdrop-blur-[2px] z-40" />}
    </div>
  );
}