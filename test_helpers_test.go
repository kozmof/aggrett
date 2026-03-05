package aggrep

import (
	"strings"
	"testing"
	"time"
)

func mustParseDate(t *testing.T, s string) time.Time {
	t.Helper()

	layout := "2006-01-02"
	if strings.Contains(s, "T") {
		layout = "2006-01-02T15:04:05"
	}

	value, err := time.Parse(layout, s)
	if err != nil {
		t.Fatalf("failed to parse %q: %v", s, err)
	}
	return value
}

func mustTimeEqual(t *testing.T, got, want time.Time) {
	t.Helper()
	if !got.Equal(want) {
		t.Fatalf("time mismatch: got=%s want=%s", got.Format(time.RFC3339Nano), want.Format(time.RFC3339Nano))
	}
}

func strPtr(v string) *string { return &v }

func floatPtr(v float64) *float64 { return &v }

func factorPtr(v Factor) *Factor { return &v }

func timePtr(v time.Time) *time.Time { return &v }
