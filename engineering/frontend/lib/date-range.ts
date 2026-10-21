export const toISODate = (date: Date) => date.toISOString().slice(0, 10);

// This module is evaluated once per server process,
// so a top-level `new Date()` would freeze "today" at server start.
export function getDataRangeEnd(): Date {
  return new Date();
}

export function getDataRangeStart(): Date {
  const end = getDataRangeEnd();
  return new Date(end.getFullYear() - 3, 0, 1);
}
