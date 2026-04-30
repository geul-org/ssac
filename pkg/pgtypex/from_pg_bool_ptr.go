//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Bool 을 *bool 로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgBoolPtr(v pgtype.Bool) *bool {
	if !v.Valid {
		return nil
	}
	out := v.Bool
	return &out
}
