//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Float4 를 *float32 로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgFloat4Ptr(v pgtype.Float4) *float32 {
	if !v.Valid {
		return nil
	}
	out := v.Float32
	return &out
}
