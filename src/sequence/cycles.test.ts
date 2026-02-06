import { describe, test, expect } from "vitest";
import { createCycles } from "./cycles";

describe("createCycles", () => {
  const mockGenId = (() => {
    let counter = 0;
    return () => `id-${++counter}`;
  })();

  test("creates daily cycles", () => {
    const start = new Date("2024-01-01");
    const end = new Date("2024-01-03");
    const result = createCycles(
      start,
      end,
      1,
      "days",
      "plus",
      100,
      "daily",
      mockGenId,
    );

    expect(result).toHaveLength(3);
    expect(result[0].time).toEqual(new Date("2024-01-01"));
    expect(result[1].time).toEqual(new Date("2024-01-02"));
    expect(result[2].time).toEqual(new Date("2024-01-03"));
    result.forEach((item) => {
      expect(item.factor).toBe("plus");
      expect(item.value).toBe(100);
      expect(item.tag).toBe("daily");
    });
  });

  test("creates weekly cycles", () => {
    const start = new Date("2024-01-01");
    const end = new Date("2024-01-22");
    const result = createCycles(
      start,
      end,
      1,
      "weeks",
      "minus",
      50,
      "weekly",
      mockGenId,
    );

    expect(result).toHaveLength(4);
    expect(result[0].time).toEqual(new Date("2024-01-01"));
    expect(result[1].time).toEqual(new Date("2024-01-08"));
    expect(result[2].time).toEqual(new Date("2024-01-15"));
    expect(result[3].time).toEqual(new Date("2024-01-22"));
    result.forEach((item) => {
      expect(item.factor).toBe("minus");
      expect(item.value).toBe(50);
      expect(item.tag).toBe("weekly");
    });
  });

  test("creates monthly cycles", () => {
    const start = new Date("2024-01-15");
    const end = new Date("2024-04-15");
    const result = createCycles(
      start,
      end,
      1,
      "months",
      "plus",
      200,
      "monthly",
      mockGenId,
    );

    expect(result).toHaveLength(4);
    expect(result[0].time).toEqual(new Date("2024-01-15"));
    expect(result[1].time).toEqual(new Date("2024-02-15"));
    expect(result[2].time).toEqual(new Date("2024-03-15"));
    expect(result[3].time).toEqual(new Date("2024-04-15"));
  });

  test("creates yearly cycles", () => {
    const start = new Date("2020-06-01");
    const end = new Date("2023-06-01");
    const result = createCycles(
      start,
      end,
      1,
      "years",
      "plus",
      1000,
      "yearly",
      mockGenId,
    );

    expect(result).toHaveLength(4);
    expect(result[0].time).toEqual(new Date("2020-06-01"));
    expect(result[1].time).toEqual(new Date("2021-06-01"));
    expect(result[2].time).toEqual(new Date("2022-06-01"));
    expect(result[3].time).toEqual(new Date("2023-06-01"));
  });

  test("creates hourly cycles", () => {
    const start = new Date("2024-01-01T10:00:00");
    const end = new Date("2024-01-01T13:00:00");
    const result = createCycles(
      start,
      end,
      1,
      "hours",
      "minus",
      25,
      "hourly",
      mockGenId,
    );

    expect(result).toHaveLength(4);
    expect(result[0].time).toEqual(new Date("2024-01-01T10:00:00"));
    expect(result[1].time).toEqual(new Date("2024-01-01T11:00:00"));
    expect(result[2].time).toEqual(new Date("2024-01-01T12:00:00"));
    expect(result[3].time).toEqual(new Date("2024-01-01T13:00:00"));
  });

  test("creates minute cycles", () => {
    const start = new Date("2024-01-01T10:00:00");
    const end = new Date("2024-01-01T10:03:00");
    const result = createCycles(
      start,
      end,
      1,
      "minutes",
      "plus",
      5,
      "minute",
      mockGenId,
    );

    expect(result).toHaveLength(4);
    expect(result[0].time).toEqual(new Date("2024-01-01T10:00:00"));
    expect(result[1].time).toEqual(new Date("2024-01-01T10:01:00"));
    expect(result[2].time).toEqual(new Date("2024-01-01T10:02:00"));
    expect(result[3].time).toEqual(new Date("2024-01-01T10:03:00"));
  });

  test("creates second cycles", () => {
    const start = new Date("2024-01-01T10:00:00");
    const end = new Date("2024-01-01T10:00:02");
    const result = createCycles(
      start,
      end,
      1,
      "seconds",
      "plus",
      1,
      "second",
      mockGenId,
    );

    expect(result).toHaveLength(3);
    expect(result[0].time).toEqual(new Date("2024-01-01T10:00:00"));
    expect(result[1].time).toEqual(new Date("2024-01-01T10:00:01"));
    expect(result[2].time).toEqual(new Date("2024-01-01T10:00:02"));
  });

  test("handles step greater than 1", () => {
    const start = new Date("2024-01-01");
    const end = new Date("2024-01-10");
    const result = createCycles(
      start,
      end,
      3,
      "days",
      "plus",
      10,
      "every3days",
      mockGenId,
    );

    expect(result).toHaveLength(4);
    expect(result[0].time).toEqual(new Date("2024-01-01"));
    expect(result[1].time).toEqual(new Date("2024-01-04"));
    expect(result[2].time).toEqual(new Date("2024-01-07"));
    expect(result[3].time).toEqual(new Date("2024-01-10"));
  });

  test("returns single item when start equals end", () => {
    const date = new Date("2024-01-01");
    const result = createCycles(
      date,
      date,
      1,
      "days",
      "plus",
      50,
      "single",
      mockGenId,
    );

    expect(result).toHaveLength(1);
    expect(result[0].time).toEqual(date);
  });

  test("generates unique ids for each item", () => {
    let idCounter = 0;
    const uniqueGenId = () => `unique-${++idCounter}`;

    const start = new Date("2024-01-01");
    const end = new Date("2024-01-03");
    const result = createCycles(
      start,
      end,
      1,
      "days",
      "plus",
      10,
      "test",
      uniqueGenId,
    );

    const ids = result.map((item) => item.id);
    const uniqueIds = new Set(ids);
    expect(uniqueIds.size).toBe(ids.length);
  });

  test("handles month overflow correctly (Jan 31 + 1 month = Feb 28/29)", () => {
    let idCounter = 0;
    const genId = () => `id-${++idCounter}`;

    const start = new Date("2024-01-31"); // 2024 is a leap year
    const end = new Date("2024-04-30");
    const result = createCycles(
      start,
      end,
      1,
      "months",
      "plus",
      100,
      "monthly",
      genId,
    );

    expect(result).toHaveLength(4);
    expect(result[0].time).toEqual(new Date("2024-01-31"));
    expect(result[1].time).toEqual(new Date("2024-02-29")); // Leap year, last day of Feb
    expect(result[2].time).toEqual(new Date("2024-03-29")); // Stays on 29th
    expect(result[3].time).toEqual(new Date("2024-04-29"));
  });

  test("handles leap year edge case (Feb 29 + 1 year = Feb 28)", () => {
    let idCounter = 0;
    const genId = () => `id-${++idCounter}`;

    const start = new Date("2024-02-29"); // Leap year
    const end = new Date("2026-02-28");
    const result = createCycles(
      start,
      end,
      1,
      "years",
      "plus",
      100,
      "yearly",
      genId,
    );

    expect(result).toHaveLength(3);
    expect(result[0].time).toEqual(new Date("2024-02-29"));
    expect(result[1].time).toEqual(new Date("2025-02-28")); // Non-leap year
    expect(result[2].time).toEqual(new Date("2026-02-28"));
  });
});
