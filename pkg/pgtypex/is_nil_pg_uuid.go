//ff:func feature=pkg-pgtypex type=util control=sequence
//ff:what pgtype.UUID 가 SQL NULL 인지 검사한다 (!v.Valid alias)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func IsNilPgUUID(v pgtype.UUID) bool {
	return !v.Valid
}
