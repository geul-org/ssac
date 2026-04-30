//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []time.Time 을 []pgtype.Timestamp 로 bulk wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgTimestamps(vs []time.Time) []pgtype.Timestamp {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Timestamp, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Timestamp{Time: v, Valid: true}
	}
	return out
}
