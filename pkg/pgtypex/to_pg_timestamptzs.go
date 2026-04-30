//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []time.Time 을 []pgtype.Timestamptz 로 bulk wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgTimestamptzs(vs []time.Time) []pgtype.Timestamptz {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Timestamptz, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Timestamptz{Time: v, Valid: true}
	}
	return out
}
