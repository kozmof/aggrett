import type { Factor } from "../types/Sequence";

export const accumulate = (
  factor: Factor,
  prevValue: number,
  value: number,
): number => {
  if (factor === "plus") {
    return prevValue + value;
  } else if (factor === "minus") {
    return prevValue - value;
  } else {
    const _exhaustive: never = factor;
    throw new Error(`Unknown factor type: ${_exhaustive}`);
  }
};
