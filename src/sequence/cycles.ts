import type { Factor, SeqFactor } from "../types/Sequence";
import type { IntervalType } from "./time";
import { generateTimeSeries } from "./time";

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
