//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Numeric 가 SQL NULL 인지 검사한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgNumeric(v pgtype.Numeric) bool {
	return !v.Valid
}
