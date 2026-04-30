//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Float4 를 float32 로 unwrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgFloat4(v pgtype.Float4) float32 {
	return v.Float32
}
