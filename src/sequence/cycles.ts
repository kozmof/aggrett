import type { Factor, SeqFactor } from "../types/Sequence";

type IntervalType =
  | "seconds"
  | "minutes"
  | "hours"
  | "days"
  | "weeks"
  | "months"
  | "years";

const generateTimeSeries = (
  start: Date,
  end: Date,
  step: number,
  intervalType: IntervalType,
): Date[] => {
  const series: Date[] = [];
  let current = new Date(start);

  while (current <= end) {
    series.push(new Date(current));

    switch (intervalType) {
      case "seconds":
        current.setSeconds(current.getSeconds() + step);
        break;
      case "minutes":
        current.setMinutes(current.getMinutes() + step);
        break;
      case "hours":
        current.setHours(current.getHours() + step);
        break;
      case "days":
        current.setDate(current.getDate() + step);
        break;
      case "weeks":
        current.setDate(current.getDate() + step * 7);
        break;
      case "months":
        current.setMonth(current.getMonth() + step);
        break;
      case "years":
        current.setFullYear(current.getFullYear() + step);
        break;
    }
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
