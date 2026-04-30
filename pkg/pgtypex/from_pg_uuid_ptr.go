//ff:func feature=pkg-pgtypex type=adapter control=sequence
//ff:what nullable pgtype.UUID 를 *OpenAPI UUID 로 unwrap 한다 (!Valid → nil)
package pgtypex

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oapi-codegen/runtime/types"
)

func FromPgUUIDPtr(v pgtype.UUID) *types.UUID {
	if !v.Valid {
		return nil
	}
	out := types.UUID(v.Bytes)
	return &out
}
