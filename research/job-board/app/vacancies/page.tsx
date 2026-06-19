// "use client";

// import { useState } from "react";
// import Header from "@/components/layout/Header";
// import { useVacancyMetadata, useVacanciesList } from "@/hooks/use-vacancies";
// import { JobVacancy } from "@/types/job";

// export default function VacanciesPage() {
//   // Query UI State Controls
//   const [search, setSearch] = useState("");
//   const [selectedSector, setSelectedSector] = useState("All Categories");
//   const [selectedProvince, setSelectedProvince] = useState<number>(0);
//   const [selectedContract, setSelectedContract] = useState("All Contracts");
//   const [selectedSeniority, setSelectedSeniority] = useState("All Seniorities");
//   const [selectedVacancy, setSelectedVacancy] = useState<JobVacancy | null>(null);

//   // Initialize TanStack BFF Data Channels
//   const { data: metadata, isLoading: isMetaLoading, isError: isMetaError } = useVacancyMetadata();

//   const { data: activeVacancies = [], isLoading: isListLoading } = useVacanciesList({
//     category: selectedSector !== "All Categories" ? selectedSector : "All",
//     province: selectedProvince !== 0 ? selectedProvince : undefined,
//     contractType: selectedContract !== "All Contracts" ? selectedContract : "All",
//     seniority: selectedSeniority !== "All Seniorities" ? selectedSeniority : "All",
//   });

//   // Client side text indexing over filtered network response items
//   const finalRenderedVacancies = activeVacancies.filter((v) => {
//     if (!search) return true;
//     const matchTarget = search.toLowerCase();
//     return (
//       v.job_role.toLowerCase().includes(matchTarget) ||
//       v.employer.toLowerCase().includes(matchTarget)
//     );
//   });

//   if (isMetaLoading || isListLoading) {
//     return (
//       <div className="min-h-screen flex items-center justify-center bg-gray-50">
//         <div className="flex flex-col items-center gap-3">
//           <div className="w-8 h-8 border-4 border-zinc-900 border-t-transparent rounded-full animate-spin" />
//           <p className="text-xs font-medium text-gray-500 tracking-wide">Syncing Workspace Data via BFF...</p>
//         </div>
//       </div>
//     );
//   }

//   if (isMetaError) {
//     return (
//       <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
//         <div className="bg-white p-6 rounded-xl border border-red-100 shadow-sm max-w-sm text-center space-y-3">
//           <div className="text-red-500 font-bold text-sm">Failed to connect to backend context</div>
//           <p className="text-xs text-gray-500">The gateway was unable to fetch validation criteria options.</p>
//         </div>
//       </div>
//     );
//   }

//   // Formatting response items safely for fallback states mapping metadata lookup objects
//   const sectorsList = ["All Categories", ...(metadata?.categories?.map(c => c.value) || [])];
//   const provincesList = [
//     { id: 0, name: "All Provinces" },
//     ...(metadata?.provinces || [])
//   ];
//   const contractsList = ["All Contracts", ...(metadata?.contractTypes?.map(ct => ct.value) || [])];
//   const senioritiesList = ["All Seniorities", ...(metadata?.seniorityLevels?.map(sl => sl.value) || [])];

//   return (
//     <div className="min-h-screen bg-gray-50">
//       <Header
//         title="Vacancy Explorer"
//         subtitle={`${finalRenderedVacancies.length} vacancies synchronized from BFF pipeline`}
//       />

//       <div className="p-8">
//         {/* Filters Panel Component Grid */}
//         <div className="bg-white p-5 rounded-xl border border-gray-200/60 shadow-sm mb-6">
//           <div className="flex flex-wrap gap-4 items-end w-full">
//             {/* Search Input Box */}
//             <div className="flex-1 min-w-[200px]">
//               <label className="block text-xs text-gray-500 mb-1 font-medium">Search</label>
//               <input
//                 type="text"
//                 placeholder="Search role or employer..."
//                 value={search}
//                 onChange={(e) => setSearch(e.target.value)}
//                 className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-zinc-900 bg-white"
//               />
//             </div>

//             {/* Category Select Controller */}
//             <div>
//               <label className="block text-xs text-gray-500 mb-1 font-medium">Category</label>
//               <select
//                 value={selectedSector}
//                 onChange={(e) => setSelectedSector(e.target.value)}
//                 className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
//               >
//                 {sectorsList.map((s) => (
//                   <option key={s} value={s}>{s}</option>
//                 ))}
//               </select>
//             </div>

//             {/* Province Select Controller */}
//             <div>
//               <label className="block text-xs text-gray-500 mb-1 font-medium">Province</label>
//               <select
//                 value={selectedProvince}
//                 onChange={(e) => setSelectedProvince(Number(e.target.value))}
//                 className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
//               >
//                 {provincesList.map((p) => (
//                   <option key={p.id} value={p.id}>{p.name}</option>
//                 ))}
//               </select>
//             </div>

//             {/* Contract Type Controller */}
//             <div>
//               <label className="block text-xs text-gray-500 mb-1 font-medium">Contract Type</label>
//               <select
//                 value={selectedContract}
//                 onChange={(e) => setSelectedContract(e.target.value)}
//                 className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
//               >
//                 {contractsList.map((ct) => (
//                   <option key={ct} value={ct}>{ct}</option>
//                 ))}
//               </select>
//             </div>

//             {/* Seniority Level Selector */}
//             <div>
//               <label className="block text-xs text-gray-500 mb-1 font-medium">Seniority Level</label>
//               <select
//                 value={selectedSeniority}
//                 onChange={(e) => setSelectedSeniority(e.target.value)}
//                 className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
//               >
//                 {senioritiesList.map((sl) => (
//                   <option key={sl} value={sl}>{sl}</option>
//                 ))}
//               </select>
//             </div>
//           </div>
//         </div>

//         {/* Dynamic Ledger Data Grid Output Table View */}
//         <div className="bg-white rounded-xl border border-gray-200/60 shadow-sm overflow-hidden w-full">
//           <div className="overflow-x-auto">
//             <table className="w-full text-left border-collapse">
//               <thead>
//                 <tr className="border-b border-gray-100 bg-gray-50/70">
//                   <th className="py-3 px-4 text-xs font-semibold text-gray-500 w-16">No.</th>
//                   <th className="py-3 px-4 text-xs font-semibold text-gray-500">Job Role</th>
//                   <th className="py-3 px-4 text-xs font-semibold text-gray-500">Employer</th>
//                   <th className="py-3 px-4 text-xs font-semibold text-gray-500">Province</th>
//                   <th className="py-3 px-4 text-xs font-semibold text-gray-500">Confidence</th>
//                   <th className="py-3 px-4 text-xs font-semibold text-gray-500">Is Remote</th>
//                 </tr>
//               </thead>
//               <tbody>
//                 {finalRenderedVacancies.length > 0 ? (
//                   finalRenderedVacancies.map((v, index) => (
//                     <tr
//                       key={v.id}
//                       onClick={() => setSelectedVacancy(v)}
//                       className="border-b border-gray-50 cursor-pointer transition-colors hover:bg-zinc-50/60"
//                     >
//                       <td className="py-3.5 px-4 text-xs text-gray-400 font-mono font-medium">{index + 1}</td>
//                       <td className="py-3.5 px-4 text-sm font-semibold text-zinc-900">{v.job_role}</td>
//                       <td className="py-3.5 px-4 text-sm text-gray-700">{v.employer}</td>
//                       <td className="py-3.5 px-4 text-sm text-gray-500">{v.meta_data.geo.province}</td>
//                       <td className="py-3.5 px-4 text-sm">
//                         <span className={`px-2 py-0.5 text-[11px] font-bold rounded-md ${v.meta_data.confidence_score >= 0.95 ? "bg-emerald-50 text-emerald-700 border border-emerald-200" : "bg-amber-50 text-amber-700 border border-amber-200"}`}>
//                           {Math.round(v.meta_data.confidence_score * 100)}%
//                         </span>
//                       </td>
//                       <td className="py-3.5 px-4 text-sm">
//                         {v.is_remote ? (
//                           <span className="bg-blue-50 text-blue-700 border border-blue-100 px-2 py-0.5 text-[11px] font-semibold rounded-md">Remote Available</span>
//                         ) : (
//                           <span className="text-xs text-gray-400 font-medium pl-2">Office Based</span>
//                         )}
//                       </td>
//                     </tr>
//                   ))
//                 ) : (
//                   <tr>
//                     <td colSpan={6} className="text-center py-12 text-sm text-gray-400">
//                       No active listings found matching active query params.
//                     </td>
//                   </tr>
//                 )}
//               </tbody>
//             </table>
//           </div>
//         </div>
//       </div>

//       {/* Detail Overlay Card Layout Panel block */}
//       {selectedVacancy && (
//         <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-zinc-950/40 backdrop-blur-sm">
//           <div className="bg-white w-full max-w-2xl rounded-xl border border-zinc-200/80 shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
//             <div className="px-6 py-4 bg-zinc-50 border-b border-zinc-100 flex items-start justify-between">
//               <div>
//                 <h3 className="text-lg font-bold text-zinc-900">{selectedVacancy.job_role}</h3>
//                 <p className="text-sm text-zinc-600 font-medium">{selectedVacancy.employer}</p>
//               </div>
//               <button onClick={() => setSelectedVacancy(null)} className="text-zinc-400 hover:text-zinc-700 p-1 rounded-md">
//                 ✕
//               </button>
//             </div>
//             <div className="p-6 space-y-5 overflow-y-auto text-sm">
//               <div className="grid grid-cols-2 gap-4 bg-zinc-50/60 p-4 rounded-lg border border-zinc-100">
//                 <DetailItem label="Location" value={selectedVacancy.location} />
//                 <DetailItem label="Province" value={selectedVacancy.meta_data.geo.province} />
//                 <DetailItem label="Category" value={selectedVacancy.meta_data.standardized_category} />
//                 <DetailItem label="Seniority Level" value={selectedVacancy.meta_data.seniority} />
//               </div>
//               <div>
//                 <h4 className="text-xs font-bold uppercase text-zinc-400 tracking-wider mb-1">Responsibilities</h4>
//                 <p className="text-zinc-700 p-3 bg-white border rounded-lg shadow-xs">{selectedVacancy.key_responsibilities}</p>
//               </div>
//               <div>
//                 <h4 className="text-xs font-bold uppercase text-zinc-400 tracking-wider mb-1">Qualifications</h4>
//                 <p className="text-zinc-700 p-3 bg-white border rounded-lg shadow-xs">{selectedVacancy.qualifications}</p>
//               </div>
//               <div>
//                 <h4 className="text-xs font-bold uppercase text-zinc-400 tracking-wider mb-2">Identified Skills</h4>
//                 <div className="flex flex-wrap gap-1.5">
//                   {selectedVacancy.skills.map((skill) => (
//                     <span key={skill} className="px-2.5 py-1 bg-zinc-100 text-zinc-800 text-xs font-medium rounded-md border">{skill}</span>
//                   ))}
//                 </div>
//               </div>
//             </div>
//           </div>
//         </div>
//       )}
//     </div>
//   );
// }

// function DetailItem({ label, value }: { label: string; value: string }) {
//   return (
//     <div>
//       <p className="text-[10px] text-zinc-400 uppercase font-bold tracking-wider">{label}</p>
//       <p className="text-sm font-semibold text-zinc-800 mt-0.5">{value}</p>
//     </div>
//   );
// }











"use client";

import { useState, useMemo } from "react";
import { useVacancyMetadata, useVacanciesList } from "@/hooks/use-vacancies";
import { JobVacancy } from "@/types/job";

const localization = {
  en: {
    title: "Vacancy Explorer",
    subtitle: (count: number) => `${count} vacancies synchronized from BFF pipeline`,
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
      <header className="bg-white border-b px-8 py-5 flex justify-between items-center sticky top-0 z-40">
        <h1 className="text-xl font-black text-gray-900">{d.title}</h1>
        <div className="flex bg-gray-100 p-1 rounded-xl">
          {(["en", "si", "ta"] as const).map((l) => (
            <button key={l} onClick={() => setCurrentLang(l)} className={`px-3 py-1 text-xs font-bold rounded-lg ${currentLang === l ? "bg-white text-blue-600 shadow-sm" : "text-gray-500"}`}>{l === "en" ? "English" : l === "si" ? "සිංහල" : "தமிழ்"}</button>
          ))}
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