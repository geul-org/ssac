//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-float
//ff:what Float8/Float4 family round-trip / nullable / slice / IsNil 검증
package pgtypex

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestFloat8_RoundTrip(t *testing.T) {
	pg := ToPgFloat8(3.14)
	if !pg.Valid || pg.Float64 != 3.14 {
		t.Fatalf("ToPgFloat8: %+v", pg)
	}
	if got := FromPgFloat8(pg); got != 3.14 {
		t.Fatalf("FromPgFloat8=%v", got)
	}
}

func TestFloat8_Ptr(t *testing.T) {
	if pg := ToPgFloat8Ptr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := -2.5
	pg := ToPgFloat8Ptr(&v)
	got := FromPgFloat8Ptr(pg)
	if got == nil || *got != -2.5 {
		t.Fatalf("FromPgFloat8Ptr=%v", got)
	}
	if FromPgFloat8Ptr(pgtype.Float8{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestFloat8_IsNil(t *testing.T) {
	if !IsNilPgFloat8(pgtype.Float8{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgFloat8(pgtype.Float8{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestFloat8_Bulk(t *testing.T) {
	out := ToPgFloat8s([]float64{1.0, 2.0, 3.0})
	if len(out) != 3 || out[2].Float64 != 3.0 {
		t.Fatalf("ToPgFloat8s: %+v", out)
	}
	if ToPgFloat8s(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}

func TestFloat4_RoundTrip(t *testing.T) {
	pg := ToPgFloat4(1.5)
	if !pg.Valid || pg.Float32 != 1.5 {
		t.Fatalf("ToPgFloat4: %+v", pg)
	}
	if got := FromPgFloat4(pg); got != 1.5 {
		t.Fatalf("FromPgFloat4=%v", got)
	}
}

func TestFloat4_Ptr(t *testing.T) {
	if pg := ToPgFloat4Ptr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := float32(0.25)
	pg := ToPgFloat4Ptr(&v)
	got := FromPgFloat4Ptr(pg)
	if got == nil || *got != 0.25 {
		t.Fatalf("FromPgFloat4Ptr=%v", got)
	}
	if FromPgFloat4Ptr(pgtype.Float4{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestFloat4_IsNil(t *testing.T) {
	if !IsNilPgFloat4(pgtype.Float4{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgFloat4(pgtype.Float4{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestFloat4_Bulk(t *testing.T) {
	out := ToPgFloat4s([]float32{1.0, 2.0})
	if len(out) != 2 || out[0].Float32 != 1.0 {
		t.Fatalf("ToPgFloat4s: %+v", out)
	}
	if ToPgFloat4s(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}
