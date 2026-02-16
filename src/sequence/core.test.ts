import { describe, test, expect } from "vitest";
import type { SeqFactor } from "../types/Sequence";
import { accumulateSequence } from "./core";

describe("accumulateSequence", () => {
  test("returns empty array for empty sequence", () => {
    expect(accumulateSequence([], 0)).toEqual([]);
  });

  test("accumulates a single factor", () => {
    const sequence: SeqFactor[] = [
      { id: "1", tag: "a", time: new Date("2024-01-01"), factor: "plus", value: 10 },
    ];
    const result = accumulateSequence(sequence, 100);

    expect(result).toHaveLength(1);
    expect(result[0].store).toBe(110);
    expect(result[0].ids).toEqual(["1"]);
  });

  test("groups factors at the same timestamp", () => {
    const time = new Date("2024-01-01");
    const sequence: SeqFactor[] = [
      { id: "1", tag: "a", time, factor: "plus", value: 10 },
      { id: "2", tag: "b", time, factor: "minus", value: 3 },
    ];
    const result = accumulateSequence(sequence, 0);

    expect(result).toHaveLength(1);
    expect(result[0].store).toBe(7);
    expect(result[0].ids).toEqual(["1", "2"]);
  });

  test("carries running total across time points", () => {
    const yesterday = new Date("2024-01-01");
    const today = new Date("2024-01-02");
    const tomorrow = new Date("2024-01-03");

    const sequence: SeqFactor[] = [
      { id: "1", tag: "a", time: yesterday, factor: "plus", value: 10 },
      { id: "2", tag: "a", time: today, factor: "plus", value: 5 },
      { id: "3", tag: "b", time: tomorrow, factor: "minus", value: 3 },
    ];
    const result = accumulateSequence(sequence, 100);

    expect(result).toHaveLength(3);
    expect(result[0].store).toBe(110);
    expect(result[1].store).toBe(115);
    expect(result[2].store).toBe(112);
  });

  test("sorts by time regardless of input order", () => {
    const earlier = new Date("2024-01-01");
    const later = new Date("2024-01-02");

    const sequence: SeqFactor[] = [
      { id: "2", tag: "a", time: later, factor: "plus", value: 5 },
      { id: "1", tag: "a", time: earlier, factor: "plus", value: 10 },
    ];
    const result = accumulateSequence(sequence, 0);

    expect(result).toHaveLength(2);
    expect(result[0].ids).toEqual(["1"]);
    expect(result[0].store).toBe(10);
    expect(result[1].ids).toEqual(["2"]);
    expect(result[1].store).toBe(15);
  });

  test("does not mutate input sequence", () => {
    const sequence: SeqFactor[] = [
      { id: "1", tag: "a", time: new Date("2024-01-02"), factor: "plus", value: 10 },
      { id: "2", tag: "a", time: new Date("2024-01-01"), factor: "plus", value: 5 },
    ];
    const original = [...sequence];
    accumulateSequence(sequence, 0);

    expect(sequence[0].id).toBe(original[0].id);
    expect(sequence[1].id).toBe(original[1].id);
  });

  test("result has no breakdown or tags fields", () => {
    const sequence: SeqFactor[] = [
      { id: "1", tag: "a", time: new Date("2024-01-01"), factor: "plus", value: 10 },
    ];
    const result = accumulateSequence(sequence, 0);

    expect(result[0]).not.toHaveProperty("breakdown");
    expect(result[0]).not.toHaveProperty("tags");
  });
});
