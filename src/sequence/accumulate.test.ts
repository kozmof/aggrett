import { describe, test, expect } from "vitest";
import { accumulate } from "./accumulate";

describe("accumulate", () => {
  test("adds value with plus factor", () => {
    expect(accumulate("plus", 10, 5)).toBe(15);
  });

  test("subtracts value with minus factor", () => {
    expect(accumulate("minus", 10, 5)).toBe(5);
  });

  test("works with zero base value", () => {
    expect(accumulate("plus", 0, 100)).toBe(100);
    expect(accumulate("minus", 0, 100)).toBe(-100);
  });

  test("works with zero value", () => {
    expect(accumulate("plus", 50, 0)).toBe(50);
    expect(accumulate("minus", 50, 0)).toBe(50);
  });

  test("handles negative results", () => {
    expect(accumulate("minus", 3, 10)).toBe(-7);
  });
});
