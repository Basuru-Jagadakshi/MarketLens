"use client";

import { useCallback } from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";

export type SearchParamUpdates = Record<string, string | null>;

/**
 * Returns a function that merges `updates` into the current query string and
 * navigates there with `router.replace` (no new history entry, no scroll
 * reset). A `null` value removes the param. This is the one place that writes
 * page state to the URL.
 */
export function useUpdateSearchParams(): (updates: SearchParamUpdates) => void {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  return useCallback(
    (updates) => {
      const params = new URLSearchParams(searchParams.toString());
      for (const [key, value] of Object.entries(updates)) {
        if (value === null) params.delete(key);
        else params.set(key, value);
      }
      const query = params.toString();
      router.replace(query ? `${pathname}?${query}` : pathname, { scroll: false });
    },
    [router, pathname, searchParams],
  );
}
