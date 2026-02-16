export type { Factor, SeqFactor, Breakdown, AccumCore, Accum } from "./types/Sequence";
export type { IntervalType } from "./sequence/time";
export { addInterval, generateTimeSeries, sliceByTimeRange } from "./sequence/time";
export { accumulate } from "./sequence/accumulate";
export { accumulateSequence } from "./sequence/core";
export { createCycles } from "./sequence/cycles";
export { aggregate } from "./sequence/aggregate";
export {
  insertFactor,
  removeFactor,
  updateFactor,
  mergeSequences,
} from "./sequence/factors";
