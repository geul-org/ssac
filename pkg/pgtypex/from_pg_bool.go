//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.Bool 을 bool 로 unwrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgBool(v pgtype.Bool) bool {
	return v.Bool
}
