"use client";

import { useCallback, useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import { getDataRangeStart, getDataRangeEnd, toISODate } from "@/lib/date-range";
import { URL_PARAMS, isDateString, isUsableDateParam } from "@/lib/url-params";
import { useUpdateSearchParams } from "@/hooks/use-update-search-params";

interface UseUrlDateRangeResult {
  fromDate: string;
  toDate: string;
  setFromDate: (value: string) => void;
  setToDate: (value: string) => void;
  rangeInvalid: boolean;
}

/**
 * The date range lives in the URL (`?from=YYYY-MM-DD&to=YYYY-MM-DD`) so any
 * snapshot of a page can be shared as a link, and a range in the URL always
 * wins. When the URL carries no (valid) range, the last range used in this
 * browser session (sessionStorage) or the current-year default is used, and
 * is then written back to the URL so the address bar is always a complete,
 * shareable snapshot.
 */
export function useUrlDateRange(storageKey: string): UseUrlDateRangeResult {
  const searchParams = useSearchParams();
  const updateSearchParams = useUpdateSearchParams();
  const rawFrom = searchParams.get(URL_PARAMS.from);
  const rawTo = searchParams.get(URL_PARAMS.to);

  const [fallback, setFallback] = useState(() => ({
    from: `${getDataRangeEnd().getFullYear()}-01-01`,
    to: toISODate(getDataRangeEnd()),
  }));
  const [restored, setRestored] = useState(false);

  const fromDate = isUsableDateParam(rawFrom) ? rawFrom.trim() : fallback.from;
  const toDate = isUsableDateParam(rawTo) ? rawTo.trim() : fallback.to;

  // Restoring from sessionStorage (a client-only API, unavailable during SSR)
  // has to happen post-mount so the first client render matches the server's.
  useEffect(() => {
    try {
      const saved = sessionStorage.getItem(storageKey);
      if (saved) {
        const parsed = JSON.parse(saved);
        if (isDateString(parsed.from) && isDateString(parsed.to)) {
          const from = parsed.from.trim();
          const to = parsed.to.trim();
          if (
            from <= to &&
            from >= toISODate(getDataRangeStart()) &&
            to <= toISODate(getDataRangeEnd())
          ) {
            // eslint-disable-next-line react-hooks/set-state-in-effect -- client-only storage, see note above
            setFallback({ from, to });
          }
        }
      }
    } catch {
      // ignore corrupt storage — defaults stay in place
    }
    setRestored(true);
  }, [storageKey]);

  // Make the address bar reflect the effective range once the fallback is known.
  useEffect(() => {
    if (!restored) return;
    const updates: Record<string, string> = {};
    if (!isUsableDateParam(rawFrom)) updates[URL_PARAMS.from] = fromDate;
    if (!isUsableDateParam(rawTo)) updates[URL_PARAMS.to] = toDate;
    if (Object.keys(updates).length > 0) updateSearchParams(updates);
  }, [restored, rawFrom, rawTo, fromDate, toDate, updateSearchParams]);

  useEffect(() => {
    if (!restored) return;
    try {
      sessionStorage.setItem(storageKey, JSON.stringify({ from: fromDate, to: toDate }));
    } catch {
      // storage unavailable — persistence is best-effort
    }
  }, [restored, storageKey, fromDate, toDate]);

  const setFromDate = useCallback(
    (value: string) => updateSearchParams({ [URL_PARAMS.from]: value }),
    [updateSearchParams],
  );
  const setToDate = useCallback(
    (value: string) => updateSearchParams({ [URL_PARAMS.to]: value }),
    [updateSearchParams],
  );

  return { fromDate, toDate, setFromDate, setToDate, rangeInvalid: fromDate > toDate };
}
