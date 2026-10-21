import { describe, expect, it } from "vitest";
import {
  isUsableDateParam,
  parseEnumParam,
  parseHierarchyIds,
  serializeHierarchyIds,
} from "./url-params";

const LEVELS = ["major-group", "sub-major-group", "minor-group"] as const;
const reader = (query: string) => new URLSearchParams(query);

describe("isUsableDateParam", () => {
  it("accepts ISO dates and the empty string, rejects the rest", () => {
    expect(isUsableDateParam("2026-01-31")).toBe(true);
    expect(isUsableDateParam("")).toBe(true);
    expect(isUsableDateParam(null)).toBe(false);
    expect(isUsableDateParam("31-01-2026")).toBe(false);
  });

  it("ignores surrounding whitespace", () => {
    expect(isUsableDateParam("      ")).toBe(true);
    expect(isUsableDateParam(" 2026-01-31 ")).toBe(true);
    expect(isUsableDateParam(" 31-01-2026 ")).toBe(false);
  });
});

describe("parseEnumParam", () => {
  it("returns allowed values and falls back otherwise", () => {
    expect(parseEnumParam("table", ["chart", "table"], "chart")).toBe("table");
    expect(parseEnumParam("grid", ["chart", "table"], "chart")).toBe("chart");
    expect(parseEnumParam(null, ["chart", "table"], "chart")).toBe("chart");
  });
});

describe("parseHierarchyIds", () => {
  it("reads one id per level in order", () => {
    expect(parseHierarchyIds(reader("major-group=2&sub-major-group=11"), LEVELS)).toEqual([2, 11]);
  });

  it("stops at the first missing level", () => {
    expect(parseHierarchyIds(reader("major-group=2&minor-group=43"), LEVELS)).toEqual([2]);
  });

  it("stops at a non-numeric id", () => {
    expect(parseHierarchyIds(reader("major-group=abc&sub-major-group=1"), LEVELS)).toEqual([]);
    expect(parseHierarchyIds(reader("major-group=-1"), LEVELS)).toEqual([]);
  });
});

describe("serializeHierarchyIds", () => {
  it("sets selected levels and clears deeper ones", () => {
    expect(serializeHierarchyIds([2, 11], LEVELS)).toEqual({
      "major-group": "2",
      "sub-major-group": "11",
      "minor-group": null,
    });
  });

  it("round-trips through parseHierarchyIds", () => {
    const params = new URLSearchParams();
    for (const [k, v] of Object.entries(serializeHierarchyIds([4, 5, 6], LEVELS))) {
      if (v !== null) params.set(k, v);
    }
    expect(parseHierarchyIds(params, LEVELS)).toEqual([4, 5, 6]);
  });
});
