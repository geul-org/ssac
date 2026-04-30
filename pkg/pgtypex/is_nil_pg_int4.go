//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Int4 가 SQL NULL 인지 검사한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgInt4(v pgtype.Int4) bool {
	return !v.Valid
}
