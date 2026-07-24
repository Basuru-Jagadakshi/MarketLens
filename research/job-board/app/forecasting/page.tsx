"use client";

import { useState } from "react";

/**
 * ROUTE: app/forecasting/page.tsx (Next.js App Router)
 * Pages Router equivalent: pages/forecasting.tsx
 *
 * TWO SECTIONS ON THIS PAGE:
 *
 * 1. "Labour Demand vs. Supply Forecast" — the client's specific ask: a
 *    simple, plain-language view of demand & supply projections across
 *    three horizons (short/medium/long-term). No filters, no dense charts —
 *    just a number, a one-line explanation, and a simple proportional bar
 *    per horizon so the widening gap is easy to see at a glance.
 *
 * 2. "Future of Work — AI Insights" — the topic cards from before
 *    (fastest-growing roles, key drivers, future skills, etc.)
 *
 * WHAT'S MOCKED
 * ----------------------------------------------------------------------
 * FORECAST_HORIZONS below holds static, national-aggregate numbers.
 * Swap the `demand` / `supply` / `explanation` fields for data from your
 * own forecasting model or an API route (e.g. /api/forecast) once that's
 * ready — the rest of the UI doesn't need to change.
 */

type HorizonId = "short" | "medium" | "long";

type ForecastHorizon = {
  id: HorizonId;
  label: string;
  subtitle: string;
  accent: string;
  demand: number;
  supply: number;
  demandGrowth: number; // %
  supplyGrowth: number; // %
  outlook: string;
  explanation: string;
};

const ACCENT: Record<string, { text: string; bg: string; lightBg: string; border: string }> = {
  blue: { text: "text-blue-600", bg: "bg-blue-600", lightBg: "bg-blue-50", border: "border-blue-100" },
  amber: { text: "text-amber-600", bg: "bg-amber-500", lightBg: "bg-amber-50", border: "border-amber-100" },
  violet: { text: "text-violet-600", bg: "bg-violet-600", lightBg: "bg-violet-50", border: "border-violet-100" },
};

// ---------------------------------------------------------------------------
// FORECAST DATA (mock — national aggregate, base year vacancies ≈ 12,450 demand / 10,800 supply)
// ---------------------------------------------------------------------------
const FORECAST_HORIZONS: ForecastHorizon[] = [
  {
    id: "short",
    label: "Short-Term",
    subtitle: "1–2 Years",
    accent: "blue",
    demand: 13450,
    supply: 11340,
    demandGrowth: 8,
    supplyGrowth: 5,
    outlook: "Emerging Shortage",
    explanation:
      "Over the next 1–2 years, vacancies are expected to keep growing steadily — especially in technology, healthcare and tourism-linked services. New graduates and returning workers should keep pace with most of this demand, though specialist digital roles will start to feel tight.",
  },
  {
    id: "medium",
    label: "Medium-Term",
    subtitle: "3–5 Years",
    accent: "amber",
    demand: 14820,
    supply: 12100,
    demandGrowth: 19,
    supplyGrowth: 12,
    outlook: "Widening Shortage",
    explanation:
      "Between three and five years out, demand growth is expected to accelerate — driven by AI adoption, automation and green-economy investment — while the supply pipeline (graduates, TVET output, returning workers) grows more slowly. Technical and vocational roles will feel this gap first.",
  },
  {
    id: "long",
    label: "Long-Term",
    subtitle: "5–10 Years",
    accent: "violet",
    demand: 17180,
    supply: 13390,
    demandGrowth: 38,
    supplyGrowth: 24,
    outlook: "Significant Shortage",
    explanation:
      "Over five to ten years, an ageing workforce, continued emigration of skilled workers, and fast-moving technology change could push demand well ahead of supply — unless the country invests heavily in education, reskilling and retention. This is the horizon where action today has the biggest payoff.",
  },
];

const MAX_VALUE = Math.max(...FORECAST_HORIZONS.map((h) => Math.max(h.demand, h.supply)));

export default function ForecastingPage() {
  const [variant, setVariant] = useState<0 | 1>(0);
  const [isGenerating, setIsGenerating] = useState(false);
  const [generatedAt, setGeneratedAt] = useState<string | null>(null);

  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");
    const d = currentLang;
    const formattedDate = new Date().toLocaleDateString(
        currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
        { year: "numeric", month: "short", day: "numeric" }
    );

  function handleRegenerate() {
    setIsGenerating(true);
    setTimeout(() => {
      setVariant((v) => (v === 0 ? 1 : 0));
      setGeneratedAt(new Date().toLocaleString());
      setIsGenerating(false);
    }, 900);
  }

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800 font-sans">
      {/* HEADER */}
      <header className="bg-white border-b border-gray-100 px-8 py-5 flex flex-col md:flex-row justify-between items-start md:items-center gap-4 sticky top-0 z-40">
        <div>
          <h1 className="text-xl font-black text-gray-900 tracking-tight">Labour Market Forecasting</h1>
          <p className="text-xs text-gray-400 mt-1 font-medium">
            Simple, AI-generated demand & supply projections and future-of-work insights
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

      <div className="p-8 max-w-6xl mx-auto pb-16 space-y-10">

        {/* ── SECTION 1: LABOUR DEMAND VS SUPPLY FORECAST ─────────────── */}
        <section>
          <div className="mb-5">
            <h2 className="text-base font-black text-gray-900">Labour Demand vs. Supply Forecast</h2>
            <p className="text-xs text-gray-400 mt-1">
              National outlook across three horizons — in plain terms, not just numbers.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {FORECAST_HORIZONS.map((h) => {
              const acc = ACCENT[h.accent];
              const gap = h.demand - h.supply;
              const demandWidth = (h.demand / MAX_VALUE) * 100;
              const supplyWidth = (h.supply / MAX_VALUE) * 100;

              return (
                <div key={h.id} className={`bg-white p-5 rounded-xl shadow-sm border ${acc.border}`}>
                  <div className="flex items-center justify-between">
                    <div>
                      <h3 className="text-sm font-black text-gray-900">{h.label}</h3>
                      <p className="text-[11px] text-gray-400 font-semibold">{h.subtitle}</p>
                    </div>
                    <span className={`text-[10px] font-bold px-2.5 py-1 rounded-full ${acc.lightBg} ${acc.text}`}>
                      {h.outlook}
                    </span>
                  </div>

                  {/* Simple proportional bars — no chart library needed */}
                  <div className="mt-4 space-y-2.5">
                    <div>
                      <div className="flex justify-between text-[10px] font-bold text-gray-500 mb-1">
                        <span>Demand</span>
                        <span className="font-mono text-blue-600">
                          {h.demand.toLocaleString()} (+{h.demandGrowth}%)
                        </span>
                      </div>
                      <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                        <div className="h-full bg-blue-600 rounded-full transition-all duration-700" style={{ width: `${demandWidth}%` }} />
                      </div>
                    </div>
                    <div>
                      <div className="flex justify-between text-[10px] font-bold text-gray-500 mb-1">
                        <span>Supply</span>
                        <span className="font-mono text-emerald-600">
                          {h.supply.toLocaleString()} (+{h.supplyGrowth}%)
                        </span>
                      </div>
                      <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                        <div className="h-full bg-emerald-600 rounded-full transition-all duration-700" style={{ width: `${supplyWidth}%` }} />
                      </div>
                    </div>
                  </div>

                  {/* Plain-language bottom line */}
                  <div className="mt-4 pt-3 border-t border-gray-50">
                    <p className="text-[11px] font-bold text-gray-700">
                      Bottom line: a projected gap of{" "}
                      <span className={acc.text}>{gap.toLocaleString()} workers</span> — demand is growing faster than supply can keep up.
                    </p>
                    <p className="text-xs text-gray-500 leading-relaxed mt-2">{h.explanation}</p>
                  </div>
                </div>
              );
            })}
          </div>
        </section>

        {/* ── SECTION 2: FUTURE-OF-WORK TOPIC INSIGHTS ────────────────── */}
        <section>
          <div className="mb-5">
            <h2 className="text-base font-black text-gray-900">Future of Work — AI Insights</h2>
            <p className="text-xs text-gray-400 mt-1">
              Short, AI-generated explanations on where jobs and skills are headed.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {TOPICS.map((topic) => {
              const acc = TOPIC_ACCENT[topic.accent];
              return (
                <div key={topic.id} className="bg-white p-5 rounded-xl shadow-sm">
                  <div className="flex items-center gap-3 border-b border-gray-50 pb-3">
                    <span className={`w-9 h-9 rounded-lg flex items-center justify-center shrink-0 ${acc.lightBg} ${acc.text}`}>
                      {topic.icon}
                    </span>
                    <h3 className="text-sm font-bold text-gray-800">{topic.title}</h3>
                  </div>

                  <p className="text-xs text-gray-600 leading-relaxed mt-3">
                    {topic.summaryVariants[variant]}
                  </p>

                  {topic.items && (
                    <div className="mt-4">
                      {topic.itemsHeader && (
                        <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-2">
                          {topic.itemsHeader}
                        </h4>
                      )}
                      <div className="flex flex-wrap gap-2">
                        {topic.items.map((item, i) => (
                          <span
                            key={i}
                            className="inline-flex items-center gap-1.5 text-[11px] font-semibold text-gray-700 bg-gray-50 border border-gray-100 px-2.5 py-1.5 rounded-lg"
                          >
                            {item.label}
                            {item.meta && (
                              <span className={`font-mono font-bold ${acc.text}`}>{item.meta}</span>
                            )}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>

          <p className="text-[11px] text-gray-400 text-center mt-8">
            These insights are AI-generated directional summaries, not guaranteed outcomes — use them as a starting point for discussion, not a forecast to plan finances around.
          </p>
        </section>

      </div>
    </div>
  );
}

// ---------------------------------------------------------------------------
// TOPIC CARDS DATA (unchanged from the previous "Future of Work" view)
// ---------------------------------------------------------------------------
type TopicItem = { label: string; meta?: string };

type Topic = {
  id: string;
  title: string;
  accent: string;
  icon: JSX.Element;
  summaryVariants: [string, string];
  itemsHeader?: string;
  items?: TopicItem[];
};

const iconProps = { className: "w-5 h-5", fill: "none", viewBox: "0 0 24 24", stroke: "currentColor", strokeWidth: 2 };

const IconTrendingUp = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18L9 11.25l4.306 4.306a11.95 11.95 0 015.814-5.518l2.74-1.22m0 0l-5.94-2.28m5.94 2.28l-2.28 5.941" /></svg>
);
const IconTrendingDown = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M2.25 6L9 12.75l4.306-4.306a11.95 11.95 0 015.814 5.518l2.74 1.22m0 0l-5.94 2.28m5.94-2.28l-2.28-5.941" /></svg>
);
const IconBolt = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M3.75 13.5l10.5-11.25L12 10.5h8.25L9.75 21.75 12 13.5H3.75z" /></svg>
);
const IconAcademic = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M4.26 10.147a60.436 60.436 0 00-.491 6.347A48.62 48.62 0 0112 20.904a48.62 48.62 0 018.232-4.41 60.46 60.46 0 00-.491-6.347m-15.482 0a50.636 50.636 0 00-2.658-.813A59.906 59.906 0 0112 3.493a59.903 59.903 0 0110.399 5.84c-.896.248-1.783.52-2.658.814m-15.482 0A50.717 50.717 0 0112 13.489a50.702 50.702 0 017.74-3.342M6.75 15a.75.75 0 100-1.5.75.75 0 000 1.5zm0 0v-3.675A55.378 55.378 0 0112 8.443" /></svg>
);
const IconShuffle = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5" /></svg>
);
const IconSparkles = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09z" /></svg>
);
const IconHeart = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M21 8.25c0-2.485-2.099-4.5-4.688-4.5-1.935 0-3.597 1.126-4.312 2.733-.715-1.607-2.377-2.733-4.313-2.733C5.1 3.75 3 5.765 3 8.25c0 7.22 9 12 9 12s9-4.78 9-12z" /></svg>
);
const IconBriefcase = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M20.25 14.15v4.25c0 1.094-.787 2.036-1.872 2.18a48.99 48.99 0 01-12.756 0C4.537 20.436 3.75 19.494 3.75 18.4v-4.25m16.5 0a2.18 2.18 0 00.75-1.653v-2.34a2.25 2.25 0 00-.622-1.55l-1.06-1.114a2.25 2.25 0 00-1.628-.693H6.31a2.25 2.25 0 00-1.628.693l-1.06 1.114A2.25 2.25 0 003 10.157v2.34c0 .633.29 1.226.75 1.653m16.5 0a2.25 2.25 0 01-1.5.573H5.25a2.25 2.25 0 01-1.5-.573m16.5 0v-1.5a3 3 0 00-3-3H8.25a3 3 0 00-3 3v1.5" /></svg>
);
const IconBook = (
  <svg {...iconProps}><path strokeLinecap="round" strokeLinejoin="round" d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25" /></svg>
);

const TOPIC_ACCENT: Record<string, { text: string; bg: string; lightBg: string; ring: string }> = {
  emerald: { text: "text-emerald-600", bg: "bg-emerald-600", lightBg: "bg-emerald-50", ring: "ring-emerald-100" },
  red: { text: "text-red-500", bg: "bg-red-500", lightBg: "bg-red-50", ring: "ring-red-100" },
  amber: { text: "text-amber-600", bg: "bg-amber-500", lightBg: "bg-amber-50", ring: "ring-amber-100" },
  indigo: { text: "text-indigo-600", bg: "bg-indigo-600", lightBg: "bg-indigo-50", ring: "ring-indigo-100" },
  cyan: { text: "text-cyan-600", bg: "bg-cyan-600", lightBg: "bg-cyan-50", ring: "ring-cyan-100" },
  violet: { text: "text-violet-600", bg: "bg-violet-600", lightBg: "bg-violet-50", ring: "ring-violet-100" },
  pink: { text: "text-pink-600", bg: "bg-pink-600", lightBg: "bg-pink-50", ring: "ring-pink-100" },
  blue: { text: "text-blue-600", bg: "bg-blue-600", lightBg: "bg-blue-50", ring: "ring-blue-100" },
  teal: { text: "text-teal-600", bg: "bg-teal-600", lightBg: "bg-teal-50", ring: "ring-teal-100" },
};

const TOPICS: Topic[] = [
  {
    id: "growing-roles",
    title: "Fastest-Growing Job Roles",
    accent: "emerald",
    icon: IconTrendingUp,
    summaryVariants: [
      "Roles tied to AI, clean energy and healthcare are expanding fastest, as employers race to build capacity in areas that can't easily be automated or outsourced.",
      "Demand is surging for roles that combine technical skill with human judgement — AI specialists, renewable energy technicians and care professionals lead the pack.",
    ],
    itemsHeader: "Top movers",
    items: [
      { label: "AI / ML Engineer", meta: "+38%" },
      { label: "Renewable Energy Technician", meta: "+31%" },
      { label: "Cybersecurity & Data Privacy Analyst", meta: "+27%" },
      { label: "Healthcare & Elder-care Specialist", meta: "+24%" },
      { label: "Digital Marketing / Growth Specialist", meta: "+19%" },
    ],
  },
  {
    id: "declining-roles",
    title: "Fastest-Declining Roles",
    accent: "red",
    icon: IconTrendingDown,
    summaryVariants: [
      "Routine, rules-based roles are shrinking the fastest, as automation and self-service tools take over repetitive data entry, transaction and manual-assembly work.",
      "Jobs built around repeatable, predictable tasks are contracting quickly — software and machines now do this work faster and cheaper than people.",
    ],
    itemsHeader: "Steepest declines",
    items: [
      { label: "Data Entry Clerk", meta: "-29%" },
      { label: "Bank Teller", meta: "-24%" },
      { label: "Manual Assembly Line Operator", meta: "-21%" },
      { label: "Print & Publishing Technician", meta: "-18%" },
      { label: "Traditional Travel Agent", meta: "-15%" },
    ],
  },
  {
    id: "key-drivers",
    title: "Key Drivers of Change",
    accent: "amber",
    icon: IconBolt,
    summaryVariants: [
      "Five forces are reshaping the job market at once: AI and automation, the shift to green energy, an ageing population, global supply-chain realignment, and changing worker expectations.",
      "The labour market is being pulled in new directions by technology, climate action, demographics, geopolitics and a workforce that now expects more flexibility.",
    ],
    itemsHeader: "Main forces",
    items: [
      { label: "Artificial intelligence & automation" },
      { label: "Climate change & the green transition" },
      { label: "Demographic shifts (ageing, urbanisation)" },
      { label: "Geopolitical & supply-chain realignment" },
      { label: "Changing worker & consumer expectations" },
    ],
  },
  {
    id: "future-skills",
    title: "Required Skills for the Future",
    accent: "indigo",
    icon: IconAcademic,
    summaryVariants: [
      "Technical fluency alone won't be enough — the most durable skills combine data and AI literacy with critical thinking, adaptability and the ability to work well with others.",
      "As tools change faster than job titles, the winning skill set is a mix of AI/data literacy, problem-solving, creativity and comfort with continuous learning.",
    ],
    itemsHeader: "Skills to build",
    items: [
      { label: "Critical thinking & complex problem-solving" },
      { label: "AI & data literacy" },
      { label: "Adaptability and continuous learning" },
      { label: "Cross-cultural collaboration" },
      { label: "Creativity and design thinking" },
    ],
  },
  {
    id: "work-structure",
    title: "Shift in Work Structure",
    accent: "cyan",
    icon: IconShuffle,
    summaryVariants: [
      "Work itself is being reorganised — hybrid arrangements, freelance and project-based contracts, and flatter, AI-assisted teams are becoming the norm rather than the exception.",
      "Traditional 9-to-5, single-employer careers are giving way to hybrid work, gig and freelance arrangements, and leaner teams supported by AI copilots.",
    ],
    itemsHeader: "What's changing",
    items: [
      { label: "Hybrid & fully-remote arrangements" },
      { label: "Rise of gig, freelance & project-based work" },
      { label: "Flatter, cross-functional teams" },
      { label: "Four-day work-week pilots" },
      { label: "AI copilots reshaping task distribution" },
    ],
  },
  {
    id: "professions-2050",
    title: "Professions in Demand in 2050",
    accent: "violet",
    icon: IconSparkles,
    summaryVariants: [
      "Looking further out, entirely new professions are likely to emerge around climate adaptation, human-AI collaboration, longevity care and even space industries.",
      "By 2050, expect whole job categories that barely exist today — built around managing AI systems, adapting to climate change, and supporting an older global population.",
    ],
    itemsHeader: "Emerging roles",
    items: [
      { label: "Climate Adaptation Specialist" },
      { label: "Human-AI Collaboration Manager" },
      { label: "Longevity & Geriatric Care Specialist" },
      { label: "Space Industry Technician" },
      { label: "Synthetic Biology Engineer" },
    ],
  },
  {
    id: "future-proof-child",
    title: "How to Future-Proof Your Child's Career",
    accent: "pink",
    icon: IconHeart,
    summaryVariants: [
      "The best preparation isn't picking the 'right' job today — it's building curiosity, strong fundamentals and the confidence to keep learning as roles keep changing.",
      "Rather than aiming for one future-proof job title, focus on skills and habits that transfer across careers: curiosity, adaptability and strong basics.",
    ],
    itemsHeader: "What helps most",
    items: [
      { label: "Encourage curiosity & lifelong learning" },
      { label: "Build strong maths, language & digital literacy" },
      { label: "Nurture creativity, empathy & collaboration" },
      { label: "Expose them to real-world problem solving" },
      { label: "Avoid locking into one narrow path too early" },
    ],
  },
  {
    id: "professions-relevance",
    title: "Professions of Relevance in the Coming Decades",
    accent: "blue",
    icon: IconBriefcase,
    summaryVariants: [
      "Healthcare, sustainability, technology maintenance, education and care-economy roles are expected to stay relevant and in demand well beyond the next decade.",
      "Look for durable demand in professions tied to human wellbeing, the green economy, and keeping technology (especially AI) running safely and effectively.",
    ],
    itemsHeader: "Enduring fields",
    items: [
      { label: "Healthcare & wellbeing professionals" },
      { label: "Green economy & sustainability roles" },
      { label: "AI, cybersecurity & cloud specialists" },
      { label: "Educators & reskilling specialists" },
      { label: "Care economy (childcare, eldercare, community)" },
    ],
  },
  {
    id: "fundamental-skills",
    title: "Strengthen Fundamental Skills",
    accent: "teal",
    icon: IconBook,
    summaryVariants: [
      "No matter how the job market shifts, literacy, numeracy, communication and emotional intelligence remain the foundation everything else is built on.",
      "Advanced tools change fast, but the basics don't: strong reading, numeracy, communication and self-management skills stay valuable in every scenario.",
    ],
    itemsHeader: "Fundamentals to prioritise",
    items: [
      { label: "Literacy & numeracy" },
      { label: "Digital & data fluency" },
      { label: "Written, verbal & cross-cultural communication" },
      { label: "Emotional intelligence & self-management" },
      { label: "Financial literacy for a freelance-driven economy" },
    ],
  },
];