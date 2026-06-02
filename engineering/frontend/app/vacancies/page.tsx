"use client";

import { useState } from "react";
import Header from "@/components/layout/Header";
import { useVacancyMetadata, useVacanciesList } from "@/hooks/use-vacancies";
import { JobVacancy } from "@/types/job";

export default function VacanciesPage() {
  // Query UI State Controls
  const [search, setSearch] = useState("");
  const [selectedSector, setSelectedSector] = useState("All Categories");
  const [selectedProvince, setSelectedProvince] = useState("All Provinces")
//   const [selectedProvince, setSelectedProvince] = useState<number>(0);
  const [selectedContract, setSelectedContract] = useState("All Contracts");
  const [selectedSeniority, setSelectedSeniority] = useState("All Seniorities");
  const [selectedVacancy, setSelectedVacancy] = useState<JobVacancy | null>(null);

  // Initialize TanStack BFF Data Channels
  const { data: metadata, isLoading: isMetaLoading, isError: isMetaError } = useVacancyMetadata();

  const { data: activeVacancies = [], isLoading: isListLoading } = useVacanciesList({
    category: selectedSector !== "All Categories" ? selectedSector : "All",
    province: selectedProvince != "All Provinces" ? selectedProvince : "All",
    //province: selectedProvince !== 0 ? selectedProvince : undefined,
    contractType: selectedContract !== "All Contracts" ? selectedContract : "All",
    seniority: selectedSeniority !== "All Seniorities" ? selectedSeniority : "All",
  });

  // Client side text indexing over filtered network response items
  const finalRenderedVacancies = activeVacancies.filter((v) => {
    if (!search) return true;
    const matchTarget = search.toLowerCase();
    return (
      v.job_role.toLowerCase().includes(matchTarget) ||
      v.employer.toLowerCase().includes(matchTarget)
    );
  });

  if (isMetaLoading || isListLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="flex flex-col items-center gap-3">
          <div className="w-8 h-8 border-4 border-zinc-900 border-t-transparent rounded-full animate-spin" />
          <p className="text-xs font-medium text-gray-500 tracking-wide">Syncing Workspace Data via BFF...</p>
        </div>
      </div>
    );
  }

  if (isMetaError) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
        <div className="bg-white p-6 rounded-xl border border-red-100 shadow-sm max-w-sm text-center space-y-3">
          <div className="text-red-500 font-bold text-sm">Failed to connect to backend context</div>
          <p className="text-xs text-gray-500">The gateway was unable to fetch validation criteria options.</p>
        </div>
      </div>
    );
  }

  // Formatting response items safely for fallback states mapping metadata lookup objects
  const sectorsList = ["All Categories", ...(metadata?.categories?.map(c => c.value) || [])];
//   const provincesList = [
//     { id: 0, name: "All Provinces" },
//     ...(metadata?.provinces || [])
//   ];
  const provincesList = ["All Provinces", ...(metadata?.provinces.map(p => p.name) || [])]
  const contractsList = ["All Contracts", ...(metadata?.contractTypes?.map(ct => ct.value) || [])];
  const senioritiesList = ["All Seniorities", ...(metadata?.seniorityLevels?.map(sl => sl.value) || [])];

  return (
    <div className="min-h-screen bg-gray-50">
      <Header
        title="Vacancy Explorer"
        subtitle={`${finalRenderedVacancies.length} vacancies synchronized from BFF pipeline`}
      />

      <div className="p-8">
        {/* Filters Panel Component Grid */}
        <div className="bg-white p-5 rounded-xl border border-gray-200/60 shadow-sm mb-6">
          <div className="flex flex-wrap gap-4 items-end w-full">
            {/* Search Input Box */}
            <div className="flex-1 min-w-[200px]">
              <label className="block text-xs text-gray-500 mb-1 font-medium">Search</label>
              <input
                type="text"
                placeholder="Search role or employer..."
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-zinc-900 bg-white"
              />
            </div>

            {/* Category Select Controller */}
            <div>
              <label className="block text-xs text-gray-500 mb-1 font-medium">Category</label>
              <select
                value={selectedSector}
                onChange={(e) => setSelectedSector(e.target.value)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
              >
                {sectorsList.map((s) => (
                  <option key={s} value={s}>{s}</option>
                ))}
              </select>
            </div>

            {/* Province Select Controller */}
            <div>
              <label className="block text-xs text-gray-500 mb-1 font-medium">Province</label>
              <select
                value={selectedProvince}
                onChange={(e) => setSelectedProvince(e.target.value)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
              >
                {provincesList.map((p) => (
                //   <option key={p.id} value={p.id}>{p.name}</option>
                    <option key={p} value={p}>{p}</option>
                ))}
              </select>
            </div>

            {/* Contract Type Controller */}
            <div>
              <label className="block text-xs text-gray-500 mb-1 font-medium">Contract Type</label>
              <select
                value={selectedContract}
                onChange={(e) => setSelectedContract(e.target.value)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
              >
                {contractsList.map((ct) => (
                  <option key={ct} value={ct}>{ct}</option>
                ))}
              </select>
            </div>

            {/* Seniority Level Selector */}
            <div>
              <label className="block text-xs text-gray-500 mb-1 font-medium">Seniority Level</label>
              <select
                value={selectedSeniority}
                onChange={(e) => setSelectedSeniority(e.target.value)}
                className="px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white outline-none focus:ring-2 focus:ring-zinc-900"
              >
                {senioritiesList.map((sl) => (
                  <option key={sl} value={sl}>{sl}</option>
                ))}
              </select>
            </div>
          </div>
        </div>

        {/* Dynamic Ledger Data Grid Output Table View */}
        <div className="bg-white rounded-xl border border-gray-200/60 shadow-sm overflow-hidden w-full">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-gray-100 bg-gray-50/70">
                  <th className="py-3 px-4 text-xs font-semibold text-gray-500 w-16">No.</th>
                  <th className="py-3 px-4 text-xs font-semibold text-gray-500">Job Role</th>
                  <th className="py-3 px-4 text-xs font-semibold text-gray-500">Employer</th>
                  <th className="py-3 px-4 text-xs font-semibold text-gray-500">Province</th>
                  <th className="py-3 px-4 text-xs font-semibold text-gray-500">Confidence</th>
                  <th className="py-3 px-4 text-xs font-semibold text-gray-500">Is Remote</th>
                </tr>
              </thead>
              <tbody>
                {finalRenderedVacancies.length > 0 ? (
                  finalRenderedVacancies.map((v, index) => (
                    <tr
                      key={v.id}
                      onClick={() => setSelectedVacancy(v)}
                      className="border-b border-gray-50 cursor-pointer transition-colors hover:bg-zinc-50/60"
                    >
                      <td className="py-3.5 px-4 text-xs text-gray-400 font-mono font-medium">{index + 1}</td>
                      <td className="py-3.5 px-4 text-sm font-semibold text-zinc-900">{v.job_role}</td>
                      <td className="py-3.5 px-4 text-sm text-gray-700">{v.employer}</td>
                      <td className="py-3.5 px-4 text-sm text-gray-500">{v.meta_data.geo.province}</td>
                      <td className="py-3.5 px-4 text-sm">
                        <span className={`px-2 py-0.5 text-[11px] font-bold rounded-md ${v.meta_data.confidence_score >= 0.95 ? "bg-emerald-50 text-emerald-700 border border-emerald-200" : "bg-amber-50 text-amber-700 border border-amber-200"}`}>
                          {Math.round(v.meta_data.confidence_score * 100)}%
                        </span>
                      </td>
                      <td className="py-3.5 px-4 text-sm">
                        {v.is_remote ? (
                          <span className="bg-blue-50 text-blue-700 border border-blue-100 px-2 py-0.5 text-[11px] font-semibold rounded-md">Remote Available</span>
                        ) : (
                          <span className="text-xs text-gray-400 font-medium pl-2">Office Based</span>
                        )}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan={6} className="text-center py-12 text-sm text-gray-400">
                      No active listings found matching active query params.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </div>
      </div>

      {/* Detail Overlay Card Layout Panel block */}
      {selectedVacancy && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-zinc-950/40 backdrop-blur-sm">
          <div className="bg-white w-full max-w-2xl rounded-xl border border-zinc-200/80 shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
            <div className="px-6 py-4 bg-zinc-50 border-b border-zinc-100 flex items-start justify-between">
              <div>
                <h3 className="text-lg font-bold text-zinc-900">{selectedVacancy.job_role}</h3>
                <p className="text-sm text-zinc-600 font-medium">{selectedVacancy.employer}</p>
              </div>
              <button onClick={() => setSelectedVacancy(null)} className="text-zinc-400 hover:text-zinc-700 p-1 rounded-md">
                ✕
              </button>
            </div>
            <div className="p-6 space-y-5 overflow-y-auto text-sm">
              <div className="grid grid-cols-2 gap-4 bg-zinc-50/60 p-4 rounded-lg border border-zinc-100">
                <DetailItem label="Location" value={selectedVacancy.location} />
                <DetailItem label="Province" value={selectedVacancy.meta_data.geo.province} />
                <DetailItem label="Category" value={selectedVacancy.meta_data.standardized_category} />
                <DetailItem label="Seniority Level" value={selectedVacancy.meta_data.seniority} />
              </div>
              <div>
                <h4 className="text-xs font-bold uppercase text-zinc-400 tracking-wider mb-1">Responsibilities</h4>
                <p className="text-zinc-700 p-3 bg-white border rounded-lg shadow-xs">{selectedVacancy.key_responsibilities}</p>
              </div>
              <div>
                <h4 className="text-xs font-bold uppercase text-zinc-400 tracking-wider mb-1">Qualifications</h4>
                <p className="text-zinc-700 p-3 bg-white border rounded-lg shadow-xs">{selectedVacancy.qualifications}</p>
              </div>
              <div>
                <h4 className="text-xs font-bold uppercase text-zinc-400 tracking-wider mb-2">Identified Skills</h4>
                <div className="flex flex-wrap gap-1.5">
                  {selectedVacancy.skills.map((skill) => (
                    <span key={skill} className="px-2.5 py-1 bg-zinc-100 text-zinc-800 text-xs font-medium rounded-md border">{skill}</span>
                  ))}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function DetailItem({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <p className="text-[10px] text-zinc-400 uppercase font-bold tracking-wider">{label}</p>
      <p className="text-sm font-semibold text-zinc-800 mt-0.5">{value}</p>
    </div>
  );
}