# SSaC v2 — AI Compact Reference

## CLI

```
ssac parse [dir]              # Print parsed sequence structure (default: specs/backend/service/)
ssac validate [dir]           # Internal validation or external SSOT cross-validation (auto-detect)
ssac gen <service-dir> <out>  # validate → codegen → gofmt (with symbol table: type conversion + model interface generation)
```

## Tech Stack

Go 1.24+, `go/ast` (parsing), `text/template` (codegen), `gopkg.in/yaml.v3` (OpenAPI), `github.com/gin-gonic/gin` (generated code target)

## DSL Syntax — One Line Per Sequence

10 sequence types. Each is a single comment line (except `@response` which is a multi-line block). Service files use `.ssac` extension (not `.go`).

### CRUD — Model Operations

```go
// @get Type var = Model.Method({Key: value, ...})        — Query (result required)
// @get Page[Type] var = Model.Method({Key: value, ...})  — Paginated query (Page or Cursor wrapper)
// @post Type var = Model.Method({Key: value, ...})       — Create (result required)
// @put Model.Method({Key: value, ...})                   — Update (no result)
// @delete Model.Method({Key: value, ...})                — Delete (no result)
```

**Package prefix model**: `pkg.Model.Method({...})` — for non-DDL models (session, cache, file, external).
- Lowercase first segment = package prefix: `session.Session.Get(...)` → Package="session", Model="Session.Get"
- No prefix (uppercase start) = DDL table model: `User.FindByID(...)` → Package="", Model="User.FindByID"
- Parser IR: `Sequence.Package` field stores package prefix (empty string if none)
- Package models are validated against Go interfaces (`st.Models["pkg.Model"]`), not DDL tables
- Package models are excluded from `models_gen.go` generation
- Parameter matching: SSaC keys ↔ interface params validated (extra → ERROR, missing → ERROR). `context.Context` excluded.

**Generic result types**: `Page[T]` and `Cursor[T]` wrappers for paginated results.
- Parser IR: `Result.Wrapper` = `"Page"` or `"Cursor"`, `Result.Type` = inner type
- Model interface returns `(*pagination.Page[T], error)` or `(*pagination.Cursor[T], error)`
- No 3-tuple return when Wrapper is used (Page/Cursor struct contains total/cursor internally)

All sequence types use unified `{Key: value}` syntax for args (CRUD, @call, @state, @auth).

Value format: `source.Field` or `"literal"`
- `request.CourseID` — from HTTP request (reserved source)
- `course.InstructorID` — from previous result variable
- `currentUser.ID` — from auth context (reserved source)
- `config.APIKey` — from environment config (reserved source)
- `query` — QueryOpts (pagination/sort/filter), explicit in inputs (reserved source)
- `"cancelled"` — string literal

Reserved sources: `request`, `currentUser`, `config`, `query` — cannot be used as result variable names.

Parser IR: all sequence types use `seq.Inputs` (map[string]string). CRUD uses `seq.Inputs` not `seq.Args`.

Required elements per type:

| Type | Required |
|---|---|
| get | Model, Result (Inputs optional) |
| post | Model, Result, Inputs |
| put | Model, Inputs |
| delete | Model, Inputs (0-input WARNING, `@delete!` suppresses) |
| empty, exists | Target, Message |
| state | DiagramID, Inputs, Transition, Message |
| auth | Action, Resource, Message |
| call | Model (pkg.Func format) |
| response | (none, Fields optional) |

### WARNING Suppression (`!` suffix)

Append `!` to any sequence type to suppress WARNINGs for that sequence. ERRORs are unaffected.

```go
// @delete! Room.DeleteAll()              — Suppresses 0-arg WARNING
// @response! { room: room }              — Suppresses stale data WARNING
```

### Guards

```go
// @empty target "message"                      — Fail if nil/zero (404)
// @exists target "message"                     — Fail if not nil/zero (409)
```

Target: variable (`course`) or variable.field (`course.InstructorID`)

### State Transition

```go
// @state diagramID {key: var.Field, ...} "transition" "message"
```

- `{inputs}`: JSON-style input mapping to state diagram package
- Codegen: `err := {id}state.CanTransition({id}state.Input{...}, "transition")` (returns error)

### Auth — OPA Permission Check

```go
// @auth "action" "resource" {key: var.Field, ...} "message"
```

- `{inputs}`: JSON-style context for OPA policy (ownership, org, etc.)
- Codegen: `authz.Check(currentUser, "action", "resource", authz.Input{...})`
- `currentUser` auto-extracted from `c.MustGet("currentUser")`

### Call — External Function

```go
// @call Type var = package.Func({Key: value, ...})       — With result
// @call package.Func({Key: value, ...})                  — Without result (guard-style error)
```

- Package name from Go import declarations in spec file
- With result: 500 on error. Without result: 401 on error.

### Response — Field Mapping Block

```go
// @response {
//   fieldName: variable,
//   fieldName: variable.Member,
//   fieldName: "literal"
// }
```

- Maps model results to OpenAPI response schema field by field
- No runtime functions (`len` etc.) — aggregation belongs in SQL
- Permission-based response differences → separate service functions (no conditionals)

**Shorthand**: `@response varName` → `c.JSON(http.StatusOK, varName)` (direct struct return, no gin.H)
- Used with Page[T]/Cursor[T] types where the struct is returned directly
- Handler skips pagination import (model handles the type internally)

## Full Example

```go
package service

import "myapp/auth"

// @auth "cancel" "reservation" {id: request.ReservationID} "권한 없음"
// @get Reservation reservation = Reservation.FindByID({reservationID: request.ReservationID})
// @empty reservation "예약을 찾을 수 없습니다"
// @state reservation {status: reservation.Status} "cancel" "취소할 수 없습니다"
// @call Refund refund = billing.CalculateRefund({id: reservation.ID, startAt: reservation.StartAt, endAt: reservation.EndAt})
// @put Reservation.UpdateStatus({reservationID: request.ReservationID, status: "cancelled"})
// @get Reservation reservation = Reservation.FindByID({reservationID: request.ReservationID})
// @response {
//   reservation: reservation,
//   refund: refund
// }
func CancelReservation() {}
```

## Directory Structure

```
cmd/ssac/main.go                 # CLI entrypoint
parser/                          # Comments → []ServiceFunc
  types.go                       #   IR structs (ServiceFunc, Sequence, Arg, Result)
  parser.go                      #   One-line expression parser
validator/                       # Internal + external SSOT validation
  validator.go                   #   Validation rules
  symbol.go                     #   Symbol table (DDL, OpenAPI, sqlc, model)
  errors.go                      #   ValidationError
generator/                       # Target interface-based codegen
  target.go                      #   Target interface + DefaultTarget()
  go_target.go                   #   GoTarget: Go+gin code generation
  go_templates.go                #   Go+gin templates
  generator.go                   #   Wrappers (Generate, GenerateWith) + utils
specs/                           # Declarations (input, SSOT)
  dummy-study/                   #   Study room reservation demo project
    service/ db/ db/queries/ api/ model/ states/ policy/
  plans/                         #   Implementation plans
artifacts/                       # Documentation
v1/                              # Archived v1 code (reference only)
```

## External Validation Project Layout

Auto-detected by `ssac validate <project-root>`:
- `<root>/service/<domain>/*.go` — Sequence specs (domain subfolder required, flat service/*.go is ERROR)
- `<root>/db/*.sql` — DDL (CREATE TABLE → column types)
- `<root>/db/queries/*.sql` — sqlc queries (filename→model, `-- name: Method :cardinality`)
- `<root>/api/openapi.yaml` — OpenAPI 3.0 (operationId=function name, x-pagination/sort/filter/include)
- `<root>/model/*.go` — Go interface→model methods, `// @dto`→DTO without DDL table
- `<root>/states/*.md` — State diagram definitions (mermaid stateDiagram-v2)
- `<root>/policy/*.rego` — OPA Rego policy files

## Codegen Features (gin framework)

Generated code uses **gin** framework (`c *gin.Context`):
- Function signature: `func Name(c *gin.Context)`
- Error responses: `c.JSON(status, gin.H{"error": "msg"})`
- Success responses: `c.JSON(http.StatusOK, gin.H{...})` with field mapping, or `c.JSON(http.StatusOK, var)` for `@response var` shorthand
- Path params: `c.Param("Name")` + type conversion
- Request body: `c.ShouldBindJSON(&req)` (2+ request params, or 1+ in POST/PUT) or `c.Query("Name")` (single GET/DELETE)
- currentUser: `c.MustGet("currentUser").(*model.CurrentUser)` — auto-generated when @auth or args reference currentUser

Additional features when symbol table (external SSOT) is available:

- **Type conversion**: DDL column type → `strconv.ParseInt`, `time.Parse` with 400 early return
- **Guard value types**: Zero value comparison based on result type (int→`== 0`/`> 0`, pointer→`== nil`/`!= nil`)
- **Stale data warning**: WARNING when response uses variable after put/delete without re-fetch (suppressed by `@response!`)
- **`:=` vs `=` tracking**: Go variable re-declaration uses `=` for already-declared variables
- **Go naming conventions**: Initialism-aware `lcFirst`/`ucFirst` (e.g. `ID`→`id`, `URL`→`url`)
- **@dto tag**: `// @dto` annotated struct → skips DDL table matching
- **DDL FK/Index parsing**: REFERENCES (inline/constraint), CREATE INDEX → `DDLTable.ForeignKeys`, `DDLTable.Indexes`
- **QueryOpts**: `query` reserved source in args → `opts := QueryOpts{}` + `c.Query()` parsing. No implicit injection.
- **List 3-tuple return**: `query` arg + `[]Type` result → `result, total, err :=` (includes count). Not used with Page[T]/Cursor[T] wrappers.
- **Query cross-validation**: OpenAPI x-extensions ↔ SSaC `query` mismatch detection (ERROR/WARNING)
- **x-pagination type validation**: `offset` ↔ `Page[T]`, `cursor` ↔ `Cursor[T]` cross-check. No x-pagination + Wrapper → ERROR
- **Wrapper field validation**: `@response var` shorthand with Page[T] → OpenAPI must have `items`, `total`. Cursor[T] → `items`, `next_cursor`, `has_next`
- **Model interface derivation**: 3 SSOT sources → `<outDir>/model/models_gen.go`
  - sqlc: method names, cardinality (:one→`*T`, :many→`[]T`, :exec→`error`)
  - SSaC: all inputs included (request, currentUser, variable refs, literals→DDL reverse-mapping, query→`opts QueryOpts`)
  - OpenAPI x-: infrastructure params validated against SSaC `query` usage
- **Domain folder structure**: `service/<domain>/*.go` required (flat service/*.go is ERROR). `service/auth/login.go` → `Domain="auth"` → `outDir/auth/login.go`, `package auth`
- **@state codegen**: `@state {id} {inputs} "transition"` → `err := {id}state.CanTransition({id}state.Input{...}, "transition")` (error return), import `"states/{id}state"`
- **@auth codegen**: `@auth "action" "resource" {inputs}` → `authz.Check(currentUser, "action", "resource", authz.Input{...})`
- **@call codegen**: `@call pkg.Func({Key: value})` → `pkg.Func(pkg.FuncRequest{Key: value, ...})`. No result → `_, err` guard-style (401), with result → value-style (500)
- **Spec file imports**: Parser collects Go import declarations from spec files and passes them to generated code
- **Package prefix model**: `pkg.Model.Method({...})` → validates against Go interface in package path. Missing interface → WARNING, missing method → ERROR with available methods list. Parameter matching: SSaC keys ↔ interface params (`context.Context` excluded). Package models skip DDL check and `models_gen.go`
- **Go reserved word validation**: DDL column names that are Go keywords (`type`, `range`, `select`, etc.) → ERROR with table name and rename suggestion. Prevents `models_gen.go` compile errors.

Singularization rules (sqlc filename → model name): `ies`→`y`, `sses`→`ss`, `xes`→`x`, otherwise remove trailing `s`

## OpenAPI x- Extensions

Infrastructure parameters declared in OpenAPI x- extensions. SSaC specs only declare business parameters.

```yaml
/api/reservations:
  get:
    operationId: ListReservations
    x-pagination:                    # style: offset|cursor, defaultLimit, maxLimit
      style: offset
      defaultLimit: 20
      maxLimit: 100
    x-sort:                          # allowed columns, default, direction
      allowed: [start_at, created_at]
      default: start_at
      direction: desc
    x-filter:                        # allowed filter columns
      allowed: [status, room_id]
    x-include:                       # FK_column:ref_table.ref_column
      allowed: [room_id:rooms.id, user_id:users.id]
```

Codegen effects:
- SSaC spec must explicitly include `query` in inputs: `Model.List({..., query: query})`
- Model methods with `query` arg get `opts QueryOpts` parameter
- `query` arg + `[]Type` result → return type includes total count: `([]T, int, error)`
- `QueryOpts` struct auto-generated (Limit, Offset, Cursor, SortCol, SortDir, Filters)
- Cross-validation: OpenAPI x- present but SSaC missing `query` → WARNING; SSaC `query` without OpenAPI x- → ERROR

## Coding Conventions

- gofmt compliant, immediate error handling (early return)
- Filenames: snake_case, variables/functions: camelCase, types: PascalCase
- Go common initialisms: `ID`, `URL`, `HTTP`, `API` etc. — all-caps (exported) or all-lowercase (unexported first word)
- Tests: `go test ./parser/... ./validator/... ./generator/... -count=1`
- 125 tests: parser 34 + validator 61 + generator 30
