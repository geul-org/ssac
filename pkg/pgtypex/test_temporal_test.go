//ff:func feature=pkg-pgtypex type=test control=sequence topic=pgtypex-temporal
//ff:what Timestamptz/Timestamp/Date family round-trip / nullable / slice / IsNil 검증
package pgtypex

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func TestTimestamptz_RoundTrip(t *testing.T) {
	src := time.Date(2026, 4, 30, 12, 30, 0, 0, time.UTC)
	pg := ToPgTimestamptz(src)
	if !pg.Valid || !pg.Time.Equal(src) {
		t.Fatalf("ToPgTimestamptz: %+v", pg)
	}
	if got := FromPgTimestamptz(pg); !got.Equal(src) {
		t.Fatalf("FromPgTimestamptz=%v", got)
	}
}

func TestTimestamptz_Ptr(t *testing.T) {
	if pg := ToPgTimestamptzPtr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	src := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	pg := ToPgTimestamptzPtr(&src)
	got := FromPgTimestamptzPtr(pg)
	if got == nil || !got.Equal(src) {
		t.Fatalf("FromPgTimestamptzPtr=%v", got)
	}
	if FromPgTimestamptzPtr(pgtype.Timestamptz{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestTimestamptz_IsNil(t *testing.T) {
	if !IsNilPgTimestamptz(pgtype.Timestamptz{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgTimestamptz(pgtype.Timestamptz{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestTimestamptz_Bulk(t *testing.T) {
	src := []time.Time{
		time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC),
	}
	out := ToPgTimestamptzs(src)
	if len(out) != 2 || !out[0].Time.Equal(src[0]) || !out[1].Time.Equal(src[1]) {
		t.Fatalf("ToPgTimestamptzs: %+v", out)
	}
	if ToPgTimestamptzs(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}

func TestTimestamp_RoundTrip(t *testing.T) {
	src := time.Date(2026, 4, 30, 12, 30, 0, 0, time.UTC)
	pg := ToPgTimestamp(src)
	if !pg.Valid {
		t.Fatal("Valid expected true")
	}
	if got := FromPgTimestamp(pg); !got.Equal(src) {
		t.Fatalf("FromPgTimestamp=%v", got)
	}
}

func TestTimestamp_Ptr(t *testing.T) {
	if pg := ToPgTimestampPtr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	src := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	pg := ToPgTimestampPtr(&src)
	got := FromPgTimestampPtr(pg)
	if got == nil || !got.Equal(src) {
		t.Fatalf("FromPgTimestampPtr=%v", got)
	}
	if FromPgTimestampPtr(pgtype.Timestamp{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestTimestamp_IsNil(t *testing.T) {
	if !IsNilPgTimestamp(pgtype.Timestamp{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgTimestamp(pgtype.Timestamp{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestTimestamp_Bulk(t *testing.T) {
	src := []time.Time{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	out := ToPgTimestamps(src)
	if len(out) != 1 || !out[0].Time.Equal(src[0]) {
		t.Fatalf("ToPgTimestamps: %+v", out)
	}
	if ToPgTimestamps(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}

func TestDate_RoundTrip(t *testing.T) {
	src := time.Date(2026, 4, 30, 0, 0, 0, 0, time.UTC)
	pg := ToPgDate(src)
	if !pg.Valid {
		t.Fatal("Valid expected true")
	}
	if got := FromPgDate(pg); !got.Equal(src) {
		t.Fatalf("FromPgDate=%v", got)
	}
}

func TestDate_Ptr(t *testing.T) {
	if pg := ToPgDatePtr(nil); pg.Valid {
		t.Fatal("nil must produce Valid:false")
	}
	src := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	pg := ToPgDatePtr(&src)
	got := FromPgDatePtr(pg)
	if got == nil || !got.Equal(src) {
		t.Fatalf("FromPgDatePtr=%v", got)
	}
	if FromPgDatePtr(pgtype.Date{}) != nil {
		t.Error("invalid must produce nil")
	}
}

func TestDate_IsNil(t *testing.T) {
	if !IsNilPgDate(pgtype.Date{}) {
		t.Error("zero must be nil")
	}
	if IsNilPgDate(pgtype.Date{Valid: true}) {
		t.Error("Valid:true must be non-nil")
	}
}

func TestDate_Bulk(t *testing.T) {
	src := []time.Time{time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)}
	out := ToPgDates(src)
	if len(out) != 1 || !out[0].Time.Equal(src[0]) {
		t.Fatalf("ToPgDates: %+v", out)
	}
	if ToPgDates(nil) != nil {
		t.Error("nil input must produce nil output")
	}
}
