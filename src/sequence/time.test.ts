import { describe, test, expect } from "vitest";
import type { SeqFactor } from "../types/Sequence";
import { addInterval, generateTimeSeries, sliceByTimeRange } from "./time";

const makeFactor = (
  id: string,
  tag: string,
  time: string,
  value: number,
  factor: "plus" | "minus" = "plus",
): SeqFactor => ({
  id,
  tag,
  time: new Date(time),
  value,
  factor,
});

describe("addInterval", () => {
  test("adds seconds", () => {
    const date = new Date("2024-01-01T10:00:00");
    const result = addInterval(date, 30, "seconds");
    expect(result).toEqual(new Date("2024-01-01T10:00:30"));
  });

  test("adds minutes", () => {
    const date = new Date("2024-01-01T10:00:00");
    const result = addInterval(date, 15, "minutes");
    expect(result).toEqual(new Date("2024-01-01T10:15:00"));
  });

  test("adds hours", () => {
    const date = new Date("2024-01-01T10:00:00");
    const result = addInterval(date, 3, "hours");
    expect(result).toEqual(new Date("2024-01-01T13:00:00"));
  });

  test("adds days", () => {
    const date = new Date("2024-01-01");
    const result = addInterval(date, 5, "days");
    expect(result).toEqual(new Date("2024-01-06"));
  });

  test("adds weeks", () => {
    const date = new Date("2024-01-01");
    const result = addInterval(date, 2, "weeks");
    expect(result).toEqual(new Date("2024-01-15"));
  });

  test("adds months", () => {
    const date = new Date("2024-01-15");
    const result = addInterval(date, 1, "months");
    expect(result).toEqual(new Date("2024-02-15"));
  });

  test("handles month overflow (Jan 31 + 1 month)", () => {
    const date = new Date("2024-01-31");
    const result = addInterval(date, 1, "months");
    expect(result).toEqual(new Date("2024-02-29")); // 2024 is a leap year
  });

  test("adds years", () => {
    const date = new Date("2024-06-01");
    const result = addInterval(date, 1, "years");
    expect(result).toEqual(new Date("2025-06-01"));
  });

  test("handles leap year (Feb 29 + 1 year)", () => {
    const date = new Date("2024-02-29");
    const result = addInterval(date, 1, "years");
    expect(result).toEqual(new Date("2025-02-28"));
  });

  test("does not mutate the original date", () => {
    const date = new Date("2024-01-01");
    addInterval(date, 5, "days");
    expect(date).toEqual(new Date("2024-01-01"));
  });
});

describe("generateTimeSeries", () => {
  test("generates daily series", () => {
    const start = new Date("2024-01-01");
    const end = new Date("2024-01-03");
    const result = generateTimeSeries(start, end, 1, "days");

    expect(result).toHaveLength(3);
    expect(result[0]).toEqual(new Date("2024-01-01"));
    expect(result[1]).toEqual(new Date("2024-01-02"));
    expect(result[2]).toEqual(new Date("2024-01-03"));
  });

  test("generates with step > 1", () => {
    const start = new Date("2024-01-01");
    const end = new Date("2024-01-10");
    const result = generateTimeSeries(start, end, 3, "days");

    expect(result).toHaveLength(4);
    expect(result[0]).toEqual(new Date("2024-01-01"));
    expect(result[1]).toEqual(new Date("2024-01-04"));
    expect(result[2]).toEqual(new Date("2024-01-07"));
    expect(result[3]).toEqual(new Date("2024-01-10"));
  });

  test("returns single item when start equals end", () => {
    const date = new Date("2024-01-01");
    const result = generateTimeSeries(date, date, 1, "days");

    expect(result).toHaveLength(1);
    expect(result[0]).toEqual(date);
  });

  test("returns empty when start > end", () => {
    const result = generateTimeSeries(
      new Date("2024-02-01"),
      new Date("2024-01-01"),
      1,
      "days",
    );

    expect(result).toHaveLength(0);
  });
});

describe("sliceByTimeRange", () => {
  const seq: SeqFactor[] = [
    makeFactor("a", "rent", "2024-01-01", 100),
    makeFactor("b", "food", "2024-02-01", 50),
    makeFactor("c", "rent", "2024-03-01", 100),
    makeFactor("d", "food", "2024-04-01", 50),
  ];

  test("returns factors within the time range inclusive", () => {
    const result = sliceByTimeRange(
      seq,
      new Date("2024-02-01"),
      new Date("2024-03-01"),
    );

    expect(result).toHaveLength(2);
    expect(result.map((f) => f.id)).toEqual(["b", "c"]);
  });

  test("includes factors exactly at boundaries", () => {
    const result = sliceByTimeRange(
      seq,
      new Date("2024-01-01"),
      new Date("2024-04-01"),
    );

    expect(result).toHaveLength(4);
  });

  test("returns empty array when no factors in range", () => {
    const result = sliceByTimeRange(
      seq,
      new Date("2025-01-01"),
      new Date("2025-12-31"),
    );

    expect(result).toHaveLength(0);
  });

  test("does not mutate the original array", () => {
    sliceByTimeRange(seq, new Date("2024-02-01"), new Date("2024-03-01"));

    expect(seq).toHaveLength(4);
  });

  test("works on an empty sequence", () => {
    const result = sliceByTimeRange(
      [],
      new Date("2024-01-01"),
      new Date("2024-12-31"),
    );

    expect(result).toHaveLength(0);
  });
});
