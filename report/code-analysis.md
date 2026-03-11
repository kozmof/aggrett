# Code Analysis: `aggrett`

**Date:** 2026-03-11
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

**Internal types** (unexported): `timeGroup`, `timeGrouping`

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

AccumCore                (flat accumulation result)
  ├─ IDs   []string
  ├─ Time  time.Time
  └─ Store float64

Accum                    (rich accumulation result)
  ├─ IDs       []string
  ├─ Tags      []string
  ├─ Time      time.Time
  ├─ Store     float64
  └─ Breakdown Breakdown

Breakdown = map[string]BreakdownEntry
BreakdownEntry
  ├─ Delta float64
  └─ IDs   []string

IntervalType (string)
  └─ Seconds / Minutes / Hours / Days / Weeks / Months / Years
  └─ IsValid() bool

timeGrouping (unexported)
  ├─ Step         int
  └─ IntervalType IntervalType
  └─ validate()

timeGroup (unexported)
  ├─ Time    time.Time
  └─ Factors []SeqFactor
```

**Key asymmetry:** `AccumCore` and `Accum` share the same top-level fields (`IDs`, `Time`, `Store`) but are not related by embedding or interface. Callers must choose the right accumulation path upfront.

---

## 3. Function Call Graph

```
Public API
├─ Aggregate(sequence, base, filter)
│   └─ aggregate(sequence, base, filter, nil)
│       ├─ FilterByTag (if filter non-empty)
│       ├─ groupByTime(filtered, nil)
│       │   └─ sort.SliceStable
│       └─ Accumulate (per factor)
│
├─ AggregateByInterval(sequence, base, filter, step, intervalType)
│   └─ aggregate(..., &timeGrouping{...})
│       └─ groupByTime(filtered, grouping)
│           └─ bucketStart(f.Time, grouping)
│               └─ timeGrouping.validate()
│
├─ AccumulateSequence(sequence, base)
│   └─ accumulateSequence(sequence, base, nil)
│       └─ groupByTime + Accumulate
│
├─ AccumulateSequenceByInterval(sequence, base, step, intervalType)
│   └─ accumulateSequence(..., &timeGrouping{...})
│
├─ AccumulateByTag(sequence, base, tag)
│   └─ FilterByTag → AccumulateSequence
│
├─ AccumulateByTagByInterval(sequence, base, tag, step, intervalType)
│   └─ FilterByTag → AccumulateSequenceByInterval
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
seq := aggrett.InsertFactor(nil, "salary", time.Now(), 5000, aggrett.FactorPlus, uuid.NewString)
seq = aggrett.InsertFactor(seq, "rent",   rentDate, 1000, aggrett.FactorMinus, uuid.NewString)

// 2. Generate recurring events
cycles := aggrett.CreateCycles(start, end, 1, aggrett.IntervalMonths,
    aggrett.FactorMinus, 1000, "rent", uuid.NewString)
seq = aggrett.MergeSequences(seq, cycles)

// 3. Aggregate
result := aggrett.Aggregate(seq, 0, nil)
// or by time bucket:
result := aggrett.AggregateByInterval(seq, 0, nil, 1, aggrett.IntervalMonths)

// 4. Tag-scoped view
rentOnly := aggrett.AccumulateByTag(seq, 0, "rent")
```

---

## 5. Pitfalls

### 5a. `Accumulate` panics on unknown `Factor`
[sequence.go:21](../sequence.go#L21) — An unvalidated `Factor` value (e.g. loaded from JSON/DB) will panic at runtime. `Factor.IsValid()` exists but is not called internally before dispatch.

```go
// If f.Factor = "plus " (note trailing space), this panics:
store = Accumulate(f.Factor, store, f.Value)
```

### 5b. `timeGrouping.validate()` panics instead of returning errors
[timeseries.go:28-35](../timeseries.go#L28) — Validation panics are called deep inside `groupByTime` → `bucketStart`. There is no way to surface these as structured errors to callers.

### 5c. Day bucket math is calendar-month-relative, not epoch-relative
[timeseries.go:83-84](../timeseries.go#L83) — `IntervalDays` buckets are calculated as `((day-1)/step*step)+1` within the current month. A `step=3` bucket in January and February start fresh each month. Cross-month multi-day buckets do **not** form a global grid.

### 5d. Week bucket math resets at year boundary
[timeseries.go:86-88](../timeseries.go#L86) — `IntervalWeeks` uses `YearDay`, so a 2-week bucket spanning Dec 31 → Jan 1 would split into two different buckets. ISO week numbering is not followed.

### 5e. Month overflow in `CreateCycles` causes date drift
[timeseries.go:100-107](../timeseries.go#L100) — Monthly cycles starting on the 31st drift: Jan-31 → Feb-29 → Mar-29 (not Mar-31). This is documented as "JS-like" but can be surprising. The drift is permanent — once clamped, subsequent months are computed from the clamped date.

### 5f. `FindByID` is O(n)
[factors.go:72](../factors.go#L72) — Linear scan over the full sequence. No indexing mechanism is provided; callers with large sequences bear the full cost.

### 5g. ID uniqueness is assumed but not enforced
`InsertFactor` delegates ID generation to the caller via `genID func() string`. `UpdateFactor` silently matches only the first occurrence if IDs are duplicated (due to early `continue` logic). `RemoveFactor` will remove **all** factors with a matching ID from the set if duplicates exist.

### 5h. `timePtr` test helper is defined but never used
[test_helpers_test.go:37](../test_helpers_test.go#L37) — `func timePtr(v time.Time) *time.Time` is declared but not referenced in any test.

### 5i. Duplicate factory helpers across test files
`makeFactor` ([factors_test.go:9](../factors_test.go#L9)) and `makeTagFactor` ([tags_test.go:8](../tags_test.go#L8)) are nearly identical. Both are in `package aggrett` (internal tests) and compile together, so they co-exist — but the duplication adds maintenance burden.

---

## 6. Improvement Points: Design Overview

1. **Return errors instead of panicking.** Both `Accumulate` and `bucketStart` should return `(float64, error)` / `(time.Time, error)`. Go's idiom is explicit error handling; panics are for programmer errors only. Consumers calling these from user-supplied data have no recovery path.

2. **Expose `TimeGrouping` as a public type.** The `step int, intervalType IntervalType` pair is repeated in every interval-based function signature. A public `TimeGrouping` struct would allow callers to build it once and reuse, and simplify the API surface.

3. **Unify result types.** `AccumCore` has fields `IDs`, `Time`, `Store`. `Accum` has all of those plus `Tags` and `Breakdown`. Either embed `AccumCore` inside `Accum`, or consolidate into one type with optional/zero-value breakdown. Currently callers are locked into one path early and cannot easily upgrade.

4. **Consider an index structure.** For use cases with large sequences and frequent `FindByID` / `UpdateFactor` calls, a companion `Index` type (`map[string]*SeqFactor`) would allow O(1) access. The current immutable-slice model is simple but doesn't scale.

5. **Week/day bucketing should use epoch-based alignment.** The current calendar-relative bucketing for days and weeks creates surprising discontinuities at month/year boundaries. An epoch-aligned approach (e.g., days since Unix epoch divided by step) would produce consistent global buckets.

---

## 7. Improvement Points: Types and Interfaces

1. **`Factor` as a typed constant with `iota`** would eliminate invalid string values entirely. Or define `Factor` as `bool` (`Plus = true`, `Minus = false`) and use a switch on the bool — prevents any string-based mistake.

2. **`SeqFactorUpdate` functional options** — the current pointer-fields pattern works but is verbose to construct. Functional options (`type FactorOption func(*SeqFactor)`) would be more idiomatic Go and easier to extend.

3. **`Breakdown` map nil-safety** — accessing `accum.Breakdown["missing-tag"]` returns a zero `BreakdownEntry` with `nil` IDs and `Delta = 0`. Document this or provide a `Get(tag string) (BreakdownEntry, bool)` accessor to make absence explicit.

4. **`fmt.Stringer` on `Factor` and `IntervalType`** — neither implements `String()`, so `fmt.Println(factor)` prints the underlying string value which happens to be readable, but it's coincidental. Implementing `Stringer` makes the intent explicit.

5. **`AccumCore` and `Accum` should share a common interface or embedding** — right now there is zero relationship between them in the type system. A reader of the API sees two unrelated result types with overlapping fields.

---

## 8. Improvement Points: Implementations

1. **`InsertFactor` pre-allocation** ([factors.go:14-16](../factors.go#L14)) — creates a new slice of `cap = len+1`, copies, then appends. This is equivalent to `append(sequence, newElem)` on a copied slice. The explicit pre-allocation is clear but verbose; using `append(append([]SeqFactor{}, sequence...), newElem)` is idiomatic.

2. **`tags` slice in `aggregate`** ([aggregate.go:33](../aggregate.go#L33)) — initialized with `make([]string, 0)` while `ids` uses `make([]string, 0, len(group.Factors))`. Since `tags` is bounded by unique tags per group (at most `len(group.Factors)`), the same capacity hint applies.

3. **`GenerateTimeSeries` cannot pre-allocate** ([timeseries.go:133](../timeseries.go#L133)) — because month/year overflow clamping means consecutive intervals may land on the same date, the result count isn't knowable upfront. This is acceptable but worth documenting.

4. **`groupByTime` appends unconditionally at end** ([sequence.go:58](../sequence.go#L58)) — `return append(groups, current)` is safe because the early-return guard ensures `current` is always populated when reached. However, the pattern diverges from typical "flush on boundary" loops and may confuse future maintainers.

5. **`lastDayOfPreviousMonth` preserves time-of-day** ([timeseries.go:118-129](../timeseries.go#L118)) — the clamped date keeps the original `Hour`, `Minute`, `Second`, `Nanosecond`. For `AddInterval`, this is intentional, but it makes the function less obviously a "date boundary" utility.

6. **`RemoveByTag` alias** ([tags.go:59-61](../tags.go#L59)) — two names for the same operation add cognitive overhead. If the alias is kept for backwards compatibility, a deprecation comment would help; otherwise, consider removing one.

---

## 9. Learning Paths

### Entry Points

| Step | File | What to learn |
|---|---|---|
| 1 | `types.go` | Core vocabulary: `Factor`, `SeqFactor`, `Accum`, `Breakdown` |
| 2 | `sequence.go` | Engine: `Accumulate`, `groupByTime`, `accumulateSequence` |
| 3 | `aggregate.go` | Full pipeline with breakdown: how `Accum` is assembled |
| 4 | `tags.go` | Tag filtering and tag-scoped accumulation |
| 5 | `timeseries.go` | Bucket math, interval arithmetic, cycle generation |
| 6 | `factors.go` | Sequence mutation CRUD pattern |

### Goals for Contributors

- **Immutability contract:** Every function returns a new slice. Understand why `groupByTime` does a `copy` before sorting ([sequence.go:31-35](../sequence.go#L31)).
- **Two result types:** Know when `AccumCore` (simple running total) vs `Accum` (with per-tag breakdown) is appropriate.
- **Bucket alignment math:** Understand how `bucketStart` maps a timestamp to its bucket start for each `IntervalType`, especially the day/week edge cases.
- **Overflow clamping:** Understand `addMonthsWithOverflowClamp` and why monthly cycles from the 31st drift.
- **Panic contract:** Currently `Accumulate` and `bucketStart` panic on invalid input. Any extension must maintain or replace this contract consistently.

---

## Summary Table

| Category | Rating | Notes |
|---|---|---|
| Correctness | High | Well-tested; immutability enforced |
| API clarity | Medium | Two result types; repeated `step/intervalType` params |
| Error handling | Low | Panics instead of errors throughout |
| Performance | Medium | O(n) find/update; no indexing |
| Test coverage | High | All public functions covered with edge cases |
| Extensibility | Medium | Unexported `timeGrouping` limits reuse |
