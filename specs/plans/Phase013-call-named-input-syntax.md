✅ 완료

# Phase 013: @call args를 JSON 형태 named field로 변경

## 목표

`@call`의 args를 positional `(arg1, arg2)`에서 `@state`, `@auth`와 동일한 `({Key: value})` 형태로 변경한다. `callFieldName()` 추론 제거.

```go
// 변경 전
// @call string accessToken = auth.IssueToken(user.ID, user.Email, user.Role)
// → auth.IssueTokenRequest{ID: user.ID, Email: user.Email, Role: user.Role}  ← ID가 UserID여야 함

// 변경 후
// @call string accessToken = auth.IssueToken({UserID: user.ID, Email: user.Email, Role: user.Role})
// → auth.IssueTokenRequest{UserID: user.ID, Email: user.Email, Role: user.Role}
```

## 변경 파일

| 파일 | 내용 |
|---|---|
| `parser/parser.go` | `@call` 파싱: `({Key: val, ...})` 형태 → Inputs map에 저장 |
| `parser/parser_test.go` | `@call` 파싱 테스트 업데이트 (named field 문법) |
| `generator/go_target.go` | `buildCallInputFields()`: Inputs map 기반으로 변경, `callFieldName()` 삭제 |
| `generator/generator_test.go` | `@call` 코드젠 테스트 업데이트 |
| `validator/validator.go` | `@call` args 변수 흐름 검증 조정 (Inputs 기반) |
| `specs/dummy-study/service/**/*.go` | 기존 @call을 새 문법으로 마이그레이션 |

## 설계

### 파서 변경

`@call`의 args 부분이 `({...})`이면 `@state`/`@auth`와 동일한 `parseInputs()`로 파싱하여 `seq.Inputs`에 저장. `seq.Args`는 비워둔다.

```go
// @call string token = auth.IssueToken({UserID: user.ID, Email: user.Email})
// → seq.Inputs = {"UserID": "user.ID", "Email": "user.Email"}
// → seq.Args = nil
```

### 코드젠 변경

`buildCallInputFields()`가 이미 `@state`/`@auth`에서 사용하는 `buildInputFieldsFromMap(seq.Inputs)`와 동일한 패턴이 된다. `buildTemplateData()`에서 `seq.Type == SeqCall`일 때 `seq.Inputs`로 InputFields 생성.

### 하위 호환

positional 문법 `(arg1, arg2)`는 더 이상 지원하지 않음. `@call`에는 반드시 `({...})` 형태 필수.

## 검증

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
ssac gen specs/dummy-study/ /tmp/ssac-phase13-check/
```

## 의존성

- 수정지시서v2/007
- Phase 011 (현재 callFieldName 추론 방식) 덮어씀
