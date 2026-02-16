import type { AccumCore, SeqFactor } from "../types/Sequence";
import { accumulateSequence } from "./core";

export const filterByTag = (
  sequence: SeqFactor[],
  tags: string[],
): SeqFactor[] => {
  const tagSet = new Set(tags);
  return sequence.filter((f) => tagSet.has(f.tag));
};

export const extractTags = (sequence: SeqFactor[]): string[] => {
  return Array.from(new Set(sequence.map((f) => f.tag)));
};

export const groupByTag = (
  sequence: SeqFactor[],
): Record<string, SeqFactor[]> => {
  const groups: Record<string, SeqFactor[]> = {};
  for (const f of sequence) {
    if (groups[f.tag]) {
      groups[f.tag].push(f);
    } else {
      groups[f.tag] = [f];
    }
  }
  return groups;
};

export const excludeByTag = (
  sequence: SeqFactor[],
  tags: string[],
): SeqFactor[] => {
  const tagSet = new Set(tags);
  return sequence.filter((f) => !tagSet.has(f.tag));
};

export const removeByTag = excludeByTag;

export const renameTag = (
  sequence: SeqFactor[],
  oldTag: string,
  newTag: string,
): SeqFactor[] => {
  return sequence.map((f) =>
    f.tag === oldTag ? { ...f, tag: newTag } : f,
  );
};

export const accumulateByTag = (
  sequence: SeqFactor[],
  baseValue: number,
  tag: string,
): AccumCore[] => {
  return accumulateSequence(filterByTag(sequence, [tag]), baseValue);
};
