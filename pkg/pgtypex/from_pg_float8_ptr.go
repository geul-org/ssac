//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Float8 을 *float64 로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgFloat8Ptr(v pgtype.Float8) *float64 {
	if !v.Valid {
		return nil
	}
	out := v.Float64
	return &out
}
