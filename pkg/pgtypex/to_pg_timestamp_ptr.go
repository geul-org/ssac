//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *time.Time 을 pgtype.Timestamp 로 wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgTimestampPtr(v *time.Time) pgtype.Timestamp {
	if v == nil {
		return pgtype.Timestamp{}
	}
	return pgtype.Timestamp{Time: *v, Valid: true}
}
