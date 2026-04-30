//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL bool 을 pgtype.Bool (Valid:true) 로 wrap 한다
package pgtypex

import "github.com/jackc/pgx/v5/pgtype"

func ToPgBool(v bool) pgtype.Bool {
	return pgtype.Bool{Bool: v, Valid: true}
}
