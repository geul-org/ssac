✅ 완료

# Phase 008: @call 코드젠 수정 + validateCurrentUserType 제거

## 목표

수정지시서v2/002의 3가지 문제를 해결한다:

1. guard-style `@call` (result 없음)이 2-value return을 무시하여 컴파일 에러
2. `@call` 구조체 리터럴의 필드명이 소스 필드명을 그대로 사용 (위치 기반이어야 함)
3. `validateCurrentUserType` 제거 (SSaC의 책임이 아님)

## 설계

### 문제 1: guard-style @call에 `_, err` 패턴

현재 `call_no_result` 템플릿:
```
if err := pkg.Func(pkg.FuncRequest{...}); err != nil {
```

func spec 함수는 항상 `(Response, error)` 2개 값을 반환하므로 `_, err` 필요:
```
if _, err := pkg.Func(pkg.FuncRequest{...}); err != nil {
```

**변경**: `go_templates.go`의 `call_no_result` 템플릿에 `_,` 추가.

### 문제 2: @call 위치 기반 필드명 매핑

현재 `buildCallInputFields`는 `Arg.Field`를 구조체 필드명으로 사용한다. 하지만 Arg는 위치 기반이므로, Request 구조체의 실제 필드명을 참조해야 한다.

SSaC가 외부 패키지의 Request 구조체 정보를 가지고 있지 않으므로, **위치 기반 초기화(unkeyed literal)**를 사용한다:

```go
// 현재 (오류)
auth.IssueTokenRequest{ID: user.ID, Email: user.Email, Role: user.Role}

// 수정 (위치 기반)
auth.IssueTokenRequest{user.ID, user.Email, user.Role}
```

**변경**: `buildCallInputFields`에서 `fieldName + ":"` 제거, 값만 나열.

**단, 리터럴 인자도 값으로 처리**:
- `"cancelled"` → `"cancelled"` (따옴표 포함)
- `user.ID` → `user.ID`

### 문제 3: validateCurrentUserType 제거

CurrentUser는 인증 인프라(fullend.yaml claims)가 담당. SSaC는 예약 소스라는 것만 알면 됨.

**변경**:
- `validator.go`: `validateCurrentUserType` 함수 + 호출 제거
- `symbol.go`: `HasCurrentUserType` 필드 + CurrentUser 파싱 로직 제거
- `validator_test.go`: `TestValidateCurrentUserTypeMissing`, `TestValidateCurrentUserTypeExists` 제거

## 변경 파일

| 파일 | 내용 |
|---|---|
| `generator/go_templates.go` | `call_no_result`: `_, err` 패턴 추가 |
| `generator/go_target.go` | `buildCallInputFields`: 위치 기반 초기화 (unkeyed) |
| `generator/generator_test.go` | @call 관련 테스트 업데이트 |
| `validator/validator.go` | `validateCurrentUserType` 함수 + 호출 제거 |
| `validator/symbol.go` | `HasCurrentUserType` 필드 + 파싱 로직 제거 |
| `validator/validator_test.go` | CurrentUser 관련 테스트 제거 |

## 검증

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
ssac gen specs/dummy-study/ /tmp/ssac-phase8-check/
```

## 의존성

- Phase 006 (reserved sources, @call 코드젠)
- 수정지시서v2/002
