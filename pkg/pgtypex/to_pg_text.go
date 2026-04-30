//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL string 을 pgtype.Text (Valid:true) 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgText(v string) pgtype.Text {
	return pgtype.Text{String: v, Valid: true}
}
