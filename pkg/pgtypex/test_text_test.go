//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-text
//ff:what Text family round-trip / nullable / slice / IsNil 검증
package pgtypex

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestText_RoundTrip(t *testing.T) {
	pg := ToPgText("hello")
	if !pg.Valid || pg.String != "hello" {
		t.Fatalf("ToPgText: %+v", pg)
	}
	if got := FromPgText(pg); got != "hello" {
		t.Fatalf("FromPgText=%q", got)
	}
}

func TestText_Ptr_Nil(t *testing.T) {
	if pg := ToPgTextPtr(nil); pg.Valid {
		t.Fatal("nil input must produce Valid:false")
	}
	if got := FromPgTextPtr(pgtype.Text{}); got != nil {
		t.Fatalf("expected nil, got %q", *got)
	}
}

func TestText_Ptr_Value(t *testing.T) {
	v := "world"
	pg := ToPgTextPtr(&v)
	if !pg.Valid || pg.String != "world" {
		t.Fatalf("ToPgTextPtr: %+v", pg)
	}
	got := FromPgTextPtr(pg)
	if got == nil || *got != "world" {
		t.Fatalf("FromPgTextPtr=%v", got)
	}
}

func TestText_IsNil(t *testing.T) {
	if !IsNilPgText(pgtype.Text{}) {
		t.Error("zero value must be nil")
	}
	if IsNilPgText(pgtype.Text{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestText_Bulk(t *testing.T) {
	out := ToPgTexts([]string{"a", "b"})
	if len(out) != 2 || out[0].String != "a" || out[1].String != "b" {
		t.Fatalf("ToPgTexts: %+v", out)
	}
	if ToPgTexts(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}
