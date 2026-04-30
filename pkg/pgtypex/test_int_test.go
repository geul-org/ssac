//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-int
//ff:what Int8/Int4/Int2 family round-trip / nullable / slice / IsNil 검증
package pgtypex

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestInt8_RoundTrip(t *testing.T) {
	pg := ToPgInt8(42)
	if !pg.Valid || pg.Int64 != 42 {
		t.Fatalf("ToPgInt8: %+v", pg)
	}
	if got := FromPgInt8(pg); got != 42 {
		t.Fatalf("FromPgInt8=%d", got)
	}
}

func TestInt8_Ptr(t *testing.T) {
	if pg := ToPgInt8Ptr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := int64(-7)
	pg := ToPgInt8Ptr(&v)
	got := FromPgInt8Ptr(pg)
	if got == nil || *got != -7 {
		t.Fatalf("FromPgInt8Ptr=%v", got)
	}
	if FromPgInt8Ptr(pgtype.Int8{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestInt8_IsNil(t *testing.T) {
	if !IsNilPgInt8(pgtype.Int8{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgInt8(pgtype.Int8{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestInt8_Bulk(t *testing.T) {
	out := ToPgInt8s([]int64{1, 2, 3})
	if len(out) != 3 || out[0].Int64 != 1 || out[2].Int64 != 3 {
		t.Fatalf("ToPgInt8s: %+v", out)
	}
	if ToPgInt8s(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}

func TestInt4_RoundTrip(t *testing.T) {
	pg := ToPgInt4(123)
	if !pg.Valid || pg.Int32 != 123 {
		t.Fatalf("ToPgInt4: %+v", pg)
	}
	if got := FromPgInt4(pg); got != 123 {
		t.Fatalf("FromPgInt4=%d", got)
	}
}

func TestInt4_Ptr(t *testing.T) {
	if pg := ToPgInt4Ptr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := int32(-1)
	pg := ToPgInt4Ptr(&v)
	got := FromPgInt4Ptr(pg)
	if got == nil || *got != -1 {
		t.Fatalf("FromPgInt4Ptr=%v", got)
	}
	if FromPgInt4Ptr(pgtype.Int4{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestInt4_IsNil(t *testing.T) {
	if !IsNilPgInt4(pgtype.Int4{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgInt4(pgtype.Int4{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestInt4_Bulk(t *testing.T) {
	out := ToPgInt4s([]int32{10, 20})
	if len(out) != 2 || out[0].Int32 != 10 {
		t.Fatalf("ToPgInt4s: %+v", out)
	}
	if ToPgInt4s(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}

func TestInt2_RoundTrip(t *testing.T) {
	pg := ToPgInt2(7)
	if !pg.Valid || pg.Int16 != 7 {
		t.Fatalf("ToPgInt2: %+v", pg)
	}
	if got := FromPgInt2(pg); got != 7 {
		t.Fatalf("FromPgInt2=%d", got)
	}
}

func TestInt2_Ptr(t *testing.T) {
	if pg := ToPgInt2Ptr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := int16(-9)
	pg := ToPgInt2Ptr(&v)
	got := FromPgInt2Ptr(pg)
	if got == nil || *got != -9 {
		t.Fatalf("FromPgInt2Ptr=%v", got)
	}
	if FromPgInt2Ptr(pgtype.Int2{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestInt2_IsNil(t *testing.T) {
	if !IsNilPgInt2(pgtype.Int2{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgInt2(pgtype.Int2{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestInt2_Bulk(t *testing.T) {
	out := ToPgInt2s([]int16{5, 6})
	if len(out) != 2 || out[1].Int16 != 6 {
		t.Fatalf("ToPgInt2s: %+v", out)
	}
	if ToPgInt2s(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}
