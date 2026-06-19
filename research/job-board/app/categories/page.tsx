// "use client";

// import { useState } from "react";
// import Header from "@/components/layout/Header";
// import { useVacancyMetadata, useCategoryAnalytics } from "@/hooks/use-vacancies";
// import {
//   BarChart,
//   Bar,
//   XAxis,
//   YAxis,
//   CartesianGrid,
//   Tooltip,
//   ResponsiveContainer,
//   Cell,
// } from "recharts";

// const CATEGORY_COLORS = (cat: string) => {
//   switch (cat) {
//     case "IT": return "#3b82f6";
//     case "Medicine": return "#ef4444";
//     case "Construction": return "#f59e0b";
//     default: return "#10b981";
//   }
// };

// export default function CategoryPage() {
//   const [selectedCategory, setSelectedCategory] = useState("All Categories");

//   // Hook 1: Fetch all standard filter categories via global metadata endpoint
//   const { data: metadata, isLoading: isMetaLoading } = useVacancyMetadata();

//   // Hook 2: Fetch aggregated metrics tailored to the current selected category parameter
//   const { data: analytics, isLoading: isAnalyticsLoading, isError } = useCategoryAnalytics(selectedCategory);

//   if (isMetaLoading || isAnalyticsLoading) {
//     return (
//       <div className="min-h-screen flex items-center justify-center bg-gray-50">
//         <div className="flex flex-col items-center gap-3">
//           <div className="w-8 h-8 border-4 border-zinc-900 border-t-transparent rounded-full animate-spin" />
//           <p className="text-xs font-medium text-gray-500">Compiling Market Demand Curves via BFF...</p>
//         </div>
//       </div>
//     );
//   }

//   if (isError || !analytics) {
//     return (
//       <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
//         <div className="bg-white p-6 rounded-xl border border-red-100 shadow-sm max-w-sm text-center">
//           <div className="text-red-500 font-bold text-sm mb-1">Analytics Disconnected</div>
//           <p className="text-xs text-gray-500">Failed to stream aggregated intelligence values.</p>
//         </div>
//       </div>
//     );
//   }

//   const categoriesOptions = ["All Categories", ...(metadata?.categories?.map(c => c.value) || [])];
//   const chartSkillsList = analytics.skills.slice(0, 15);

//   return (
//     <div className="min-h-screen bg-gray-50">
//       <Header 
//         title="Category Analysis" 
//         subtitle="Tracking talent demand curves, regional job distribution metrics, and domain skill clusters" 
//       />

//       <div className="p-8 space-y-6">
//         {/* Dynamic Filter Dropdown */}
//         <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
//           <div className="max-w-xs w-full">
//             <label className="block text-xs text-gray-500 mb-1 font-medium">Filter by Core Category</label>
//             <select
//               value={selectedCategory}
//               onChange={(e) => setSelectedCategory(e.target.value)}
//               className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm bg-white cursor-pointer outline-none focus:ring-2 focus:ring-zinc-900"
//             >
//               {categoriesOptions.map((cat) => (
//                 <option key={cat} value={cat}>{cat}</option>
//               ))}
//             </select>
//           </div>
//           <div className="text-sm text-gray-500 font-medium">
//             Displaying <span className="text-zinc-900 font-bold">{analytics.skills.length}</span> critical competencies
//           </div>
//         </div>

//         {/* Charts Section Block Layout */}
//         <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
//           {/* Top Skills Volume Chart */}
//           <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm">
//             <h3 className="text-sm font-semibold text-gray-700 mb-4">
//               Top Skills by Volume {selectedCategory !== "All Categories" && `— ${selectedCategory}`}
//             </h3>
//             <ResponsiveContainer width="100%" height={380}>
//               <BarChart data={chartSkillsList} layout="vertical" margin={{ left: 140, right: 10 }}>
//                 <CartesianGrid strokeDasharray="3 3" stroke="#f8f9fa" />
//                 <XAxis type="number" tick={{ fontSize: 11 }} />
//                 <YAxis type="category" dataKey="skill" tick={{ fontSize: 11 }} width={140} />
//                 <Tooltip formatter={(value) => [`${value} Open Vacancies`, "Demand"]} />
//                 <Bar dataKey="demand" radius={[0, 4, 4, 0]}>
//                   {chartSkillsList.map((s, i) => (
//                     <Cell key={i} fill={CATEGORY_COLORS(s.category)} />
//                   ))}
//                 </Bar>
//               </BarChart>
//             </ResponsiveContainer>
//           </div>

//           {/* Regional Job Distribution Chart */}
//           <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm">
//             <h3 className="text-sm font-semibold text-gray-700 mb-4">
//               Province-Wise Job Vacancies {selectedCategory !== "All Categories" && `— ${selectedCategory}`}
//             </h3>
//             <ResponsiveContainer width="100%" height={380}>
//               <BarChart data={analytics.provinces} margin={{ bottom: 10, left: 10, right: 10 }}>
//                 <CartesianGrid strokeDasharray="3 3" stroke="#f8f9fa" />
//                 <XAxis dataKey="name" tick={{ fontSize: 11 }} />
//                 <YAxis tick={{ fontSize: 11 }} />
//                 <Tooltip formatter={(value) => [`${value} Open Positions`, "Vacancies"]} />
//                 <Bar dataKey="vacancies" radius={[4, 4, 0, 0]}>
//                   {analytics.provinces.map((entry, index) => (
//                     <Cell 
//                       key={`cell-${index}`} 
//                       fill={selectedCategory === "All Categories" ? "#10b981" : CATEGORY_COLORS(selectedCategory)} 
//                     />
//                   ))}
//                 </Bar>
//               </BarChart>
//             </ResponsiveContainer>
//           </div>
//         </div>

//         {/* Data Grid & Target Hiring Context Block Cards */}
//         <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
//           {/* Detailed Skills Analysis Table */}
//           <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm lg:col-span-2">
//             <h3 className="text-sm font-semibold text-gray-700 mb-4">Complete Category Skills Breakdown</h3>
//             <div className="overflow-x-auto">
//               <table className="w-full text-left border-collapse">
//                 <thead>
//                   <tr className="border-b border-gray-200 bg-gray-50/70">
//                     <th className="py-3 px-4 text-xs font-semibold uppercase tracking-wider text-gray-500 w-16">No</th>
//                     <th className="py-3 px-4 text-xs font-semibold uppercase tracking-wider text-gray-500">Tracked Competency / Skill</th>
//                     <th className="py-3 px-4 text-xs font-semibold uppercase tracking-wider text-gray-500 w-44 text-right">Active Vacancies</th>
//                   </tr>
//                 </thead>
//                 <tbody>
//                   {analytics.skills.map((s, i) => (
//                     <tr key={s.skill} className="border-b border-gray-100 hover:bg-gray-50/40 transition-colors">
//                       <td className="py-3.5 px-4 text-gray-400 text-xs font-mono">{i + 1}</td>
//                       <td className="py-3.5 px-4 font-semibold text-gray-900 text-sm">{s.skill}</td>
//                       <td className="py-3.5 px-4 font-mono text-sm text-gray-700 text-right">{s.demand.toLocaleString()}</td>
//                     </tr>
//                   ))}
//                 </tbody>
//               </table>
//             </div>
//           </div>

//           {/* Hiring Institutions Dashboard Panel Card */}
//           <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex flex-col justify-between">
//             <div>
//               <div className="flex items-center justify-between mb-4">
//                 <h3 className="text-sm font-semibold text-gray-700">Top Hiring Employers</h3>
//                 <span 
//                   className="text-[10px] font-bold uppercase tracking-wider px-2 py-0.5 rounded border"
//                   style={{ 
//                     borderColor: CATEGORY_COLORS(selectedCategory), 
//                     color: CATEGORY_COLORS(selectedCategory),
//                     backgroundColor: `${CATEGORY_COLORS(selectedCategory)}08`
//                   }}
//                 >
//                   {selectedCategory === "All Categories" ? "All Sectors" : selectedCategory}
//                 </span>
//               </div>
//               <p className="text-xs text-gray-400 mb-4">
//                 Enterprise institutions and corporations exhibiting the highest volume of ongoing requisition pipelines.
//               </p>
              
//               <div className="space-y-3">
//                 {analytics.employers.map((employer) => (
//                   <div key={employer.name} className="p-3.5 border border-gray-100 rounded-lg bg-gray-50/30 flex items-center justify-between">
//                     <div className="space-y-0.5">
//                       <div className="text-sm font-semibold text-gray-900">{employer.name}</div>
//                       <div className="text-[11px] text-gray-400 font-medium">{employer.location}</div>
//                     </div>
//                     <div className="text-right">
//                       <div className="text-sm font-bold text-gray-800 font-mono">{employer.openRoles}</div>
//                       <div className="text-[10px] uppercase font-bold text-gray-400 tracking-tight">Openings</div>
//                     </div>
//                   </div>
//                 ))}
//               </div>
//             </div>

//             <div className="mt-6 pt-4 border-t border-gray-100 flex items-center justify-between text-xs text-gray-400">
//               <span>Data source: Live Feed Integration</span>
//               <span className="w-2 h-2 rounded-full bg-emerald-500 animate-pulse" />
//             </div>
//           </div>
//         </div>
//       </div>
//     </div>
//   );
// }















"use client";

import { useState, useMemo } from "react";
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

const INDUSTRY_COLORS = (ind: string) => {
  switch (ind) {
    case "Information Technology": return "#3b82f6";
    case "Healthcare & Medicine": return "#ef4444";
    case "Construction & Real Estate": return "#f59e0b";
    case "Finance & Banking": return "#8b5cf6";
    case "Education & Academics": return "#ec4899";
    default: return "#10b981";
  }
};

// 21 Core Macro Industries mapped with localized labels
const INDUSTRIES_21 = [
  { id: "All Industries", en: "All Industries", si: "සියලුම කර්මාන්ත", ta: "அனைத்து தொழில்துறைகள்" },
  { id: "Information Technology", en: "Information Technology", si: "තොරතුරු තාක්ෂණය", ta: "தகவல் தொழில்நுட்பம்" },
  { id: "Healthcare & Medicine", en: "Healthcare & Medicine", si: "සෞඛ්‍ය හා වෛද්‍ය", ta: "சுகாதாரம் & மருத்துவம்" },
  { id: "Construction & Real Estate", en: "Construction & Real Estate", si: "ඉදිකිරීම් හා දේපළ වෙළඳාම්", ta: "கட்டுமானம் & ரியல் எஸ்டேட்" },
  { id: "Finance & Banking", en: "Finance & Banking", si: "මුදල් හා බැංකු", ta: "நிதி & வங்கி" },
  { id: "Education & Academics", en: "Education & Academics", si: "අධ්‍යාපන හා ශාස්ත්‍රීය", ta: "கல்வி & கல்வித்துறை" },
  { id: "Manufacturing & Industrial", en: "Manufacturing & Industrial", si: "නිෂ්පාදන හා කාර්මික", ta: "உற்பத்தி & தொழில்முறை" },
  { id: "Retail & E-Commerce", en: "Retail & E-Commerce", si: "සිල්ලර වෙළඳාම සහ විද්‍යුත් වාණිජ්‍යය", ta: "சில்லறை & மின்-வணிகம்" },
  { id: "Hospitality & Tourism", en: "Hospitality & Tourism", si: "ආගන්තුක සත්කාරය සහ සංචාරක", ta: "உபசரிப்பு & சுற்றுலா" },
  { id: "Transportation & Logistics", en: "Transportation & Logistics", si: "ප්‍රවාහන හා ලොජිස්ටික්ස්", ta: "போக்குவரத்து & லாஜிஸ்டிக்ஸ்" },
  { id: "Agriculture & Farming", en: "Agriculture & Farming", si: "කෘෂිකර්මාන්තය", ta: "விவசாயம்" },
  { id: "Legal Services", en: "Legal Services", si: "නීතිමය සේවා", ta: "சட்ட சேவைகள்" },
  { id: "Marketing & Advertising", en: "Marketing & Advertising", si: "අලෙවිකරණය සහ ප්‍රචාරණය", ta: "சந்தைப்படுத்தல் & விளம்பரம்" },
  { id: "Entertainment & Media", en: "Entertainment & Media", si: "විනෝදාස්වාදය සහ මාධ්‍ය", ta: "பொழுதுபோக்கு & ஊடகம்" },
  { id: "Energy & Utilities", en: "Energy & Utilities", si: "බලශක්ති හා උපයෝගිතා", ta: "எரிசக்தி & பயன்பாடுகள்" },
  { id: "Telecommunications", en: "Telecommunications", si: "විදුලි සංදේශ", ta: "தொலைத்தொடர்பு" },
  { id: "Automotive", en: "Automotive", si: "මෝටර් රථ කර්මාන්තය", ta: "தானியங்கி" },
  { id: "Aerospace & Defense", en: "Aerospace & Defense", si: "අභ්‍යවකාශ හා ආරක්ෂක", ta: "விண்வெளி & பாதுகாப்பு" },
  { id: "Pharmaceuticals", en: "Pharmaceuticals", si: "ඖෂධ නිෂ්පාදනය", ta: "மருந்துகள்" },
  { id: "Non-Profit & NGO", en: "Non-Profit & NGO", si: "ලාභ නොලබන සහ රාජ්‍ය නොවන සංවිධාන", ta: "அறக்கட்டளை & தன்னார்வ தொண்டு நிறுவனம்" },
  { id: "Government & Public Sector", en: "Government & Public Sector", si: "රාජ්‍ය අංශය", ta: "அரசு & பொதுத்துறை" },
  { id: "Human Resources", en: "Human Resources", si: "මානව සම්පත්", ta: "மனித வளம்" }
];

const MOCK_DATA_BY_LANG = {
  en: {
    "All Industries": {
      skills: [
        { skill: "Project Management", demand: 1420 },
        { skill: "Data Analysis", demand: 1250 },
        { skill: "Java / Spring Boot", demand: 1100 },
        { skill: "Cloud Architecture", demand: 980 },
        { skill: "SQL Databases", demand: 890 },
        { skill: "React / Frontend", demand: 820 },
        { skill: "DevOps & Docker", demand: 750 },
        { skill: "Financial Auditing", demand: 680 },
        { skill: "Digital Marketing", demand: 610 },
        { skill: "UI/UX Design", demand: 540 },
        { skill: "Go Language", demand: 490 },
        { skill: "Cybersecurity", demand: 430 },
        { skill: "Machine Learning", demand: 390 },
        { skill: "Agile Scrum", demand: 320 },
        { skill: "Technical Writing", demand: 280 }
      ],
      employers: [
        { name: "WSO2 Lanka", location: "Colombo 03", openRoles: 45 },
        { name: "Dialog Axiata", location: "Colombo 02", openRoles: 38 },
        { name: "Sysco LABS", location: "Colombo 07", openRoles: 32 },
        { name: "John Keells Holdings", location: "Colombo 02", openRoles: 27 },
        { name: "London Stock Exchange Group", location: "Colombo 04", openRoles: 22 }
      ]
    },
    "Information Technology": {
      skills: [
        { skill: "Java / Spring Boot", demand: 950 },
        { skill: "React / Frontend", demand: 840 },
        { skill: "Cloud Architecture", demand: 790 },
        { skill: "Go Language", demand: 680 },
        { skill: "DevOps & Docker", demand: 620 },
        { skill: "Microservices Architecture", demand: 590 },
        { skill: "Python Data Science", demand: 480 },
        { skill: "Kubernetes Orchestration", demand: 410 },
        { skill: "GraphQL APIs", demand: 320 },
        { skill: "TypeScript Ecosystem", demand: 290 },
        { skill: "Node.js Backends", demand: 250 },
        { skill: "CI/CD Automations", demand: 210 },
        { skill: "Next.js Framework", demand: 180 },
        { skill: "Redis Caching Systems", demand: 140 },
        { skill: "Ballerina Integrations", demand: 120 }
      ],
      employers: [
        { name: "WSO2 Lanka", location: "Colombo 03", openRoles: 85 },
        { name: "Sysco LABS", location: "Colombo 07", openRoles: 64 },
        { name: "Virtusa Private Ltd", location: "Colombo 05", openRoles: 59 },
        { name: "99x Technologies", location: "Colombo 03", openRoles: 42 },
        { name: "IFS Sri Lanka", location: "Colombo 04", openRoles: 31 }
      ]
    }
  },
  si: {
    "All Industries": {
      skills: [
        { skill: "ව්‍යාපෘති කළමනාකරණය (Project Management)", demand: 1390 },
        { skill: "දත්ත විශ්ලේෂණය (Data Analysis)", demand: 1180 },
        { skill: "ජාවා මෘදුකාංග සංවර්ධනය (Java)", demand: 1050 },
        { skill: "ගිණුම්කරණය (Accounting)", demand: 920 },
        { skill: "ඩිජිටල් අලෙවිකරණය (Digital Marketing)", demand: 810 },
        { skill: "වෙබ් අඩවි නිර්මාණය (Web Design)", demand: 740 },
        { skill: "පාරිභෝගික සේවා (Customer Service)", demand: 680 },
        { skill: "බාහිර සන්නිවේදනය (Public Relations)", demand: 590 },
        { skill: "මිනිස් බල කළමනාකරණය (HR)", demand: 520 },
        { skill: "ජාලකරණ ඉංජිනේරු විද්‍යාව (Networking)", demand: 460 },
        { skill: "දත්ත ගබඩා පද්ධති (SQL)", demand: 410 },
        { skill: "ග්‍රැෆික් නිර්මාණකරණය (UI/UX)", demand: 380 },
        { skill: "තත්ත්ව පාලනය (Quality Assurance)", demand: 310 },
        { skill: "සැපයුම් දාම කළමනාකරණය", demand: 260 },
        { skill: "තාක්ෂණික ලේඛනකරණය", demand: 190 }
      ],
      employers: [
        { name: "ඩයලොග් ආසියාටා", location: "කොළඹ 02", openRoles: 52 },
        { name: "ජෝන් කීල්ස් සමූහය", location: "කොළඹ 02", openRoles: 44 },
        { name: "සංවර්ධන බැංකුව", location: "කොළඹ 01", openRoles: 29 },
        { name: "කොමර්ෂල් බැංකුව", location: "කොළඹ 03", openRoles: 26 },
        { name: "හේලීස් සමාගම", location: "කොළඹ 10", openRoles: 19 }
      ]
    },
    "Information Technology": {
      skills: [
        { skill: "ජාවා / ස්ප්‍රින්ග් බූට් (Java / Spring)", demand: 920 },
        { skill: "ප්‍රතික්‍රියාශීලී වෙබ් (React.js)", demand: 810 },
        { skill: "වලාකුළු පරිගණකකරණය (Cloud)", demand: 760 },
        { skill: "ගෝ නිරූපණ භාෂාව (Go Lang)", demand: 640 },
        { skill: "ඩොකර් සහ මෙහෙයුම් (DevOps)", demand: 600 },
        { skill: "දත්ත විද්‍යාව (Python Data)", demand: 490 },
        { skill: "මයික්‍රොසර්විස් වාස්තු විද්‍යාව", demand: 440 },
        { skill: "කුබර්නෙටීස් පද්ධති (Kubernetes)", demand: 390 },
        { skill: "ඒපීඅයි කළමනාකරණය (APIs)", demand: 310 },
        { skill: "ටයිප්ස්ක්‍රිප්ට් පරිසරය", demand: 260 },
        { skill: "නෝඩ් මෘදුකාංග (Node.js)", demand: 220 },
        { skill: "ස්වයංක්‍රීයකරණ පද්ධති (CI/CD)", demand: 190 },
        { skill: "නෙක්ස්ට් මෘදුකාංග (Next.js)", demand: 150 },
        { skill: "මතක ගබඩාකරණය (Redis)", demand: 120 },
        { skill: "බැලරිනා ක්‍රමලේඛනය (Ballerina)", demand: 110 }
      ],
      employers: [
        { name: "ඩබ්ලිව්එස්ඕටූ ලංකා (WSO2)", location: "කොළඹ 03", openRoles: 88 },
        { name: "සිස්කෝ ලැබ්ස් (Sysco LABS)", location: "කොළඹ 07", openRoles: 61 },
        { name: "වර්ටූසා පුද්ගලික සමාගම", location: "කොළඹ 05", openRoles: 55 },
        { name: "99එක්ස් ටෙක්නොලොජීස්", location: "කොළඹ 03", openRoles: 39 },
        { name: "අයිඑෆ්එස් ශ්‍රී ලංකා", location: "කොළඹ 04", openRoles: 28 }
      ]
    }
  },
  ta: {
    "All Industries": {
      skills: [
        { skill: "திட்ட மேலாண்மை (Project Management)", demand: 1310 },
        { skill: "தரவு பகுப்பாய்வு (Data Analysis)", demand: 1200 },
        { skill: "ஜாவா மென்பொருள் (Java)", demand: 1020 },
        { skill: "நிதி கணக்கியல் (Accounting)", demand: 940 },
        { skill: "டிஜிட்டல் சந்தைப்படுத்தல்", demand: 830 },
        { skill: "வலைத்தள வடிவமைப்பு (Web)", demand: 710 },
        { skill: "வாடிக்கையாளர் சேவை", demand: 660 },
        { skill: "மனித வள மேலாண்மை (HRM)", demand: 540 },
        { skill: "தரவுத்தள அமைப்புகள் (SQL)", demand: 480 },
        { skill: "வலைப்பின்னல் பொறியியல்", demand: 420 },
        { skill: "வரைகலை வடிவமைப்பு (UI/UX)", demand: 390 },
        { skill: "சைபர் பாதுகாப்பு (Cybersecurity)", demand: 340 },
        { skill: "தரக் கட்டுப்பாடு (QA Testing)", demand: 290 },
        { skill: "விநியோக சங்கிலி மேலாண்மை", demand: 230 },
        { skill: "தொழில்நுட்ப ஆவணமாக்கல்", demand: 170 }
      ],
      employers: [
        { name: "டயலாக் ஆக்சியாட்டா", location: "கொழும்பு 02", openRoles: 49 },
        { name: "ஜான் கீல்ஸ் ஹோல்டிங்ஸ்", location: "கொழும்பு 02", openRoles: 41 },
        { name: "இலங்கை மத்திய வங்கி", location: "கொழும்பு 01", openRoles: 33 },
        { name: "கொமர்ஷல் வங்கி", location: "கொழும்பு 03", openRoles: 28 },
        { name: "ஹெய்லீஸ் நிறுவனம்", location: "கொழும்பு 10", openRoles: 21 }
      ]
    },
    "Information Technology": {
      skills: [
        { skill: "ஜாவா / ஸ்பிரிங் பூட் (Java Boot)", demand: 940 },
        { skill: "ரியாக்ட் மென்பொருள் (React.js)", demand: 830 },
        { skill: "கிளவுட் கட்டிடக்கலை (Cloud Architecture)", demand: 740 },
        { skill: "கோ மொழி நிரலாக்கம் (Go Lang)", demand: 660 },
        { skill: "டெவொப்ஸ் மற்றும் டாக்கர் (DevOps)", demand: 610 },
        { skill: "நுண்ணிய சேவை கட்டிடக்கலை", demand: 530 },
        { skill: "பைதான் தரவு அறிவியல் (Python)", demand: 470 },
        { skill: "குபெர்னெட்டீஸ் (Kubernetes)", demand: 380 },
        { skill: "ஏபிஐ மேலாண்மை (GraphQL)", demand: 320 },
        { skill: "டைப்ஸ்கிரிப்ட் சுற்றுச்சூழல்", demand: 280 },
        { skill: "நோட் மென்பொருள் (Node.js)", demand: 240 },
        { skill: "தானியங்கி அமைப்புகள் (CI/CD)", demand: 190 },
        { skill: "நெக்ஸ்ட் கட்டமைப்பு (Next.js)", demand: 160 },
        { skill: "ரெடிஸ் சேமிப்பகம் (Redis)", demand: 130 },
        { skill: "பலேரினா நிரலாக்கம் (Ballerina)", demand: 105 }
      ],
      employers: [
        { name: "டபிள்யூஎஸ்ஓ2 லங்கா (WSO2)", location: "கொழும்பு 03", openRoles: 82 },
        { name: "சிஸ்கோ லேப்ஸ் (Sysco LABS)", location: "கொழும்பு 07", openRoles: 66 },
        { name: "வெர்டூசா நிறுவனம்", location: "கொழும்பு 05", openRoles: 51 },
        { name: "99எக்ஸ் டெக்னாலஜிஸ்", location: "கொழும்பு 03", openRoles: 37 },
        { name: "ஐஎப்எஸ் ஸ்ரீலங்கா", location: "கொழும்பு 04", openRoles: 29 }
      ]
    }
  }
};

const localization = {
  en: {
    title: "Industry Skill Analytics",
    subtitle: "Tracking talent demand curves and domain competency metrics across active job segments",
    filterLabel: "Filter by Core Industry",
    matrixText: "Sector Matrix: 21 Verticals Available",
    kpiUniqueTitle: "Unique Skills Tracked",
    kpiUniqueDesc: "Distinct competencies indexed inside selected pipeline segment.",
    kpiDemandTitle: "Most Demanding Competency",
    kpiDemandSuffix: "Requisitions Pending",
    chartTitle: "Top 15 Competencies Framework Volume Graph",
    chartDesc: "Aggregated unique live requisition postings tracking highest operational density thresholds.",
    tableTitle: "Complete Industry Skills Breakdown",
    tableDesc: "Select a skill row target below to partition regional hiring employers metrics dynamically.",
    thNo: "No",
    thSkill: "Tracked Competency / Skill",
    thVacancies: "Active Vacancies",
    sidebarTitle: "Top Hiring Employers",
    sidebarDesc: "Enterprise institutions with top volume requirements specifically highlighting this node in active requisition structures.",
    sidebarSuffix: "Pipelines",
    sidebarEmpty: "No matching hiring records.",
    bffChannel: "BFF Live Metrics Channel"
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
    tableDesc: "ආයතනවල බඳවා ගැනීමේ දත්ත සක්‍රීයව බැලීමට පහතින් කුසලතා පේළියක් තෝරන්න.",
    thNo: "අංකය",
    thSkill: "ලුහුබැඳි නිපුණතාවය / කුසලතාව",
    thVacancies: "සක්‍රීය පුරප්පාඩු",
    sidebarTitle: "ප්‍රමුඛ බඳවා ගන්නන්",
    sidebarDesc: "තෝරාගත් කුසලතාවය සඳහා ඉහළම අවශ්‍යතා සහිත ප්‍රමුඛ පෙළේ සමාගම් සහ ආයතන.",
    sidebarSuffix: "නාලිකා",
    sidebarEmpty: "ගැලපෙන බඳවා ගැනීමේ වාර්තා කිසිවක් හමු නොවීය.",
    bffChannel: "BFF සජීවී දත්ත නාලිකාව"
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
    bffChannel: "BFF நேரடி அளவீட்டு சேனல்"
  }
};

export default function IndustryPage() {
  const [selectedIndustry, setSelectedIndustry] = useState("All Industries");
  const [selectedSkillName, setSelectedSkillName] = useState<string | null>(null);
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");

  const d = useMemo(() => localization[currentLang], [currentLang]);

  const analytics = useMemo(() => {
    const languageBucket = MOCK_DATA_BY_LANG[currentLang];
    return languageBucket[selectedIndustry as keyof typeof languageBucket] || languageBucket["Information Technology"];
  }, [selectedIndustry, currentLang]);

  const handleIndustryChange = (ind: string) => {
    setSelectedIndustry(ind);
    setSelectedSkillName(null);
  };

  const uniqueSkillsCount = analytics.skills.length;
  const mostDemandingSkill = useMemo(() => {
    const topNode = analytics.skills.reduce((prev, current) => (prev.demand > current.demand) ? prev : current);
    return { name: topNode.skill, count: topNode.demand };
  }, [analytics.skills]);

  const activeSkillMetrics = useMemo(() => {
    if (selectedSkillName) {
      return analytics.skills.find(s => s.skill === selectedSkillName) || analytics.skills[0];
    }
    return analytics.skills[0];
  }, [analytics.skills, selectedSkillName]);

  const formattedDate = new Date().toLocaleDateString(currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK", {
    year: "numeric",
    month: "short",
    day: "numeric"
  });

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

          <div className="flex bg-gray-100 p-1 rounded-xl shadow-sm">
            {(["en", "si", "ta"] as const).map((l) => (
              <button key={l} onClick={() => { setCurrentLang(l); setSelectedSkillName(null); }} className={`px-3 py-1 text-xs font-bold rounded-lg transition-all ${currentLang === l ? "bg-white text-blue-600 shadow-sm" : "text-gray-500 hover:text-gray-900"}`}>
                {l === "en" ? "English" : l === "si" ? "සිංහල" : "தமிழ்"}
              </button>
            ))}
          </div>
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-blue-600 to-indigo-600 flex items-center justify-center text-white text-xs font-black shadow-md border border-white">BJ</div>
        </div>
      </header>

      <div className="p-8 space-y-6">
        
        {/* Industry Filter Dropdown Control */}
        <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div className="max-w-md w-full">
            <label className="block text-xs text-gray-500 mb-1 font-semibold uppercase tracking-wider">{d.filterLabel}</label>
            <select
              value={selectedIndustry}
              onChange={(e) => handleIndustryChange(e.target.value)}
              className="w-full px-3 py-2.5 border border-gray-200 rounded-xl text-sm bg-white cursor-pointer font-medium outline-none focus:ring-2 focus:ring-zinc-900"
            >
              {INDUSTRIES_21.map((ind) => (
                <option key={ind.id} value={ind.id}>{ind[currentLang]}</option>
              ))}
            </select>
          </div>
          <div className="text-xs text-gray-400 font-medium font-mono bg-gray-50 px-3 py-1.5 rounded-lg border">
            {d.matrixText}
          </div>
        </div>

        {/* Dynamic Metrics KPI Cards Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex items-center justify-between">
            <div>
              <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider">{d.kpiUniqueTitle}</p>
              <h3 className="text-3xl font-black text-zinc-900 mt-1">{uniqueSkillsCount}</h3>
              <p className="text-xs text-gray-500 mt-1">{d.kpiUniqueDesc}</p>
            </div>
            <div className="p-3.5 bg-blue-50 text-blue-600 rounded-xl font-bold text-lg">💡</div>
          </div>

          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm flex items-center justify-between">
            <div>
              <p className="text-[10px] font-black uppercase text-gray-400 tracking-wider">{d.kpiDemandTitle}</p>
              <h3 className="text-xl font-black text-zinc-900 mt-2 truncate max-w-[280px]">{mostDemandingSkill.name}</h3>
              <p className="text-xs text-emerald-600 font-semibold mt-1">
                🔥 {mostDemandingSkill.count.toLocaleString()} {d.kpiDemandSuffix}
              </p>
            </div>
            <div className="p-3.5 bg-emerald-50 text-emerald-600 rounded-xl font-bold text-lg">📈</div>
          </div>
        </div>

        {/* --- VERTICAL BAR CHART: TOP 15 SKILLS VOLUME --- */}
        <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm">
          <div className="mb-4">
            <h3 className="text-xs font-black uppercase text-gray-400 tracking-wider">{d.chartTitle}</h3>
            <p className="text-xs text-gray-400 mt-0.5">{d.chartDesc}</p>
          </div>
          <ResponsiveContainer width="100%" height={340}>
            <BarChart data={analytics.skills} margin={{ top: 10, bottom: 25, left: 10, right: 10 }}>
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
              <Tooltip formatter={(value) => [`${value}`, `${d.thVacancies}`]} />
              <Bar dataKey="demand" radius={[4, 4, 0, 0]} maxBarSize={45}>
                {analytics.skills.map((s, i) => (
                  <Cell key={i} fill={INDUSTRY_COLORS(selectedIndustry)} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Data Grid Master-Detail Interaction Layout */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 items-start">
          
          {/* Complete Skills Table Breakdown (Master) */}
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
                  {analytics.skills.map((s, i) => {
                    const isSelected = activeSkillMetrics?.skill === s.skill;
                    return (
                      <tr 
                        key={s.skill} 
                        onClick={() => setSelectedSkillName(s.skill)}
                        className={`border-b border-gray-50 cursor-pointer transition-all ${isSelected ? "bg-zinc-900 text-white hover:bg-zinc-800" : "hover:bg-gray-50/80"}`}
                      >
                        <td className="py-3.5 px-4 text-xs font-mono text-gray-400">{i + 1}</td>
                        <td className={`py-3.5 px-4 text-sm font-semibold ${isSelected ? "text-white" : "text-gray-900"}`}>{s.skill}</td>
                        <td className={`py-3.5 px-4 font-mono text-sm text-right ${isSelected ? "text-emerald-400 font-bold" : "text-gray-700"}`}>
                          {s.demand.toLocaleString()}
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>

          {/* Reactive Hiring Institutions Sidebar (Detail Panel) */}
          <div className="bg-white p-6 rounded-xl border border-gray-200/60 shadow-sm sticky top-[100px]">
            <div>
              <div className="flex flex-col gap-1 mb-4">
                <h3 className="text-xs font-black uppercase text-gray-400 tracking-wider">{d.sidebarTitle}</h3>
                <div className="mt-1">
                  <span 
                    className="inline-block text-[11px] font-black uppercase tracking-wide px-2.5 py-1 rounded-lg border max-w-full truncate"
                    style={{ 
                      borderColor: INDUSTRY_COLORS(selectedIndustry), 
                      color: INDUSTRY_COLORS(selectedIndustry),
                      backgroundColor: `${INDUSTRY_COLORS(selectedIndustry)}08`
                    }}
                  >
                    🎯 {activeSkillMetrics?.skill || "Default Aggregate"}
                  </span>
                </div>
              </div>
              <p className="text-xs text-gray-400 mb-4 leading-relaxed">{d.sidebarDesc}</p>
              
              <div className="space-y-3 min-h-[280px]">
                {analytics.employers && analytics.employers.length > 0 ? (
                  analytics.employers.map((employer, idx) => {
                    const scalingFactor = activeSkillMetrics ? Math.max(1, Math.round(activeSkillMetrics.demand / (idx + 1.8))) : employer.openRoles;
                    
                    return (
                      <div key={employer.name} className="p-3.5 border border-gray-100 rounded-xl bg-gray-50/50 flex items-center justify-between shadow-xs">
                        <div className="space-y-0.5 pr-2">
                          <div className="text-sm font-bold text-gray-900 tracking-tight truncate max-w-[160px]">{employer.name}</div>
                          <div className="text-[11px] text-gray-400 font-medium">{employer.location}</div>
                        </div>
                        <div className="text-right min-w-[65px]">
                          <div className="text-sm font-black text-blue-600 font-mono">{scalingFactor}</div>
                          <div className="text-[9px] uppercase font-black text-gray-400 tracking-tighter">{d.sidebarSuffix}</div>
                        </div>
                      </div>
                    );
                  })
                ) : (
                  <div className="text-center py-12 text-xs text-gray-400">
                    {d.sidebarEmpty}
                  </div>
                )}
              </div>
            </div>

            <div className="mt-6 pt-4 border-t border-gray-100 flex items-center justify-between text-[11px] text-gray-400 font-medium">
              <span>{d.bffChannel}</span>
              <span className="w-1.5 h-1.5 rounded-full bg-emerald-500 animate-pulse" />
            </div>
          </div>

        </div>
      </div>
    </div>
  );
}