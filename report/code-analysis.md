# Code Analysis: `aggrett`

**Date:** 2026-03-13
**Module:** `github.com/kozmof/aggrett`
**Language:** Go 1.26

---

## 1. Code Organization and Structure

The library is a pure Go package with no external dependencies. It implements a **time-series factor accumulation system** — accumulating add/subtract events (factors) over ordered time data.

| File | Role |
|---|---|
| `types.go` | Core data types: `Factor`, `SeqFactor`, `Accum`, `AccumCore`, `Breakdown`, `BreakdownEntry` |
| `factors.go` | Sequence mutation: Insert, Remove, Update, Find, Merge |
| `sequence.go` | Core engine: `Accumulate`, `groupByTime`, `AccumulateSequence` |
| `tags.go` | Tag-based filtering and tag-scoped accumulation |
| `timeseries.go` | Interval types, bucket math, time series generation |
| `aggregate.go` | Thin facade: `Aggregate` / `AggregateByInterval` returning rich `Accum` |

**Internal types** (unexported): `timeGroup`

**Exported types added**: `TimeGrouping`, `Index`, `Accumulated` interface

The library enforces an **immutable-sequence contract**: every mutation function returns a new `[]SeqFactor` without modifying the input. This is consistently validated across the test suite.

---

## 2. Type and Interface Relations

```
Factor (string)
  └─ FactorPlus / FactorMinus
  └─ IsValid() bool

SeqFactor
  ├─ ID     string
  ├─ Tag    string
  ├─ Time   time.Time
  ├─ Value  float64
  └─ Factor Factor

SeqFactorUpdate          (partial patch; pointer fields = optional)
  ├─ Tag    *string
  ├─ Time   *time.Time
  ├─ Value  *float64
  └─ Factor *Factor

AccumCore                (flat accumulation result — implements Accumulated)
  ├─ IDs   []string
  ├─ Time  time.Time
  └─ Store float64

Accum                    (rich accumulation result — implements Accumulated)
  ├─ AccumCore (embedded)
  ├─ Tags      []string
  └─ Breakdown Breakdown

Accumulated (interface)
  ├─ GetIDs()   []string
  ├─ GetTime()  time.Time
  └─ GetStore() float64

Breakdown = map[string]BreakdownEntry
BreakdownEntry
  ├─ Delta float64
  └─ IDs   []string

IntervalType (string)
  └─ Seconds / Minutes / Hours / Days / Weeks / Months / Years
  └─ IsValid() bool

TimeGrouping (exported)
  ├─ Step         int
  └─ IntervalType IntervalType
  └─ validate() error

Index (exported)
  └─ BuildIndex([]SeqFactor) Index
  └─ Lookup([]SeqFactor, id) (SeqFactor, bool)

timeGroup (unexported)
  ├─ Time    time.Time
  └─ Factors []SeqFactor
```

Both `AccumCore` and `Accum` implement the `Accumulated` interface. `Accum` embeds `AccumCore`, so field access (`accum.Store`, `accum.IDs`) is promoted. Callers can write functions accepting `Accumulated` to work with either type.

---

## 3. Function Call Graph

```
Public API
├─ Aggregate(sequence, base, filter) ([]Accum, error)
│   └─ aggregate(sequence, base, filter, nil)
│       ├─ FilterByTag (if filter non-empty)
│       ├─ groupByTime(filtered, nil)
│       │   └─ sort.SliceStable
│       └─ Accumulate (per factor) → error on unknown Factor
│
├─ AggregateByInterval(sequence, base, filter, step, intervalType) ([]Accum, error)
│   └─ aggregate(..., &TimeGrouping{...})
│       └─ groupByTime(filtered, grouping)
│           └─ bucketStart(f.Time, grouping) → error on invalid grouping
│               └─ TimeGrouping.validate() → error
│
├─ AggregateByGrouping(sequence, base, filter, grouping TimeGrouping) ([]Accum, error)
│   └─ aggregate(..., &grouping)
│
├─ AccumulateSequence(sequence, base) ([]AccumCore, error)
│   └─ accumulateSequence(sequence, base, nil)
│       └─ groupByTime + Accumulate
│
├─ AccumulateSequenceByInterval(sequence, base, step, intervalType) ([]AccumCore, error)
│   └─ accumulateSequence(..., &TimeGrouping{...})
│
├─ AccumulateSequenceByGrouping(sequence, base, grouping TimeGrouping) ([]AccumCore, error)
│   └─ accumulateSequence(..., &grouping)
│
├─ AccumulateByTag(sequence, base, tag) ([]AccumCore, error)
│   └─ FilterByTag → AccumulateSequence
│
├─ AccumulateByTagByInterval(sequence, base, tag, step, intervalType) ([]AccumCore, error)
│   └─ FilterByTag → AccumulateSequenceByInterval
│
├─ InsertFactor(sequence, tag, time, value, factor, genID) ([]SeqFactor, error)
│   └─ factor.IsValid() → error on invalid Factor
│
├─ UpdateFactor(sequence, id, fields) ([]SeqFactor, error)
│   └─ fields.Factor.IsValid() → error on invalid Factor (when non-nil)
│
└─ CreateCycles(start, end, step, intervalType, factor, value, tag, genID)
    └─ GenerateTimeSeries(start, end, step, intervalType)
        └─ AddInterval(current, step, intervalType) [loop]
            ├─ addMonthsWithOverflowClamp (for months)
            └─ addYearsWithOverflowClamp  (for years)
                └─ lastDayOfPreviousMonth
```

---

## 4. Specific Contexts and Usages

The library is designed for **personal finance / budgeting scenarios** (test data uses tags: `rent`, `food`, `salary`, `utilities`). Typical usage pattern:

```go
// 1. Build a sequence
seq, err := aggrett.InsertFactor(nil, "salary", time.Now(), 5000, aggrett.FactorPlus, uuid.NewString)
seq, err = aggrett.InsertFactor(seq, "rent", rentDate, 1000, aggrett.FactorMinus, uuid.NewString)

// 2. Generate recurring events
cycles := aggrett.CreateCycles(start, end, 1, aggrett.IntervalMonths,
    aggrett.FactorMinus, 1000, "rent", uuid.NewString)
seq = aggrett.MergeSequences(seq, cycles)

// 3. Aggregate
result, err := aggrett.Aggregate(seq, 0, nil)
// or by time bucket:
result, err := aggrett.AggregateByInterval(seq, 0, nil, 1, aggrett.IntervalMonths)

// 4. Tag-scoped view
rentOnly, err := aggrett.AccumulateByTag(seq, 0, "rent")
```

---

## 5. Pitfalls

### ~~5a. `Accumulate` panics on unknown `Factor`~~ — Fixed
`Accumulate` now returns `(float64, error)` instead of panicking. `InsertFactor` and `UpdateFactor` validate the `Factor` field at the boundary and return `([]SeqFactor, error)`, so well-formed sequences never reach `Accumulate` with an invalid factor.

### ~~5b. `timeGrouping.validate()` panics instead of returning errors~~ — Fixed
`validate()` now returns `error`. `bucketStart` returns `(time.Time, error)`. The error propagates up through `groupByTime` → `accumulateSequence` / `aggregate` → all public interval-based functions, which now return errors.

### ~~5c. Day bucket math is calendar-month-relative, not epoch-relative~~ — Fixed
`IntervalDays` now uses epoch-aligned integer division: days since Unix epoch divided by step gives a consistent global grid. A `step=3` bucket crossing a month boundary (e.g., Jan-31 and Feb-1) correctly lands in the same bucket.

### ~~5d. Week bucket math resets at year boundary~~ — Fixed
`IntervalWeeks` now uses epoch-aligned integer division: days since Unix epoch divided by 7 gives stable week indices. A week bucket spanning Dec 31 → Jan 1 no longer splits.

### 5e. Month overflow in `CreateCycles` causes date drift
[timeseries.go](../timeseries.go) — Monthly cycles starting on the 31st drift: Jan-31 → Feb-29 → Mar-29 (not Mar-31). This is intentional ("JS-like" overflow clamping) and is documented in the `AddInterval` and `CreateCycles` doc comments. The drift is permanent once it begins.

### ~~5f. `FindByID` is O(n)~~ — Mitigated
An `Index` companion type has been added ([index.go](../index.go)). `BuildIndex([]SeqFactor)` builds a `map[string]int` in O(n); subsequent `Lookup` calls are O(1). The index includes stale-detection: if the slice is mutated after the index was built, `Lookup` returns `false` rather than a wrong result.

### ~~5g. ID uniqueness is assumed but not enforced~~ — Fixed
`InsertFactor` now checks the generated ID against all existing IDs before inserting and returns an error if a duplicate is found.

### ~~5h. `timePtr` test helper is defined but never used~~ — Fixed
Removed from `test_helpers_test.go`.

### ~~5i. Duplicate factory helpers across test files~~ — Fixed
`makeFactor` and `makeTagFactor` have been removed and consolidated into a single `makeSeqFactor` in `test_helpers_test.go`, used by all test files.

---

## 6. Improvement Points: Design Overview

1. ~~**Return errors instead of panicking.**~~ **Done.** `Accumulate`, `bucketStart`, and `validate` now return errors. All public accumulation and aggregation functions return `(result, error)`. Factor validation is enforced at `InsertFactor` / `UpdateFactor`.

2. ~~**Expose `TimeGrouping` as a public type.**~~ **Done.** `TimeGrouping` is now exported. `AccumulateSequenceByGrouping` and `AggregateByGrouping` accept it directly, letting callers build a grouping once and reuse it. The original `ByInterval` signatures are retained for convenience.

3. ~~**Unify result types.**~~ **Done.** `Accum` now embeds `AccumCore`. All field accesses are promoted. Both types implement the `Accumulated` interface, enabling polymorphic code.

4. ~~**Consider an index structure.**~~ **Done.** `Index` companion type added in [index.go](../index.go). O(1) lookup via `BuildIndex` + `Lookup`. Includes stale-detection.

5. ~~**Week/day bucketing should use epoch-based alignment.**~~ **Done.** See items 5c and 5d above.

---

## 7. Improvement Points: Types and Interfaces

1. **`Factor` as a typed constant with `iota`** would eliminate invalid string values entirely. Or define `Factor` as `bool` (`Plus = true`, `Minus = false`) and use a switch on the bool — prevents any string-based mistake.

2. **`SeqFactorUpdate` functional options** — the current pointer-fields pattern works but is verbose to construct. Functional options (`type FactorOption func(*SeqFactor)`) would be more idiomatic Go and easier to extend.

3. ~~**`Breakdown` map nil-safety**~~ **Done.** `Breakdown.Get(tag string) (BreakdownEntry, bool)` added. Returns `false` for nil map or missing key, making absence explicit without a map lookup.

4. ~~**`fmt.Stringer` on `Factor` and `IntervalType`**~~ **Done.** Both types now implement `String() string`, making their output in `fmt.Println` and error messages intentional.

5. ~~**`AccumCore` and `Accum` should share a common interface or embedding**~~ **Done.** `Accum` embeds `AccumCore`; both implement `Accumulated`. See item 6.3 above.

---

## 8. Improvement Points: Implementations

1. ~~**`InsertFactor` pre-allocation**~~ **Done.** Simplified to `append(append([]SeqFactor{}, sequence...), newElem)` — idiomatic and equivalent.

2. ~~**`tags` slice in `aggregate`**~~ **Done.** `tags` now uses `make([]string, 0, len(group.Factors))` matching the capacity hint already used for `ids`.

3. ~~**`GenerateTimeSeries` cannot pre-allocate**~~ **Documented.** A doc comment on `GenerateTimeSeries` now explains why pre-allocation is not possible: month/year overflow clamping means consecutive intervals may produce the same date, making the final count unknowable upfront.

4. ~~**`groupByTime` appends unconditionally at end**~~ **Documented.** A comment at the final `append(groups, current)` explains that the empty-sequence guard above guarantees `current` is always populated at this point.

5. ~~**`lastDayOfPreviousMonth` preserves time-of-day**~~ **Documented.** A doc comment on `lastDayOfPreviousMonth` clarifies that preserving time-of-day is intentional for `AddInterval` semantics.

6. ~~**`RemoveByTag` alias**~~ **Done.** A `// Deprecated: Use ExcludeByTag instead.` comment has been added to `RemoveByTag`.

---

## 9. Learning Paths

### Entry Points

| Step | File | What to learn |
|---|---|---|
| 1 | `types.go` | Core vocabulary: `Factor`, `SeqFactor`, `Accum`, `AccumCore`, `Breakdown`, `Accumulated` |
| 2 | `sequence.go` | Engine: `Accumulate`, `groupByTime`, `accumulateSequence` |
| 3 | `aggregate.go` | Full pipeline with breakdown: how `Accum` is assembled |
| 4 | `tags.go` | Tag filtering and tag-scoped accumulation |
| 5 | `timeseries.go` | Bucket math, interval arithmetic, cycle generation |
| 6 | `factors.go` | Sequence mutation CRUD pattern |

### Goals for Contributors

- **Immutability contract:** Every function returns a new slice. Understand why `groupByTime` does a `copy` before sorting ([sequence.go:31-35](../sequence.go#L31)).
- **Two result types:** Know when `AccumCore` (simple running total) vs `Accum` (with per-tag breakdown) is appropriate. Both implement `Accumulated` and `Accum` embeds `AccumCore`.
- **Bucket alignment math:** `bucketStart` uses epoch-aligned integer division for `IntervalDays` and `IntervalWeeks`, and calendar-relative logic for months/years. Understand the cross-month consistency guarantee for days/weeks.
- **Overflow clamping:** Understand `addMonthsWithOverflowClamp` and why monthly cycles from the 31st drift.
- **Error contract:** `Accumulate`, `bucketStart`, `InsertFactor`, and `UpdateFactor` return errors on invalid input. All public accumulation and aggregation functions propagate these errors.

---

## Summary Table

| Category | Rating | Notes |
|---|---|---|
| Correctness | High | Well-tested; immutability enforced; epoch-aligned day/week buckets |
| API clarity | **High** | `TimeGrouping` exported; `ByGrouping` variants added; `Accumulated` interface unifies result types |
| Error handling | **High** | Panics replaced with errors throughout; Factor validated at insertion boundaries |
| Performance | **High** | O(1) lookup via `Index`; O(n) `FindByID` still available for one-off use |
| Test coverage | High | All public functions covered; error paths, epoch bucket tests, interface tests added |
| Extensibility | **High** | `TimeGrouping` exported; `Accumulated` interface; `Index` companion type |
