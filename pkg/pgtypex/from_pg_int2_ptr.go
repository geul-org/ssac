//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Int2 를 *int16 으로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgInt2Ptr(v pgtype.Int2) *int16 {
	if !v.Valid {
		return nil
	}
	out := v.Int16
	return &out
}
