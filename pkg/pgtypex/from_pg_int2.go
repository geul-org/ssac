//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Int2 를 int16 으로 unwrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgInt2(v pgtype.Int2) int16 {
	return v.Int16
}
