"use client";

import { useState, useMemo } from "react";
import { useVacancyMetadata, useVacanciesList } from "@/hooks/use-vacancies";
import { JobVacancy } from "@/types/job";

const localization = {
  en: {
    title: "Vacancy Explorer",
    subtitle: "vacancies synchronized from BFF pipeline",
    searchLabel: "Search",
    searchPlaceholder: "Search role or employer...",
    industryLabel: "Industry",
    provinceLabel: "Province",
    contractLabel: "Contract Type",
    seniorityLabel: "Seniority Level",
    thNo: "No.",
    thRole: "Job Role",
    thEmployer: "Employer",
    thProvince: "Province",
    thRemote: "Is Remote",
    remoteText: "Remote Available",
    officeText: "Office Based",
    noRecords: "No active listings found.",
    syncText: "Syncing Workspace Data via BFF...",
    errorTitle: "Failed to connect to backend",
    allIndustries: "All Industries",
    allProvinces: "All Provinces",
    allContracts: "All Contracts",
    allSeniorities: "All Seniorities",
    // Panel Labels
    panelRole: "Job Role",
    panelEmployer: "Employer",
    panelDescription: "Description",
    panelRemote: "Work Mode",
    panelLocation: "Location",
    panelType: "Job Type",
    panelSeniority: "Seniority",
    panelSkills: "Skills"
  },
  si: {
    title: "රැකියා ගවේෂකය",
    subtitle: (count: number) => `රැකියා ඇබෑර්තු ${count} ක් ඇත`,
    searchLabel: "සෙවීම",
    searchPlaceholder: "තනතුර හෝ ආයතනය සොයන්න...",
    industryLabel: "කර්මාන්තය",
    provinceLabel: "පළාත",
    contractLabel: "ගිවිසුම් වර්ගය",
    seniorityLabel: "ජ්‍යෙෂ්ඨත්ව මට්ටම",
    thNo: "අංකය",
    thRole: "තනතුර",
    thEmployer: "ආයතනය",
    thProvince: "පළාත",
    thRemote: "සේවා ක්‍රමය",
    remoteText: "දුරස්ථ සේවය",
    officeText: "කාර්යාලීය",
    noRecords: "ප්‍රතිඵල නොමැත.",
    syncText: "දත්ත සම්බන්ධ වෙමින් පවතී...",
    errorTitle: "සම්බන්ධ විය නොහැක",
    allIndustries: "සියලුම කර්මාන්ත",
    allProvinces: "සියලුම පළාත්",
    allContracts: "සියලුම ගිවිසුම්",
    allSeniorities: "සියලුම ජ්‍යෙෂ්ඨත්වයන්",
    panelRole: "රැකියා තනතුර",
    panelEmployer: "ආයතනය",
    panelDescription: "විස්තරය",
    panelRemote: "සේවා ක්‍රමය",
    panelLocation: "ස්ථානය",
    panelType: "වර්ගය",
    panelSeniority: "ජ්‍යෙෂ්ඨත්වය",
    panelSkills: "නිපුණතා"
  },
  ta: {
    title: "வேலைவாய்ப்பு ஆய்வாளர்",
    subtitle: (count: number) => `${count} வேலைவாய்ப்புகள்`,
    searchLabel: "தேடு",
    searchPlaceholder: "பதவி அல்லது நிறுவனத்தைத் தேடு...",
    industryLabel: "தொழிற்துறை",
    provinceLabel: "மாகாணம்",
    contractLabel: "ஒப்பந்த வகை",
    seniorityLabel: "முதுநிலை நிலை",
    thNo: "எண்.",
    thRole: "வேலைப்பணி",
    thEmployer: "நிறுவனம்",
    thProvince: "மாகாணம்",
    thRemote: "வேலை முறை",
    remoteText: "தொலைதூர வேலை",
    officeText: "அலுவலகம்",
    noRecords: "தேடல் முடிவுகள் இல்லை.",
    syncText: "தரவு ஒத்திசைக்கப்படுகிறது...",
    errorTitle: "இணைக்க முடியவில்லை",
    allIndustries: "அனைத்துத் தொழில்துறைகள்",
    allProvinces: "அனைத்து மாகாணங்கள்",
    allContracts: "அனைத்து ஒப்பந்தங்கள்",
    allSeniorities: "அனைத்து நிலைகள்",
    panelRole: "வேலைப்பணி",
    panelEmployer: "நிறுவனம்",
    panelDescription: "விவரம்",
    panelRemote: "வேலை முறை",
    panelLocation: "இடம்",
    panelType: "வகை",
    panelSeniority: "முதுநிலை",
    panelSkills: "திறன்கள்"
  }
};

export default function VacanciesPage() {
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");
  const d = useMemo(() => localization[currentLang], [currentLang]);

  const [search, setSearch] = useState("");
  const [selectedSector, setSelectedSector] = useState("All Industries");
  const [selectedProvince, setSelectedProvince] = useState<number>(0);
  const [selectedContract, setSelectedContract] = useState("All Contracts");
  const [selectedSeniority, setSelectedSeniority] = useState("All Seniorities");
  const [selectedVacancy, setSelectedVacancy] = useState<JobVacancy | null>(null);

  const formattedDate = new Date().toLocaleDateString(
        currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
        { year: "numeric", month: "short", day: "numeric" }
    );

  const { data: metadata, isLoading: isMetaLoading, isError: isMetaError } = useVacancyMetadata();
  const { data: activeVacancies = [], isLoading: isListLoading } = useVacanciesList({
    category: selectedSector !== "All Industries" ? selectedSector : "All",
    province: selectedProvince !== 0 ? selectedProvince : undefined,
    contractType: selectedContract !== "All Contracts" ? selectedContract : "All",
    seniority: selectedSeniority !== "All Seniorities" ? selectedSeniority : "All",
  });

  const finalRenderedVacancies = activeVacancies.filter((v) => {
    if (!search) return true;
    const matchTarget = search.toLowerCase();
    return v.job_role.toLowerCase().includes(matchTarget) || v.employer.toLowerCase().includes(matchTarget);
  });

  if (isMetaLoading || isListLoading) return <div className="p-10 text-center">{d.syncText}</div>;

  const sectorsList = [{ value: "All Industries", label: d.allIndustries }, ...(metadata?.categories?.map(c => ({ value: c.value, label: c.value })) || [])];
  const provincesList = [{ id: 0, name: d.allProvinces }, ...(metadata?.provinces || [])];
  const contractsList = [{ value: "All Contracts", label: d.allContracts }, ...(metadata?.contractTypes?.map(ct => ({ value: ct.value, label: ct.value })) || [])];
  const senioritiesList = [{ value: "All Seniorities", label: d.allSeniorities }, ...(metadata?.seniorityLevels?.map(sl => ({ value: sl.value, label: sl.value })) || [])];

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-100 px-8 py-5 flex flex-col md:flex-row justify-between items-start md:items-center gap-4 sticky top-0 z-40">
        <div>
          <h1 className="text-xl font-black text-gray-900 tracking-tight">{d.title}</h1>
          <p className="text-xs text-gray-400 mt-1 font-medium">
            {d.subtitle}
          </p>
        </div>
        <div className="flex flex-wrap items-center gap-4 ml-auto md:ml-0 w-full md:w-auto justify-end">
          <div className="flex items-center gap-2 bg-gray-50 px-3 py-1.5 rounded-xl text-xs text-gray-500 font-medium shadow-inner">
            <svg className="w-3.5 h-3.5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth="2">
              <path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 002-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
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

      <div className="p-8">
        <div className="bg-white p-5 rounded-xl shadow-sm mb-6 flex flex-wrap gap-4">
          <input type="text" placeholder={d.searchPlaceholder} value={search} onChange={(e) => setSearch(e.target.value)} className="flex-1 px-3 py-2 border rounded-lg text-sm outline-none" />
          <select value={selectedSector} onChange={(e) => setSelectedSector(e.target.value)} className="px-3 py-2 border rounded-lg text-sm">{sectorsList.map(s => <option key={s.value} value={s.value}>{s.label}</option>)}</select>
          <select value={selectedProvince} onChange={(e) => setSelectedProvince(Number(e.target.value))} className="px-3 py-2 border rounded-lg text-sm">{provincesList.map(p => <option key={p.id} value={p.id}>{p.name}</option>)}</select>
          <select value={selectedContract} onChange={(e) => setSelectedContract(e.target.value)} className="px-3 py-2 border rounded-lg text-sm">{contractsList.map(ct => <option key={ct.value} value={ct.value}>{ct.label}</option>)}</select>
          <select value={selectedSeniority} onChange={(e) => setSelectedSeniority(e.target.value)} className="px-3 py-2 border rounded-lg text-sm">{senioritiesList.map(sl => <option key={sl.value} value={sl.value}>{sl.label}</option>)}</select>
        </div>

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
              {finalRenderedVacancies.map((v, index) => (
                <tr key={v.id} onClick={() => setSelectedVacancy(v)} className="border-b hover:bg-zinc-50 cursor-pointer">
                  <td className="py-3.5 px-4 text-xs text-gray-400 font-mono">{index + 1}</td>
                  <td className="py-3.5 px-4 text-sm font-semibold text-zinc-900">{v.job_role}</td>
                  <td className="py-3.5 px-4 text-sm text-gray-700">{v.employer}</td>
                  <td className="py-3.5 px-4 text-sm text-gray-500">{v.meta_data.geo.province}</td>
                  <td className="py-3.5 px-4 text-sm">{v.is_remote ? <span className="text-blue-600 font-bold">{d.remoteText}</span> : <span>{d.officeText}</span>}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      {/* Detail Panel */}
      <div className={`fixed inset-0 z-50 transition-opacity ${selectedVacancy ? "pointer-events-auto opacity-100" : "pointer-events-none opacity-0"}`}>
        <div onClick={() => setSelectedVacancy(null)} className="absolute inset-0 bg-black/20 backdrop-blur-sm" />
        <div className="absolute inset-y-0 right-0 w-full max-w-lg bg-white shadow-2xl p-8 overflow-y-auto">
          <button onClick={() => setSelectedVacancy(null)} className="text-gray-400 font-bold mb-4">✕ Close</button>
          <div className="space-y-6">
            <h2 className="text-2xl font-black">{selectedVacancy?.job_role}</h2>
            <div className="grid grid-cols-2 gap-4 text-sm">
                <DetailItem label={d.panelEmployer} value={selectedVacancy?.employer || ""} />
                <DetailItem label={d.panelRemote} value={selectedVacancy?.is_remote ? d.remoteText : d.officeText} />
                <DetailItem label={d.panelLocation} value={selectedVacancy?.location || ""} />
                <DetailItem label={d.panelSeniority} value={selectedVacancy?.meta_data.seniority || ""} />
            </div>
            <div>
                <h4 className="text-[10px] font-black uppercase text-gray-400 mb-2">{d.panelDescription}</h4>
                <p className="text-sm text-gray-700 bg-gray-50 p-4 rounded-xl">{selectedVacancy?.key_responsibilities} {selectedVacancy?.qualifications}</p>
            </div>
            <div>
                <h4 className="text-[10px] font-black uppercase text-gray-400 mb-3">{d.panelSkills}</h4>
                <div className="flex flex-wrap gap-2">
                    {selectedVacancy?.skills.map(s => <span key={s} className="px-3 py-1 bg-zinc-100 text-[11px] font-bold rounded-full">{s}</span>)}
                </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return <div><p className="text-[9px] text-gray-400 uppercase font-black">{label}</p><p className="text-xs font-bold text-zinc-800 mt-0.5">{value || "N/A"}</p></div>;
}