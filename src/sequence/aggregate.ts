import type { Accum, Breakdown, Factor, SeqFactor } from "../types/Sequence";

const accumulate = (factor: Factor, prevValue: number, value: number) => {
  if (factor === "plus") {
    return prevValue + value;
  } else if (factor === "minus") {
    return prevValue - value;
  } else {
    throw new Error();
  }
};

const addBreakdown = (breakdown: Breakdown, factor: SeqFactor) => {
  if (factor.tag in breakdown) {
    breakdown[factor.tag] = {
      ids: [...breakdown[factor.tag].ids, factor.id],
      store: accumulate(
        factor.factor,
        breakdown[factor.tag].store,
        factor.value,
      )
    };
  } else {
    breakdown[factor.tag] = { ids: [factor.id], store: accumulate(factor.factor, 0, factor.value) };
  }
  return breakdown;
};

export const aggregate = (
  sequence: SeqFactor[],
  baseValue: number,
  filter: string[],
): Accum[] => {
  if (sequence.length < 1) return [];

  const seq = sequence;
  const accums: Accum[] = [];
  const sorted = seq.sort((a, b) => a.time.getTime() - b.time.getTime());

  const [firstFactor, ...restFactors] = sorted;

  let timePos = firstFactor.time;
  let accum: Accum = {
    ids: [firstFactor.id],
    tags: [firstFactor.tag],
    time: timePos,
    store: accumulate(firstFactor.factor, baseValue, firstFactor.value),
    breakdown: {
      [firstFactor.tag]: { ids: [firstFactor.id], store: accumulate(firstFactor.factor, 0, firstFactor.value) },
    },
  };

  for (const seqFactor of restFactors) {
    if (filter.length > 0 && !filter.includes(seqFactor.tag)) continue;

    const factor = seqFactor.factor;
    const value = seqFactor.value;
    const prevValue = accum.store;

    if (timePos === seqFactor.time) {
      accum = {
        ids: [...accum.ids, seqFactor.id],
        tags: Array.from(new Set([...accum.tags, seqFactor.tag])),
        time: timePos,
        store: accumulate(factor, prevValue, value),
        breakdown: addBreakdown(accum.breakdown, seqFactor),
      };
    } else {
      accums.push(accum);
      const newTimePos = seqFactor.time;
      accum = {
        ids: [seqFactor.id],
        tags: [seqFactor.tag],
        time: newTimePos,
        store: accumulate(factor, prevValue, value),
        breakdown: {
          [seqFactor.tag]: { ids: [seqFactor.id], store: accumulate(seqFactor.factor, 0, seqFactor.value) },
        },
      };
      timePos = newTimePos;
    }
  }
  if (accum !== null) {
    accums.push(accum);
  }
  return accums;
};
