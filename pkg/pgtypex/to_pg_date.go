//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL time.Time 을 pgtype.Date (Valid:true) 로 wrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func ToPgDate(v time.Time) pgtype.Date {
	return pgtype.Date{Time: v, Valid: true}
}
