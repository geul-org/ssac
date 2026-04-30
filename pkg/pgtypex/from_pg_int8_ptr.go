//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Int8 을 *int64 로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgInt8Ptr(v pgtype.Int8) *int64 {
	if !v.Valid {
		return nil
	}
	out := v.Int64
	return &out
}
