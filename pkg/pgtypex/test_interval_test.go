//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-interval
//ff:what Interval family ISO 8601 round-trip / nullable / slice / IsNil 검증
package pgtypex

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestInterval_RoundTrip(t *testing.T) {
	cases := []struct {
		name string
		v    pgtype.Interval
		s    string
	}{
		{"zero", pgtype.Interval{Valid: true}, "PT0S"},
		{"1 hour", pgtype.Interval{Microseconds: 3600 * 1_000_000, Valid: true}, "PT1H"},
		{"1 day", pgtype.Interval{Days: 1, Valid: true}, "P1D"},
		{"1 month", pgtype.Interval{Months: 1, Valid: true}, "P1M"},
		{"1 year", pgtype.Interval{Months: 12, Valid: true}, "P1Y"},
		{"complex", pgtype.Interval{Months: 14, Days: 5, Microseconds: int64(3661)*1_000_000 + 500_000, Valid: true}, "P1Y2M5DT1H1M1.5S"},
	}
	for _, tc := range cases {
		got := formatIntervalISO(tc.v)
		if got != tc.s {
			t.Errorf("%s: format=%q want %q", tc.name, got, tc.s)
		}
		parsed, err := parseIntervalISO(tc.s)
		if err != nil {
			t.Errorf("%s: parse(%q): %v", tc.name, tc.s, err)
			continue
		}
		if parsed.Months != tc.v.Months || parsed.Days != tc.v.Days || parsed.Microseconds != tc.v.Microseconds {
			t.Errorf("%s: parse mismatch: %+v want %+v", tc.name, parsed, tc.v)
		}
	}
}

func TestInterval_PublicAPI(t *testing.T) {
	pg := ToPgInterval("PT2H")
	if !pg.Valid || pg.Microseconds != 2*3600*1_000_000 {
		t.Fatalf("ToPgInterval: %+v", pg)
	}
	if got := FromPgInterval(pg); got != "PT2H" {
		t.Errorf("FromPgInterval=%q", got)
	}
}

func TestInterval_ParseFailureFallsBackToInvalid(t *testing.T) {
	if pg := ToPgInterval("not-iso"); pg.Valid {
		t.Fatal("malformed input must produce Valid:false")
	}
}

func TestInterval_Ptr(t *testing.T) {
	if pg := ToPgIntervalPtr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := "PT30M"
	pg := ToPgIntervalPtr(&v)
	got := FromPgIntervalPtr(pg)
	if got == nil || *got != "PT30M" {
		t.Fatalf("FromPgIntervalPtr=%v", got)
	}
	if FromPgIntervalPtr(pgtype.Interval{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestInterval_IsNil(t *testing.T) {
	if !IsNilPgInterval(pgtype.Interval{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgInterval(pgtype.Interval{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestInterval_Bulk(t *testing.T) {
	out := ToPgIntervals([]string{"PT1H", "P1D"})
	if len(out) != 2 || !out[0].Valid || !out[1].Valid {
		t.Fatalf("ToPgIntervals: %+v", out)
	}
	if ToPgIntervals(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}
