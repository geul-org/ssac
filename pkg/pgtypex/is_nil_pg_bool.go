//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Bool 이 SQL NULL 인지 검사한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgBool(v pgtype.Bool) bool {
	return !v.Valid
}
