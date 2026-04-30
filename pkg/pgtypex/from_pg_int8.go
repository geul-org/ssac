//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Int8 을 int64 로 unwrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgInt8(v pgtype.Int8) int64 {
	return v.Int64
}
