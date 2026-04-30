//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what []time.Time 을 []pgtype.Date 로 bulk wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgDates(vs []time.Time) []pgtype.Date {
	if vs == nil {
		return nil
	}
	out := make([]pgtype.Date, len(vs))
	for i, v := range vs {
		out[i] = pgtype.Date{Time: v, Valid: true}
	}
	return out
}
