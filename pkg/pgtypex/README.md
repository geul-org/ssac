# pkg/pgtypex

OpenAPI 표면 타입 (`runtime/types.UUID`, `time.Time`, `string`, `int64`, …) 과 sqlc pgx/v5 pgtype wrapper (`pgtype.UUID`, `pgtype.Timestamptz`, …) 사이를 닫는 양방향 bridge 라이브러리.

## 변경이력

- 2026-04-30: 초기 패키지. 12 PG wrapper 타입 (UUID / Text / Int8 / Int4 / Int2 / Bool / Timestamptz / Timestamp / Date / Numeric / Interval / Float8 / Float4) × 5 함수 family + slice / UUIDToString.

## 개요

yongol 이 코드젠한 SaaS 백엔드는 oapi-codegen 으로부터 OpenAPI 표면 타입을, sqlc pgx/v5 로부터 pgtype wrapper 타입을 동시에 받는다. 두 표면이 만나는 경계 — request → sqlc 인자, sqlc row → response, OPA Owners 입력, nil 가드 — 에서 wrap/unwrap 이 필요하다. 본 패키지는 그 변환을 결정론적 함수 셋으로 닫는다. 함수 본체에 프로젝트별 결정이 없으므로 SaaS 별 emit 이 아닌 ssac 라이브러리에 둔다 (BUG-041 처방).

## 공개 API (타입별 5+1 함수 family)

| 함수 | 시그니처 (UUID 예) | 용도 |
|---|---|---|
| `ToPg<T>` | `ToPgUUID(v types.UUID) pgtype.UUID` | NOT NULL 컬럼: OpenAPI → pgtype (`Valid: true`) |
| `ToPg<T>Ptr` | `ToPgUUIDPtr(v *types.UUID) pgtype.UUID` | NULLABLE 컬럼: nil → `Valid: false`, 값 → `Valid: true` |
| `FromPg<T>` | `FromPgUUID(v pgtype.UUID) types.UUID` | NOT NULL row: pgtype → OpenAPI (caller 가 Valid 보장) |
| `FromPg<T>Ptr` | `FromPgUUIDPtr(v pgtype.UUID) *types.UUID` | NULLABLE row: `!Valid` → nil |
| `IsNilPg<T>` | `IsNilPgUUID(v pgtype.UUID) bool` | nil 가드 (`@empty` SSaC 시퀀스 등) |
| `ToPg<T>s` | `ToPgUUIDs(vs []types.UUID) []pgtype.UUID` | bulk insert 인자용 |
| `UUIDToString` | `UUIDToString(v pgtype.UUID) string` | OPA Owners-map 등 string 직렬화 사이트 (UUID 전용) |

지원 타입 12 종 — UUID · Text · Int8 · Int4 · Int2 · Bool · Timestamptz · Timestamp · Date · Numeric · Interval · Float8 · Float4. JSONB / BYTEA 는 sqlc 와 OpenAPI 양쪽이 `[]byte` 로 동일하므로 bridge 불필요. 사용자 정의 ENUM / typed JSONB struct 는 SaaS 내부 `internal/typesx/` 별 Phase.

NUMERIC ↔ string 은 정밀도 보존을 위해 `pgtype.Numeric.MarshalJSON` 의 텍스트 출력을 그대로 통과시킨다 (NaN 은 `"NaN"`).

INTERVAL ↔ string 은 ISO 8601 duration (`PT1H30M0S`) 표기를 사용한다. `pgtype.Interval` 의 `Microseconds + Days + Months` 를 합쳐 변환하되, `Months` 는 days 로 환산할 수 없으므로 별도 컴포넌트로 직렬화한다.

## 사용 예시

```go
import (
    "github.com/jackc/pgx/v5/pgtype"
    "github.com/oapi-codegen/runtime/types"
    "github.com/park-jun-woo/ssac/pkg/pgtypex"
)

// request → sqlc 인자
row, err := qtx.OwnerLookupWorkflow(ctx, pgtypex.ToPgUUID(request.Id))

// sqlc row → response
resp.OwnerId = pgtypex.FromPgUUID(row.OwnerID)

// nullable
resp.UpdatedAt = pgtypex.FromPgTimestamptzPtr(row.UpdatedAt)

// nil 가드 (SSaC @empty 시퀀스)
if pgtypex.IsNilPgUUID(row.WorkflowID) {
    return c.JSON(404, ...)
}

// OPA Owners-map
authz.Check(authz.CheckRequest{
    Action:   "read",
    Resource: "workflow",
    Owners:   map[string]any{"workflow": pgtypex.UUIDToString(row.OwnerID)},
})
```

## 외부 의존성

- `github.com/jackc/pgx/v5/pgtype` — sqlc pgx/v5 wrapper 타입 출처. 메이저 업그레이드 시 SaaS go.mod 와 ssac go.mod 동시 마이그레이션 필요.
- `github.com/oapi-codegen/runtime/types` — `types.UUID` (= `[16]byte`) 출처. 메이저 업그레이드는 동일 절차.
