//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable OpenAPI UUID 포인터를 pgtype.UUID 로 wrap 한다 (nil → Valid:false)
package pgtypex

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
)

func ToPgUUIDPtr(v *types.UUID) pgtype.UUID {
	if v == nil {
		return pgtype.UUID{}
	}
	return pgtype.UUID{Bytes: *v, Valid: true}
}
