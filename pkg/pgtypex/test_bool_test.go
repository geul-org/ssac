//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-bool
//ff:what Bool family round-trip / nullable / slice / IsNil 검증
package pgtypex

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestBool_RoundTrip(t *testing.T) {
	pg := ToPgBool(true)
	if !pg.Valid || !pg.Bool {
		t.Fatalf("ToPgBool: %+v", pg)
	}
	if got := FromPgBool(pg); !got {
		t.Fatalf("FromPgBool=%v", got)
	}
}

func TestBool_Ptr(t *testing.T) {
	if pg := ToPgBoolPtr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	v := false
	pg := ToPgBoolPtr(&v)
	got := FromPgBoolPtr(pg)
	if got == nil || *got != false {
		t.Fatalf("FromPgBoolPtr=%v", got)
	}
	if FromPgBoolPtr(pgtype.Bool{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestBool_IsNil(t *testing.T) {
	if !IsNilPgBool(pgtype.Bool{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgBool(pgtype.Bool{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestBool_Bulk(t *testing.T) {
	out := ToPgBools([]bool{true, false, true})
	if len(out) != 3 || !out[0].Bool || out[1].Bool {
		t.Fatalf("ToPgBools: %+v", out)
	}
	if ToPgBools(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}
