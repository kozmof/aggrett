export type { Factor, SeqFactor, Breakdown, Accum } from "./types/Sequence";
export { createCycles } from "./sequence/cycles";
export { aggregate } from "./sequence/aggregate";
export {
  insertFactor,
  removeFactor,
  updateFactor,
  mergeSequences,
  sliceByTimeRange,
} from "./sequence/factors";
