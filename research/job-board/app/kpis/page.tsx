"use client";

import { useState } from "react";

/**
 * ROUTE: app/admin/kpis/page.tsx  (Next.js App Router)
 * If you use the Pages Router instead, put this file at pages/admin/kpis.tsx
 * and it will work unchanged (it has no App-Router-only APIs).
 *
 * This page provides full CRUD (Create, Read, Update, Delete) management
 * for every KPI reference list used by the main labour-market dashboard:
 *  - Employment Sector
 *  - Experience
 *  - Education Level
 *  - Formal / Informal
 *  - Gender
 *  - Vocational Education (NVQ)
 *  - Remote / On-Site
 *  - Contract Type
 *
 * ...plus two hierarchical, parent/child KPIs, each with 5 levels:
 *  - Occupation:
 *      Major Group -> Sub Major Group -> Minor Group -> Unit Group -> Occupation
 *  - Industry:
 *      Industry Sector -> Industry Division -> Industry Group -> Industry Class
 *      -> Industry Sub Class
 *
 *  Every node at every level of both hierarchies is { name, code, parent_id }.
 *  Only the top level of each hierarchy (Major Group / Industry Sector) has no
 *  parent_id — every level below it must have one. Both hierarchies share one
 *  generic engine (HIERARCHIES / HierarchyItem / the handleH* functions), so
 *  adding a third hierarchical KPI later is just another entry in
 *  HIERARCHIES + INITIAL_HIERARCHY_DATA.
 *
 * Data lives in local React state (mock data) so it works standalone.
 * Swap the setData()/setHData() calls for API calls (fetch/axios to your
 * backend) when you're ready to wire this up to a real data source — the
 * shapes (KPIItem / HierarchyItem) and handlers are written so that swap is
 * a drop-in change. For the backend, each HierarchyItem maps directly onto
 * { name, code, parent_id } per level — parent_id should be omitted/null for
 * the top level of a hierarchy and required for every level below it.
 */

// ---------------------------------------------------------------------------
// TYPES — flat, single-level KPIs
// ---------------------------------------------------------------------------
type KPIItem = {
  id: string;
  name: string;
};

type KPICategoryKey =
  | "employmentSector"
  | "experience"
  | "educationLevel"
  | "formalInformal"
  | "gender"
  | "vocationalEducation"
  | "remoteOnsite"
  | "contractType";

type KPICategoryConfig = {
  title: string;
  description: string;
  itemLabel: string;
  accent: string; // tailwind color token, e.g. "blue"
};

// ---------------------------------------------------------------------------
// TYPES — hierarchical KPIs (Occupation, Industry, ...)
// ---------------------------------------------------------------------------
type HierarchyId = "occupation" | "industry";

type HierarchyItem = {
  id: string;
  name: string;
  code: string;
  parentId: string | null; // null only for the top level of the hierarchy
};

type HierarchyLevelDef = {
  key: string; // unique within this hierarchy, e.g. "majorGroup" / "sector"
  title: string; // e.g. "Major Group" / "Industry Sector"
  itemLabel: string; // used in buttons/labels, usually same as title
};

type HierarchyDef = {
  label: string; // sidebar nav label, e.g. "Occupation"
  accent: string;
  levels: HierarchyLevelDef[]; // ordered top (no parent) -> bottom (leaf)
};

// ---------------------------------------------------------------------------
// CATEGORY CONFIG — flat KPIs
// ---------------------------------------------------------------------------
const KPI_CATEGORIES: Record<KPICategoryKey, KPICategoryConfig> = {
  employmentSector: {
    title: "Employment Sector",
    description: "Government, Semi-Government, Private and NGO sector breakdown",
    itemLabel: "Sector",
    accent: "blue",
  },
  experience: {
    title: "Experience",
    description: "Experience-band breakdown used across occupation & industry charts",
    itemLabel: "Experience Level",
    accent: "amber",
  },
  educationLevel: {
    title: "Education Level",
    description: "Academic qualification thresholds",
    itemLabel: "Education Level",
    accent: "violet",
  },
  formalInformal: {
    title: "Formal / Informal",
    description: "Formal vs. informal sector classification",
    itemLabel: "Type",
    accent: "emerald",
  },
  gender: {
    title: "Gender",
    description: "Gender breakdown of current vacancies",
    itemLabel: "Gender",
    accent: "pink",
  },
  vocationalEducation: {
    title: "Vocational Education (NVQ)",
    description: "National Vocational Qualification levels",
    itemLabel: "NVQ Level",
    accent: "indigo",
  },
  remoteOnsite: {
    title: "Remote / On-Site",
    description: "Work-mode configuration share",
    itemLabel: "Work Mode",
    accent: "cyan",
  },
  contractType: {
    title: "Contract Type",
    description: "Employment contract type share",
    itemLabel: "Contract Type",
    accent: "teal",
  },
};

const CATEGORY_ORDER: KPICategoryKey[] = [
  "employmentSector",
  "experience",
  "educationLevel",
  "formalInformal",
  "gender",
  "vocationalEducation",
  "remoteOnsite",
  "contractType",
];

// ---------------------------------------------------------------------------
// HIERARCHY CONFIG — Occupation + Industry, top level -> leaf level
// ---------------------------------------------------------------------------
const HIERARCHY_ORDER: HierarchyId[] = ["occupation", "industry"];

const HIERARCHIES: Record<HierarchyId, HierarchyDef> = {
  occupation: {
    label: "Occupation",
    accent: "rose",
    levels: [
      { key: "majorGroup", title: "Major Group", itemLabel: "Major Group" },
      { key: "subMajorGroup", title: "Sub Major Group", itemLabel: "Sub Major Group" },
      { key: "minorGroup", title: "Minor Group", itemLabel: "Minor Group" },
      { key: "unitGroup", title: "Unit Group", itemLabel: "Unit Group" },
      { key: "occupation", title: "Occupation", itemLabel: "Occupation" },
    ],
  },
  industry: {
    label: "Industry",
    accent: "orange",
    levels: [
      { key: "sector", title: "Industry Sector", itemLabel: "Industry Sector" },
      { key: "division", title: "Industry Division", itemLabel: "Industry Division" },
      { key: "group", title: "Industry Group", itemLabel: "Industry Group" },
      { key: "class", title: "Industry Class", itemLabel: "Industry Class" },
      { key: "subClass", title: "Industry Sub Class", itemLabel: "Industry Sub Class" },
    ],
  },
};

// ---------------------------------------------------------------------------
// MOCK DATA — flat KPIs
// ---------------------------------------------------------------------------
const INITIAL_KPI_DATA: Record<KPICategoryKey, KPIItem[]> = {
  employmentSector: [
    { id: "sec-1", name: "Government" },
    { id: "sec-2", name: "Semi Government" },
    { id: "sec-3", name: "Private" },
    { id: "sec-4", name: "NGO" },
  ],
  experience: [
    { id: "exp-1", name: "Entry Level" },
    { id: "exp-2", name: "Junior" },
    { id: "exp-3", name: "Mid-Level" },
    { id: "exp-4", name: "Senior" },
  ],
  educationLevel: [
    { id: "edu-1", name: "Degree" },
    { id: "edu-2", name: "A/L" },
    { id: "edu-3", name: "O/L" },
    { id: "edu-4", name: "Below O/L" },
    { id: "edu-5", name: "Not Specified" },
  ],
  formalInformal: [
    { id: "fi-1", name: "Formal" },
    { id: "fi-2", name: "Informal" },
  ],
  gender: [
    { id: "gen-1", name: "Male" },
    { id: "gen-2", name: "Female" },
    { id: "gen-3", name: "Not Specified" },
  ],
  vocationalEducation: [
    { id: "nvq-1", name: "NVQ 1" },
    { id: "nvq-2", name: "NVQ 2" },
    { id: "nvq-3", name: "NVQ 3" },
    { id: "nvq-4", name: "NVQ 4" },
    { id: "nvq-5", name: "NVQ 5" },
    { id: "nvq-6", name: "NVQ 6" },
    { id: "nvq-7", name: "NVQ 7" },
  ],
  remoteOnsite: [
    { id: "rem-1", name: "On-Site" },
    { id: "rem-2", name: "Remote" },
    { id: "rem-3", name: "Hybrid" },
  ],
  contractType: [
    { id: "ct-1", name: "Full-Time" },
    { id: "ct-2", name: "Part-Time" },
    { id: "ct-3", name: "Contract / Temporary" },
  ],
};

// ---------------------------------------------------------------------------
// MOCK DATA — hierarchical KPIs (small ISCO/ISIC-style sample trees)
// ---------------------------------------------------------------------------
const INITIAL_HIERARCHY_DATA: Record<HierarchyId, Record<string, HierarchyItem[]>> = {
  occupation: {
    majorGroup: [
      { id: "mg-1", name: "Managers", code: "1", parentId: null },
      { id: "mg-2", name: "Professionals", code: "2", parentId: null },
    ],
    subMajorGroup: [
      { id: "smg-11", name: "Chief Executives, Senior Officials and Legislators", code: "11", parentId: "mg-1" },
      { id: "smg-12", name: "Administrative and Commercial Managers", code: "12", parentId: "mg-1" },
      { id: "smg-21", name: "Science and Engineering Professionals", code: "21", parentId: "mg-2" },
    ],
    minorGroup: [
      { id: "mng-111", name: "Legislators and Senior Officials", code: "111", parentId: "smg-11" },
      { id: "mng-211", name: "Physical and Earth Science Professionals", code: "211", parentId: "smg-21" },
    ],
    unitGroup: [
      { id: "ug-1111", name: "Legislators", code: "1111", parentId: "mng-111" },
      { id: "ug-2111", name: "Physicists and Astronomers", code: "2111", parentId: "mng-211" },
    ],
    occupation: [
      { id: "occ-11111", name: "Member of Parliament", code: "11111", parentId: "ug-1111" },
      { id: "occ-21111", name: "Physicist", code: "21111", parentId: "ug-2111" },
    ],
  },
  industry: {
    sector: [
      { id: "is-a", name: "Agriculture, Forestry and Fishing", code: "A", parentId: null },
      { id: "is-c", name: "Manufacturing", code: "C", parentId: null },
    ],
    division: [
      { id: "id-01", name: "Crop and Animal Production, Hunting", code: "01", parentId: "is-a" },
      { id: "id-10", name: "Manufacture of Food Products", code: "10", parentId: "is-c" },
    ],
    group: [
      { id: "ig-011", name: "Growing of Non-Perennial Crops", code: "011", parentId: "id-01" },
      { id: "ig-107", name: "Manufacture of Other Food Products", code: "107", parentId: "id-10" },
    ],
    class: [
      { id: "ic-0111", name: "Growing of Cereals, Leguminous Crops and Oil Seeds", code: "0111", parentId: "ig-011" },
      { id: "ic-1071", name: "Manufacture of Bakery Products", code: "1071", parentId: "ig-107" },
    ],
    subClass: [
      { id: "isc-01111", name: "Growing of Cereals (Except Rice)", code: "01111", parentId: "ic-0111" },
      { id: "isc-10711", name: "Manufacture of Bread, Fresh Pastry Goods and Cakes", code: "10711", parentId: "ic-1071" },
    ],
  },
};

// ---------------------------------------------------------------------------
// ACCENT COLOR CLASS MAP (Tailwind needs literal class names, not template strings)
// ---------------------------------------------------------------------------
const ACCENT_CLASSES: Record<string, { bg: string; text: string; ring: string; solidBg: string; lightBg: string }> = {
  blue: { bg: "bg-blue-600", text: "text-blue-600", ring: "ring-blue-500", solidBg: "bg-blue-600", lightBg: "bg-blue-50" },
  amber: { bg: "bg-amber-600", text: "text-amber-600", ring: "ring-amber-500", solidBg: "bg-amber-500", lightBg: "bg-amber-50" },
  violet: { bg: "bg-violet-600", text: "text-violet-600", ring: "ring-violet-500", solidBg: "bg-violet-600", lightBg: "bg-violet-50" },
  emerald: { bg: "bg-emerald-600", text: "text-emerald-600", ring: "ring-emerald-500", solidBg: "bg-emerald-600", lightBg: "bg-emerald-50" },
  pink: { bg: "bg-pink-600", text: "text-pink-600", ring: "ring-pink-500", solidBg: "bg-pink-600", lightBg: "bg-pink-50" },
  indigo: { bg: "bg-indigo-600", text: "text-indigo-600", ring: "ring-indigo-500", solidBg: "bg-indigo-600", lightBg: "bg-indigo-50" },
  cyan: { bg: "bg-cyan-600", text: "text-cyan-600", ring: "ring-cyan-500", solidBg: "bg-cyan-600", lightBg: "bg-cyan-50" },
  teal: { bg: "bg-teal-600", text: "text-teal-600", ring: "ring-teal-500", solidBg: "bg-teal-600", lightBg: "bg-teal-50" },
  rose: { bg: "bg-rose-600", text: "text-rose-600", ring: "ring-rose-500", solidBg: "bg-rose-600", lightBg: "bg-rose-50" },
  orange: { bg: "bg-orange-600", text: "text-orange-600", ring: "ring-orange-500", solidBg: "bg-orange-600", lightBg: "bg-orange-50" },
};

// ---------------------------------------------------------------------------
// SMALL ID GENERATOR (works without relying on crypto.randomUUID typings)
// ---------------------------------------------------------------------------
function makeId(prefix: string) {
  return `${prefix}-${Date.now()}-${Math.floor(Math.random() * 100000)}`;
}

// ---------------------------------------------------------------------------
// HIERARCHY TREE HELPERS (pure, generic across any hierarchy's level list —
// easy to unit test / swap for API calls)
// ---------------------------------------------------------------------------

// Recursively count how many descendant nodes (across all lower levels) sit
// underneath a given node. Used to warn the user before a cascading delete.
function countHierarchyDescendants(
  data: Record<string, HierarchyItem[]>,
  levels: HierarchyLevelDef[],
  levelIndex: number,
  id: string
): number {
  if (levelIndex === levels.length - 1) return 0;
  const childLevel = levels[levelIndex + 1];
  const children = data[childLevel.key].filter((c) => c.parentId === id);
  return children.reduce(
    (sum, child) => sum + 1 + countHierarchyDescendants(data, levels, levelIndex + 1, child.id),
    0
  );
}

// Removes a node and every descendant beneath it, across all lower levels.
function cascadeDeleteHierarchy(
  data: Record<string, HierarchyItem[]>,
  levels: HierarchyLevelDef[],
  levelIndex: number,
  id: string
): Record<string, HierarchyItem[]> {
  let result = { ...data };
  if (levelIndex < levels.length - 1) {
    const childLevel = levels[levelIndex + 1];
    const childIds = data[childLevel.key].filter((c) => c.parentId === id).map((c) => c.id);
    for (const childId of childIds) {
      result = cascadeDeleteHierarchy(result, levels, levelIndex + 1, childId);
    }
  }
  const level = levels[levelIndex];
  result = { ...result, [level.key]: result[level.key].filter((i) => i.id !== id) };
  return result;
}

type SectionKey = "flat" | HierarchyId;

export default function ManageKpisPage() {
  // ---- SECTION SWITCH: flat single-value KPIs vs. one of the hierarchical KPIs
  const [activeSection, setActiveSection] = useState<SectionKey>("flat");

  // =========================================================================
  // FLAT KPI STATE
  // =========================================================================
  const [data, setData] = useState<Record<KPICategoryKey, KPIItem[]>>(INITIAL_KPI_DATA);
  const [activeCategory, setActiveCategory] = useState<KPICategoryKey>("employmentSector");

  const [newName, setNewName] = useState("");
  const [formError, setFormError] = useState<string | null>(null);

  const [editingId, setEditingId] = useState<string | null>(null);
  const [editName, setEditName] = useState("");
  const [editError, setEditError] = useState<string | null>(null);

  const [deleteTarget, setDeleteTarget] = useState<KPIItem | null>(null);
  const [toast, setToast] = useState<string | null>(null);

  const config = KPI_CATEGORIES[activeCategory];
  const accent = ACCENT_CLASSES[config.accent];
  const items = data[activeCategory];

  const [currentLang, setCurrentLang] = useState<"en" | "si" | "ta">("en");
    const d = currentLang;
    const formattedDate = new Date().toLocaleDateString(
        currentLang === "en" ? "en-US" : currentLang === "si" ? "si-LK" : "ta-LK",
        { year: "numeric", month: "short", day: "numeric" }
    );

  function flashToast(msg: string) {
    setToast(msg);
    setTimeout(() => setToast(null), 2200);
  }

  function switchCategory(cat: KPICategoryKey) {
    setActiveSection("flat");
    setActiveCategory(cat);
    setEditingId(null);
    setEditError(null);
    setFormError(null);
    setNewName("");
  }

  // ---- CREATE -------------------------------------------------------------
  function handleAdd(e: React.FormEvent) {
    e.preventDefault();
    const trimmed = newName.trim();
    if (!trimmed) {
      setFormError(`${config.itemLabel} name is required.`);
      return;
    }
    if (items.some((i) => i.name.toLowerCase() === trimmed.toLowerCase())) {
      setFormError(`"${trimmed}" already exists in ${config.title}.`);
      return;
    }

    const newItem: KPIItem = { id: makeId(activeCategory), name: trimmed };
    setData((prev) => ({ ...prev, [activeCategory]: [...prev[activeCategory], newItem] }));
    setNewName("");
    setFormError(null);
    flashToast(`Added "${trimmed}"`);
  }

  // ---- UPDATE ---------------------------------------------------------------
  function startEdit(item: KPIItem) {
    setEditingId(item.id);
    setEditName(item.name);
    setEditError(null);
  }

  function cancelEdit() {
    setEditingId(null);
    setEditError(null);
  }

  function saveEdit(id: string) {
    const trimmed = editName.trim();
    if (!trimmed) {
      setEditError(`${config.itemLabel} name is required.`);
      return;
    }
    if (items.some((i) => i.id !== id && i.name.toLowerCase() === trimmed.toLowerCase())) {
      setEditError(`"${trimmed}" already exists in ${config.title}.`);
      return;
    }

    setData((prev) => ({
      ...prev,
      [activeCategory]: prev[activeCategory].map((i) =>
        i.id === id ? { ...i, name: trimmed } : i
      ),
    }));
    setEditingId(null);
    setEditError(null);
    flashToast(`Updated "${trimmed}"`);
  }

  // ---- DELETE ---------------------------------------------------------------
  function handleDelete() {
    if (!deleteTarget) return;
    setData((prev) => ({
      ...prev,
      [activeCategory]: prev[activeCategory].filter((i) => i.id !== deleteTarget.id),
    }));
    flashToast(`Deleted "${deleteTarget.name}"`);
    setDeleteTarget(null);
  }

  // =========================================================================
  // HIERARCHICAL KPI STATE (shared engine for Occupation + Industry)
  // =========================================================================
  const [hData, setHData] = useState<Record<HierarchyId, Record<string, HierarchyItem[]>>>(
    INITIAL_HIERARCHY_DATA
  );
  const [hActiveLevelIndex, setHActiveLevelIndex] = useState(0);

  const [hNewName, setHNewName] = useState("");
  const [hNewCode, setHNewCode] = useState("");
  const [hNewParentId, setHNewParentId] = useState("");
  const [hFormError, setHFormError] = useState<string | null>(null);

  const [hEditingId, setHEditingId] = useState<string | null>(null);
  const [hEditName, setHEditName] = useState("");
  const [hEditCode, setHEditCode] = useState("");
  const [hEditParentId, setHEditParentId] = useState("");
  const [hEditError, setHEditError] = useState<string | null>(null);

  const [hDeleteTarget, setHDeleteTarget] = useState<HierarchyItem | null>(null);

  const isHierarchySection = activeSection === "occupation" || activeSection === "industry";
  const hierarchyId: HierarchyId | null = isHierarchySection ? (activeSection as HierarchyId) : null;
  const hierarchyDef = hierarchyId ? HIERARCHIES[hierarchyId] : null;
  const hAccent = hierarchyDef ? ACCENT_CLASSES[hierarchyDef.accent] : ACCENT_CLASSES.rose;
  const hLevels = hierarchyDef ? hierarchyDef.levels : [];
  const hLevelDef = hierarchyDef ? hLevels[hActiveLevelIndex] : null;
  const hParentLevelDef = hierarchyDef && hActiveLevelIndex > 0 ? hLevels[hActiveLevelIndex - 1] : null;
  const hItems = hierarchyId && hLevelDef ? hData[hierarchyId][hLevelDef.key] : [];
  const hParentItems = hierarchyId && hParentLevelDef ? hData[hierarchyId][hParentLevelDef.key] : [];

  function hierarchyTotalCount(id: HierarchyId): number {
    return HIERARCHIES[id].levels.reduce((sum, lvl) => sum + hData[id][lvl.key].length, 0);
  }

  function hParentDisplay(item: HierarchyItem): string {
    if (!hierarchyId || !hParentLevelDef || !item.parentId) return "—";
    const parent = hData[hierarchyId][hParentLevelDef.key].find((p) => p.id === item.parentId);
    return parent ? `${parent.code} — ${parent.name}` : "—";
  }

  function resetHForms() {
    setHEditingId(null);
    setHEditError(null);
    setHFormError(null);
    setHNewName("");
    setHNewCode("");
    setHNewParentId("");
  }

  function switchHierarchy(id: HierarchyId) {
    setActiveSection(id);
    setHActiveLevelIndex(0);
    resetHForms();
  }

  function switchHLevel(index: number) {
    setHActiveLevelIndex(index);
    resetHForms();
  }

  // ---- CREATE -------------------------------------------------------------
  function handleHAdd(e: React.FormEvent) {
    e.preventDefault();
    if (!hierarchyId || !hLevelDef) return;

    const trimmedName = hNewName.trim();
    const trimmedCode = hNewCode.trim();

    if (!trimmedName) {
      setHFormError(`${hLevelDef.itemLabel} name is required.`);
      return;
    }
    if (!trimmedCode) {
      setHFormError(`${hLevelDef.itemLabel} code is required.`);
      return;
    }
    if (hParentLevelDef && !hNewParentId) {
      setHFormError(`Select a parent ${hParentLevelDef.itemLabel}.`);
      return;
    }
    if (hItems.some((i) => i.code.toLowerCase() === trimmedCode.toLowerCase())) {
      setHFormError(`Code "${trimmedCode}" already exists in ${hLevelDef.title}.`);
      return;
    }
    if (hItems.some((i) => i.name.toLowerCase() === trimmedName.toLowerCase())) {
      setHFormError(`"${trimmedName}" already exists in ${hLevelDef.title}.`);
      return;
    }

    const newItem: HierarchyItem = {
      id: makeId(hLevelDef.key),
      name: trimmedName,
      code: trimmedCode,
      parentId: hParentLevelDef ? hNewParentId : null,
    };
    setHData((prev) => ({
      ...prev,
      [hierarchyId]: { ...prev[hierarchyId], [hLevelDef.key]: [...prev[hierarchyId][hLevelDef.key], newItem] },
    }));
    setHNewName("");
    setHNewCode("");
    setHNewParentId("");
    setHFormError(null);
    flashToast(`Added "${trimmedName}"`);
  }

  // ---- UPDATE ---------------------------------------------------------------
  function startHEdit(item: HierarchyItem) {
    setHEditingId(item.id);
    setHEditName(item.name);
    setHEditCode(item.code);
    setHEditParentId(item.parentId ?? "");
    setHEditError(null);
  }

  function cancelHEdit() {
    setHEditingId(null);
    setHEditError(null);
  }

  function saveHEdit(id: string) {
    if (!hierarchyId || !hLevelDef) return;

    const trimmedName = hEditName.trim();
    const trimmedCode = hEditCode.trim();

    if (!trimmedName) {
      setHEditError(`${hLevelDef.itemLabel} name is required.`);
      return;
    }
    if (!trimmedCode) {
      setHEditError(`${hLevelDef.itemLabel} code is required.`);
      return;
    }
    if (hParentLevelDef && !hEditParentId) {
      setHEditError(`Select a parent ${hParentLevelDef.itemLabel}.`);
      return;
    }
    if (hItems.some((i) => i.id !== id && i.code.toLowerCase() === trimmedCode.toLowerCase())) {
      setHEditError(`Code "${trimmedCode}" already exists in ${hLevelDef.title}.`);
      return;
    }
    if (hItems.some((i) => i.id !== id && i.name.toLowerCase() === trimmedName.toLowerCase())) {
      setHEditError(`"${trimmedName}" already exists in ${hLevelDef.title}.`);
      return;
    }

    setHData((prev) => ({
      ...prev,
      [hierarchyId]: {
        ...prev[hierarchyId],
        [hLevelDef.key]: prev[hierarchyId][hLevelDef.key].map((i) =>
          i.id === id
            ? {
                ...i,
                name: trimmedName,
                code: trimmedCode,
                parentId: hParentLevelDef ? hEditParentId : null,
              }
            : i
        ),
      },
    }));
    setHEditingId(null);
    setHEditError(null);
    flashToast(`Updated "${trimmedName}"`);
  }

  // ---- DELETE ---------------------------------------------------------------
  function handleHDelete() {
    if (!hierarchyId || !hierarchyDef || !hDeleteTarget) return;
    setHData((prev) => ({
      ...prev,
      [hierarchyId]: cascadeDeleteHierarchy(prev[hierarchyId], hierarchyDef.levels, hActiveLevelIndex, hDeleteTarget.id),
    }));
    flashToast(`Deleted "${hDeleteTarget.name}"`);
    setHDeleteTarget(null);
  }

  const hDeleteDescendantCount =
    hierarchyId && hierarchyDef && hDeleteTarget
      ? countHierarchyDescendants(hData[hierarchyId], hierarchyDef.levels, hActiveLevelIndex, hDeleteTarget.id)
      : 0;

  return (
    <div className="min-h-screen bg-gray-50 text-gray-800 font-sans">
      {/* HEADER */}
      <header className="bg-white border-b border-gray-100 px-8 py-5 flex flex-col md:flex-row justify-between items-start md:items-center gap-4 sticky top-0 z-40">
        <div>
          <h1 className="text-xl font-black text-gray-900 tracking-tight">KPI Reference Data Manager</h1>
          <p className="text-xs text-gray-400 mt-1 font-medium">
            Create, edit, and remove the reference values that power the labour market dashboard charts
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

      <div className="flex">
        {/* ── LEFT: CATEGORY NAV ─────────────────────────────────────── */}
        <aside className="w-72 min-w-72 border-r border-gray-100 bg-white h-[calc(100vh-89px)] overflow-y-auto sticky top-[89px] py-4 px-3 space-y-1">
          {CATEGORY_ORDER.map((key) => {
            const cfg = KPI_CATEGORIES[key];
            const acc = ACCENT_CLASSES[cfg.accent];
            const count = data[key].length;
            const isActive = activeSection === "flat" && activeCategory === key;
            return (
              <button
                key={key}
                onClick={() => switchCategory(key)}
                className={`w-full text-left px-4 py-3 rounded-xl transition-all flex items-center justify-between gap-2 ${
                  isActive ? `${acc.lightBg} ring-1 ${acc.ring}` : "hover:bg-gray-50"
                }`}
              >
                <div>
                  <p className={`text-xs font-bold ${isActive ? acc.text : "text-gray-700"}`}>{cfg.title}</p>
                  <p className="text-[10px] text-gray-400 mt-0.5">{count} items</p>
                </div>
                <span
                  className={`w-2 h-2 rounded-full shrink-0 ${isActive ? acc.solidBg : "bg-gray-200"}`}
                />
              </button>
            );
          })}

          {/* Hierarchical KPIs */}
          <div className="pt-3 mt-2 border-t border-gray-100">
            <p className="px-4 pb-1.5 text-[10px] font-black uppercase text-gray-300 tracking-wider">
              Hierarchical
            </p>
            {HIERARCHY_ORDER.map((id) => {
              const def = HIERARCHIES[id];
              const acc = ACCENT_CLASSES[def.accent];
              const isActive = activeSection === id;
              return (
                <button
                  key={id}
                  onClick={() => switchHierarchy(id)}
                  className={`w-full text-left px-4 py-3 rounded-xl transition-all flex items-center justify-between gap-2 ${
                    isActive ? `${acc.lightBg} ring-1 ${acc.ring}` : "hover:bg-gray-50"
                  }`}
                >
                  <div>
                    <p className={`text-xs font-bold ${isActive ? acc.text : "text-gray-700"}`}>{def.label}</p>
                    <p className="text-[10px] text-gray-400 mt-0.5">
                      {hierarchyTotalCount(id)} items · {def.levels.length} levels
                    </p>
                  </div>
                  <span
                    className={`w-2 h-2 rounded-full shrink-0 ${isActive ? acc.solidBg : "bg-gray-200"}`}
                  />
                </button>
              );
            })}
          </div>
        </aside>

        {/* ── RIGHT: CRUD WORKSPACE ──────────────────────────────────── */}
        {activeSection === "flat" ? (
          <main className="flex-1 p-8 max-w-4xl mx-auto space-y-6">
            <div>
              <h2 className="text-lg font-black text-gray-900">{config.title}</h2>
              <p className="text-xs text-gray-400 mt-1">{config.description}</p>
            </div>

            {/* SUMMARY CARD */}
            <div className="bg-white p-5 rounded-xl shadow-sm w-fit min-w-[200px]">
              <p className="text-[11px] text-gray-400 font-bold uppercase tracking-wider">Total Items</p>
              <p className="text-2xl font-black text-gray-900 mt-1">{items.length}</p>
            </div>

            {/* ADD (CREATE) FORM */}
            <form
              onSubmit={handleAdd}
              className="bg-white p-5 rounded-xl shadow-sm flex flex-col sm:flex-row gap-3 sm:items-end"
            >
              <div className="flex-1">
                <label className="text-[10px] font-black uppercase text-gray-400 tracking-wider block mb-1.5">
                  {config.itemLabel}
                </label>
                <input
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  placeholder={`e.g. New ${config.itemLabel}`}
                  className="w-full bg-gray-50 rounded-lg p-2.5 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                />
              </div>
              <button
                type="submit"
                className={`${accent.bg} text-white text-xs font-bold px-5 py-2.5 rounded-lg hover:opacity-90 transition-all shrink-0`}
              >
                + Add {config.itemLabel}
              </button>
            </form>
            {formError && (
              <p className="text-[11px] font-bold text-red-500 -mt-3 px-1">{formError}</p>
            )}

            {/* TABLE (READ / UPDATE / DELETE) */}
            <div className="bg-white rounded-xl shadow-sm overflow-hidden">
              <table className="w-full text-xs">
                <thead>
                  <tr className="bg-gray-50 text-gray-400 text-[10px] font-black uppercase tracking-wider">
                    <th className="text-left px-5 py-3">{config.itemLabel}</th>
                    <th className="text-right px-5 py-3 w-40">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {items.length === 0 && (
                    <tr>
                      <td colSpan={2} className="px-5 py-10 text-center text-gray-400 font-medium">
                        No items yet — add the first {config.itemLabel.toLowerCase()} above.
                      </td>
                    </tr>
                  )}
                  {items.map((item) => {
                    const isEditing = editingId === item.id;
                    return (
                      <tr key={item.id} className="border-t border-gray-50 hover:bg-gray-50/60">
                        <td className="px-5 py-3 align-top">
                          {isEditing ? (
                            <>
                              <input
                                value={editName}
                                onChange={(e) => setEditName(e.target.value)}
                                className="w-full bg-gray-50 rounded-lg p-2 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                                autoFocus
                              />
                              {editError && (
                                <p className="text-[10px] font-bold text-red-500 mt-1">{editError}</p>
                              )}
                            </>
                          ) : (
                            <span className="font-bold text-gray-800">{item.name}</span>
                          )}
                        </td>
                        <td className="px-5 py-3 align-top">
                          <div className="flex justify-end gap-2">
                            {isEditing ? (
                              <>
                                <button
                                  onClick={() => saveEdit(item.id)}
                                  className="text-[11px] font-bold text-white bg-emerald-600 hover:bg-emerald-700 px-3 py-1.5 rounded-lg transition-all"
                                >
                                  Save
                                </button>
                                <button
                                  onClick={cancelEdit}
                                  className="text-[11px] font-bold text-gray-500 bg-gray-100 hover:bg-gray-200 px-3 py-1.5 rounded-lg transition-all"
                                >
                                  Cancel
                                </button>
                              </>
                            ) : (
                              <>
                                <button
                                  onClick={() => startEdit(item)}
                                  className={`text-[11px] font-bold ${accent.text} ${accent.lightBg} hover:opacity-80 px-3 py-1.5 rounded-lg transition-all`}
                                >
                                  Edit
                                </button>
                                <button
                                  onClick={() => setDeleteTarget(item)}
                                  className="text-[11px] font-bold text-red-500 bg-red-50 hover:bg-red-100 px-3 py-1.5 rounded-lg transition-all"
                                >
                                  Delete
                                </button>
                              </>
                            )}
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </main>
        ) : hierarchyDef && hLevelDef ? (
          <main className="flex-1 p-8 max-w-5xl mx-auto space-y-6">
            <div>
              <h2 className="text-lg font-black text-gray-900">{hierarchyDef.label} Classification</h2>
              <p className="text-xs text-gray-400 mt-1">
                {hierarchyDef.levels.map((l) => l.title).join(" → ")}. Each item has a name, a code, and
                (except the top level) a parent from the level above.
              </p>
            </div>

            {/* LEVEL TABS */}
            <div className="flex flex-wrap gap-2">
              {hierarchyDef.levels.map((level, idx) => {
                const isActive = hActiveLevelIndex === idx;
                return (
                  <button
                    key={level.key}
                    onClick={() => switchHLevel(idx)}
                    className={`flex items-center gap-2 px-4 py-2 rounded-xl text-xs font-bold transition-all ${
                      isActive
                        ? `${hAccent.solidBg} text-white`
                        : "bg-white text-gray-500 hover:bg-gray-100 shadow-sm"
                    }`}
                  >
                    <span
                      className={`w-4 h-4 rounded-full flex items-center justify-center text-[9px] ${
                        isActive ? "bg-white/20" : "bg-gray-100"
                      }`}
                    >
                      {idx + 1}
                    </span>
                    {level.title}
                    <span className={isActive ? "text-white/70" : "text-gray-300"}>
                      ({hData[hierarchyId as HierarchyId][level.key].length})
                    </span>
                  </button>
                );
              })}
            </div>

            {/* SUMMARY CARD */}
            <div className="bg-white p-5 rounded-xl shadow-sm w-fit min-w-[200px]">
              <p className="text-[11px] text-gray-400 font-bold uppercase tracking-wider">
                Total {hLevelDef.title} Items
              </p>
              <p className="text-2xl font-black text-gray-900 mt-1">{hItems.length}</p>
            </div>

            {/* ADD (CREATE) FORM */}
            <form
              onSubmit={handleHAdd}
              className="bg-white p-5 rounded-xl shadow-sm flex flex-col sm:flex-row gap-3 sm:items-end sm:flex-wrap"
            >
              <div className="flex-1 min-w-[160px]">
                <label className="text-[10px] font-black uppercase text-gray-400 tracking-wider block mb-1.5">
                  {hLevelDef.itemLabel} Name
                </label>
                <input
                  value={hNewName}
                  onChange={(e) => setHNewName(e.target.value)}
                  placeholder={`e.g. New ${hLevelDef.itemLabel}`}
                  className="w-full bg-gray-50 rounded-lg p-2.5 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                />
              </div>
              <div className="w-full sm:w-32">
                <label className="text-[10px] font-black uppercase text-gray-400 tracking-wider block mb-1.5">
                  Code
                </label>
                <input
                  value={hNewCode}
                  onChange={(e) => setHNewCode(e.target.value)}
                  placeholder="e.g. 1112"
                  className="w-full bg-gray-50 rounded-lg p-2.5 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                />
              </div>
              {hParentLevelDef && (
                <div className="w-full sm:w-64">
                  <label className="text-[10px] font-black uppercase text-gray-400 tracking-wider block mb-1.5">
                    Parent {hParentLevelDef.itemLabel}
                  </label>
                  <select
                    value={hNewParentId}
                    onChange={(e) => setHNewParentId(e.target.value)}
                    className="w-full bg-gray-50 rounded-lg p-2.5 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                  >
                    <option value="">Select {hParentLevelDef.itemLabel}…</option>
                    {hParentItems.map((p) => (
                      <option key={p.id} value={p.id}>
                        {p.code} — {p.name}
                      </option>
                    ))}
                  </select>
                </div>
              )}
              <button
                type="submit"
                className={`${hAccent.bg} text-white text-xs font-bold px-5 py-2.5 rounded-lg hover:opacity-90 transition-all shrink-0`}
              >
                + Add {hLevelDef.itemLabel}
              </button>
            </form>
            {hParentLevelDef && hParentItems.length === 0 && (
              <p className="text-[11px] font-bold text-amber-600 -mt-3 px-1">
                Add at least one {hParentLevelDef.itemLabel} first before creating a{" "}
                {hLevelDef.itemLabel.toLowerCase()}.
              </p>
            )}
            {hFormError && (
              <p className="text-[11px] font-bold text-red-500 -mt-3 px-1">{hFormError}</p>
            )}

            {/* TABLE (READ / UPDATE / DELETE) */}
            <div className="bg-white rounded-xl shadow-sm overflow-hidden">
              <table className="w-full text-xs">
                <thead>
                  <tr className="bg-gray-50 text-gray-400 text-[10px] font-black uppercase tracking-wider">
                    <th className="text-left px-5 py-3">{hLevelDef.itemLabel}</th>
                    <th className="text-left px-5 py-3 w-32">Code</th>
                    {hParentLevelDef && (
                      <th className="text-left px-5 py-3 w-64">Parent {hParentLevelDef.itemLabel}</th>
                    )}
                    <th className="text-right px-5 py-3 w-40">Actions</th>
                  </tr>
                </thead>
                <tbody>
                  {hItems.length === 0 && (
                    <tr>
                      <td
                        colSpan={hParentLevelDef ? 4 : 3}
                        className="px-5 py-10 text-center text-gray-400 font-medium"
                      >
                        No items yet — add the first {hLevelDef.itemLabel.toLowerCase()} above.
                      </td>
                    </tr>
                  )}
                  {hItems.map((item) => {
                    const isEditing = hEditingId === item.id;
                    return (
                      <tr key={item.id} className="border-t border-gray-50 hover:bg-gray-50/60">
                        <td className="px-5 py-3 align-top">
                          {isEditing ? (
                            <input
                              value={hEditName}
                              onChange={(e) => setHEditName(e.target.value)}
                              className="w-full bg-gray-50 rounded-lg p-2 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                              autoFocus
                            />
                          ) : (
                            <span className="font-bold text-gray-800">{item.name}</span>
                          )}
                        </td>
                        <td className="px-5 py-3 align-top">
                          {isEditing ? (
                            <input
                              value={hEditCode}
                              onChange={(e) => setHEditCode(e.target.value)}
                              className="w-full bg-gray-50 rounded-lg p-2 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                            />
                          ) : (
                            <span className="font-mono font-bold text-gray-600">{item.code}</span>
                          )}
                        </td>
                        {hParentLevelDef && (
                          <td className="px-5 py-3 align-top">
                            {isEditing ? (
                              <select
                                value={hEditParentId}
                                onChange={(e) => setHEditParentId(e.target.value)}
                                className="w-full bg-gray-50 rounded-lg p-2 text-xs font-semibold border border-gray-200 outline-none focus:ring-2 focus:ring-offset-0"
                              >
                                <option value="">Select {hParentLevelDef.itemLabel}…</option>
                                {hParentItems.map((p) => (
                                  <option key={p.id} value={p.id}>
                                    {p.code} — {p.name}
                                  </option>
                                ))}
                              </select>
                            ) : (
                              <span className="text-gray-500">{hParentDisplay(item)}</span>
                            )}
                          </td>
                        )}
                        <td className="px-5 py-3 align-top">
                          {isEditing && hEditError && (
                            <p className="text-[10px] font-bold text-red-500 mb-1 text-right">{hEditError}</p>
                          )}
                          <div className="flex justify-end gap-2">
                            {isEditing ? (
                              <>
                                <button
                                  onClick={() => saveHEdit(item.id)}
                                  className="text-[11px] font-bold text-white bg-emerald-600 hover:bg-emerald-700 px-3 py-1.5 rounded-lg transition-all"
                                >
                                  Save
                                </button>
                                <button
                                  onClick={cancelHEdit}
                                  className="text-[11px] font-bold text-gray-500 bg-gray-100 hover:bg-gray-200 px-3 py-1.5 rounded-lg transition-all"
                                >
                                  Cancel
                                </button>
                              </>
                            ) : (
                              <>
                                <button
                                  onClick={() => startHEdit(item)}
                                  className={`text-[11px] font-bold ${hAccent.text} ${hAccent.lightBg} hover:opacity-80 px-3 py-1.5 rounded-lg transition-all`}
                                >
                                  Edit
                                </button>
                                <button
                                  onClick={() => setHDeleteTarget(item)}
                                  className="text-[11px] font-bold text-red-500 bg-red-50 hover:bg-red-100 px-3 py-1.5 rounded-lg transition-all"
                                >
                                  Delete
                                </button>
                              </>
                            )}
                          </div>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </main>
        ) : null}
      </div>

      {/* DELETE CONFIRMATION MODAL — flat KPIs */}
      {deleteTarget && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 px-4">
          <div className="bg-white rounded-xl shadow-xl max-w-sm w-full p-6">
            <h3 className="text-sm font-black text-gray-900">Delete {config.itemLabel}?</h3>
            <p className="text-xs text-gray-500 mt-2">
              This will permanently remove{" "}
              <span className="font-bold text-gray-800">"{deleteTarget.name}"</span> from{" "}
              {config.title}. This action can't be undone.
            </p>
            <div className="flex justify-end gap-2 mt-6">
              <button
                onClick={() => setDeleteTarget(null)}
                className="text-xs font-bold text-gray-500 bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded-lg transition-all"
              >
                Cancel
              </button>
              <button
                onClick={handleDelete}
                className="text-xs font-bold text-white bg-red-600 hover:bg-red-700 px-4 py-2 rounded-lg transition-all"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}

      {/* DELETE CONFIRMATION MODAL — hierarchical KPIs (warns about cascading child deletes) */}
      {hDeleteTarget && hLevelDef && (
        <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50 px-4">
          <div className="bg-white rounded-xl shadow-xl max-w-sm w-full p-6">
            <h3 className="text-sm font-black text-gray-900">Delete {hLevelDef.itemLabel}?</h3>
            <p className="text-xs text-gray-500 mt-2">
              This will permanently remove{" "}
              <span className="font-bold text-gray-800">
                "{hDeleteTarget.code} — {hDeleteTarget.name}"
              </span>
              . This action can't be undone.
            </p>
            {hDeleteDescendantCount > 0 && (
              <p className="text-xs font-bold text-red-500 mt-2">
                It also has {hDeleteDescendantCount} item{hDeleteDescendantCount === 1 ? "" : "s"} nested
                underneath it across the lower levels — those will be deleted too.
              </p>
            )}
            <div className="flex justify-end gap-2 mt-6">
              <button
                onClick={() => setHDeleteTarget(null)}
                className="text-xs font-bold text-gray-500 bg-gray-100 hover:bg-gray-200 px-4 py-2 rounded-lg transition-all"
              >
                Cancel
              </button>
              <button
                onClick={handleHDelete}
                className="text-xs font-bold text-white bg-red-600 hover:bg-red-700 px-4 py-2 rounded-lg transition-all"
              >
                Delete
              </button>
            </div>
          </div>
        </div>
      )}

      {/* TOAST */}
      {toast && (
        <div className="fixed bottom-6 right-6 bg-gray-900 text-white text-xs font-bold px-4 py-3 rounded-xl shadow-lg z-50">
          {toast}
        </div>
      )}
    </div>
  );
}