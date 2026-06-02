"use client";

import { useState } from "react";
import Header from "@/components/layout/Header";
import { useVacancyMetadata, useCategoryAnalytics } from "@/hooks/use-vacancies";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";

// A professional, modern dashboard color scheme array
const COLOR_PALETTE = [
  "#3b82f6", // Blue
  "#ef4444", // Red
  "#f59e0b", // Amber
  "#10b981", // Emerald
  "#8b5cf6", // Purple
  "#ec4899", // Pink
  "#06b6d4", // Cyan
  "#f97316"  // Orange
];

const categoryColorRegistry: Record<string, string> = {};
let colorIndexPointer = 0;

const getDynamicCategoryColor = (categoryName: string): string => {
  if (!categoryName) return "#64748b"; 

  if (categoryColorRegistry[categoryName]) {
    return categoryColorRegistry[categoryName];
  }

  const assignedColor = COLOR_PALETTE[colorIndexPointer % COLOR_PALETTE.length];
  categoryColorRegistry[categoryName] = assignedColor;

  colorIndexPointer++;

  return assignedColor;
};

const CATEGORY_COLORS = (cat: string) => {
  switch (cat) {
    case "IT": return "#3b82f6";
    case "Medicine": return "#ef4444";
    case "Construction": return "#f59e0b";
    default: return "#10b981";
  }
};

export default function CategoryPage() {
  const [selectedCategory, setSelectedCategory] = useState("All Categories");

  // Hook 1: Fetch all standard filter categories via global metadata endpoint
  const { data: metadata, isLoading: isMetaLoading } = useVacancyMetadata();

  // Hook 2: Fetch aggregated metrics tailored to the current selected category parameter
  const { data: analytics, isLoading: isAnalyticsLoading, isError } = useCategoryAnalytics(selectedCategory);

  if (isMetaLoading || isAnalyticsLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="flex flex-col items-center gap-3">
          <div className="w-8 h-8 border-4 border-zinc-900 border-t-transparent rounded-full animate-spin" />
          <p className="text-xs font-medium text-gray-500">Compiling Market Demand Curves via BFF...</p>
        </div>
      </div>
    );
  }

  if (isError || !analytics) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
        <div className="bg-white p-6 rounded-xl border border-red-100 shadow-sm max-w-sm text-center">
          <div className="text-red-500 font-bold text-sm mb-1">Analytics Disconnected</div>
          <p className="text-xs text-gray-500">Failed to stream aggregated intelligence values.</p>
        </div>
      </div>
    );
  }

  const categoriesOptions = ["All Categories", ...(metadata?.categories?.map(c => c.value) || [])];
  const chartSkillsList = analytics.skills.slice(0, 15);

  return (
    <div className="min-h-screen bg-gray-50">
      <Header 
        title="Category Analysis" 
        subtitle="Tracking talent demand curves, regional job distribution metrics, and domain skill clusters" 
      />

      <div className="p-8 space-y-6">
        {/* Dynamic Filter Dropdown */}
        <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="max-w-xs w-full">
            <label className="block text-xs text-gray-500 mb-1 font-medium">Filter by Core Category</label>
            <select
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
              className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white cursor-pointer outline-none focus:ring-2 focus:ring-zinc-900"
            >
              {categoriesOptions.map((cat) => (
                <option key={cat} value={cat}>{cat}</option>
              ))}
            </select>
          </div>
          <div className="text-sm text-gray-500 font-medium">
            Displaying <span className="text-zinc-900 font-bold">{analytics.skills.length}</span> critical competencies
          </div>
        </div>

        {/* Charts Section Block Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Top Skills Volume Chart */}
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm">
            <h3 className="text-sm font-semibold text-gray-700 mb-4">
              Top Skills by Volume {selectedCategory !== "All Categories" && `— ${selectedCategory}`}
            </h3>
            <ResponsiveContainer width="100%" height={380}>
              <BarChart data={chartSkillsList} layout="vertical" margin={{ left: 140, right: 10 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f8f9fa" />
                <XAxis type="number" tick={{ fontSize: 11 }} />
                <YAxis type="category" dataKey="skill" tick={{ fontSize: 11 }} width={140} />
                <Tooltip formatter={(value) => [`${value} Open Vacancies`, "Demand"]} />
                <Bar dataKey="demand" radius={[0, 4, 4, 0]}>
                  {chartSkillsList.map((s, i) => (
                    //<Cell key={i} fill={CATEGORY_COLORS(s.category)} 
                    <Cell key={i} fill={getDynamicCategoryColor(s.category)} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>

          {/* Regional Job Distribution Chart */}
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm">
            <h3 className="text-sm font-semibold text-gray-700 mb-4">
              Province-Wise Job Vacancies {selectedCategory !== "All Categories" && `— ${selectedCategory}`}
            </h3>
            <ResponsiveContainer width="100%" height={380}>
              <BarChart data={analytics.provinces} margin={{ bottom: 10, left: 10, right: 10 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#f8f9fa" />
                <XAxis dataKey="name" tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} />
                <Tooltip formatter={(value) => [`${value} Open Positions`, "Vacancies"]} />
                <Bar dataKey="vacancies" radius={[4, 4, 0, 0]}>
                  {analytics.provinces.map((entry, index) => (
                    <Cell 
                      key={`cell-${index}`} 
                      fill={selectedCategory === "All Categories" ? "#10b981" : getDynamicCategoryColor(selectedCategory)}
                      //fill={selectedCategory === "All Categories" ? "#10b981" : CATEGORY_COLORS(selectedCategory)} 
                    />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Data Grid & Target Hiring Context Block Cards */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Detailed Skills Analysis Table */}
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm lg:col-span-2">
            <h3 className="text-sm font-semibold text-gray-700 mb-4">Complete Category Skills Breakdown</h3>
            <div className="overflow-x-auto">
              <table className="w-full text-left border-collapse">
                <thead>
                  <tr className="border-b border-gray-200 bg-gray-50/70">
                    <th className="py-3 px-4 text-xs font-semibold uppercase tracking-wider text-gray-500 w-16">No</th>
                    <th className="py-3 px-4 text-xs font-semibold uppercase tracking-wider text-gray-500">Tracked Competency / Skill</th>
                    <th className="py-3 px-4 text-xs font-semibold uppercase tracking-wider text-gray-500 w-44 text-right">Active Vacancies</th>
                  </tr>
                </thead>
                <tbody>
                  {analytics.skills.map((s, i) => (
                    <tr key={s.skill} className="border-b border-gray-100 hover:bg-gray-50/40 transition-colors">
                      <td className="py-3.5 px-4 text-gray-400 text-xs font-mono">{i + 1}</td>
                      <td className="py-3.5 px-4 font-semibold text-gray-900 text-sm">{s.skill}</td>
                      <td className="py-3.5 px-4 font-mono text-sm text-gray-700 text-right">{s.demand.toLocaleString()}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Hiring Institutions Dashboard Panel Card */}
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex flex-col justify-between">
            <div>
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-semibold text-gray-700">Top Hiring Employers</h3>
                <span 
                  className="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded border"
                    style={{ 
                        borderColor: getDynamicCategoryColor(selectedCategory), 
                        color: getDynamicCategoryColor(selectedCategory),
                        backgroundColor: `${getDynamicCategoryColor(selectedCategory)}08`
                    }}
                //   style={{ 
                //     borderColor: CATEGORY_COLORS(selectedCategory), 
                //     color: CATEGORY_COLORS(selectedCategory),
                //     backgroundColor: `${CATEGORY_COLORS(selectedCategory)}08`
                //   }}
                >
                  {selectedCategory === "All Categories" ? "All Sectors" : selectedCategory}
                </span>
              </div>
              <p className="text-xs text-gray-400 mb-4">
                Enterprise institutions and corporations exhibiting the highest volume of ongoing requisition pipelines.
              </p>
              
              <div className="space-y-3">
                {analytics.employers.map((employer) => (
                  <div key={employer.name} className="p-3.5 border border-gray-100 rounded-lg bg-gray-50/30 flex items-center justify-between">
                    <div className="space-y-0.5">
                      <div className="text-sm font-semibold text-gray-900">{employer.name}</div>
                      <div className="text-[11px] text-gray-400 font-medium">{employer.location}</div>
                    </div>
                    <div className="text-right">
                      <div className="text-sm font-bold text-gray-800 font-mono">{employer.openRoles}</div>
                      <div className="text-[10px] uppercase font-bold text-gray-400 tracking-tight">Openings</div>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="mt-6 pt-4 border-t border-gray-100 flex items-center justify-between text-xs text-gray-400">
              <span>Data source: Live Feed Integration</span>
              <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}