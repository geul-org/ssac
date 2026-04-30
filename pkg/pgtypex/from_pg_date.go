//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Date 를 time.Time 으로 unwrap 한다
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FromPgDate(v pgtype.Date) time.Time {
	return v.Time
}
