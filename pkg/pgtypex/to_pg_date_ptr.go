//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *time.Time 을 pgtype.Date 로 wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgDatePtr(v *time.Time) pgtype.Date {
	if v == nil {
		return pgtype.Date{}
	}
	return pgtype.Date{Time: *v, Valid: true}
}
