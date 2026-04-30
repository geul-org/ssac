//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Float8 이 SQL NULL 인지 검사한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgFloat8(v pgtype.Float8) bool {
	return !v.Valid
}
