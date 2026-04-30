//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Int4 를 *int32 로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgInt4Ptr(v pgtype.Int4) *int32 {
	if !v.Valid {
		return nil
	}
	out := v.Int32
	return &out
}
