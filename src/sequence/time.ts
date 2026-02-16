import type { SeqFactor } from "../types/Sequence";

export type IntervalType =
  | "seconds"
  | "minutes"
  | "hours"
  | "days"
  | "weeks"
  | "months"
  | "years";

export const addInterval = (
  date: Date,
  step: number,
  intervalType: IntervalType,
): Date => {
  const result = new Date(date);

  switch (intervalType) {
    case "seconds":
      result.setSeconds(result.getSeconds() + step);
      break;
    case "minutes":
      result.setMinutes(result.getMinutes() + step);
      break;
    case "hours":
      result.setHours(result.getHours() + step);
      break;
    case "days":
      result.setDate(result.getDate() + step);
      break;
    case "weeks":
      result.setDate(result.getDate() + step * 7);
      break;
    case "months": {
      const originalDay = result.getDate();
      result.setMonth(result.getMonth() + step);
      // Handle month overflow (e.g., Jan 31 + 1 month should be Feb 28/29, not Mar 2/3)
      if (result.getDate() !== originalDay) {
        result.setDate(0); // Set to last day of previous month
      }
      break;
    }
    case "years": {
      const originalDay = result.getDate();
      const originalMonth = result.getMonth();
      result.setFullYear(result.getFullYear() + step);
      // Handle leap year edge case (Feb 29 + 1 year should be Feb 28)
      if (
        result.getMonth() !== originalMonth ||
        result.getDate() !== originalDay
      ) {
        result.setDate(0);
      }
      break;
    }
  }

  return result;
};

export const generateTimeSeries = (
  start: Date,
  end: Date,
  step: number,
  intervalType: IntervalType,
): Date[] => {
  const series: Date[] = [];
  let current = new Date(start);

  while (current <= end) {
    series.push(current);
    current = addInterval(current, step, intervalType);
  }

  return series;
};

export const sliceByTimeRange = (
  sequence: SeqFactor[],
  start: Date,
  end: Date,
): SeqFactor[] => {
  const startMs = start.getTime();
  const endMs = end.getTime();
  return sequence.filter((f) => {
    const t = f.time.getTime();
    return t >= startMs && t <= endMs;
  });
};
