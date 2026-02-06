import type { Factor, SeqFactor } from "../types/Sequence";

type IntervalType =
  | "seconds"
  | "minutes"
  | "hours"
  | "days"
  | "weeks"
  | "months"
  | "years";

const addInterval = (
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

const generateTimeSeries = (
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

export const createCycles = (
  start: Date,
  end: Date,
  step: number,
  intervalType: IntervalType,
  factor: Factor,
  value: number,
  tag: string,
  genId: () => string,
): SeqFactor[] => {
  const times = generateTimeSeries(start, end, step, intervalType);
  const sequence: SeqFactor[] = [];
  for (const time of times) {
    sequence.push({
      id: genId(),
      tag: tag,
      value: value,
      factor: factor,
      time: time,
    });
  }
  return sequence;
};
