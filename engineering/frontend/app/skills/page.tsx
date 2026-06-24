"use client";

import { useState, useMemo } from "react";
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid,
  Tooltip, ResponsiveContainer, Cell,
} from "recharts";
import { useIndustries, useIndustrySkillsAnalytics } from "@/hooks/use-industry";
import { SkillDemand } from "@/types/industry";

const CHART_COLOR = "#10b981";

const localization = {
  en: {
    title: "Industry Skill Analytics",
    subtitle: "Tracking talent demand curves and domain skills across active job segments",
    filterLabel: "Filter by Core Industry",
    matrixText: "Sector Matrix: 21 Verticals Available",
    kpiUniqueTitle: "Unique Skills Tracked",
    kpiUniqueDesc: "Distinct skills indexed inside selected industry.",
    kpiDemandTitle: "Most Demanding Skill",
    kpiDemandSuffix: "Requisitions Pending",
    chartTitle: "Top 15 Skills Framework Volume Graph",
    chartDesc: "This bar chart represents the highly demanding skills for the selected industry.",
    tableTitle: "Complete Industry Skills Breakdown",
    tableDesc: "This table view represents all skills for the selected industry.",
    thNo: "No",
    thSkill: "Tracked Skill",
    thVacancies: "Active Vacancies",
    sidebarTitle: "Top Hiring Employers",
    sidebarDesc: "Enterprise institutions with most hirings for the selected industry.",
    sidebarSuffix: "Vacancies",
    sidebarEmpty: "No matching hiring records.",
    bffChannel: "BFF Live Metrics Channel",
    selectIndustry: "Select an industry to view analytics",
    loading: "Loading analytics...",
    error: "Failed to load data.",
    noSkills: "No skills data available.",
  },
  si: {
    title: "කර්මාන්ත නිපුණතා විශ්ලේෂණය",
    subtitle: "සක්‍රීය රැකියා අංශ හරහා දක්ෂතා ඉල්ලුම සහ නිපුණතා මිනුම් ලුහුබැඳීම",
    filterLabel: "ප්‍රධාන කර්මාන්තය අනුව තෝරන්න",
    matrixText: "අංශ න්‍යාසය: සිරස් අංශ 21 ක් පවතී",
    kpiUniqueTitle: "හඳුනාගත් සුවිශේෂී නිපුණතා",
    kpiUniqueDesc: "තෝරාගත් කර්මාන්ත ක්ෂේත්‍රය තුළ සුචිගත කර ඇති වෙනස්ම කුසලතා ප්‍රමාණය.",
    kpiDemandTitle: "වැඩිම ඉල්ලුමක් ඇති නිපුණතාවය",
    kpiDemandSuffix: "බලාපොරොත්තු වන පුරප්පාඩු",
    chartTitle: "ඉහළම නිපුණතා 15 හි පරිමාව ප්‍රස්ථාරය",
    chartDesc: "ඉහළම මෙහෙයුම් ඝනත්ව සීමාවන් ලුහුබඳින සක්‍රීය රැකියා පුරප්පාඩු පළ කිරීම්.",
    tableTitle: "සම්පූර්ණ කර්මාන්ත නිපුණතා බිඳවැටීම",
    tableDesc: "ආයතනවල බඳවා ගැනීමේ දත්ත සක්‍රීයව බැලීම.",
    thNo: "අංකය",
    thSkill: "ලුහුබැඳි නිපුණතාවය / කුසලතාව",
    thVacancies: "සක්‍රීය පුරප්පාඩු",
    sidebarTitle: "ප්‍රමුඛ බඳවා ගන්නන්",
    sidebarDesc: "කර්මාන්තය සඳහා ඉහළම අවශ්‍යතා සහිත ප්‍රමුඛ පෙළේ සමාගම් සහ ආයතන.",
    sidebarSuffix: "නාලිකා",
    sidebarEmpty: "ගැලපෙන බඳවා ගැනීමේ වාර්තා කිසිවක් හමු නොවීය.",
    bffChannel: "BFF සජීවී දත්ත නාලිකාව",
    selectIndustry: "විශ්ලේෂණ බැලීමට කර්මාන්තයක් තෝරන්න",
    loading: "දත්ත පූරණය වෙමින්...",
    error: "දත්ත පූරණය අසාර්ථකයි.",
    noSkills: "නිපුණතා දත්ත නොමැත.",
  },
  ta: {
    title: "தொழிற்துறை திறன் பகுப்பாய்வு",
    subtitle: "செயலில் உள்ள வேலைப் பிரிவுகளில் திறமை தேவை வளைவுகள் மற்றும் திறன் அளவீடுகளைக் கண்காணித்தல்",
    filterLabel: "முதன்மை தொழிற்துறை மூலம் வடிகட்டவும்",
    matrixText: "துறை அணி: 21 பிரிவுகள் கிடைக்கின்றன",
    kpiUniqueTitle: "கண்காணிக்கப்பட்ட தனித்துவமான திறன்கள்",
    kpiUniqueDesc: "தேர்ந்தெடுக்கப்பட்ட துறையினுள் குறியிடப்பட்ட தனித்துவமான தகுதிகள்.",
    kpiDemandTitle: "அதிக தேவை உள்ள திறன்",
    kpiDemandSuffix: "நிலுவையில் உள்ள கோரிக்கைகள்",
    chartTitle: "முதல் 15 திறன்களின் கட்டமைப்பு தொகுதி வரைபடம்",
    chartDesc: "அதிக செயல்பாட்டு அடர்த்தி வரம்புகளைக் கண்காணிக்கும் செயலில் உள்ள வேலை இடுகைகளின் மொத்தத் தரவு.",
    tableTitle: "முழுதுமையான தொழிற்துறை திறன்களின் முறிவு",
    tableDesc: "நிறுவனங்களின் வேலைவாய்ப்பு அளவீடுகளை மாறும் வகையில் பார்க்க கீழே உள்ள திறன் வரிசையைத் தேர்ந்தெடுக்கவும்.",
    thNo: "எண்",
    thSkill: "கண்காணிக்கப்பட்ட தகுதி / திறன்",
    thVacancies: "செயலில் உள்ள காலியிடங்கள்",
    sidebarTitle: "முன்னணி வேலை வழங்குநர்கள்",
    sidebarDesc: "செயலில் உள்ள தேவைகளில் இந்த திறனைக் குறிப்பிடும் அதிக அளவு தேவைகளைக் கொண்ட நிறுவனங்கள்.",
    sidebarSuffix: "பகிர்வுகள்",
    sidebarEmpty: "பொருத்தமான வேலைவாய்ப்பு பதிவுகள் எதுவும் இல்லை.",
    bffChannel: "BFF நேரடி அளவீட்டு சேனல்",
    selectIndustry: "பகுப்பாய்வைக் காண தொழிற்துறையைத் தேர்ந்தெடுக்கவும்",
    loading: "தரவு ஏற்றப்படுகிறது...",
    error: "தரவு ஏற்றுவதில் தோல்வி.",
    noSkills: "திறன் தரவு இல்லை.",
  },
};

export default function IndustryPage() {
  const [selectedIndustryId, setSelectedIndustryId] = useState<number | null>(null);
  const [selectedSkillId,    setSelectedSkillId]    = useState<number | null>(null);
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");

  const d = useMemo(() => localization[currentLang], [currentLang]);

  // ── Data fetching ───────────────────────────────────────────────────────────
  const { data: industryList, isLoading: isIndustriesLoading } = useIndustries();

  const {
    data: analytics,
    isLoading: isAnalyticsLoading,
    isError: isAnalyticsError,
  } = useIndustrySkillsAnalytics(selectedIndustryId);

  // ── Derived data ────────────────────────────────────────────────────────────
  const top15Skills   = analytics?.top15_skills ?? [];
  const allSkills     = analytics?.all_skills   ?? [];
  const topEmployers  = analytics?.top_employers ?? [];

  const activeSkill: SkillDemand | null = useMemo(() => {
    if (!selectedSkillId) return allSkills[0] ?? null;
    return allSkills.find((s) => s.id === selectedSkillId) ?? allSkills[0] ?? null;
  }, [allSkills, selectedSkillId]);

  const selectedIndustryName = useMemo(() => {
    return industryList?.industries.find((i) => i.id === selectedIndustryId)?.name ?? "";
  }, [industryList, selectedIndustryId]);

  const formattedDate = new Date().toLocaleDateString(
    currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
    { year: "numeric", month: "short", day: "numeric" }
  );

  const handleIndustryChange = (value: string) => {
    setSelectedIndustryId(value ? Number(value) : null);
    setSelectedSkillId(null);
  };

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 antialiased">

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
              <button
                key={l}
                onClick={() => { setCurrentLang(l); setSelectedSkillId(null); }}
                className={`px-3 py-1 text-xs font-bold rounded-lg transition-all ${currentLang === l ? "bg-white text-blue-600 shadow-sm" : "text-gray-500 hover:text-gray-900"}`}
              >
                {l === "en" ? "English" : l === "si" ? "සිංහල" : "தமிழ்"}
              </button>
            ))}
          </div>
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white text-xs font-black shadow-md border border-white">BJ</div>
        </div>
      </header>

      <div className="p-8 space-y-6">

        {/* INDUSTRY FILTER */}
        <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="max-w-md w-full">
            <label className="block text-xs text-gray-500 mb-1 font-semibold uppercase tracking-wider">{d.filterLabel}</label>
            <select
              value={selectedIndustryId ?? ""}
              onChange={(e) => handleIndustryChange(e.target.value)}
              disabled={isIndustriesLoading}
              className="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm bg-white cursor-pointer font-medium outline-none focus:ring-2 focus:ring-zinc-900 disabled:opacity-50"
            >
              <option value="">{d.selectIndustry}</option>
              {industryList?.industries.map((ind) => (
                <option key={ind.id} value={ind.id}>{ind.name}</option>
              ))}
            </select>
          </div>
          <div className="text-xs text-gray-400 font-medium font-mono bg-gray-50 px-3 py-1.5 rounded-lg border">
            {d.matrixText}
          </div>
        </div>

        {/* NO INDUSTRY SELECTED STATE */}
        {!selectedIndustryId && (
          <div className="bg-white rounded-xl border border-gray-200/60 shadow-sm p-16 text-center">
            <p className="text-4xl mb-4">📊</p>
            <p className="text-sm text-gray-400 font-medium">{d.selectIndustry}</p>
          </div>
        )}

        {/* LOADING STATE */}
        {selectedIndustryId && isAnalyticsLoading && (
          <div className="bg-white rounded-xl border border-gray-200/60 shadow-sm p-16 flex items-center justify-center">
            <div className="text-center space-y-3">
              <div className="w-8 h-8 border-2 border-emerald-500 border-t-transparent rounded-full animate-spin mx-auto" />
              <p className="text-sm text-gray-400 font-medium">{d.loading}</p>
            </div>
          </div>
        )}

        {/* ERROR STATE */}
        {selectedIndustryId && isAnalyticsError && (
          <div className="bg-white rounded-xl border border-red-200 shadow-sm p-12 text-center">
            <p className="text-sm text-red-500 font-medium">{d.error}</p>
          </div>
        )}

        {/* ANALYTICS CONTENT */}
        {selectedIndustryId && analytics && !isAnalyticsLoading && (
          <>
            {/* KPI CARDS */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex items-center justify-between">
                <div>
                  <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider">{d.kpiUniqueTitle}</p>
                  <h3 className="text-3xl font-black text-zinc-900 mt-1">{analytics.unique_skills_count}</h3>
                  <p className="text-xs text-gray-500 mt-1">{d.kpiUniqueDesc}</p>
                </div>
                <div className="p-3.5 bg-blue-50 text-blue-600 rounded-xl font-bold text-lg">💡</div>
              </div>

              <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex items-center justify-between">
                <div>
                  <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider">{d.kpiDemandTitle}</p>
                  <h3 className="text-xl font-black text-zinc-900 mt-2 truncate max-w-[280px]">
                    {analytics.most_in_demand_skill?.skill ?? "—"}
                  </h3>
                  <p className="text-xs text-emerald-600 font-semibold mt-1">
                    🔥 {(analytics.most_in_demand_skill?.open_job_count ?? 0).toLocaleString()} {d.kpiDemandSuffix}
                  </p>
                </div>
                <div className="p-3.5 bg-emerald-50 text-emerald-600 rounded-xl font-bold text-lg">📈</div>
              </div>
            </div>

            {/* TOP 15 BAR CHART */}
            <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm">
              <div className="mb-4">
                <h3 className="text-xs font-black uppercase text-gray-400 tracking-wider">{d.chartTitle}</h3>
                <p className="text-xs text-gray-400 mt-0.5">{d.chartDesc}</p>
              </div>
              {top15Skills.length === 0 ? (
                <p className="text-center text-sm text-gray-400 py-12">{d.noSkills}</p>
              ) : (
                <ResponsiveContainer width="100%" height={340}>
                  <BarChart data={top15Skills} margin={{ top: 10, bottom: 25, left: 10, right: 10 }}>
                    <CartesianGrid strokeDasharray="3 3" stroke="#f8f9fa" vertical={false} />
                    <XAxis
                      dataKey="skill"
                      tick={{ fontSize: 10 }}
                      interval={0}
                      angle={-20}
                      textAnchor="end"
                      height={65}
                    />
                    <YAxis tick={{ fontSize: 11 }} />
                    <Tooltip formatter={(value) => [`${value}`, d.thVacancies]} />
                    <Bar dataKey="open_job_count" radius={[4, 4, 0, 0]} maxBarSize={45}>
                      {top15Skills.map((_, i) => (
                        <Cell key={i} fill={CHART_COLOR} />
                      ))}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              )}
            </div>

            {/* MASTER-DETAIL: ALL SKILLS TABLE + EMPLOYERS SIDEBAR */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">

              {/* ALL SKILLS TABLE */}
              <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm lg:col-span-2">
                <div className="mb-4">
                  <h3 className="text-sm font-bold text-zinc-900">{d.tableTitle}</h3>
                  <p className="text-xs text-gray-400 mt-0.5">{d.tableDesc}</p>
                </div>
                <div className="overflow-y-auto max-h-[480px] border border-gray-100 rounded-xl shadow-inner">
                  <table className="w-full text-left border-collapse">
                    <thead className="sticky top-0 z-10">
                      <tr className="border-b border-gray-200 bg-gray-50/90 backdrop-blur-xs">
                        <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-16">{d.thNo}</th>
                        <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400">{d.thSkill}</th>
                        <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-44 text-right">{d.thVacancies}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {allSkills.length === 0 ? (
                        <tr>
                          <td colSpan={3} className="py-12 text-center text-sm text-gray-400">{d.noSkills}</td>
                        </tr>
                      ) : (
                        allSkills.map((s, i) => {
                          const isSelected = activeSkill?.id === s.id;
                          return (
                            <tr
                              key={s.id}
                              onClick={() => setSelectedSkillId(s.id)}
                              className={`border-b border-gray-50 cursor-pointer transition-all ${isSelected ? "bg-zinc-900 text-white hover:bg-zinc-800" : "hover:bg-gray-50/80"}`}
                            >
                              <td className="py-3.5 px-4 text-xs font-mono text-gray-400">{i + 1}</td>
                              <td className={`py-3.5 px-4 text-sm font-semibold ${isSelected ? "text-white" : "text-gray-900"}`}>{s.skill}</td>
                              <td className={`py-3.5 px-4 font-mono text-sm text-right ${isSelected ? "text-emerald-400 font-bold" : "text-gray-700"}`}>
                                {s.open_job_count.toLocaleString()}
                              </td>
                            </tr>
                          );
                        })
                      )}
                    </tbody>
                  </table>
                </div>
              </div>

              {/* EMPLOYERS SIDEBAR */}
              <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm sticky top-[100px]">
                <div className="flex flex-col gap-1 mb-4">
                  <h3 className="text-xs font-black uppercase text-gray-400 tracking-wider">{d.sidebarTitle}</h3>
                  {activeSkill && (
                    <div className="mt-1">
                      <span className="inline-block text-[11px] font-black uppercase tracking-wide px-2.5 py-1 rounded-lg border max-w-full truncate border-emerald-500 text-emerald-600 bg-emerald-500/5">
                        🎯 {activeSkill.skill}
                      </span>
                    </div>
                  )}
                </div>
                <p className="text-xs text-gray-400 mb-4 leading-relaxed">{d.sidebarDesc}</p>

                <div className="space-y-3 min-h-[280px]">
                  {topEmployers.length === 0 ? (
                    <div className="text-center py-12 text-xs text-gray-400">{d.sidebarEmpty}</div>
                  ) : (
                    topEmployers.map((employer) => (
                      <div key={employer.id} className="p-3.5 border border-gray-100 rounded-xl bg-gray-50/50 flex items-center justify-between shadow-xs">
                        <div className="space-y-0.5 pr-2">
                          <div className="text-sm font-bold text-gray-900 tracking-tight truncate max-w-[160px]">{employer.name}</div>
                        </div>
                        <div className="text-right min-w-[65px]">
                          <div className="text-sm font-black text-blue-600 font-mono">{employer.open_job_count}</div>
                          <div className="text-[9px] uppercase font-black text-gray-400 tracking-tighter">{d.sidebarSuffix}</div>
                        </div>
                      </div>
                    ))
                  )}
                </div>

                <div className="mt-6 pt-4 border-t border-gray-100 flex items-center justify-between text-[11px] text-gray-400 font-medium">
                  <span>{d.bffChannel}</span>
                  <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
                </div>
              </div>

            </div>
          </>
        )}
      </div>
    </div>
  );
}