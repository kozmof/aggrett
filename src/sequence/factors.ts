import type { Factor, SeqFactor } from "../types/Sequence";

export const insertFactor = (
  sequence: SeqFactor[],
  tag: string,
  time: Date,
  value: number,
  factor: Factor,
  genId: () => string,
): SeqFactor[] => {
  return [...sequence, { id: genId(), tag, time, value, factor }];
};

export const removeFactor = (
  sequence: SeqFactor[],
  ids: string[],
): SeqFactor[] => {
  const idSet = new Set(ids);
  return sequence.filter((f) => !idSet.has(f.id));
};

export const updateFactor = (
  sequence: SeqFactor[],
  id: string,
  fields: Partial<Pick<SeqFactor, "tag" | "time" | "value" | "factor">>,
): SeqFactor[] => {
  return sequence.map((f) => (f.id === id ? { ...f, ...fields } : f));
};

export const mergeSequences = (
  a: SeqFactor[],
  b: SeqFactor[],
): SeqFactor[] => {
  return [...a, ...b];
};
