//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Text 를 *string 으로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgTextPtr(v pgtype.Text) *string {
	if !v.Valid {
		return nil
	}
	out := v.String
	return &out
}
