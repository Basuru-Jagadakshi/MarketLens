"use client";

import { useMemo, useState } from "react";
import Link from "next/link";

// ─────────────────────────────────────────────────────────────
// OCCUPATION DEMAND FORECAST
// Shows the most in-demand occupations for three horizons:
// Short Term, Mid Term, Long Term. Top 100 occupation names per
// horizon, selected/ranked by the forecasting pipeline. The mock
// ranking below is deterministic per term — replace topForTerm()
// with the API call (e.g. GET /api/v1/forecast/occupations?term=
// short|mid|long) when the backend endpoint is ready.
// ─────────────────────────────────────────────────────────────

type TermKey = "short" | "mid" | "long";

const TERMS: { key: TermKey; label: string; desc: string }[] = [
  { key: "short", label: "Short Term", desc: "Projected demand over the next 6 months" },
  { key: "mid", label: "Mid Term", desc: "Projected demand over the next 1 year" },
  { key: "long", label: "Long Term", desc: "Projected demand over the next 3 year" },
];

const TOP_N = 100;

// Mock pool of realistic occupation titles across the Sri Lankan
// labour market. The real page will receive titles from the
// forecaster's precomputed results.
const OCCUPATION_POOL: string[] = [
  // ICT
  "Software Engineer",
  "Backend Developer",
  "Frontend Developer",
  "Full Stack Developer",
  "Mobile Application Developer",
  "DevOps Engineer",
  "QA Engineer",
  "Data Analyst",
  "Data Scientist",
  "AI / Machine Learning Engineer",
  "Cybersecurity Analyst",
  "Network Engineer",
  "Systems Administrator",
  "Database Administrator",
  "IT Support Technician",
  "UI/UX Designer",
  "Cloud Engineer",
  "Solutions Architect",
  "Business Analyst",
  "IT Project Manager",
  "Product Manager",
  // Health
  "Registered Nurse",
  "Midwife",
  "Pharmacist",
  "Medical Laboratory Technician",
  "Physiotherapist",
  "Caregiver",
  "Medical Officer",
  "Dental Assistant",
  "Public Health Inspector",
  "Radiographer",
  // Finance & admin
  "Accountant",
  "Audit Associate",
  "Finance Executive",
  "Bank Teller",
  "Credit Officer",
  "Insurance Advisor",
  "Investment Analyst",
  "Payroll Officer",
  "Bookkeeper",
  "Tax Consultant",
  "HR Executive",
  "Recruitment Officer",
  "Administrative Assistant",
  "Office Clerk",
  "Data Entry Operator",
  "Legal Officer",
  // Education
  "Primary School Teacher",
  "Secondary School Teacher",
  "Vocational Instructor",
  "University Lecturer",
  "Preschool Teacher",
  "English Language Instructor",
  "ICT Instructor",
  // Sales, marketing & service
  "Digital Marketing Executive",
  "Marketing Manager",
  "Sales Executive",
  "Sales Manager",
  "Merchandiser",
  "Brand Executive",
  "Content Writer",
  "Graphic Designer",
  "Social Media Coordinator",
  "SEO Specialist",
  "Customer Service Representative",
  "Call Center Agent",
  "Receptionist",
  "Cashier",
  "Retail Sales Assistant",
  // Apparel & manufacturing
  "Garment Machine Operator",
  "Quality Checker (Apparel)",
  "Fashion Designer",
  "Pattern Maker",
  "Textile Technician",
  "Production Supervisor",
  "Industrial Engineer",
  "Machine Operator",
  "Assembly Line Worker",
  // Hospitality & tourism
  "Chef",
  "Cook",
  "Baker",
  "Barista",
  "Waiter / Steward",
  "Hotel Manager",
  "Housekeeping Supervisor",
  "Front Office Executive",
  "Tour Guide",
  "Travel Consultant",
  "Event Coordinator",
  // Construction & trades
  "Mason",
  "Carpenter",
  "Electrician",
  "Plumber",
  "Welder",
  "Painter",
  "Quantity Surveyor",
  "Civil Engineer",
  "Site Supervisor",
  "Draughtsman",
  "Architect",
  "Mechanical Engineer",
  "Electrical Engineer",
  "Maintenance Technician",
  "Air Conditioning Technician",
  "Solar PV Technician",
  "CCTV Technician",
  "Automobile Mechanic",
  "Motorcycle Mechanic",
  // Transport & logistics
  "Heavy Vehicle Driver",
  "Light Vehicle Driver",
  "Forklift Operator",
  "Crane Operator",
  "Warehouse Assistant",
  "Logistics Coordinator",
  "Supply Chain Executive",
  "Procurement Officer",
  "Store Keeper",
  "Delivery Rider",
  // Agriculture, security & other services
  "Security Officer",
  "Farm Supervisor",
  "Agricultural Technician",
  "Fisheries Worker",
  "Landscaper",
  "Veterinary Assistant",
  "Translator",
  "Journalist",
  "Photographer",
  "Video Editor",
];

// Deterministic mock ranking per term (FNV-style hash), so each
// tab shows a stable but distinct top-100 ordering.
function hash01(s: string): number {
  let h = 2166136261;
  for (let i = 0; i < s.length; i++) {
    h ^= s.charCodeAt(i);
    h = Math.imul(h, 16777619);
  }
  return ((h >>> 0) % 10000) / 10000;
}

function topForTerm(term: TermKey): string[] {
  return OCCUPATION_POOL.map((name) => ({
    name,
    score: hash01(term + "::" + name),
  }))
    .sort((a, b) => b.score - a.score)
    .slice(0, TOP_N)
    .map((o) => o.name);
}

// ─────────────────────────────────────────────────────────────
// PAGE
// ─────────────────────────────────────────────────────────────
export default function ForecastPage() {
  const [activeTerm, setActiveTerm] = useState<TermKey>("short");

  const occupations = useMemo(() => topForTerm(activeTerm), [activeTerm]);
  const term = TERMS.find((t) => t.key === activeTerm)!;

  return (
    <div className="min-h-screen w-full min-w-0 bg-zinc-50 text-zinc-800 font-sans">
      {/* HEADER — same as the main dashboard page */}
      <header className="bg-white border-b border-zinc-200 px-4 md:px-8 py-5 sticky top-0 z-40">
        <div className="flex justify-between items-start gap-4">
          <div className="min-w-0">
            <h1 className="text-lg md:text-xl font-black text-zinc-900 tracking-tight truncate">
              Occupation Demand Forecast
            </h1>
            <p className="text-xs text-zinc-400 mt-1 font-medium">
              The most in demand occupations projected for the short, mid and long term
            </p>
          </div>
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white text-xs font-black shadow-md border border-white shrink-0">
            BJ
          </div>
        </div>
      </header>

      <div className="p-4 md:p-8 space-y-6 md:space-y-8 w-full max-w-[1400px] mx-auto pb-20">

        {/* FORECAST PANEL */}
        <div className="bg-white border border-zinc-200 rounded-xl shadow-sm">
          {/* TABS — term selector */}
          <div className="px-5 border-b border-zinc-200 flex gap-6">
            {TERMS.map((t) => (
              <button
                key={t.key}
                onClick={() => setActiveTerm(t.key)}
                className={`cursor-pointer py-3 text-xs font-bold border-b-2 -mb-px transition-colors ${
                  activeTerm === t.key
                    ? "border-zinc-900 text-zinc-900"
                    : "border-transparent text-zinc-400 hover:text-zinc-600"
                }`}
              >
                {t.label}
              </button>
            ))}
          </div>

          {/* TERM HEADER */}
          <div className="px-5 py-4 border-b border-zinc-200 bg-zinc-50/70 flex flex-wrap justify-between items-end gap-4">
            <div className="min-w-0">
              <p className="text-[10px] font-black uppercase tracking-wider text-zinc-400 mb-1">
                {term.label} Forecast
              </p>
              <p className="text-sm font-bold text-zinc-700">{term.desc}</p>
            </div>
            <div className="text-right shrink-0">
              <p className="text-[10px] font-black uppercase tracking-wider text-zinc-400">
                Occupations Listed
              </p>
              <p className="text-lg font-black text-zinc-900">
                Top {occupations.length}
              </p>
            </div>
          </div>

          {/* TOP 100 OCCUPATIONS — names only, ranked */}
          <div className="p-4 md:p-5">
            <ol className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-2">
              {occupations.map((name, i) => (
                <li
                  key={name}
                  className="flex items-center gap-3 bg-zinc-50 border border-zinc-100 rounded-lg px-3 py-2.5"
                >
                  <span
                    className={`shrink-0 w-7 h-7 rounded-full text-[11px] font-black flex items-center justify-center ${
                      i < 10
                        ? "bg-zinc-900 text-white"
                        : "bg-zinc-200 text-zinc-600"
                    }`}
                  >
                    {i + 1}
                  </span>
                  <span className="text-xs font-bold text-zinc-800 truncate">
                    {name}
                  </span>
                </li>
              ))}
            </ol>
          </div>
        </div>
      </div>
    </div>
  );
}