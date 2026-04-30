//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable *string 을 pgtype.Text 로 wrap 한다 (nil → Valid:false)
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgTextPtr(v *string) pgtype.Text {
	if v == nil {
		return pgtype.Text{}
	}
	return pgtype.Text{String: *v, Valid: true}
}
