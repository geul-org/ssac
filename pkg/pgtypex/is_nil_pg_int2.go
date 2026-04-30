//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Int2 가 SQL NULL 인지 검사한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgInt2(v pgtype.Int2) bool {
	return !v.Valid
}
