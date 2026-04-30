//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what NOT NULL pgtype.UUID 를 OpenAPI UUID 로 unwrap 한다 (caller 가 Valid 보장)
package pgtypex

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
)

func FromPgUUID(v pgtype.UUID) types.UUID {
	return types.UUID(v.Bytes)
}
