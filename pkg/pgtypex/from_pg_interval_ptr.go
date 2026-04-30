//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Interval 을 *string 으로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgIntervalPtr(v pgtype.Interval) *string {
	if !v.Valid {
		return nil
	}
	out := formatIntervalISO(v)
	return &out
}
