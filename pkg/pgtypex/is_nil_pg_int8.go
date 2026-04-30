//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Int8 이 SQL NULL 인지 검사한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgInt8(v pgtype.Int8) bool {
	return !v.Valid
}
