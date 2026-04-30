//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Timestamp 를 *time.Time 으로 unwrap 한다 (!Valid → nil)
package pgtypex

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func FromPgTimestampPtr(v pgtype.Timestamp) *time.Time {
	if !v.Valid {
		return nil
	}
	out := v.Time
	return &out
}
