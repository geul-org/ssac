# ✅ 완료 Phase032: Handler struct 도입 + 자동 트랜잭션 래핑

## 배경

1. 현재 코드젠은 bare function (`func Name(c *gin.Context)`)과 package-level 모델 변수(`courseModel.FindByID(...)`)를 사용. 테스트 불가, DI 불가.
2. 쓰기 시퀀스가 여러 개인 핸들러에서 트랜잭션이 없어 데이터 불일치 발생 가능.

두 문제를 한 번에 해결: **Handler struct 도입 → 트랜잭션 래핑**.

## 설계 결정

| 항목 | 결정 |
|---|---|
| DB 참조 | Handler struct의 `DB *sql.DB` 필드 |
| 함수 시그니처 | `func (h *Handler) Name(c *gin.Context)` |
| 모델 참조 | `h.CourseModel.Method(...)` |
| tx 범위 | 쓰기(@post/@put/@delete) 하나라도 있으면 함수 전체 tx |
| tx 적용 대상 | 모든 CRUD 모델 호출 (읽기+쓰기 정합성) |
| @call | tx 밖 — 리소스 접근 금지 계약 |
| @publish | tx 안 — 사용자 설계 순서 유지 |
| @state, @auth | tx 밖 — DB 무관 |
| subscribe | 쓰기 있으면 동일하게 tx 적용 |
| 단일 DB | multi-DB / 2PC는 범위 밖 |

## 생성 코드 변경

### 현재

```go
func AcceptProposal(c *gin.Context) {
    proposal, err := proposalModel.FindByID(proposalID)
    // ...
    err = proposalModel.UpdateStatus(proposalID, "accepted")
    // ...
    tx, err := transactionModel.Create(gigID, freelancerID, amount)
    // ...
    c.JSON(200, gin.H{"proposal": proposal})
}
```

### 변경 후 (쓰기 있는 함수)

```go
func (h *Handler) AcceptProposal(c *gin.Context) {
    // ... param parsing, currentUser ...

    tx, err := h.DB.BeginTx(c.Request.Context(), nil)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "transaction failed"})
        return
    }
    defer tx.Rollback()

    proposal, err := h.ProposalModel.WithTx(tx).FindByID(proposalID)
    // ...
    err = h.ProposalModel.WithTx(tx).UpdateStatus(proposalID, "accepted")
    // ...
    transaction, err := h.TransactionModel.WithTx(tx).Create(gigID, freelancerID, amount)
    // ...

    if err := tx.Commit(); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "commit failed"})
        return
    }
    c.JSON(200, gin.H{"proposal": proposal})
}
```

### 변경 후 (읽기 전용 함수)

```go
func (h *Handler) GetCourse(c *gin.Context) {
    course, err := h.CourseModel.FindByID(courseID)
    // ...
    c.JSON(200, gin.H{"course": course})
}
```

### Handler struct (도메인별 생성)

```go
// <outDir>/<domain>/handler.go
package reservation

import (
    "database/sql"
    "model"
)

type Handler struct {
    DB               *sql.DB
    ProposalModel    model.ProposalModel
    TransactionModel model.TransactionModel
}
```

### Model interface (WithTx 추가)

```go
// <outDir>/model/models_gen.go
type ProposalModel interface {
    WithTx(tx *sql.Tx) ProposalModel
    FindByID(proposalID int64) (*Proposal, error)
    UpdateStatus(proposalID int64, status string) error
}
```

### Subscribe (쓰기 있는 경우)

```go
func (h *Handler) ProcessOrder(ctx context.Context, message OnOrderMessage) error {
    tx, err := h.DB.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("transaction failed: %w", err)
    }
    defer tx.Rollback()

    // model calls with .WithTx(tx)

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit failed: %w", err)
    }
    return nil
}
```

## 변경 파일

### Part A: Handler struct 도입

| 파일 | 변경 |
|---|---|
| `generator/go_handler.go` | 함수 시그니처 `func (h *Handler) Name(...)` 생성 |
| `generator/go_helpers.go` | `buildTemplateData`의 ModelCall 생성을 `h.` prefix로 변경 |
| `generator/go_target.go` | `GenerateHandlerStruct` 메서드 추가 — 도메인별 Handler struct 생성 |
| `generator/target.go` | `Target` interface에 `GenerateHandlerStruct` 추가 |
| `generator/generator.go` | `GenerateWith`에서 `GenerateHandlerStruct` 호출 추가 |
| `generator/go_params.go` | `collectImports`에 `database/sql` 조건부 추가 |

### Part B: 트랜잭션 래핑

| 파일 | 변경 |
|---|---|
| `generator/go_helpers.go` | `hasWriteSequence` 헬퍼 추가, `buildTemplateData`에 `useTx` 파라미터 추가 → ModelCall에 `.WithTx(tx)` 삽입 |
| `generator/go_handler.go` | `generateHTTPFunc` — BeginTx 블록 생성, @response 직전 Commit 삽입 |
| `generator/go_handler.go` | `generateSubscribeFunc` — BeginTx 블록 생성, `return nil` 직전 Commit 삽입 |
| `generator/go_templates.go` | 변경 없음 — ModelCall 필드가 이미 전체 호출 문자열을 담고 있으므로 WithTx는 ModelCall 생성 시 처리 |
| `generator/go_interface.go` | `renderInterfaces`에서 모든 interface에 `WithTx(*sql.Tx) <InterfaceName>` 메서드 추가, `database/sql` import 추가 |

### Part C: 테스트

| 파일 | 변경 |
|---|---|
| `generator/go_handler_test.go` | 기존 전체 테스트 assertion 수정 (`h.` prefix, `func (h *Handler)` 시그니처) + 트랜잭션 테스트 추가 |
| `generator/go_subscribe_test.go` | 기존 전체 테스트 assertion 수정 + subscribe tx 테스트 추가 |
| `generator/go_interface_test.go` | WithTx 메서드 포함 assertion 추가 |
| `generator/go_args_test.go` | assertion 수정 (ModelCall 변경에 따라) |

## 상세 구현

### hasWriteSequence

```go
func hasWriteSequence(seqs []parser.Sequence) bool {
    for _, seq := range seqs {
        switch seq.Type {
        case parser.SeqPost, parser.SeqPut, parser.SeqDelete:
            return true
        }
    }
    return false
}
```

### ModelCall 생성 변경 (go_helpers.go buildTemplateData)

현재:
```go
d.ModelCall = strcase.ToGoCamel(parts[0]) + "Model." + parts[1]
// → "courseModel.FindByID"
```

변경:
```go
modelRef := "h." + strcase.ToGoPascal(parts[0]) + "Model"
if useTx {
    modelRef += ".WithTx(tx)"
}
d.ModelCall = modelRef + "." + parts[1]
// 읽기 전용: "h.CourseModel.FindByID"
// tx 함수:  "h.CourseModel.WithTx(tx).FindByID"
```

### BeginTx 블록 삽입 위치 (go_handler.go)

HTTP 함수:
```
path params → currentUser → request params → QueryOpts → [BeginTx] → sequences → [Commit] → response
```

Subscribe 함수:
```
[BeginTx] → sequences → [Commit] → return nil
```

Commit은 시퀀스 루프에서 `@response`를 만났을 때 직전에 삽입. Subscribe는 루프 종료 후 `return nil` 직전에 삽입.

### errDeclared 영향

`BeginTx` 코드가 `err`를 선언하므로, tx 함수에서는 `errDeclared`를 `true`로 초기화:

```go
errDeclared := hasConversionErr(requestParams)
if useTx {
    errDeclared = true // BeginTx에서 err := 선언
}
```

### Handler struct 생성

도메인별로 사용되는 모든 모델을 수집:

```go
func collectDomainModels(funcs []parser.ServiceFunc) map[string][]string {
    // domain → []modelName (중복 제거, 정렬)
}
```

도메인이 없는 함수는 `"service"` 패키지. 각 도메인별 `handler.go` 파일 생성.

패키지 접두사 모델 (`session.Token.Create`)도 동일하게 `h.TokenModel`로 참조. Handler struct 필드 타입은 `model.TokenModel` (현재 models_gen.go에서 제외되는 패키지 접두사 모델은 변경 없음 — 사용자가 직접 interface 제공).

### collectImports 변경

tx 함수에서 `database/sql` import 추가:

```go
if hasWriteSequence(sf.Sequences) {
    imports = append(imports, "database/sql")
}
```

### WithTx interface 메서드 (go_interface.go)

```go
for _, iface := range interfaces {
    fmt.Fprintf(&buf, "type %s interface {\n", iface.Name)
    fmt.Fprintf(&buf, "\tWithTx(tx *sql.Tx) %s\n", iface.Name)
    for _, m := range iface.Methods {
        // 기존 메서드들
    }
    buf.WriteString("}\n\n")
}
```

`renderInterfaces`에서 `database/sql` import 무조건 추가 (모든 interface에 WithTx가 있으므로).

## 의존성

- 수정지시서001 (파일 분리) 완료 — 완료됨
- 외부 패키지 추가 없음

## 검증

```bash
go test ./generator/... -count=1
```

### 테스트 케이스

| 테스트 | 검증 |
|---|---|
| 기존 전체 테스트 | assertion 수정 후 통과 (h. prefix, receiver 시그니처) |
| TestGenerateTxWriteFunction | 쓰기 함수 → BeginTx + defer Rollback + Commit 코드 존재 |
| TestGenerateTxReadOnlyFunction | 읽기 전용 → BeginTx/Commit 없음, `h.Model.Method()` (WithTx 없음) |
| TestGenerateTxModelCallWithTx | tx 함수 내 @get도 `.WithTx(tx)` 포함 |
| TestGenerateTxCommitBeforeResponse | Commit이 @response 직전에 위치 |
| TestGenerateTxSubscribe | subscribe 쓰기 함수 → tx 래핑 |
| TestGenerateHandlerStruct | 도메인별 Handler struct에 DB + 사용 모델 필드 포함 |
| TestGenerateModelInterfaceWithTx | WithTx 메서드가 interface에 포함 |
