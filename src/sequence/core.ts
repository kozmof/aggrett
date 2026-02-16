import type { AccumCore, SeqFactor } from "../types/Sequence";
import { accumulate } from "./accumulate";

export const accumulateSequence = (
  sequence: SeqFactor[],
  baseValue: number,
): AccumCore[] => {
  if (sequence.length < 1) return [];

  const sorted = [...sequence].sort(
    (a, b) => a.time.getTime() - b.time.getTime(),
  );

  const [firstFactor, ...restFactors] = sorted;

  const accums: AccumCore[] = [];
  let timePos = firstFactor.time;
  let accum: AccumCore = {
    ids: [firstFactor.id],
    time: timePos,
    store: accumulate(firstFactor.factor, baseValue, firstFactor.value),
  };

  for (const seqFactor of restFactors) {
    const prevValue = accum.store;

    if (timePos.getTime() === seqFactor.time.getTime()) {
      accum = {
        ids: [...accum.ids, seqFactor.id],
        time: timePos,
        store: accumulate(seqFactor.factor, prevValue, seqFactor.value),
      };
    } else {
      accums.push(accum);
      timePos = seqFactor.time;
      accum = {
        ids: [seqFactor.id],
        time: timePos,
        store: accumulate(seqFactor.factor, prevValue, seqFactor.value),
      };
    }
  }
  accums.push(accum);
  return accums;
};
