"use client";

import { useState, useMemo } from "react";

// Localized Mock Data matching requested properties precisely
const CRAWLING_MOCK_DATA = {
  en: {
    kpis: {
      totalCrawled: 3420,
      activeSources: 4,
      lastRuntime: "14 mins ago"
    },
    sources: [
      { id: 1, name: "LinkedIn Jobs", job_count: 1450 },
      { id: 2, name: "TopJobs.lk", job_count: 920 },
      { id: 3, name: "Roar Jobs", job_count: 610 },
      { id: 4, name: "Government Gazette", job_count: 440 }
    ],
    history: [
      { id: "CR-9082", start_time: "2026-06-18 10:00:00", end_time: "2026-06-18 10:02:04", status: "Complete" },
      { id: "CR-9081", start_time: "2026-06-18 10:35:12", end_time: "--:--:--", status: "Running" },
      { id: "CR-9080", start_time: "2026-06-18 08:00:00", end_time: "2026-06-18 08:01:50", status: "Complete" },
      { id: "CR-9079", start_time: "2026-06-18 06:00:00", end_time: "2026-06-18 06:01:29", status: "Complete" },
      { id: "CR-9078", start_time: "2026-06-17 23:00:00", end_time: "2026-06-17 23:10:12", status: "Failed" }
    ]
  },
  si: {
    kpis: {
      totalCrawled: 3150,
      activeSources: 4,
      lastRuntime: "මිනිත්තු 14 කට පෙර"
    },
    sources: [
      { id: 1, name: "LinkedIn Jobs", job_count: 1320 },
      { id: 2, name: "TopJobs.lk", job_count: 880 },
      { id: 3, name: "Roar Jobs", job_count: 560 },
      { id: 4, name: "රජයේ ගැසට් පත්‍රය", job_count: 390 }
    ],
    history: [
      { id: "CR-9082", start_time: "2026-06-18 10:00:00", end_time: "2026-06-18 10:02:04", status: "සාර්ථකයි" },
      { id: "CR-9081", start_time: "2026-06-18 10:35:12", end_time: "--:--:--", status: "ධාවනය වේ" },
      { id: "CR-9080", start_time: "2026-06-18 08:00:00", end_time: "2026-06-18 08:01:50", status: "සාර්ථකයි" },
      { id: "CR-9079", start_time: "2026-06-18 06:00:00", end_time: "2026-06-18 06:01:29", status: "සාර්ථකයි" },
      { id: "CR-9078", start_time: "2026-06-17 23:00:00", end_time: "2026-06-17 23:10:12", status: "අසාර්ථකයි" }
    ]
  },
  ta: {
    kpis: {
      totalCrawled: 3240,
      activeSources: 4,
      lastRuntime: "14 நிமிடங்களுக்கு முன்பு"
    },
    sources: [
      { id: 1, name: "LinkedIn Jobs", job_count: 1380 },
      { id: 2, name: "TopJobs.lk", job_count: 900 },
      { id: 3, name: "Roar Jobs", job_count: 580 },
      { id: 4, name: "அரசு வர்த்தமானி", job_count: 380 }
    ],
    history: [
      { id: "CR-9082", start_time: "2026-06-18 10:00:00", end_time: "2026-06-18 10:02:04", status: "நிறைவடைந்தது" },
      { id: "CR-9081", start_time: "2026-06-18 10:35:12", end_time: "--:--:--", status: "இயங்குகிறது" },
      { id: "CR-9080", start_time: "2026-06-18 08:00:00", end_time: "2026-06-18 08:01:50", status: "நிறைவடைந்தது" },
      { id: "CR-9079", start_time: "2026-06-18 06:00:00", end_time: "2026-06-18 06:01:29", status: "நிறைவடைந்தது" },
      { id: "CR-9078", start_time: "2026-06-17 23:00:00", end_time: "2026-06-17 23:10:12", status: "தோல்வியடைந்தது" }
    ]
  }
};

const localization = {
  en: {
    title: "Crawling Operations & Details",
    subtitle: "Real-time visibility into engine scraping nodes, live pipelines, and database integration structures",
    kpi1: "Latest Total Crawled",
    kpi1Desc: "Total job records synced during the latest operational batch interval.",
    kpi2: "Monitored Sources",
    kpi4: "Engine Last Running",
    sectionSources: "Sources Breakdown",
    thSourceId: "ID",
    thSourceName: "Source Name",
    thJobCount: "Crawled Job Count",
    sectionHistory: "Crawling Activity Log",
    thLogId: "ID",
    thStartTime: "Start Time",
    thEndTime: "End Time",
    thStatus: "Status",
    bffNode: "BFF Aggregator Node Active"
  },
  si: {
    title: "දත්ත එකතු කිරීමේ විස්තර (Crawling)",
    subtitle: "බාහිර රැකියා වෙබ් අඩවිවලින් දත්ත ලබාගන්නා පද්ධති සහ සජීවී නාලිකාවල තත්ත්වය",
    kpi1: "අලුත්ම මුළු රැකියා සංඛ්‍යාව",
    kpi1Desc: "අවසන් වරට පද්ධතිය ක්‍රියාත්මක වීමේදී එකතු කරන ලද මුළු රැකියා වාර්තා ප්‍රමාණය.",
    kpi2: "සක්‍රීය මූලාශ්‍ර",
    kpi4: "අවසන් ධාවන කාලය",
    sectionSources: "දත්ත මූලාශ්‍ර විස්තරය",
    thSourceId: "අංකය",
    thSourceName: "මූලාශ්‍ර නම",
    thJobCount: "එකතු කළ රැකියා ගණන",
    sectionHistory: "දත්ත එකතු කිරීමේ ඉතිහාසය",
    thLogId: "අංකය",
    thStartTime: "ආරම්භක වේලාව",
    thEndTime: "අවසන් වූ වේලාව",
    thStatus: "තත්ත්වය",
    bffNode: "BFF දත්ත සම්බන්ධකය සක්‍රීයයි"
  },
  ta: {
    title: "தரவு சேகரிப்பு விவரங்கள் (Crawling)",
    subtitle: "தரவு பிரித்தெடுத்தல் முனைகள், நேரடி தரவுக் குழாய்கள் மற்றும் தரவுத்தள ஒருங்கிணைப்பு பற்றிய நிகழ்நேர பார்வை",
    kpi1: "சமீபத்திய மொத்த வேலைகள்",
    kpi1Desc: "சமீபத்திய செயல்பாட்டு தொகுதியின் போது ஒத்திசைக்கப்பட்ட மொத்த வேலை பதிவுகள்.",
    kpi2: "கண்காணிக்கப்படும் ஆதாரங்கள்",
    kpi4: "இயந்திரத்தின் கடைசி இயக்கம்",
    sectionSources: "மூலங்கள் முறிவு",
    thSourceId: "ஐடி",
    thSourceName: "மூலப் பெயர்",
    thJobCount: "சேகரிக்கப்பட்ட வேலைகள்",
    sectionHistory: "தரவு சேகரிப்பு செயல்பாட்டு பதிவு",
    thLogId: "ஐடி",
    thStartTime: "ஆரம்ப நேரம்",
    thEndTime: "முடிவு நேரம்",
    thStatus: "நிலை",
    bffNode: "BFF ஒருங்கிணைப்பு முனை செயலில் உள்ளது"
  }
};

export default function CrawlingDetailsPage() {
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");

  const d = useMemo(() => localization[currentLang], [currentLang]);
  const data = useMemo(() => CRAWLING_MOCK_DATA[currentLang], [currentLang]);

  const formattedDate = new Date().toLocaleDateString(currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK", {
    year: "numeric",
    month: "short",
    day: "numeric"
  });

  const getStatusStyle = (status: string) => {
    const s = status.toLowerCase();
    if (s.includes("run") || s.includes("ධාවන") || s.includes("இயங்கு")) {
      return "bg-amber-50 text-amber-700 border border-amber-200 animate-pulse";
    }
    if (s.includes("fail") || s.includes("අසාර්ථක") || s.includes("தோல்வி")) {
      return "bg-red-50 text-red-700 border border-red-200";
    }
    return "bg-emerald-50 text-emerald-700 border border-emerald-200";
  };

  return (
    <div className="min-h-screen bg-gray-50 text-gray-900 antialiased">
      
      {/* --- HEADER --- */}
      <header className="bg-white border-b border-gray-100 px-8 py-5 flex flex-col md:flex-row justify-between items-start md:items-center gap-4 sticky top-0 z-40">
        <div>
          <h1 className="text-xl font-black text-gray-900 tracking-tight">{d.title}</h1>
          <p className="text-xs text-gray-400 mt-1 font-medium">{d.subtitle}</p>
        </div>
        
        <div className="flex flex-wrap items-center gap-4 ml-auto md:ml-0 w-full md:w-auto justify-end">
          <div className="flex items-center gap-2 bg-gray-50 px-3 py-1.5 rounded-xl text-xs text-gray-500 font-medium shadow-inner">
            <svg className="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2"><path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 002-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
            <span>{formattedDate}</span>
          </div>

          {/* <div className="flex bg-gray-100 p-1 rounded-xl shadow-sm">
            {(["en", "si", "ta"] as const).map((l) => (
              <button key={l} onClick={() => setCurrentLang(l)} className={`px-3 py-1 text-xs font-bold rounded-lg transition-all ${currentLang === l ? "bg-white text-blue-600 shadow-sm" : "text-gray-500 hover:text-gray-900"}`}>
                {l === "en" ? "English" : l === "si" ? "සිංහල" : "தமிழ்"}
              </button>
            ))}
          </div> */}
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white text-xs font-black shadow-md border border-white">BJ</div>
        </div>
      </header>

      <div className="p-8 space-y-6">
        
        {/* --- KPI SUMMARY GRID (3 Columns Only Now) --- */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-6">
          <div className="bg-white p-5 rounded-xl border border-gray-200/60 shadow-sm">
            <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider">{d.kpi1}</p>
            <h3 className="text-3xl font-black text-blue-600 mt-1">{data.kpis.totalCrawled.toLocaleString()}</h3>
          </div>

          <div className="bg-white p-5 rounded-xl border border-gray-200/60 shadow-sm flex flex-col justify-between">
            <div>
              <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider">{d.kpi2}</p>
              <h3 className="text-3xl font-black text-zinc-900 mt-1">{data.kpis.activeSources}</h3>
            </div>
            <p className="text-[11px] text-emerald-600 font-bold mt-2">✔️ Channels Active</p>
          </div>

          <div className="bg-white p-5 rounded-xl border border-gray-200/60 shadow-sm flex flex-col justify-between">
            <div>
              <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider">{d.kpi4}</p>
              <h3 className="text-lg font-black text-zinc-800 mt-2 truncate">{data.kpis.lastRuntime}</h3>
            </div>
            <div className="flex items-center gap-1.5 text-[11px] text-gray-400 mt-2">
              <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
              <span>Cron Daemon Healthy</span>
            </div>
          </div>
        </div>

        {/* --- SPLIT SIDE-BY-SIDE TWO COLUMN LAYOUT --- */}
        <div className="grid grid-cols-1 lg:grid-cols-5 gap-6 items-start">
          
          {/* LEFT SIDE COLUMN: SOURCES TABLE (Takes 2/5 width) */}
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm lg:col-span-2">
            <div className="mb-4">
              <h3 className="text-sm font-bold text-zinc-900 tracking-tight">{d.sectionSources}</h3>
            </div>
            
            <div className="overflow-x-auto border border-gray-100 rounded-xl shadow-inner">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-gray-200 bg-gray-50">
                    <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-16">{d.thSourceId}</th>
                    <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400">{d.thSourceName}</th>
                    <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-36 text-right">{d.thJobCount}</th>
                  </tr>
                </thead>
                <tbody>
                  {data.sources.map((src) => (
                    <tr key={src.id} className="border-b border-gray-50 hover:bg-gray-50/50 transition-colors">
                      <td className="py-3.5 px-4 text-xs font-mono font-bold text-gray-400">#{src.id}</td>
                      <td className="py-3.5 px-4 text-sm font-bold text-gray-900">{src.name}</td>
                      <td className="py-3.5 px-4 text-sm font-mono font-black text-gray-700 text-right">{src.job_count.toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* RIGHT SIDE COLUMN: CRAWLING RESTRUCTURING LOG TABLE (Takes 3/5 width) */}
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm lg:col-span-3">
            <div className="mb-4">
              <h3 className="text-sm font-bold text-zinc-900 tracking-tight">{d.sectionHistory}</h3>
            </div>

            <div className="overflow-x-auto border border-gray-100 rounded-xl shadow-inner">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-gray-200 bg-gray-50">
                    <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-24">{d.thLogId}</th>
                    <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-44">{d.thStartTime}</th>
                    <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-44">{d.thEndTime}</th>
                    <th className="py-3 px-4 text-xs font-bold uppercase tracking-wider text-gray-400 w-28 text-right">{d.thStatus}</th>
                  </tr>
                </thead>
                <tbody>
                  {data.history.map((log) => (
                    <tr key={log.id} className="border-b border-gray-50 hover:bg-gray-50/80 transition-colors">
                      <td className="py-3.5 px-4 text-xs font-mono font-bold text-gray-400">{log.id}</td>
                      <td className="py-3.5 px-4 text-xs font-mono text-gray-500">{log.start_time}</td>
                      <td className="py-3.5 px-4 text-xs font-mono text-gray-500">{log.end_time}</td>
                      <td className="py-3.5 px-4 text-xs text-right">
                        <span className={`px-2 py-0.5 rounded-md font-bold text-[10px] uppercase inline-block ${getStatusStyle(log.status)}`}>
                          {log.status}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            <div className="mt-5 pt-4 border-t border-gray-100 flex items-center justify-between text-[11px] text-gray-400 font-medium">
              <span>{d.bffNode}</span>
              <span className="w-1.5 h-1.5 rounded-full bg-blue-500 animate-pulse" />
            </div>
          </div>

        </div>
      </div>
    </div>
  );
}