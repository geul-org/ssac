//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL OpenAPI UUID 를 pgtype.UUID (Valid:true) 로 wrap 한다
package pgtypex

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
)

func ToPgUUID(v types.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: v, Valid: true}
}
