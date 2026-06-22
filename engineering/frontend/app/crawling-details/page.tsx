"use client";

import { useState, useMemo } from "react";
import { useCrawlerOverview } from "@/hooks/use-crawler";

const localization = {
  en: {
    title: "Crawling Operations & Details",
    subtitle: "Real-time visibility into engine scraping nodes, live pipelines, and database integration.",
    kpi1: "Latest Total Crawled",
    kpi2: "Active Sources",
    kpi4: "Engine Last Run",
    sectionSources: "Sources Breakdown",
    thSourceId: "ID",
    thSourceName: "Source Name",
    thJobCount: "Jobs",
    sectionHistory: "Crawling Activity Log",
    thLogId: "Run ID",
    thStartTime: "Start Time",
    thEndTime: "End Time",
    thStatus: "Status",
    bffNode: "BFF Aggregator Node Active"
  },
  si: {
    title: "දත්ත එකතු කිරීමේ විස්තර (Crawling)",
    subtitle: "බාහිර රැකියා වෙබ් අඩවිවලින් දත්ත ලබාගන්නා පද්ධති සහ සජීවී නාලිකාවල තත්ත්වය.",
    kpi1: "අවසන් මුළු රැකියා",
    kpi2: "සක්‍රීය මූලාශ්‍ර",
    kpi4: "අවසන් ධාවනය",
    sectionSources: "දත්ත මූලාශ්‍ර",
    thSourceId: "අංකය",
    thSourceName: "මූලාශ්‍ර නම",
    thJobCount: "රැකියා ගණන",
    sectionHistory: "ඉතිහාසය",
    thLogId: "ID",
    thStartTime: "ආරම්භය",
    thEndTime: "අවසන්",
    thStatus: "තත්ත්වය",
    bffNode: "BFF දත්ත සම්බන්ධකය සක්‍රීයයි"
  },
  ta: {
    title: "தரவு சேகரிப்பு விவரங்கள் (Crawling)",
    subtitle: "தரவு பிரித்தெடுத்தல் முனைகள் மற்றும் நேரடி தரவுக் குழாய்கள்.",
    kpi1: "சமீபத்திய மொத்த வேலைகள்",
    kpi2: "செயலில் உள்ள ஆதாரங்கள்",
    kpi4: "இயந்திரத்தின் கடைசி இயக்கம்",
    sectionSources: "மூலங்கள்",
    thSourceId: "ஐடி",
    thSourceName: "மூலப் பெயர்",
    thJobCount: "வேலைகள்",
    sectionHistory: "செயல்பாட்டு பதிவு",
    thLogId: "ID",
    thStartTime: "ஆரம்ப நேரம்",
    thEndTime: "முடிவு நேரம்",
    thStatus: "நிலை",
    bffNode: "BFF ஒருங்கிணைப்பு முனை செயலில் உள்ளது"
  }
};

export default function CrawlingDetailsPage() {
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");
  const { data, isLoading, isError } = useCrawlerOverview();
  
  const d = useMemo(() => localization[currentLang], [currentLang]);

  const getStatusStyle = (status: string) => {
    if (status === "RUNNING") return "bg-amber-50 text-amber-700 border-amber-200 animate-pulse";
    if (status === "FAILED") return "bg-red-50 text-red-700 border-red-200";
    return "bg-emerald-50 text-emerald-700 border-emerald-200";
  };

  const formattedDate = new Date().toLocaleDateString(
    currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
    { year: "numeric", month: "short", day: "numeric" }
  );

  if (isLoading) return <div className="p-10 text-center">Loading...</div>;
  if (isError || !data) return <div className="p-10 text-center text-red-500">Error loading data.</div>;

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800 font-sans pb-20">
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

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div className="bg-white p-5 rounded-xl border border-gray-200 shadow-sm">
            <p className="text-[10px] font-black uppercase text-gray-400">{d.kpi1}</p>
            <h3 className="text-3xl font-black text-blue-600">{data.kpis.last_crawl_job_count.toLocaleString()}</h3>
            </div>
            <div className="bg-white p-5 rounded-xl border border-gray-200 shadow-sm">
            <p className="text-[10px] font-black uppercase text-gray-400">{d.kpi2}</p>
            <h3 className="text-3xl font-black text-zinc-900">{data.sources.length}</h3>
            </div>
            <div className="bg-white p-5 rounded-xl border border-gray-200 shadow-sm">
            <p className="text-[10px] font-black uppercase text-gray-400">{d.kpi4}</p>
            <h3 className="text-lg font-black text-zinc-800">{data.kpis.gap_human}</h3>
            </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-5 gap-6">
            <div className="bg-white p-6 rounded-xl border border-gray-200 lg:col-span-2">
            <h3 className="text-sm font-bold mb-4">{d.sectionSources}</h3>
            <table className="w-full text-left">
                <tbody>
                {data.sources.map((src) => (
                    <tr key={src.id} className="border-b border-gray-50">
                    <td className="py-3 text-xs text-gray-400 font-mono">#{src.id}</td>
                    <td className="py-3 text-sm font-bold">{src.source}</td>
                    <td className="py-3 text-sm text-right font-black">{src.open_job_count.toLocaleString()}</td>
                    </tr>
                ))}
                </tbody>
            </table>
            </div>

            <div className="bg-white p-6 rounded-xl border border-gray-200 lg:col-span-3">
            <h3 className="text-sm font-bold mb-4">{d.sectionHistory}</h3>
            <table className="w-full text-left">
                <thead>
                <tr className="text-gray-400 text-xs uppercase border-b border-gray-100">
                    <th className="pb-3">{d.thLogId}</th>
                    <th className="pb-3">{d.thStartTime}</th>
                    <th className="pb-3 text-right">{d.thStatus}</th>
                </tr>
                </thead>
                <tbody>
                {data.crawler_runs.map((log) => (
                    <tr key={log.id} className="border-b border-gray-50">
                    <td className="py-3 text-xs font-mono">{log.id}</td>
                    <td className="py-3 text-xs text-gray-500">{new Date(log.started_at).toLocaleString()}</td>
                    <td className="py-3 text-right">
                        <span className={`px-2 py-0.5 rounded text-[10px] font-bold border ${getStatusStyle(log.status)}`}>
                        {log.status}
                        </span>
                    </td>
                    </tr>
                ))}
                </tbody>
            </table>
            </div>
        </div>
      </div>
    </div>
  );
}