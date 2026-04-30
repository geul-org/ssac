//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *bool 을 pgtype.Bool 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgBoolPtr(v *bool) pgtype.Bool {
	if v == nil {
		return pgtype.Bool{}
	}
	return pgtype.Bool{Bool: *v, Valid: true}
}
