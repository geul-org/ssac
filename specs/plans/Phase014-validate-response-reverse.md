✅ 완료

# Phase 014: validateRequest·validateResponse 역방향 검증 추가

## 목표

`validateRequest`와 `validateResponse`에 역방향 검증을 추가한다.
OpenAPI에 정의된 필드가 SSaC에 없는 불일치를 감지한다.

## 변경 사항

### 1. validator/validator.go — `validateResponse` 함수 교체 (line 87-106)

기존 단방향 루프를 유지하면서 역방향 블록 추가:

- response sequence의 `@var` 목록을 `responseVars map[string]bool`로 수집 + `responseSeqIdx` 기록
- 기존 정방향 검증은 동일 (SSaC `@var` → OpenAPI `ResponseFields`)
- 루프 후 역방향 블록: `op.ResponseFields`의 각 필드가 `responseVars`에 있는지 확인
- `x-pagination != nil`이면 `"total"` 필드 스킵 (코드젠 자동 생성)
- 누락 시 `errCtx.err("@var", ...)` — ERROR (기본 Level)

### 2. validator/validator.go — `validateRequest` 함수 교체 (line 65-84)

기존 단방향 루프를 유지하면서 역방향 블록 추가:

- 정방향 루프에서 `usedRequestFields map[string]bool`에 `source == "request"` 파라미터 수집
- 루프 후 역방향 블록: `op.RequestFields`의 각 필드가 `usedRequestFields`에 있는지 확인
- `op.PathParams`에 있는 필드는 스킵 (라우팅 자동 바인딩)
- 누락 시 `ValidationError` 직접 구성 (`errCtx.err`는 Level 설정 불가):
  - `SeqIndex: -1` (특정 sequence에 귀속되지 않음)
  - `Level: "WARNING"` (optional 필드 가능성)

### 3. validator/validator_test.go — 테스트 4개 추가

기존 테스트 영향: `TestValidateWithSymbolsMissingRequestField`는 `assertContainsError`(포함 검사)만 하므로 역방향 WARNING이 추가되어도 깨지지 않음.

#### 3-1. TestValidateReverseResponseMissing
- SymbolTable: `Operations["Test"].ResponseFields = {"user": true, "instructor": true}`
- ServiceFunc: `response json` + `Vars: ["user"]` (instructor 누락)
- 기대: ERROR, `@var`, `"instructor"`

#### 3-2. TestValidateReverseResponsePaginationTotal
- SymbolTable: `Operations["Test"].ResponseFields = {"items": true, "total": true}`, `XPagination != nil`
- ServiceFunc: `response json` + `Vars: ["items"]` (total 누락)
- 기대: 에러 없음 (total은 x-pagination 시 자동 생성)

#### 3-3. TestValidateReverseRequestMissing
- SymbolTable: `Operations["Test"].RequestFields = {"Email": true, "Description": true}`
- ServiceFunc: `@param Email request`만 사용 (Description 누락)
- 기대: WARNING, `@param`, `"Description"`

#### 3-4. TestValidateReverseRequestPathParamSkip
- SymbolTable: `Operations["Test"].RequestFields = {"CourseID": true}`, `PathParams: [{Name: "CourseID"}]`
- ServiceFunc: `@param` request 없음
- 기대: 에러 없음 (path param은 역방향 검증 제외)

## 변경 파일 목록

| 파일 | 변경 |
|---|---|
| `validator/validator.go` | `validateResponse` 함수 교체 (line 87-106), `validateRequest` 함수 교체 (line 65-84) |
| `validator/validator_test.go` | 역방향 검증 테스트 4개 추가 |

## 의존성

- 없음

## 검증

```bash
go test ./validator/... -count=1
```
