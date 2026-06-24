"use client";

import { useState, useMemo } from "react";
import { useVacancyMetadata, useVacancies } from "@/hooks/use-vacancies";
import { Vacancy } from "@/types/vacancy";

const localization = {
  en: {
    title: "Vacancy Explorer",
    subtitle: "vacancies synchronized from job postings web sites",
    searchPlaceholder: "Search role or employer...",
    allIndustries: "All Industries",
    allProvinces: "All Provinces",
    allJobTypes: "All Job Types",
    allExperiences: "All Experience Levels",
    thNo: "No.",
    thRole: "Job Role",
    thEmployer: "Employer",
    thProvince: "Province",
    thRemote: "Work Mode",
    remoteText: "Remote Available",
    officeText: "Office Based",
    noRecords: "No active listings found.",
    syncText: "Syncing Workspace Data via BFF...",
    errorText: "Failed to connect to backend.",
    panelEmployer: "Employer",
    panelRemote: "Work Mode",
    panelLocation: "Location",
    panelExperience: "Experience",
    panelIndustry: "Industry",
    panelJobType: "Job Type",
    panelDescription: "Description",
    panelSkills: "Skills",
    panelPostedAt: "Posted",
    closePanel: "Close",
    noSkills: "No skills listed",
  },
  si: {
    title: "රැකියා ගවේෂකය",
    subtitle: "රැකියා දැන්වීම් වෙබ් අඩවි වලින් සමමුහුර්ත කරන ලද පුරප්පාඩු",
    searchPlaceholder: "තනතුර හෝ ආයතනය සොයන්න...",
    allIndustries: "සියලුම කර්මාන්ත",
    allProvinces: "සියලුම පළාත්",
    allJobTypes: "සියලුම රැකියා වර්ග",
    allExperiences: "සියලුම අත්දැකීම්",
    thNo: "අංකය",
    thRole: "තනතුර",
    thEmployer: "ආයතනය",
    thProvince: "පළාත",
    thRemote: "සේවා ක්‍රමය",
    remoteText: "දුරස්ථ සේවය",
    officeText: "කාර්යාලීය",
    noRecords: "ප්‍රතිඵල නොමැත.",
    syncText: "දත්ත සම්බන්ධ වෙමින් පවතී...",
    errorText: "සම්බන්ධ විය නොහැක.",
    panelEmployer: "ආයතනය",
    panelRemote: "සේවා ක්‍රමය",
    panelLocation: "ස්ථානය",
    panelExperience: "අත්දැකීම",
    panelIndustry: "කර්මාන්තය",
    panelJobType: "රැකියා වර්ගය",
    panelDescription: "විස්තරය",
    panelSkills: "නිපුණතා",
    panelPostedAt: "පළ කළ දිනය",
    closePanel: "වසන්න",
    noSkills: "නිපුණතා නොමැත",
  },
  ta: {
    title: "வேலைவாய்ப்பு ஆய்வாளர்",
    subtitle: "வேலைவாய்ப்பு இணையதளங்களிலிருந்து ஒத்திசைக்கப்பட்ட காலியிடங்கள்",
    searchPlaceholder: "பதவி அல்லது நிறுவனத்தைத் தேடு...",
    allIndustries: "அனைத்துத் தொழில்துறைகள்",
    allProvinces: "அனைத்து மாகாணங்கள்",
    allJobTypes: "அனைத்து வேலை வகைகள்",
    allExperiences: "அனைத்து அனுபவ நிலைகள்",
    thNo: "எண்.",
    thRole: "வேலைப்பணி",
    thEmployer: "நிறுவனம்",
    thProvince: "மாகாணம்",
    thRemote: "வேலை முறை",
    remoteText: "தொலைதூர வேலை",
    officeText: "அலுவலகம்",
    noRecords: "தேடல் முடிவுகள் இல்லை.",
    syncText: "தரவு ஒத்திசைக்கப்படுகிறது...",
    errorText: "இணைக்க முடியவில்லை.",
    panelEmployer: "நிறுவனம்",
    panelRemote: "வேலை முறை",
    panelLocation: "இடம்",
    panelExperience: "அனுபவம்",
    panelIndustry: "தொழிற்துறை",
    panelJobType: "வேலை வகை",
    panelDescription: "விவரம்",
    panelSkills: "திறன்கள்",
    panelPostedAt: "இடுகையிட்டது",
    closePanel: "மூடு",
    noSkills: "திறன்கள் இல்லை",
  },
};

export default function VacanciesPage() {
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");
  const d = useMemo(() => localization[currentLang], [currentLang]);

  const formattedDate = new Date().toLocaleDateString(
    currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
    { year: "numeric", month: "short", day: "numeric" },
  );

  const [search,       setSearch]       = useState("");
  const [industryId,   setIndustryId]   = useState<number | undefined>();
  const [provinceId,   setProvinceId]   = useState<number | undefined>();
  const [jobTypeId,    setJobTypeId]    = useState<number | undefined>();
  const [experienceId, setExperienceId] = useState<number | undefined>();
  const [selectedVacancy, setSelectedVacancy] = useState<Vacancy | null>(null);

  const { data: metadata,    isLoading: isMetaLoading, isError: isMetaError } = useVacancyMetadata();
  const { data: vacancyData, isLoading: isListLoading, isError: isListError  } = useVacancies({
    industry_id:   industryId,
    geo_data_id:   provinceId,
    job_type_id:   jobTypeId,
    experience_id: experienceId,
  });

  const filteredVacancies = useMemo(() => {
    const all = vacancyData?.jobs ?? [];
    if (!search.trim()) return all;
    const q = search.toLowerCase();
    return all.filter(
      (v) => v.job_role.toLowerCase().includes(q) || v.employer.name.toLowerCase().includes(q),
    );
  }, [vacancyData?.jobs, search]);

  const isLoading = isMetaLoading || isListLoading;
  const isError   = isMetaError   || isListError;

  if (isLoading) return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <div className="text-center space-y-3">
        <div className="w-8 h-8 border-2 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto" />
        <p className="text-sm text-gray-400 font-medium">{d.syncText}</p>
      </div>
    </div>
  );

  if (isError) return (
    <div className="min-h-screen bg-gray-50 flex items-center justify-center">
      <p className="text-sm text-red-500 font-medium">{d.errorText}</p>
    </div>
  );

  return (
    <div className="h-screen bg-gray-50 text-gray-800 font-sans flex flex-col overflow-hidden">

      {/* HEADER */}
      <header className="bg-white border-b border-gray-100 px-8 py-5 flex flex-col md:flex-row justify-between items-start md:items-center gap-4 sticky top-0 z-40 shrink-0">
        <div>
          <h1 className="text-xl font-black text-gray-900 tracking-tight">{d.title}</h1>
          <p className="text-xs text-gray-400 mt-1 font-medium">
            <span className="font-bold text-gray-600">{filteredVacancies.length}</span> {d.subtitle}
          </p>
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
              <button key={l} onClick={() => setCurrentLang(l)}
                className={`px-3 py-1 text-xs font-bold rounded-lg transition-all ${currentLang === l ? "bg-white text-blue-600 shadow-sm" : "text-gray-500 hover:text-gray-900"}`}>
                {l === "en" ? "English" : l === "si" ? "සිංහල" : "தமிழ்"}
              </button>
            ))}
          </div>
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white text-xs font-black shadow-md border border-white">BJ</div>
        </div>
      </header>

      {/* ── ROOT SPLIT LAYOUT ──────────────────────────────────────────────── */}
      <div className="flex flex-1 min-h-0">

        {/* ── LEFT: MAIN SCROLLABLE CONTENT ─────────────────────────────── */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-8 space-y-6">

            {/* FILTER BAR */}
            <div className="bg-white p-4 rounded-xl shadow-sm flex flex-wrap gap-3">
              <input
                type="text"
                placeholder={d.searchPlaceholder}
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="flex-1 min-w-[200px] px-3 py-2 border border-gray-200 rounded-lg text-sm outline-none focus:ring-2 focus:ring-blue-100"
              />
              <select
                value={industryId ?? ""}
                onChange={(e) => setIndustryId(e.target.value ? Number(e.target.value) : undefined)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm max-w-[180px]"
              >
                <option value="">{d.allIndustries}</option>
                {metadata?.industries.map((ind) => (
                  <option key={ind.id} value={ind.id}>{ind.name}</option>
                ))}
              </select>
              <select
                value={provinceId ?? ""}
                onChange={(e) => setProvinceId(e.target.value ? Number(e.target.value) : undefined)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm"
              >
                <option value="">{d.allProvinces}</option>
                {metadata?.provinces.map((p) => (
                  <option key={p.id} value={p.id}>{p.province}</option>
                ))}
              </select>
              <select
                value={jobTypeId ?? ""}
                onChange={(e) => setJobTypeId(e.target.value ? Number(e.target.value) : undefined)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm"
              >
                <option value="">{d.allJobTypes}</option>
                {metadata?.job_types.map((jt) => (
                  <option key={jt.id} value={jt.id}>{jt.type}</option>
                ))}
              </select>
              <select
                value={experienceId ?? ""}
                onChange={(e) => setExperienceId(e.target.value ? Number(e.target.value) : undefined)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm"
              >
                <option value="">{d.allExperiences}</option>
                {metadata?.experiences.map((exp) => (
                  <option key={exp.id} value={exp.id}>{exp.name}</option>
                ))}
              </select>
            </div>

            {/* VACANCIES TABLE */}
            <div className="bg-white rounded-xl shadow-sm overflow-hidden">
              <table className="w-full text-left">
                <thead>
                  <tr className="border-b bg-gray-50/70">
                    <th className="py-3 px-4 text-xs font-semibold text-gray-500">{d.thNo}</th>
                    <th className="py-3 px-4 text-xs font-semibold text-gray-500">{d.thRole}</th>
                    <th className="py-3 px-4 text-xs font-semibold text-gray-500">{d.thEmployer}</th>
                    <th className="py-3 px-4 text-xs font-semibold text-gray-500">{d.thProvince}</th>
                    <th className="py-3 px-4 text-xs font-semibold text-gray-500">{d.thRemote}</th>
                  </tr>
                </thead>
                <tbody>
                  {filteredVacancies.length === 0 ? (
                    <tr>
                      <td colSpan={5} className="py-16 text-center text-sm text-gray-400">{d.noRecords}</td>
                    </tr>
                  ) : (
                    filteredVacancies.map((v, index) => (
                      <tr
                        key={v.id}
                        onClick={() => setSelectedVacancy(v)}
                        className={`border-b cursor-pointer transition-colors ${
                          selectedVacancy?.id === v.id
                            ? "bg-blue-50 border-l-2 border-l-blue-500"
                            : "hover:bg-zinc-50"
                        }`}
                      >
                        <td className="py-3.5 px-4 text-xs text-gray-400 font-mono">{index + 1}</td>
                        <td className="py-3.5 px-4 text-sm font-semibold text-zinc-900">{v.job_role}</td>
                        <td className="py-3.5 px-4 text-sm text-gray-700">{v.employer.name}</td>
                        <td className="py-3.5 px-4 text-sm text-gray-500">{v.meta_data.geo_data.province}</td>
                        <td className="py-3.5 px-4 text-sm">
                          {v.is_remote
                            ? <span className="text-blue-600 font-bold text-xs bg-blue-50 px-2 py-0.5 rounded-full">{d.remoteText}</span>
                            : <span className="text-gray-400 text-xs">{d.officeText}</span>
                          }
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

          </div>
        </div>

        {/* ── RIGHT: DETAIL PANEL (inline, not overlay) ─────────────────── */}
        {selectedVacancy && (
          <div className="w-[420px] min-w-[420px] border-l border-gray-200 bg-white flex flex-col overflow-hidden shadow-lg">

            {/* Panel Header */}
            <div className="px-6 py-4 border-b border-gray-100 flex justify-between items-start shrink-0">
              <div className="flex-1 pr-3">
                <h2 className="text-sm font-black text-gray-900 leading-snug">{selectedVacancy.job_role}</h2>
                <p className="text-[11px] text-gray-400 mt-0.5">
                  {d.panelPostedAt}: {new Date(selectedVacancy.meta_data.posted_at).toLocaleDateString()}
                </p>
              </div>
              <button
                onClick={() => setSelectedVacancy(null)}
                className="shrink-0 text-xs font-bold text-gray-400 hover:text-gray-700 bg-gray-100 hover:bg-gray-200 px-3 py-1.5 rounded-lg transition-all"
              >
                {d.closePanel} ✕
              </button>
            </div>

            {/* Panel Body */}
            <div className="flex-1 overflow-y-auto p-5 space-y-5">

              {/* Meta grid */}
              <div className="grid grid-cols-2 gap-3">
                <DetailItem label={d.panelEmployer}   value={selectedVacancy.employer.name} />
                <DetailItem label={d.panelRemote}     value={selectedVacancy.is_remote ? d.remoteText : d.officeText} />
                <DetailItem label={d.panelLocation}   value={selectedVacancy.location} />
                <DetailItem label={d.panelJobType}    value={selectedVacancy.job_type.type} />
                <DetailItem label={d.panelIndustry}   value={selectedVacancy.meta_data.industry.name} />
                <DetailItem label={d.panelExperience} value={selectedVacancy.meta_data.experience.name} />
              </div>

              {/* Divider */}
              <div className="border-t border-gray-100" />

              {/* Description */}
              <div>
                <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-2">{d.panelDescription}</p>
                <p className="text-xs text-gray-700 bg-gray-50 p-4 rounded-xl leading-relaxed">
                  {selectedVacancy.job_description || "N/A"}
                </p>
              </div>

              {/* Skills */}
              <div>
                <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.panelSkills}</p>
                <div className="flex flex-wrap gap-2">
                  {selectedVacancy.skills.length > 0
                    ? selectedVacancy.skills.map((s) => (
                        <span key={s.id} className="px-3 py-1 bg-zinc-100 text-[11px] font-bold rounded-full text-zinc-700">
                          {s.skill}
                        </span>
                      ))
                    : <span className="text-xs text-gray-400">{d.noSkills}</span>
                  }
                </div>
              </div>

            </div>
          </div>
        )}

      </div>
    </div>
  );
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div className="bg-gray-50 px-3 py-2.5 rounded-lg">
      <p className="text-[9px] text-gray-400 uppercase font-black tracking-wider">{label}</p>
      <p className="text-xs font-bold text-zinc-800 mt-0.5 truncate">{value || "N/A"}</p>
    </div>
  );
}