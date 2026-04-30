//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-numeric
//ff:what Numeric family round-trip / nullable / slice / IsNil + 정밀도 보존 검증
package pgtypex

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestNumeric_RoundTripPrecision(t *testing.T) {
	for _, tc := range []string{"0", "1", "-1", "123.456", "9999999999999999999999.999999"} {
		pg := ToPgNumeric(tc)
		if !pg.Valid {
			t.Errorf("%q: ToPgNumeric Valid expected true, got %+v", tc, pg)
			continue
		}
		got := FromPgNumeric(pg)
		if got != tc {
			t.Errorf("%q round-trip: got %q", tc, got)
		}
	}
}

func TestNumeric_ParseFailureFallsBackToInvalid(t *testing.T) {
	if pg := ToPgNumeric("not-a-number"); pg.Valid {
		t.Fatal("malformed input must produce Valid:false")
	}
}

func TestNumeric_Ptr(t *testing.T) {
	if pg := ToPgNumericPtr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := "42.5"
	pg := ToPgNumericPtr(&v)
	got := FromPgNumericPtr(pg)
	if got == nil || *got != "42.5" {
		t.Fatalf("FromPgNumericPtr=%v", got)
	}
	if FromPgNumericPtr(pgtype.Numeric{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestNumeric_IsNil(t *testing.T) {
	if !IsNilPgNumeric(pgtype.Numeric{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgNumeric(pgtype.Numeric{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestNumeric_Bulk(t *testing.T) {
	out := ToPgNumerics([]string{"1", "2.5", "-3"})
	if len(out) != 3 || !out[0].Valid || !out[1].Valid {
		t.Fatalf("ToPgNumerics: %+v", out)
	}
	if FromPgNumeric(out[0]) != "1" || FromPgNumeric(out[1]) != "2.5" {
		t.Errorf("round-trip mismatch: %v", out)
	}
	if ToPgNumerics(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}

func TestNumeric_NaN(t *testing.T) {
	pg := pgtype.Numeric{NaN: true, Valid: true}
	if got := FromPgNumeric(pg); got != "NaN" {
		t.Fatalf("FromPgNumeric NaN=%q", got)
	}
}
