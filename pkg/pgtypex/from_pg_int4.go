//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Int4 를 int32 로 unwrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgInt4(v pgtype.Int4) int32 {
	return v.Int32
}
