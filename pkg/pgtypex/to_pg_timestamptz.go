//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL time.Time 을 pgtype.Timestamptz (Valid:true) 로 wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgTimestamptz(v time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: v, Valid: true}
}
