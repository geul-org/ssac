//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.Float4 가 SQL NULL 인지 검사한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgFloat4(v pgtype.Float4) bool {
	return !v.Valid
}
