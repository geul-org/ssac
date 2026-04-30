//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *time.Time 을 pgtype.Timestamptz 로 wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgTimestamptzPtr(v *time.Time) pgtype.Timestamptz {
	if v == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *v, Valid: true}
}
