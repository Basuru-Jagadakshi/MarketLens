"use client";

import { useState } from "react";
import {
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  Legend,
  BarChart,
  Bar,
  LineChart,
  Line,
} from "recharts";

// --- DYNAMIC YEAR GENERATOR (Future-Proofing) ---
const currentYearNum = new Date().getFullYear();
const DYNAMIC_YEARS = [
  String(currentYearNum),
  String(currentYearNum - 1),
  String(currentYearNum - 2)
];

// --- PURE JAVASCRIPT DICTIONARY MATRIX ---
const TRANSLATIONS = {
  en: {
    title: "Labour Market Demand Dashboard",
    subtitle: "National Strategic Overview Driven by SLSO & SLSIC Registries",
    vacancies: "Current Vacancies",
    mom: "+4.8% Month over Month",
    occFramework: "Occupations Framework",
    indFramework: "Industries Framework",
    slsoBased: "Based on SLSO",
    slsicBased: "Based on SLSIC",
    collapse: "Click to collapse",
    expand: "Click to view full breakdown",
    matrixTitle: "Registered Framework Classifications Matrix",
    open: "open",
    occChartTitle: "Current Job Distribution by Occupation (SLSO)",
    occChartSub: "Horizontal mapping showing all 10 standard occupation bands",
    indChartTitle: "Current Job Distribution by Industry (SLSIC)",
    indChartSub: "Vertical bar chart projection featuring rotated X-axis headers for all 21 divisions",
    analyticsBtn: "See Analytics",
    expChartTitle: "Current Job Distribution by Experience",
    eduChartTitle: "Current Job Distribution by Education Level",
    remoteTitle: "Remote / On-Site Configuration",
    contractTitle: "Contract Type Share",
    occModalTitle: "Occupation Analytics Breakdown",
    occModalSub: "Macro trend patterns & top jobs mapped by tier",
    indModalTitle: "Industry Segment Workspace Matrix",
    indModalSub: "Past year performance indices categorized by specific market sectors",
    selectOccLabel: "Select Occupation",
    trendChartHeader: "Variation Over Past Years (Historical Demand Trend)",
    demandingJobsHeader: "Demanding Jobs for",
    sectorTrendHeader: "Sector Variant Level Across Years",
    expAllocHeader: "Experience Allocation Distribution",
    regionalShareHeader: "Regional Province Share Allocation",
    eduThresholdHeader: "Minimum Educational Level Threshold",
    topEnterpriseHeader: "Top Sector Enterprise Hiring Groups",
    positionsCount: "open positions",
    closePanel: "Close",
    // --- NEW KEYS ---
    sectorChartTitle: "Current Job Distribution by Employment Sector",
    sectorChartSub: "Government, Semi-Government, Private and NGO share of current vacancies",
    formalChartTitle: "Current Job Distribution by Formal / Informal Sector",
    genderChartTitle: "Current Job Distribution by Gender",
    vocationalChartTitle: "Current Job Distribution by Vocational Education (NVQ Level)",
    sectorModalTitle: "Employment Sector Analytics Breakdown",
    sectorModalSub: "Yearly vacancy trend for each employment sector",
    selectSectorLabel: "Select Sector",
    sectorPanelTrendHeader: "Yearly Trend for Selected Sector",
    formalInformalHeader: "Formal / Informal Job Count",
    genderHeader: "Gender Wise Job Count",
    vocationalIndHeader: "Vocational Education Wise Job Count (NVQ)",
    yearFilterHeader: "Filter by Year",
  },
  si: {
    title: "ශ්‍රම වෙළඳපල ඉල්ලුම උපකරණ පුවරුව",
    subtitle: "SLSO සහ SLSIC ලේඛන මගින් මෙහෙයවන ජාතික උපායමාර්ගික දළ විශ්ලේෂණය",
    vacancies: "වත්මන් පුරප්පාඩු",
    mom: "පෙර මාසයට සාපේක්ෂව +4.8%",
    occFramework: "වෘත්තීය රාමුව",
    indFramework: "කර්මාන්ත රාමුව",
    slsoBased: "SLSO මත පදනම්ව",
    slsicBased: "SLSIC මත පදනම්ව",
    collapse: "හකුලන්න ක්ලික් කරන්න",
    expand: "සම්පූර්ණ විස්තරය බැලීමට ක්ලික් කරන්න",
    matrixTitle: "ලියාපදිංචි රාමු වර්ගීකරණ අනුකෘතිය",
    open: "විවෘතයි",
    occChartTitle: "වෘත්තිය අනුව වත්මන් රැකියා ව්‍යාප්තිය (SLSO)",
    occChartSub: "ප්‍රමිතිගත වෘත්තීය කාණ්ඩ 10ම පෙන්වන තිරස් සිතියම්කරණය",
    indChartTitle: "කර්මාන්තය අනුව වත්මන් රැකියා ව්‍යාප්තිය (SLSIC)",
    indChartSub: "අංශ 21 සඳහාම භ්‍රමණය වූ X-අක්ෂ ශීර්ෂයන් සහිත සිරස් තීරු ප්‍රස්තාරය",
    analyticsBtn: "විශ්ලේෂණ බලන්න",
    expChartTitle: "අත්දැකීම් අනුව වත්මන් රැකියා ව්‍යාප්තිය",
    eduChartTitle: "අධ්‍යාපන මට්ටම අනුව වත්මන් රැකියා ව්‍යාප්තිය",
    remoteTitle: "දුරස්ථ / සේවා ස්ථානගත වින්‍යාසය",
    contractTitle: "කොන්ත්‍රාත්තු වර්ගයේ කොටස",
    occModalTitle: "වෘත්තීය විශ්ලේෂණ බිඳවැටීම",
    occModalSub: "ස්ථර අනුව සිතියම්ගත කරන ලද සාර්ව ප්‍රවණතා රටා සහ ඉහළම රැකියා",
    indModalTitle: "කර්මාන්ත අංශ වැඩ අවකාශ අනුකෘතිය",
    indModalSub: "නිශ්චිත වෙළඳපල අංශ අනුව වර්ගීකරණය කරන ලද පසුගිය වසරේ කාර්ය සාධන දර්ශක",
    selectOccLabel: "වෘත්තිය තෝරන්න",
    trendChartHeader: "පසුගිය වසරවල විචලනය (ඓතිහාසික ඉල්ලුමේ ප්‍රවණතාවය)",
    demandingJobsHeader: "සඳහා ඉල්ලුමක් ඇති රැකියා",
    sectorTrendHeader: "වසර පුරා අංශ විචල්‍ය මට්ටම",
    expAllocHeader: "අත්දැකීම් වෙන් කිරීමේ ව්‍යාප්තිය",
    regionalShareHeader: "ප්‍රාදේශීය පළාත් කොටස් වෙන් කිරීම",
    eduThresholdHeader: "අවම අධ්‍යාපන මට්ටමේ සීමාව",
    topEnterpriseHeader: "ඉහළම අංශයේ ව්‍යවසාය බඳවා ගැනීමේ කණ්ඩායම්",
    positionsCount: "ඇබෑර්තු සංඛ්‍යාව",
    closePanel: "වසන්න",
    // --- NEW KEYS ---
    sectorChartTitle: "රැකියා අංශය අනුව වත්මන් රැකියා ව්‍යාප්තිය",
    sectorChartSub: "රජය, අර්ධ රාජ්‍ය, පුද්ගලික සහ රාජ්‍ය නොවන සංවිධාන අංශවල වත්මන් පුරප්පාඩු කොටස",
    formalChartTitle: "විධිමත් / අවිධිමත් අංශය අනුව වත්මන් රැකියා ව්‍යාප්තිය",
    genderChartTitle: "ස්ත්‍රී පුරුෂ භාවය අනුව වත්මන් රැකියා ව්‍යාප්තිය",
    vocationalChartTitle: "වෘත්තීය අධ්‍යාපනය අනුව වත්මන් රැකියා ව්‍යාප්තිය (NVQ මට්ටම)",
    sectorModalTitle: "රැකියා අංශ විශ්ලේෂණ බිඳවැටීම",
    sectorModalSub: "එක් එක් රැකියා අංශය සඳහා වාර්ෂික පුරප්පාඩු ප්‍රවණතාව",
    selectSectorLabel: "අංශය තෝරන්න",
    sectorPanelTrendHeader: "තෝරාගත් අංශය සඳහා වාර්ෂික ප්‍රවණතාව",
    formalInformalHeader: "විධිමත් / අවිධිමත් රැකියා ගණන",
    genderHeader: "ස්ත්‍රී පුරුෂ භාවය අනුව රැකියා ගණන",
    vocationalIndHeader: "වෘත්තීය අධ්‍යාපනය අනුව රැකියා ගණන (NVQ)",
    yearFilterHeader: "වර්ෂය අනුව පෙරහන් කරන්න",
  },
  ta: {
    title: "தொழில் சந்தை தேவை தகவல் பலகை",
    subtitle: "SLSO & SLSIC பதிவேடுகளால் இயக்கப்படும் தேசிய மூலோபாய கண்ணோட்டம்",
    vacancies: "தற்போதைய காலியிடங்கள்",
    mom: "கடந்த மாதத்தை விட +4.8%",
    occFramework: "தொழில் கட்டமைப்பு",
    indFramework: "தொழில்துறை கட்டமைப்பு",
    slsoBased: "SLSO இன் அடிப்படையில்",
    slsicBased: "SLSIC இன் அடிப்படையில்",
    collapse: "சுருக்க கிளிக் செய்யவும்",
    expand: "முழு விபரங்களையும் பார்க்க கிளிக் செய்யவும்",
    matrixTitle: "பதிவுசெய்யப்பட்ட கட்டமைப்பு வகைப்பாடு அணி",
    open: "காலியிடம்",
    occChartTitle: "தொழில் வாரியான தற்போதைய வேலை விநியோகம் (SLSO)",
    occChartSub: "அனைத்து 10 நிலையான தொழில் குழுக்களையும் காட்டும் கிடைமட்ட வரைபடம்",
    indChartTitle: "தொழில்துறை வாரியான தற்போதைய வேலை விநியோகம் (SLSIC)",
    indChartSub: "அனைத்து 21 பிரிவுகளுக்கான சுழற்றப்பட்ட X-அச்சு தலைப்புகளைக் கொண்ட செங்குத்து பட்டை வரைபடம்",
    analyticsBtn: "பகுப்பாய்வைக் காண்க",
    expChartTitle: "அனுபவ வாரியான தற்போதைய வேலை விநியோகம்",
    eduChartTitle: "கல்வித் தகுதி வாரியான தற்போதைய வேலை விநியோகம்",
    remoteTitle: "தொலைதூர / தள வேலை கட்டமைப்பு",
    contractTitle: "ஒப்பந்த வகை பங்கீடு",
    occModalTitle: "தொழில் பகுப்பாய்வு முறிவு",
    occModalSub: "மேக்ரோ போக்கு வடிவங்கள் மற்றும் அடுக்கு வாரியாக வரைபடமாக்கப்பட்ட சிறந்த வேலைகள்",
    indModalTitle: "தொழில்துறை பிரிவு பணிவெளி அணி",
    indModalSub: "குறிப்பிட்ட சந்தைத் துறைகளால் வகைப்படுத்தப்பட்ட கடந்த ஆண்டு செயல்திறன் குறியீடுகள்",
    selectOccLabel: "தொழிலைத் தேர்ந்தெடுக்கவும்",
    trendChartHeader: "கடந்த ஆண்டுகளின் மாறுபாடு (வரலாற்று தேவை போக்கு)",
    demandingJobsHeader: "விருப்பமுள்ள வேலைகள்",
    sectorTrendHeader: "ஆண்டுகள் முழுவதும் துறை மாறுபாட்டின் அளவு",
    expAllocHeader: "அனுபவ ஒதுக்கீடு விநியோகம்",
    regionalShareHeader: "பிராந்திய மாகாணப் பங்கு ஒதுக்கீடு",
    eduThresholdHeader: "கையெழுத்து கல்வித் தகுதி வரம்பு",
    topEnterpriseHeader: "முன்னணி துறை நிறுவன வேலைவாய்ப்பு குழுக்கள்",
    positionsCount: "காலியிடங்கள்",
    closePanel: "மூடு",
    // --- NEW KEYS ---
    sectorChartTitle: "வேலைவாய்ப்பு துறை வாரியான தற்போதைய வேலை விநியோகம்",
    sectorChartSub: "அரசு, அரை அரசு, தனியார் மற்றும் அரசு சாரா நிறுவன துறைகளின் தற்போதைய காலியிடப் பங்கு",
    formalChartTitle: "முறைசார் / முறைசாரா துறை வாரியான தற்போதைய வேலை விநியோகம்",
    genderChartTitle: "பாலினம் வாரியான தற்போதைய வேலை விநியோகம்",
    vocationalChartTitle: "தொழில் சார் கல்வி வாரியான தற்போதைய வேலை விநியோகம் (NVQ நிலை)",
    sectorModalTitle: "வேலைவாய்ப்பு துறை பகுப்பாய்வு முறிவு",
    sectorModalSub: "ஒவ்வொரு வேலைவாய்ப்பு துறைக்கான ஆண்டு காலியிடப் போக்கு",
    selectSectorLabel: "துறையைத் தேர்ந்தெடுக்கவும்",
    sectorPanelTrendHeader: "தேர்ந்தெடுக்கப்பட்ட துறைக்கான ஆண்டு போக்கு",
    formalInformalHeader: "முறைசார் / முறைசாரா வேலை எண்ணிக்கை",
    genderHeader: "பாலினம் வாரியான வேலை எண்ணிக்கை",
    vocationalIndHeader: "தொழில் சார் கல்வி வாரியான வேலை எண்ணிக்கை (NVQ)",
    yearFilterHeader: "ஆண்டு வாரியாக வடிகட்டவும்",
  }
};

const CHART_COLORS = ["#2563eb", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#6366f1", "#ec4899"];

// --- 10 SLSO OCCUPATIONS ---
const MOCK_SLSO_FULL = [
  { id: 1, name: "1. Managers", count: 1240, demanding: ["Operations Director", "Branch Manager", "HR Generalist"] },
  { id: 2, name: "2. Professionals", count: 3450, demanding: ["Software Engineer", "Accountant", "Data Analyst"] },
  { id: 3, name: "3. Technicians & Associate Professionals", count: 2100, demanding: ["Network Admin", "Lab Assistant", "Safety Inspector"] },
  { id: 4, name: "4. Clerical Support Workers", count: 950, demanding: ["Data Entry Specialist", "Receptionist", "Executive Assistant"] },
  { id: 5, name: "5. Service and Sales Workers", count: 1850, demanding: ["Sales Representative", "Customer Care Agent", "Cashier"] },
  { id: 6, name: "6. Skilled Agricultural, Forestry & Fishery Workers", count: 420, demanding: ["Farm Supervisor", "Agri-Consultant", "Fisheries Officer"] },
  { id: 7, name: "7. Craft and Related Trades Workers", count: 1100, demanding: ["Electrician", "Machinist", "Welding Supervisor"] },
  { id: 8, name: "8. Plant & Machine Operators & Assemblers", count: 680, demanding: ["CNC Operator", "Delivery Driver", "Assembly Line Lead"] },
  { id: 9, name: "9. Elementary Occupations", count: 450, demanding: ["Office Helper", "Store Keeper", "Security Associate"] },
  { id: 10, name: "10. Armed Forces Occupations", count: 210, demanding: ["Security Analyst", "Logistics Officer", "Field Guard"] }
];

// --- 21 SLSIC INDUSTRIES ---
const MOCK_SLSIC_FULL = [
  { id: 1, name: "A - Agriculture", count: 420 }, { id: 2, name: "B - Mining & Quarrying", count: 110 },
  { id: 3, name: "C - Manufacturing", count: 2450 }, { id: 4, name: "D - Electricity & Gas", count: 190 },
  { id: 5, name: "E - Water & Waste Mgmt", count: 130 }, { id: 6, name: "F - Construction", count: 680 },
  { id: 7, name: "G - Wholesale & Retail", count: 1150 }, { id: 8, name: "H - Transport & Storage", count: 540 },
  { id: 9, name: "I - Accommodation & Food", count: 980 }, { id: 10, name: "J - Info & Tech Comms", count: 3100 },
  { id: 11, name: "K - Finance & Insurance", count: 1420 }, { id: 12, name: "L - Real Estate", count: 180 },
  { id: 13, name: "M - Professional & Sci", count: 490 }, { id: 14, name: "N - Admin & Support", count: 310 },
  { id: 15, name: "O - Public Admin", count: 220 }, { id: 16, name: "P - Education", count: 380 },
  { id: 17, name: "Q - Health & Social Work", count: 290 }, { id: 18, name: "R - Arts, Ent & Rec", count: 150 },
  { id: 19, name: "S - Other Services", count: 120 }, { id: 20, name: "T - Private Households", count: 80 },
  { id: 21, name: "U - Extraterritorial Org", count: 50 }
];

// --- NEW: 4 EMPLOYMENT SECTORS ---
const MOCK_SECTOR_FULL = [
  { id: 1, name: "Government", count: 2450 },
  { id: 2, name: "Semi Government", count: 1180 },
  { id: 3, name: "Private", count: 8120 },
  { id: 4, name: "NGO", count: 700 },
];

// --- NEW: FORMAL / INFORMAL DISTRIBUTION ---
const MOCK_FORMAL_DIST = [
  { name: "Formal", value: 9800 },
  { name: "Informal", value: 2650 },
];

// --- NEW: GENDER DISTRIBUTION ---
const MOCK_GENDER_DIST = [
  { name: "Male", value: 7100 },
  { name: "Female", value: 4950 },
  { name: "Not Specified", value: 400 },
];

// --- NEW: VOCATIONAL EDUCATION (NVQ) DISTRIBUTION ---
const MOCK_VOCATIONAL_DIST = [
  { name: "NVQ 1", value: 420 },
  { name: "NVQ 2", value: 610 },
  { name: "NVQ 3", value: 890 },
  { name: "NVQ 4", value: 1120 },
  { name: "NVQ 5", value: 760 },
  { name: "NVQ 6", value: 310 },
  { name: "NVQ 7", value: 140 },
];

const MOCK_OCCUPATION_HISTORICAL_TRENDS: Record<string, { year: string; vacancies: number }[]> = {
  "1. Managers": [{ year: DYNAMIC_YEARS[2], vacancies: 950 }, { year: DYNAMIC_YEARS[1], vacancies: 1050 }, { year: DYNAMIC_YEARS[0], vacancies: 1240 }],
  "2. Professionals": [{ year: DYNAMIC_YEARS[2], vacancies: 1900 }, { year: DYNAMIC_YEARS[1], vacancies: 2400 }, { year: DYNAMIC_YEARS[0], vacancies: 3450 }],
  default: [{ year: DYNAMIC_YEARS[2], vacancies: 400 }, { year: DYNAMIC_YEARS[1], vacancies: 410 }, { year: DYNAMIC_YEARS[0], vacancies: 420 }]
};

// --- NEW: FORMAL/INFORMAL + GENDER, PER OCCUPATION, PER YEAR ---
const MOCK_OCC_YEARLY_DETAIL: Record<string, Record<string, { formal: { name: string; value: number }[]; gender: { name: string; value: number }[] }>> = {
  "1. Managers": {
    [DYNAMIC_YEARS[0]]: {
      formal: [{ name: "Formal", value: 1080 }, { name: "Informal", value: 160 }],
      gender: [{ name: "Male", value: 780 }, { name: "Female", value: 440 }, { name: "Not Specified", value: 20 }],
    },
    [DYNAMIC_YEARS[1]]: {
      formal: [{ name: "Formal", value: 900 }, { name: "Informal", value: 150 }],
      gender: [{ name: "Male", value: 650 }, { name: "Female", value: 390 }, { name: "Not Specified", value: 10 }],
    },
    [DYNAMIC_YEARS[2]]: {
      formal: [{ name: "Formal", value: 810 }, { name: "Informal", value: 140 }],
      gender: [{ name: "Male", value: 590 }, { name: "Female", value: 350 }, { name: "Not Specified", value: 10 }],
    },
  },
  "2. Professionals": {
    [DYNAMIC_YEARS[0]]: {
      formal: [{ name: "Formal", value: 3100 }, { name: "Informal", value: 350 }],
      gender: [{ name: "Male", value: 1950 }, { name: "Female", value: 1450 }, { name: "Not Specified", value: 50 }],
    },
    [DYNAMIC_YEARS[1]]: {
      formal: [{ name: "Formal", value: 2150 }, { name: "Informal", value: 250 }],
      gender: [{ name: "Male", value: 1380 }, { name: "Female", value: 990 }, { name: "Not Specified", value: 30 }],
    },
    [DYNAMIC_YEARS[2]]: {
      formal: [{ name: "Formal", value: 1720 }, { name: "Informal", value: 180 }],
      gender: [{ name: "Male", value: 1080 }, { name: "Female", value: 790 }, { name: "Not Specified", value: 30 }],
    },
  },
  default: {
    [DYNAMIC_YEARS[0]]: {
      formal: [{ name: "Formal", value: 340 }, { name: "Informal", value: 80 }],
      gender: [{ name: "Male", value: 260 }, { name: "Female", value: 150 }, { name: "Not Specified", value: 10 }],
    },
    [DYNAMIC_YEARS[1]]: {
      formal: [{ name: "Formal", value: 330 }, { name: "Informal", value: 80 }],
      gender: [{ name: "Male", value: 250 }, { name: "Female", value: 150 }, { name: "Not Specified", value: 10 }],
    },
    [DYNAMIC_YEARS[2]]: {
      formal: [{ name: "Formal", value: 320 }, { name: "Informal", value: 80 }],
      gender: [{ name: "Male", value: 245 }, { name: "Female", value: 145 }, { name: "Not Specified", value: 10 }],
    },
  },
};

// --- NEW: TRENDS PER EMPLOYMENT SECTOR ---
const MOCK_SECTOR_TRENDS: Record<string, { year: string; vacancies: number }[]> = {
  "Government": [{ year: DYNAMIC_YEARS[2], vacancies: 2100 }, { year: DYNAMIC_YEARS[1], vacancies: 2280 }, { year: DYNAMIC_YEARS[0], vacancies: 2450 }],
  "Semi Government": [{ year: DYNAMIC_YEARS[2], vacancies: 980 }, { year: DYNAMIC_YEARS[1], vacancies: 1060 }, { year: DYNAMIC_YEARS[0], vacancies: 1180 }],
  "Private": [{ year: DYNAMIC_YEARS[2], vacancies: 6400 }, { year: DYNAMIC_YEARS[1], vacancies: 7250 }, { year: DYNAMIC_YEARS[0], vacancies: 8120 }],
  "NGO": [{ year: DYNAMIC_YEARS[2], vacancies: 540 }, { year: DYNAMIC_YEARS[1], vacancies: 610 }, { year: DYNAMIC_YEARS[0], vacancies: 700 }],
  default: [{ year: DYNAMIC_YEARS[2], vacancies: 500 }, { year: DYNAMIC_YEARS[1], vacancies: 540 }, { year: DYNAMIC_YEARS[0], vacancies: 580 }],
};

const MOCK_INDUSTRY_ANALYTICS: Record<string, any> = {
  default: {
    trends: [{ year: DYNAMIC_YEARS[2], vacancies: 1800 }, { year: DYNAMIC_YEARS[1], vacancies: 2200 }, { year: DYNAMIC_YEARS[0], vacancies: 2450 }],
    experience: {
      [DYNAMIC_YEARS[0]]: [{ label: "Entry Level", value: 400 }, { label: "Junior", value: 1200 }, { label: "Mid-Level", value: 650 }, { label: "Senior", value: 200 }],
      [DYNAMIC_YEARS[1]]: [{ label: "Entry Level", value: 350 }, { label: "Junior", value: 1100 }, { label: "Mid-Level", value: 550 }, { label: "Senior", value: 200 }],
      [DYNAMIC_YEARS[2]]: [{ label: "Entry Level", value: 300 }, { label: "Junior", value: 900 }, { label: "Mid-Level", value: 450 }, { label: "Senior", value: 150 }]
    },
    employers: {
      [DYNAMIC_YEARS[0]]: [{ name: "Enterprise Corp A", open_job_count: 140 }, { name: "Global Solutions", open_job_count: 95 }, { name: "National Industries", open_job_count: 80 }],
      [DYNAMIC_YEARS[1]]: [{ name: "Enterprise Corp A", open_job_count: 110 }, { name: "Global Solutions", open_job_count: 80 }, { name: "National Industries", open_job_count: 75 }],
      [DYNAMIC_YEARS[2]]: [{ name: "Enterprise Corp A", open_job_count: 90 }, { name: "Global Solutions", open_job_count: 70 }, { name: "National Industries", open_job_count: 60 }]
    },
    provinces: {
      [DYNAMIC_YEARS[0]]: [{ name: "Western", value: 1500 }, { name: "Central", value: 500 }, { name: "Southern", value: 450 }],
      [DYNAMIC_YEARS[1]]: [{ name: "Western", value: 1300 }, { name: "Central", value: 480 }, { name: "Southern", value: 420 }],
      [DYNAMIC_YEARS[2]]: [{ name: "Western", value: 1100 }, { name: "Central", value: 400 }, { name: "Southern", value: 300 }]
    },
    education: {
      [DYNAMIC_YEARS[0]]: [{ label: "Degree", value: 1400 }, { label: "A/L", value: 800 }, { label: "O/L", value: 250 }],
      [DYNAMIC_YEARS[1]]: [{ label: "Degree", value: 1200 }, { label: "A/L", value: 750 }, { label: "O/L", value: 250 }],
      [DYNAMIC_YEARS[2]]: [{ label: "Degree", value: 1000 }, { label: "A/L", value: 600 }, { label: "O/L", value: 200 }]
    },
    // --- NEW: VOCATIONAL EDUCATION (NVQ) PER YEAR, FOR SELECTED INDUSTRY ---
    vocational: {
      [DYNAMIC_YEARS[0]]: [
        { label: "NVQ 1", value: 90 }, { label: "NVQ 2", value: 150 }, { label: "NVQ 3", value: 240 },
        { label: "NVQ 4", value: 330 }, { label: "NVQ 5", value: 190 }, { label: "NVQ 6", value: 95 }, { label: "NVQ 7", value: 45 }
      ],
      [DYNAMIC_YEARS[1]]: [
        { label: "NVQ 1", value: 75 }, { label: "NVQ 2", value: 120 }, { label: "NVQ 3", value: 200 },
        { label: "NVQ 4", value: 270 }, { label: "NVQ 5", value: 150 }, { label: "NVQ 6", value: 70 }, { label: "NVQ 7", value: 30 }
      ],
      [DYNAMIC_YEARS[2]]: [
        { label: "NVQ 1", value: 60 }, { label: "NVQ 2", value: 95 }, { label: "NVQ 3", value: 160 },
        { label: "NVQ 4", value: 210 }, { label: "NVQ 5", value: 115 }, { label: "NVQ 6", value: 50 }, { label: "NVQ 7", value: 20 }
      ]
    }
  }
};

const MOCK_EXPERIENCE_DIST = [{ name: "Entry Level", value: 2500 }, { name: "Junior", value: 4800 }, { name: "Mid-Level", value: 3100 }, { name: "Senior", value: 1450 }];
const MOCK_EDUCATION_DIST = [{ name: "Degree", value: 5200 }, { name: "A/L", value: 3100 }, { name: "O/L", value: 1800 }, { name: "Below O/L", value: 650 }, { name: "Not Specified", value: 1700 }];

export default function DashboardPage() {
  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");
  const [activeKpiRow, setActiveKpiRow] = useState<"SLSO" | "SLSIC" | null>(null);
  const [activePanel, setActivePanel] = useState<"OCCUPATION" | "INDUSTRY" | "SECTOR" | null>(null);

  const [selectedOccId, setSelectedOccId] = useState<number>(MOCK_SLSO_FULL[0].id);
  const [selectedOccName, setSelectedOccName] = useState<string>(MOCK_SLSO_FULL[0].name);
  const [selectedIndId, setSelectedIndId] = useState<number>(MOCK_SLSIC_FULL[0].id);
  const [selectedIndName, setSelectedIndName] = useState<string>(MOCK_SLSIC_FULL[0].name);
  const [analyticsYear, setAnalyticsYear] = useState<number>(Number(DYNAMIC_YEARS[0]));

  // --- NEW STATE: employment sector panel selection ---
  const [selectedSectorId, setSelectedSectorId] = useState<number>(MOCK_SECTOR_FULL[0].id);
  const [selectedSectorName, setSelectedSectorName] = useState<string>(MOCK_SECTOR_FULL[0].name);

  // --- NEW STATE: year filter dedicated to the occupation panel ---
  const [occAnalyticsYear, setOccAnalyticsYear] = useState<number>(Number(DYNAMIC_YEARS[0]));

  const d = TRANSLATIONS[currentLang];
  const formattedDate = new Date().toLocaleDateString(
    currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
    { year: "numeric", month: "short", day: "numeric" }
  );

  const activeOccData = MOCK_SLSO_FULL.find(o => o.name === selectedOccName);
  const occTrendData = MOCK_OCCUPATION_HISTORICAL_TRENDS[selectedOccName] || MOCK_OCCUPATION_HISTORICAL_TRENDS.default;
  const indData = MOCK_INDUSTRY_ANALYTICS.default;

  // --- NEW DERIVED DATA ---
  const occYearDetail =
    (MOCK_OCC_YEARLY_DETAIL[selectedOccName] && MOCK_OCC_YEARLY_DETAIL[selectedOccName][String(occAnalyticsYear)]) ||
    MOCK_OCC_YEARLY_DETAIL.default[String(occAnalyticsYear)];
  const sectorTrendData = MOCK_SECTOR_TRENDS[selectedSectorName] || MOCK_SECTOR_TRENDS.default;
  const vocationalIndData = indData.vocational[analyticsYear];

  const handleOpenOccPanel = () => {
    setActivePanel("OCCUPATION");
  };

  const handleOpenIndPanel = () => {
    setActivePanel("INDUSTRY");
  };

  // --- NEW HANDLER ---
  const handleOpenSectorPanel = () => {
    setActivePanel("SECTOR");
  };

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800 font-sans">

      {/* HEADER */}
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

      {/* ── ROOT SPLIT LAYOUT ─────────────────────────────────────────────── */}
      <div className="flex h-[calc(100vh-73px)]">

        {/* ── LEFT: MAIN SCROLLABLE CONTENT ─────────────────────────────── */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-8 space-y-8 max-w-[1200px] mx-auto pb-20">

            {/* TIER 1: KPI CARDS */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              <div className="bg-white p-6 rounded-xl shadow-sm flex flex-col justify-between">
                <div>
                  <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">{d.vacancies}</p>
                  <h3 className="text-3xl font-black text-gray-900 mt-2">12,450</h3>
                </div>
                <div className="mt-4 flex items-center text-emerald-600 text-xs font-bold"><span>{d.mom}</span></div>
              </div>

              <div onClick={() => setActiveKpiRow(activeKpiRow === "SLSO" ? null : "SLSO")} className={`bg-white p-6 rounded-xl transition-all cursor-pointer shadow-sm flex flex-col justify-between ${activeKpiRow === "SLSO" ? "ring-2 ring-blue-500" : ""}`}>
                <div>
                  <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">{d.occFramework}</p>
                  <h3 className="text-3xl font-black text-gray-900 mt-2">10</h3>
                  <p className="text-[11px] font-semibold text-gray-400 mt-1">{d.slsoBased}</p>
                </div>
                <p className="mt-4 text-[11px] text-blue-600 font-medium">{activeKpiRow === "SLSO" ? d.collapse : d.expand}</p>
              </div>

              <div onClick={() => setActiveKpiRow(activeKpiRow === "SLSIC" ? null : "SLSIC")} className={`bg-white p-6 rounded-xl transition-all cursor-pointer shadow-sm flex flex-col justify-between ${activeKpiRow === "SLSIC" ? "ring-2 ring-emerald-500" : ""}`}>
                <div>
                  <p className="text-xs text-gray-400 font-bold uppercase tracking-wider">{d.indFramework}</p>
                  <h3 className="text-3xl font-black text-gray-900 mt-2">21</h3>
                  <p className="text-[11px] font-semibold text-gray-400 mt-1">{d.slsicBased}</p>
                </div>
                <p className="mt-4 text-[11px] text-emerald-600 font-medium">{activeKpiRow === "SLSIC" ? d.collapse : d.expand}</p>
              </div>
            </div>

            {/* EXPANDABLE MATRIX */}
            {activeKpiRow && (
              <div className="bg-white p-6 rounded-xl shadow-inner">
                <h4 className="text-xs font-bold text-gray-500 uppercase tracking-wider mb-4 font-mono">{d.matrixTitle} ({activeKpiRow})</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                  {(activeKpiRow === "SLSO" ? MOCK_SLSO_FULL : MOCK_SLSIC_FULL).map((item, idx) => (
                    <div key={idx} className="flex justify-between items-center p-3 rounded-lg bg-gray-50 text-xs">
                      <span className="font-semibold text-gray-700 truncate mr-2">{item.name}</span>
                      <span className="font-mono font-bold bg-white px-2.5 py-1 rounded text-gray-600">{item.count} {d.open}</span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* OCCUPATION CHART */}
            <div className="bg-white p-5 rounded-xl shadow-sm h-[520px] flex flex-col">
              <div className="flex justify-between items-center border-b border-gray-50 pb-2">
                <div>
                  <h4 className="text-sm font-bold text-gray-800">{d.occChartTitle}</h4>
                  <p className="text-[11px] text-gray-400">{d.occChartSub}</p>
                </div>
                <button
                  onClick={handleOpenOccPanel}
                  className={`text-xs font-bold px-3 py-1.5 rounded-lg transition-all ${activePanel === "OCCUPATION" ? "bg-blue-600 text-white" : "text-blue-600 bg-blue-50 hover:bg-blue-100"}`}
                >
                  {d.analyticsBtn}
                </button>
              </div>
              <div className="flex-1 mt-4">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={MOCK_SLSO_FULL} layout="vertical" margin={{ left: 20, right: 20 }}>
                    <XAxis type="number" tick={{ fontSize: 10 }} />
                    <YAxis dataKey="name" type="category" tick={{ fontSize: 10 }} width={160} />
                    <Tooltip />
                    <Bar dataKey="count" fill="#2563eb" radius={[0, 4, 4, 0]} barSize={16} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            {/* INDUSTRY CHART */}
            <div className="bg-white p-5 rounded-xl shadow-sm h-[620px] flex flex-col">
              <div className="flex justify-between items-center border-b border-gray-50 pb-2">
                <div>
                  <h4 className="text-sm font-bold text-gray-800">{d.indChartTitle}</h4>
                  <p className="text-[11px] text-gray-400">{d.indChartSub}</p>
                </div>
                <button
                  onClick={handleOpenIndPanel}
                  className={`text-xs font-bold px-3 py-1.5 rounded-lg transition-all ${activePanel === "INDUSTRY" ? "bg-emerald-600 text-white" : "text-emerald-600 bg-emerald-50 hover:bg-emerald-100"}`}
                >
                  {d.analyticsBtn}
                </button>
              </div>
              <div className="flex-1 mt-4 pb-16">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={MOCK_SLSIC_FULL}>
                    <XAxis dataKey="name" tick={{ fontSize: 9 }} interval={0} angle={-45} textAnchor="end" height={100} tickFormatter={(v) => v.length > 20 ? `${v.substring(0, 20)}...` : v} />
                    <YAxis type="number" tick={{ fontSize: 10 }} />
                    <Tooltip />
                    <Bar dataKey="count" fill="#10b981" radius={[4, 4, 0, 0]} barSize={22} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            {/* NEW: EMPLOYMENT SECTOR CHART */}
            <div className="bg-white p-5 rounded-xl shadow-sm h-[440px] flex flex-col">
              <div className="flex justify-between items-center border-b border-gray-50 pb-2">
                <div>
                  <h4 className="text-sm font-bold text-gray-800">{d.sectorChartTitle}</h4>
                  <p className="text-[11px] text-gray-400">{d.sectorChartSub}</p>
                </div>
                <button
                  onClick={handleOpenSectorPanel}
                  className={`text-xs font-bold px-3 py-1.5 rounded-lg transition-all ${activePanel === "SECTOR" ? "bg-indigo-600 text-white" : "text-indigo-600 bg-indigo-50 hover:bg-indigo-100"}`}
                >
                  {d.analyticsBtn}
                </button>
              </div>
              <div className="flex-1 mt-4">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={MOCK_SECTOR_FULL} margin={{ top: 10, right: 10, left: -10, bottom: 5 }}>
                    <XAxis dataKey="name" tick={{ fontSize: 10 }} />
                    <YAxis tick={{ fontSize: 10 }} />
                    <Tooltip />
                    <Bar dataKey="count" radius={[4, 4, 0, 0]} barSize={48}>
                      {MOCK_SECTOR_FULL.map((_, i) => <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                    </Bar>
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            {/* EXPERIENCE & EDUCATION */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="bg-white p-5 rounded-xl shadow-sm h-[400px] flex flex-col">
                <h4 className="text-sm font-bold text-gray-800 border-b border-gray-50 pb-2">{d.expChartTitle}</h4>
                <div className="flex-1 mt-4">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={MOCK_EXPERIENCE_DIST} margin={{ top: 10, right: 10, left: -20, bottom: 5 }}>
                      <XAxis dataKey="name" tick={{ fontSize: 10 }} />
                      <YAxis tick={{ fontSize: 10 }} />
                      <Tooltip />
                      <Bar dataKey="value" fill="#f59e0b" radius={[4, 4, 0, 0]} barSize={30} />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </div>

              <div className="bg-white p-5 rounded-xl shadow-sm h-[400px] flex flex-col">
                <h4 className="text-sm font-bold text-gray-800 border-b border-gray-50 pb-2">{d.eduChartTitle}</h4>
                <div className="flex-1 flex items-center justify-center">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={MOCK_EDUCATION_DIST} cx="50%" cy="45%" innerRadius={65} outerRadius={90} paddingAngle={3} dataKey="value">
                        {MOCK_EDUCATION_DIST.map((_, i) => <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                      </Pie>
                      <Legend verticalAlign="bottom" iconSize={8} iconType="circle" wrapperStyle={{ fontSize: 11 }} />
                      <Tooltip />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              </div>
            </div>

            {/* NEW: FORMAL/INFORMAL & GENDER */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
              <div className="bg-white p-5 rounded-xl shadow-sm h-[380px] flex flex-col">
                <h4 className="text-sm font-bold text-gray-800 border-b border-gray-50 pb-2">{d.formalChartTitle}</h4>
                <div className="flex-1 flex items-center justify-center">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={MOCK_FORMAL_DIST} cx="50%" cy="45%" innerRadius={60} outerRadius={85} paddingAngle={3} dataKey="value">
                        {MOCK_FORMAL_DIST.map((_, i) => <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                      </Pie>
                      <Legend verticalAlign="bottom" iconSize={8} iconType="circle" wrapperStyle={{ fontSize: 11 }} />
                      <Tooltip />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              </div>

              <div className="bg-white p-5 rounded-xl shadow-sm h-[380px] flex flex-col">
                <h4 className="text-sm font-bold text-gray-800 border-b border-gray-50 pb-2">{d.genderChartTitle}</h4>
                <div className="flex-1 flex items-center justify-center">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={MOCK_GENDER_DIST} cx="50%" cy="45%" innerRadius={60} outerRadius={85} paddingAngle={3} dataKey="value">
                        {MOCK_GENDER_DIST.map((_, i) => <Cell key={i} fill={CHART_COLORS[(i + 2) % CHART_COLORS.length]} />)}
                      </Pie>
                      <Legend verticalAlign="bottom" iconSize={8} iconType="circle" wrapperStyle={{ fontSize: 11 }} />
                      <Tooltip />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
              </div>
            </div>

            {/* NEW: VOCATIONAL EDUCATION (NVQ) */}
            <div className="bg-white p-5 rounded-xl shadow-sm h-[380px] flex flex-col">
              <h4 className="text-sm font-bold text-gray-800 border-b border-gray-50 pb-2">{d.vocationalChartTitle}</h4>
              <div className="flex-1 mt-4">
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart data={MOCK_VOCATIONAL_DIST} margin={{ top: 10, right: 10, left: -10, bottom: 5 }}>
                    <XAxis dataKey="name" tick={{ fontSize: 10 }} />
                    <YAxis tick={{ fontSize: 10 }} />
                    <Tooltip />
                    <Bar dataKey="value" fill="#8b5cf6" radius={[4, 4, 0, 0]} barSize={28} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </div>

            {/* REMOTE & JOB TYPE */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-8 pb-8">
              <div className="bg-white p-6 rounded-xl shadow-sm">
                <h3 className="text-sm font-bold text-gray-800 mb-4 flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-indigo-600" />{d.remoteTitle}
                </h3>
                <div className="space-y-4">
                  {[{ label: "On-Site", share: 64 }, { label: "Remote", share: 15 }].map((item, i) => (
                    <div key={i}>
                      <div className="flex justify-between text-xs mb-1 font-medium text-gray-600">
                        <span>{item.label}</span><span className="font-mono">{item.share}%</span>
                      </div>
                      <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                        <div className="h-full bg-indigo-600 transition-all duration-1000" style={{ width: `${item.share}%` }} />
                      </div>
                    </div>
                  ))}
                </div>
              </div>

              <div className="bg-white p-6 rounded-xl shadow-sm">
                <h3 className="text-sm font-bold text-gray-800 mb-4 flex items-center gap-2">
                  <span className="w-2 h-2 rounded-full bg-emerald-600" />{d.contractTitle}
                </h3>
                <div className="space-y-4">
                  {[{ type: "Full-Time", share: 72 }, { type: "Part-Time", share: 14 }].map((jt, i) => (
                    <div key={i}>
                      <div className="flex justify-between text-xs mb-1 font-medium text-gray-600">
                        <span>{jt.type}</span><span className="font-mono">{jt.share}%</span>
                      </div>
                      <div className="w-full h-2 bg-gray-100 rounded-full overflow-hidden">
                        <div className="h-full bg-emerald-600 transition-all duration-1000" style={{ width: `${jt.share}%` }} />
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>

          </div>
        </div>

        {/* ── RIGHT: ANALYTICS PANEL (inline) ──────────────── */}
        {activePanel && (
          <div className="w-[480px] min-w-[480px] border-l border-gray-200 bg-white flex flex-col overflow-hidden shadow-lg">

            {/* Panel Header */}
            <div className="px-6 py-4 border-b border-gray-100 flex justify-between items-start bg-white shrink-0">
              <div>
                <h2 className="text-sm font-black text-gray-900">
                  {activePanel === "OCCUPATION" ? d.occModalTitle : activePanel === "INDUSTRY" ? d.indModalTitle : d.sectorModalTitle}
                </h2>
                <p className="text-[11px] text-gray-400 mt-0.5">
                  {activePanel === "OCCUPATION" ? d.occModalSub : activePanel === "INDUSTRY" ? d.indModalSub : d.sectorModalSub}
                </p>
              </div>
              <button
                onClick={() => setActivePanel(null)}
                className="ml-4 shrink-0 text-xs font-bold text-gray-400 hover:text-gray-700 bg-gray-100 hover:bg-gray-200 px-3 py-1.5 rounded-lg transition-all"
              >
                {d.closePanel} ✕
              </button>
            </div>

            {/* Panel Body */}
            <div className="flex-1 overflow-y-auto p-5 space-y-5">

              {/* ── OCCUPATION PANEL CONTENT ────────────────────────────── */}
              {activePanel === "OCCUPATION" && (
                <>
                  <div className="flex flex-wrap gap-2">
                    {MOCK_SLSO_FULL.map((occ) => (
                      <button
                        key={occ.id}
                        onClick={() => { setSelectedOccId(occ.id); setSelectedOccName(occ.name); }}
                        className={`text-[11px] px-3 py-1.5 rounded-lg font-bold transition-all border ${selectedOccName === occ.name ? "bg-blue-600 text-white border-blue-600" : "bg-white text-gray-600 border-gray-200 hover:border-blue-300"}`}
                      >
                        {occ.name}
                      </button>
                    ))}
                  </div>

                  {/* NEW: year filter for the occupation panel */}
                  <div>
                    <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-2">{d.yearFilterHeader}</h4>
                    <div className="flex bg-gray-100 p-1 rounded-xl w-fit">
                      {DYNAMIC_YEARS.map((year) => (
                        <button
                          key={year}
                          onClick={() => setOccAnalyticsYear(Number(year))}
                          className={`px-4 py-1.5 text-[11px] font-bold rounded-lg transition-all ${occAnalyticsYear === Number(year) ? "bg-white text-blue-600 shadow-sm" : "text-gray-500 hover:text-gray-800"}`}
                        >
                          {year}
                        </button>
                      ))}
                    </div>
                  </div>

                  <div className="bg-gray-50 p-4 rounded-xl">
                    <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.trendChartHeader}</h4>
                    <div className="h-40">
                      <ResponsiveContainer>
                        <LineChart data={occTrendData}>
                          <XAxis dataKey="year" tick={{ fontSize: 10 }} />
                          <YAxis tick={{ fontSize: 10 }} />
                          <Tooltip />
                          <Line type="monotone" dataKey="vacancies" stroke="#2563eb" strokeWidth={2} dot={{ r: 3 }} />
                        </LineChart>
                      </ResponsiveContainer>
                    </div>
                  </div>

                  {/* NEW: Formal / Informal job count, filtered by occupation + year */}
                  <div className="bg-gray-50 p-4 rounded-xl">
                    <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.formalInformalHeader}</h4>
                    <div className="h-36">
                      <ResponsiveContainer>
                        <BarChart data={occYearDetail.formal} margin={{ left: -20 }}>
                          <XAxis dataKey="name" tick={{ fontSize: 10 }} />
                          <YAxis tick={{ fontSize: 10 }} />
                          <Tooltip />
                          <Bar dataKey="value" radius={[4, 4, 0, 0]} barSize={40}>
                            {occYearDetail.formal.map((_: any, i: number) => <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                          </Bar>
                        </BarChart>
                      </ResponsiveContainer>
                    </div>
                  </div>

                  {/* NEW: Gender wise job count, filtered by occupation + year */}
                  <div className="bg-gray-50 p-4 rounded-xl">
                    <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.genderHeader}</h4>
                    <div className="h-40">
                      <ResponsiveContainer>
                        <PieChart>
                          <Pie data={occYearDetail.gender} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={55} paddingAngle={2}>
                            {occYearDetail.gender.map((_: any, i: number) => <Cell key={i} fill={CHART_COLORS[(i + 2) % CHART_COLORS.length]} />)}
                          </Pie>
                          <Legend iconSize={8} iconType="circle" wrapperStyle={{ fontSize: 10 }} />
                          <Tooltip />
                        </PieChart>
                      </ResponsiveContainer>
                    </div>
                  </div>

                  <div>
                    <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">
                      {d.demandingJobsHeader} <span className="text-blue-600">{selectedOccName}</span>
                    </h4>
                    <div className="space-y-2">
                      {activeOccData?.demanding.map((role, i) => (
                        <div key={i} className="flex items-center justify-between bg-gray-50 px-4 py-3 rounded-xl border border-gray-100">
                          <div className="flex items-center gap-3">
                            <span className="text-[10px] font-black text-gray-300 font-mono w-4">{i + 1}</span>
                            <span className="text-xs font-bold text-gray-800">{role}</span>
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </>
              )}

              {/* ── INDUSTRY PANEL CONTENT ──────────────────────────────── */}
              {activePanel === "INDUSTRY" && (
                <>
                  <div className="space-y-3">
                    <select
                      value={selectedIndId}
                      onChange={(e) => {
                        const found = MOCK_SLSIC_FULL.find((i) => i.id === Number(e.target.value));
                        if (found) { setSelectedIndId(found.id); setSelectedIndName(found.name); }
                      }}
                      className="w-full bg-gray-50 rounded-xl p-2.5 text-xs font-bold border border-gray-200 outline-none focus:ring-2 focus:ring-emerald-200"
                    >
                      {MOCK_SLSIC_FULL.map((ind) => (
                        <option key={ind.id} value={ind.id}>{ind.name}</option>
                      ))}
                    </select>

                    <div className="flex bg-gray-100 p-1 rounded-xl w-fit">
                      {DYNAMIC_YEARS.map((year) => (
                        <button
                          key={year}
                          onClick={() => setAnalyticsYear(Number(year))}
                          className={`px-4 py-1.5 text-[11px] font-bold rounded-lg transition-all ${analyticsYear === Number(year) ? "bg-white text-emerald-600 shadow-sm" : "text-gray-500 hover:text-gray-800"}`}
                        >
                          {year}
                        </button>
                      ))}
                    </div>
                  </div>

                  <div className="space-y-4">
                    <div className="bg-gray-50 p-4 rounded-xl">
                      <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.sectorTrendHeader}</h4>
                      <div className="h-36">
                        <ResponsiveContainer>
                          <LineChart data={indData.trends}>
                            <XAxis dataKey="year" tick={{ fontSize: 10 }} />
                            <YAxis tick={{ fontSize: 10 }} />
                            <Tooltip />
                            <Line type="monotone" dataKey="vacancies" stroke="#10b981" strokeWidth={2} dot={{ r: 3 }} />
                          </LineChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div className="bg-gray-50 p-4 rounded-xl">
                      <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.expAllocHeader}</h4>
                      <div className="h-36">
                        <ResponsiveContainer>
                          <BarChart data={indData.experience[analyticsYear]} margin={{ left: -20 }}>
                            <XAxis dataKey="label" tick={{ fontSize: 9 }} />
                            <YAxis tick={{ fontSize: 9 }} />
                            <Tooltip />
                            <Bar dataKey="value" fill="#3b82f6" radius={[3, 3, 0, 0]} />
                          </BarChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div className="bg-gray-50 p-4 rounded-xl">
                      <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.regionalShareHeader}</h4>
                      <div className="h-44">
                        <ResponsiveContainer>
                          <PieChart>
                            <Pie data={indData.provinces[analyticsYear]} dataKey="value" nameKey="name" cx="50%" cy="50%" outerRadius={55} paddingAngle={2}>
                              {indData.provinces[analyticsYear]?.map((_: any, i: number) => <Cell key={i} fill={CHART_COLORS[i % CHART_COLORS.length]} />)}
                            </Pie>
                            <Legend iconSize={8} iconType="circle" wrapperStyle={{ fontSize: 10 }} />
                            <Tooltip />
                          </PieChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div className="bg-gray-50 p-4 rounded-xl">
                      <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.eduThresholdHeader}</h4>
                      <div className="h-36">
                        <ResponsiveContainer>
                          <BarChart data={indData.education[analyticsYear]} margin={{ left: -20 }}>
                            <XAxis dataKey="label" tick={{ fontSize: 9 }} />
                            <YAxis tick={{ fontSize: 9 }} />
                            <Tooltip />
                            <Bar dataKey="value" fill="#f59e0b" radius={[3, 3, 0, 0]} />
                          </BarChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    {/* NEW: Vocational education wise job count, for selected industry + year */}
                    <div className="bg-gray-50 p-4 rounded-xl">
                      <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.vocationalIndHeader}</h4>
                      <div className="h-36">
                        <ResponsiveContainer>
                          <BarChart data={vocationalIndData} margin={{ left: -20 }}>
                            <XAxis dataKey="label" tick={{ fontSize: 9 }} />
                            <YAxis tick={{ fontSize: 9 }} />
                            <Tooltip />
                            <Bar dataKey="value" fill="#8b5cf6" radius={[3, 3, 0, 0]} />
                          </BarChart>
                        </ResponsiveContainer>
                      </div>
                    </div>

                    <div>
                      <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">{d.topEnterpriseHeader}</h4>
                      <div className="space-y-2">
                        {indData.employers[analyticsYear]?.map((emp: any, i: number) => (
                          <div key={i} className="flex items-center justify-between bg-gray-50 px-4 py-3 rounded-xl border border-gray-100">
                            <div className="flex items-center gap-3">
                              <span className="text-[10px] font-black text-gray-300 font-mono w-4">{i + 1}</span>
                              <span className="text-xs font-semibold text-gray-800 truncate max-w-[240px]">{emp.name}</span>
                            </div>
                            <span className="text-[11px] font-black text-emerald-600 bg-emerald-50 px-2.5 py-1 rounded-lg shrink-0">{emp.open_job_count}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                  </div>
                </>
              )}

              {/* ── NEW: EMPLOYMENT SECTOR PANEL CONTENT ────────────────── */}
              {activePanel === "SECTOR" && (
                <>
                  <div>
                    <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-2">{d.selectSectorLabel}</h4>
                    <div className="flex flex-wrap gap-2">
                      {MOCK_SECTOR_FULL.map((sec) => (
                        <button
                          key={sec.id}
                          onClick={() => { setSelectedSectorId(sec.id); setSelectedSectorName(sec.name); }}
                          className={`text-[11px] px-3 py-1.5 rounded-lg font-bold transition-all border ${selectedSectorName === sec.name ? "bg-indigo-600 text-white border-indigo-600" : "bg-white text-gray-600 border-gray-200 hover:border-indigo-300"}`}
                        >
                          {sec.name}
                        </button>
                      ))}
                    </div>
                  </div>

                  <div className="bg-gray-50 p-4 rounded-xl">
                    <h4 className="text-[10px] font-black uppercase text-gray-400 tracking-wider mb-3">
                      {d.sectorPanelTrendHeader} — <span className="text-indigo-600">{selectedSectorName}</span>
                    </h4>
                    <div className="h-44">
                      <ResponsiveContainer>
                        <LineChart data={sectorTrendData}>
                          <XAxis dataKey="year" tick={{ fontSize: 10 }} />
                          <YAxis tick={{ fontSize: 10 }} />
                          <Tooltip />
                          <Line type="monotone" dataKey="vacancies" stroke="#6366f1" strokeWidth={2} dot={{ r: 3 }} />
                        </LineChart>
                      </ResponsiveContainer>
                    </div>
                  </div>
                </>
              )}

            </div>
          </div>
        )}
      </div>
    </div>
  );
}