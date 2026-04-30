//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.Numeric 를 *string 으로 unwrap 한다 (!Valid → nil)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func FromPgNumericPtr(v pgtype.Numeric) *string {
	if !v.Valid {
		return nil
	}
	out := FromPgNumeric(v)
	return &out
}
