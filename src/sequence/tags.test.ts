import { describe, test, expect } from "vitest";
import type { SeqFactor } from "../types/Sequence";
import {
  filterByTag,
  extractTags,
  groupByTag,
  removeByTag,
  renameTag,
  excludeByTag,
  accumulateByTag,
} from "./tags";

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

const seq: SeqFactor[] = [
  makeFactor("1", "rent", "2024-01-01", 1000, "minus"),
  makeFactor("2", "food", "2024-01-05", 200, "minus"),
  makeFactor("3", "salary", "2024-01-10", 5000, "plus"),
  makeFactor("4", "food", "2024-01-15", 150, "minus"),
  makeFactor("5", "utilities", "2024-01-20", 100, "minus"),
  makeFactor("6", "rent", "2024-02-01", 1000, "minus"),
];

describe("filterByTag", () => {
  test("filters to matching tags", () => {
    const result = filterByTag(seq, ["food"]);
    expect(result).toHaveLength(2);
    expect(result.every((f) => f.tag === "food")).toBe(true);
  });

  test("filters to multiple tags", () => {
    const result = filterByTag(seq, ["rent", "utilities"]);
    expect(result).toHaveLength(3);
    expect(result.map((f) => f.id)).toEqual(["1", "5", "6"]);
  });

  test("returns empty array when no tags match", () => {
    expect(filterByTag(seq, ["nonexistent"])).toHaveLength(0);
  });

  test("returns empty array for empty sequence", () => {
    expect(filterByTag([], ["food"])).toHaveLength(0);
  });

  test("does not mutate the original array", () => {
    filterByTag(seq, ["food"]);
    expect(seq).toHaveLength(6);
  });
});

describe("extractTags", () => {
  test("returns unique tags", () => {
    const tags = extractTags(seq);
    expect(tags).toHaveLength(4);
    expect(new Set(tags)).toEqual(
      new Set(["rent", "food", "salary", "utilities"]),
    );
  });

  test("returns empty array for empty sequence", () => {
    expect(extractTags([])).toHaveLength(0);
  });

  test("returns single tag when all items share it", () => {
    const uniform = [
      makeFactor("1", "rent", "2024-01-01", 100),
      makeFactor("2", "rent", "2024-02-01", 100),
    ];
    expect(extractTags(uniform)).toEqual(["rent"]);
  });
});

describe("groupByTag", () => {
  test("groups factors by tag", () => {
    const groups = groupByTag(seq);
    expect(Object.keys(groups)).toHaveLength(4);
    expect(groups["rent"]).toHaveLength(2);
    expect(groups["food"]).toHaveLength(2);
    expect(groups["salary"]).toHaveLength(1);
    expect(groups["utilities"]).toHaveLength(1);
  });

  test("preserves factor data in groups", () => {
    const groups = groupByTag(seq);
    expect(groups["salary"][0].id).toBe("3");
    expect(groups["salary"][0].value).toBe(5000);
  });

  test("returns empty object for empty sequence", () => {
    expect(groupByTag([])).toEqual({});
  });
});

describe("excludeByTag", () => {
  test("excludes matching tags", () => {
    const result = excludeByTag(seq, ["food"]);
    expect(result).toHaveLength(4);
    expect(result.every((f) => f.tag !== "food")).toBe(true);
  });

  test("excludes multiple tags", () => {
    const result = excludeByTag(seq, ["food", "utilities"]);
    expect(result).toHaveLength(3);
    expect(result.map((f) => f.id)).toEqual(["1", "3", "6"]);
  });

  test("returns all when no tags match", () => {
    const result = excludeByTag(seq, ["nonexistent"]);
    expect(result).toHaveLength(6);
  });

  test("returns empty array for empty sequence", () => {
    expect(excludeByTag([], ["food"])).toHaveLength(0);
  });

  test("does not mutate the original array", () => {
    excludeByTag(seq, ["food"]);
    expect(seq).toHaveLength(6);
  });
});

describe("removeByTag", () => {
  test("is the same operation as excludeByTag", () => {
    const excluded = excludeByTag(seq, ["food", "rent"]);
    const removed = removeByTag(seq, ["food", "rent"]);
    expect(removed).toEqual(excluded);
  });
});

describe("renameTag", () => {
  test("renames all matching factors", () => {
    const result = renameTag(seq, "food", "dining");
    const diningItems = result.filter((f) => f.tag === "dining");
    const foodItems = result.filter((f) => f.tag === "food");
    expect(diningItems).toHaveLength(2);
    expect(foodItems).toHaveLength(0);
  });

  test("preserves other fields", () => {
    const result = renameTag(seq, "food", "dining");
    const renamed = result.find((f) => f.id === "2")!;
    expect(renamed.tag).toBe("dining");
    expect(renamed.value).toBe(200);
    expect(renamed.factor).toBe("minus");
  });

  test("leaves non-matching factors unchanged", () => {
    const result = renameTag(seq, "food", "dining");
    expect(result.find((f) => f.id === "1")!.tag).toBe("rent");
  });

  test("returns unchanged copy when tag not found", () => {
    const result = renameTag(seq, "nonexistent", "new");
    expect(result).toHaveLength(6);
    expect(result).not.toBe(seq);
    expect(result).toEqual(seq);
  });

  test("does not mutate the original array", () => {
    renameTag(seq, "food", "dining");
    expect(seq[1].tag).toBe("food");
  });
});

describe("accumulateByTag", () => {
  test("accumulates only the specified tag", () => {
    const result = accumulateByTag(seq, 0, "food");
    expect(result).toHaveLength(2);
    expect(result[0].store).toBe(-200);
    expect(result[1].store).toBe(-350);
  });

  test("returns empty array when tag not found", () => {
    expect(accumulateByTag(seq, 0, "nonexistent")).toEqual([]);
  });

  test("applies base value", () => {
    const result = accumulateByTag(seq, 1000, "rent");
    expect(result[0].store).toBe(0); // 1000 - 1000
    expect(result[1].store).toBe(-1000); // 0 - 1000
  });

  test("returns empty array for empty sequence", () => {
    expect(accumulateByTag([], 0, "food")).toEqual([]);
  });
});
