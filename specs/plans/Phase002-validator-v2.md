# Phase 002: v2 밸리데이터 — 내부 + 외부 검증

## 목표

v2 IR(`[]ServiceFunc`)에 대한 내부 검증과 외부 SSOT 교차 검증을 구현한다.
`validator/symbol.go`(심볼 테이블 로딩)는 v1에서 그대로 복사하여 재사용.

## 변경 사항

### 1. validator/symbol.go — v1에서 복사

v1의 `validator/symbol.go`를 그대로 가져온다. DDL/OpenAPI/sqlc/model 파싱은 문법과 무관.

변경 없음:
- `SymbolTable`, `ModelSymbol`, `MethodInfo`, `DDLTable`, `OperationSymbol` 등 구조체
- `LoadSymbolTable()`, `loadDDL()`, `loadSqlcQueries()`, `loadOpenAPI()`, `loadGoInterfaces()`
- `parseDDLTables()`, `parseCreateIndex()`, `parseInlineFK()` 등 유틸

### 2. validator/errors.go — v1에서 복사

`ValidationError` 구조체 그대로 재사용.

### 3. validator/validator.go — 신규 작성

#### 내부 검증 (`Validate`)

| 시퀀스 타입 | 검증 규칙 |
|---|---|
| `get`, `post` | Model 필수, Result 필수, Args 1개 이상 |
| `put`, `delete` | Model 필수, Result nil이어야 함, Args 1개 이상 |
| `empty`, `exists` | Target 필수, Message 필수 |
| `state` | DiagramID 필수, Inputs 1개 이상, Transition 필수, Message 필수 |
| `auth` | Action 필수, Resource 필수, Message 필수 |
| `call` | Model 필수 (package.Func 형식), Args 0개 이상 |
| `response` | Fields 1개 이상 |

변수 흐름 검증:
- `@get`/`@post`/`@call`의 `Result.Var`가 선언
- `@empty`/`@exists`의 `Target`이 선언된 변수 참조
- Args의 `Source`가 선언된 변수이거나 `request`/`currentUser`
- `@response` Fields의 값이 선언된 변수 참조
- `@state`/`@auth` Inputs의 값이 선언된 변수 참조

#### 외부 검증 (`ValidateWithSymbols`)

v1에서 로직 이관 + v2 구조에 맞게 적응:

- Model/Method 존재 확인 (sqlc 쿼리, Go interface)
- Request 필드 존재 (OpenAPI request ↔ Args에서 source=="request")
- Response 필드 매핑 (OpenAPI response ↔ @response Fields 키)
- 정방향 + 역방향 검증 (Phase 014 로직 이관)
- @call 패키지 함수 존재 확인
- Stale 데이터 경고 (put/delete 후 response에서 재조회 없이 사용)

## 생성 파일

| 파일 | 내용 |
|---|---|
| `validator/errors.go` | ValidationError (v1 복사) |
| `validator/symbol.go` | 심볼 테이블 (v1 복사) |
| `validator/validator.go` | 내부 + 외부 검증 로직 |
| `validator/validator_test.go` | 테스트 |

## 테스트 케이스

1. 타입별 필수 필드 누락 검증
2. 변수 흐름: 미선언 변수 참조 에러
3. 외부: Model/Method 미존재
4. 외부: Request/Response 필드 불일치 (정방향 + 역방향)
5. 외부: Stale 데이터 경고
6. dummy-study 전체 검증 통과

## 의존성

- Phase 001 (parser IR)

## 검증

```bash
go test ./validator/... -count=1
```
