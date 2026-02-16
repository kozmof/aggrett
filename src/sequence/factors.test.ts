import { describe, test, expect } from "vitest";
import type { SeqFactor } from "../types/Sequence";
import {
  insertFactor,
  removeFactor,
  updateFactor,
  mergeSequences,
} from "./factors";

const mockGenId = (() => {
  let counter = 0;
  return () => `fac-${++counter}`;
})();

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

describe("insertFactor", () => {
  test("appends a new factor to the sequence", () => {
    const seq: SeqFactor[] = [makeFactor("a", "rent", "2024-01-01", 100)];
    const result = insertFactor(
      seq,
      "food",
      new Date("2024-02-01"),
      50,
      "minus",
      mockGenId,
    );

    expect(result).toHaveLength(2);
    expect(result[1].tag).toBe("food");
    expect(result[1].value).toBe(50);
    expect(result[1].factor).toBe("minus");
    expect(result[1].id).toMatch(/^fac-/);
  });

  test("does not mutate the original array", () => {
    const seq: SeqFactor[] = [makeFactor("a", "rent", "2024-01-01", 100)];
    insertFactor(seq, "food", new Date("2024-02-01"), 50, "plus", mockGenId);

    expect(seq).toHaveLength(1);
  });

  test("works on an empty sequence", () => {
    const result = insertFactor(
      [],
      "salary",
      new Date("2024-03-01"),
      3000,
      "plus",
      mockGenId,
    );

    expect(result).toHaveLength(1);
    expect(result[0].tag).toBe("salary");
  });
});

describe("removeFactor", () => {
  const seq: SeqFactor[] = [
    makeFactor("a", "rent", "2024-01-01", 100),
    makeFactor("b", "food", "2024-01-15", 50),
    makeFactor("c", "rent", "2024-02-01", 100),
  ];

  test("removes a single factor by ID", () => {
    const result = removeFactor(seq, ["b"]);

    expect(result).toHaveLength(2);
    expect(result.map((f) => f.id)).toEqual(["a", "c"]);
  });

  test("removes multiple factors by ID", () => {
    const result = removeFactor(seq, ["a", "c"]);

    expect(result).toHaveLength(1);
    expect(result[0].id).toBe("b");
  });

  test("returns unchanged copy when ID not found", () => {
    const result = removeFactor(seq, ["nonexistent"]);

    expect(result).toHaveLength(3);
    expect(result).not.toBe(seq);
  });

  test("does not mutate the original array", () => {
    removeFactor(seq, ["a"]);

    expect(seq).toHaveLength(3);
  });

  test("works on an empty sequence", () => {
    const result = removeFactor([], ["a"]);

    expect(result).toHaveLength(0);
  });
});

describe("updateFactor", () => {
  const seq: SeqFactor[] = [
    makeFactor("a", "rent", "2024-01-01", 100, "minus"),
    makeFactor("b", "food", "2024-01-15", 50, "minus"),
  ];

  test("updates a single field", () => {
    const result = updateFactor(seq, "a", { value: 200 });

    expect(result[0].value).toBe(200);
    expect(result[0].tag).toBe("rent");
    expect(result[0].factor).toBe("minus");
  });

  test("updates multiple fields", () => {
    const result = updateFactor(seq, "b", {
      tag: "dining",
      factor: "plus",
      value: 75,
    });

    expect(result[1].tag).toBe("dining");
    expect(result[1].factor).toBe("plus");
    expect(result[1].value).toBe(75);
  });

  test("returns unchanged copy when ID not found", () => {
    const result = updateFactor(seq, "nonexistent", { value: 999 });

    expect(result).toHaveLength(2);
    expect(result[0].value).toBe(100);
    expect(result[1].value).toBe(50);
    expect(result).not.toBe(seq);
  });

  test("does not mutate the original array", () => {
    updateFactor(seq, "a", { value: 999 });

    expect(seq[0].value).toBe(100);
  });
});

describe("mergeSequences", () => {
  test("concatenates two sequences", () => {
    const a: SeqFactor[] = [makeFactor("a", "rent", "2024-01-01", 100)];
    const b: SeqFactor[] = [makeFactor("b", "food", "2024-02-01", 50)];
    const result = mergeSequences(a, b);

    expect(result).toHaveLength(2);
    expect(result[0].id).toBe("a");
    expect(result[1].id).toBe("b");
  });

  test("merging with an empty array returns copy of the other", () => {
    const a: SeqFactor[] = [makeFactor("a", "rent", "2024-01-01", 100)];

    expect(mergeSequences(a, [])).toHaveLength(1);
    expect(mergeSequences([], a)).toHaveLength(1);
  });

  test("does not mutate the original arrays", () => {
    const a: SeqFactor[] = [makeFactor("a", "rent", "2024-01-01", 100)];
    const b: SeqFactor[] = [makeFactor("b", "food", "2024-02-01", 50)];
    mergeSequences(a, b);

    expect(a).toHaveLength(1);
    expect(b).toHaveLength(1);
  });
});
