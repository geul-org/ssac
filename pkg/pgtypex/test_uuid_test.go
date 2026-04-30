//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-uuid
//ff:what UUID family 의 round-trip / nullable / slice / IsNil / UUIDToString 검증
package pgtypex

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
)

func TestUUID_RoundTrip(t *testing.T) {
	src := types.UUID{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10}
	pg := ToPgUUID(src)
	if !pg.Valid {
		t.Fatal("Valid expected true")
	}
	got := FromPgUUID(pg)
	if got != src {
		t.Fatalf("round-trip mismatch: %v != %v", got, src)
	}
}

func TestUUID_NullablePtr_NilInput(t *testing.T) {
	pg := ToPgUUIDPtr(nil)
	if pg.Valid {
		t.Fatal("nil input must produce Valid:false")
	}
	if got := FromPgUUIDPtr(pg); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestUUID_NullablePtr_ValueInput(t *testing.T) {
	src := types.UUID{0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a}
	pg := ToPgUUIDPtr(&src)
	if !pg.Valid {
		t.Fatal("Valid expected true")
	}
	got := FromPgUUIDPtr(pg)
	if got == nil {
		t.Fatal("expected non-nil")
	}
	if *got != src {
		t.Fatalf("round-trip mismatch: %v != %v", *got, src)
	}
}

func TestUUID_IsNil(t *testing.T) {
	if !IsNilPgUUID(pgtype.UUID{}) {
		t.Error("zero value must be nil")
	}
	if IsNilPgUUID(pgtype.UUID{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestUUID_Bulk(t *testing.T) {
	src := []types.UUID{{0x01}, {0x02}, {0x03}}
	pgs := ToPgUUIDs(src)
	if len(pgs) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(pgs))
	}
	for i, pg := range pgs {
		if !pg.Valid {
			t.Errorf("element %d: Valid expected true", i)
		}
		if got := FromPgUUID(pg); got != src[i] {
			t.Errorf("element %d: %v != %v", i, got, src[i])
		}
	}
	if ToPgUUIDs(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}

func TestUUIDToString(t *testing.T) {
	pg := pgtype.UUID{Bytes: [16]byte{0x55, 0x0e, 0x84, 0x00, 0xe2, 0x9b, 0x41, 0xd4, 0xa7, 0x16, 0x44, 0x66, 0x55, 0x44, 0x00, 0x00}, Valid: true}
	got := UUIDToString(pg)
	want := "550e8400-e29b-41d4-a716-446655440000"
	if got != want {
		t.Fatalf("UUIDToString=%q, want %q", got, want)
	}
	if got := UUIDToString(pgtype.UUID{}); got != "" {
		t.Errorf("NULL must render as empty: got %q", got)
	}
}
