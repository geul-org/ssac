# Phase 018: Go 예약어 파라미터명 검증

## 목표

DDL 컬럼명이 Go 예약어(`type`, `range`, `select` 등)인 경우, models_gen.go 코드젠 전에 validator에서 ERROR를 출력하여 컴파일 에러를 사전 차단.

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `validator/validator.go` | `validateGoReservedWords()` 함수 추가, `ValidateWithSymbols()`에서 호출 |
| `validator/validator_test.go` | 예약어 검증 테스트 추가 |

## 상세 설계

### 1. Go 예약어 목록

```go
var goReservedWords = map[string]bool{
    "break": true, "case": true, "chan": true, "const": true,
    "continue": true, "default": true, "defer": true, "else": true,
    "fallthrough": true, "for": true, "func": true, "go": true,
    "goto": true, "if": true, "import": true, "interface": true,
    "map": true, "package": true, "range": true, "return": true,
    "select": true, "struct": true, "switch": true, "type": true,
    "var": true,
}
```

### 2. 검증 위치

`ValidateWithSymbols()`에서 새 함수 호출:

```go
errs = append(errs, validateGoReservedWords(funcs, st)...)
```

### 3. 검증 로직

`validateGoReservedWords(funcs, st)`:
- SSaC spec의 모든 CRUD 시퀀스에서 `seq.Inputs` 키를 순회
- `lcFirst(key)`로 변환한 파라미터명이 Go 예약어이면 ERROR
- DDL 컬럼명(`toSnakeCase(key)`)을 역추적하여 에러 메시지에 테이블명 포함

실제 컴파일 에러가 발생하는 경로:
1. SSaC spec: `{type: request.Type}` → Inputs key = `type`
2. `deriveInterfaces()`: `lcFirst("type")` = `type` → 파라미터명
3. models_gen.go: `Create(amount int64, type string)` → **컴파일 에러**

따라서 Inputs 키를 `lcFirst()` 변환한 결과가 예약어인지 확인.

### 4. ERROR 메시지 형식

```
ERROR: DDL column "type" in table "transactions" is a Go reserved word — rename the column (e.g. "tx_type")
```

테이블을 특정할 수 없는 경우:
```
ERROR: parameter name "type" is a Go reserved word — rename the DDL column
```

## 테스트 계획

| 테스트 | 검증 내용 |
|---|---|
| `TestValidateGoReservedWordInInputs` | `{type: request.Type}` + DDL "type" 컬럼 → ERROR |
| `TestValidateGoReservedWordNoConflict` | `{txType: request.TxType}` → OK |
| `TestValidateGoReservedWordRange` | `{range: request.Range}` → ERROR |

## 의존성

- `generator/generator.go`의 `lcFirst()`, `toSnakeCase()` — validator에서 import 필요 (또는 로직 인라인)

## 검증 방법

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
```
