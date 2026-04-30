//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Timestamptz 를 *time.Time 으로 unwrap 한다 (!Valid → nil)
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FromPgTimestamptzPtr(v pgtype.Timestamptz) *time.Time {
	if !v.Valid {
		return nil
	}
	out := v.Time
	return &out
}
