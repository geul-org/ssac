//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL time.Time 을 pgtype.Timestamp (Valid:true) 로 wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgTimestamp(v time.Time) pgtype.Timestamp {
	return pgtype.Timestamp{Time: v, Valid: true}
}
