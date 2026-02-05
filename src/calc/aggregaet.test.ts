import { describe, test, expect } from "vitest";
import type { SeqFactor } from "../types/Sequence";
import { aggregate } from "./aggregate";
describe("basic test", () => {
  test("aggregate", () => {
    const today = new Date();
    const yesterday = new Date();
    yesterday.setDate(yesterday.getDate() - 1);
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);

    const sequence: SeqFactor[] = [
      {
        tag: "test",
        time: today,
        factor: "plus",
        value: 4,
      },
      {
        tag: "test",
        time: today,
        factor: "minus",
        value: 2,
      },
      {
        tag: "test",
        time: yesterday,
        factor: "plus",
        value: 10,
      },
      {
        tag: "test",
        time: tomorrow,
        factor: "minus",
        value: 3,
      },
      {
        tag: "other tag",
        time: tomorrow,
        factor: "minus",
        value: 5,
      },
    ];
    const accum = aggregate(sequence, 10, []);
    expect(accum[0].store).toBe(20);
    expect(accum[0].breakdown["test"]).toBe(10);
    expect(accum[1].store).toBe(22);
    expect(accum[1].breakdown["test"]).toBe(2);
    expect(accum[2].store).toBe(14);
    expect(accum[2].breakdown["test"]).toBe(-3);
    expect(accum[2].breakdown["other tag"]).toBe(-5);
  });
});
