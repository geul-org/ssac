# ✅ 완료 Phase 022: @auth 코드젠 변경 + @call 타입 검증 + 미사용 변수 _ 처리 + config.* 코드젠

## 목표

A. `@auth` 코드젠을 `@call` 방식으로 변경 — `authz.Check(authz.CheckRequest{...})`
B. `@call` inputs의 필드 타입을 func Request struct와 비교 검증
C. 미참조 result 변수를 `_`로 생성
D. `config.Field` → `config.Get("UPPER_SNAKE")` 코드젠 변환

## 변경 파일 목록

### 1. 코드젠

| 파일 | 변경 |
|---|---|
| `generator/go_templates.go` | `auth`/`sub_auth` 템플릿 변경, `get`/`post`/`call_with_result`/`sub_get`/`sub_post`/`sub_call_with_result` 템플릿에 Unused 분기 추가 |
| `generator/go_target.go` | `needsCurrentUser()`: `SeqAuth` 무조건 true 제거, `inputValueToCode()`: config 변환, `toUpperSnake()` 헬퍼, config import, 미사용 변수 스캔, `templateData.Unused` 필드 |
| `generator/generator.go` | `toUpperSnake()` 유틸 (또는 go_target.go에 배치) |

### 2. 검증

| 파일 | 변경 |
|---|---|
| `validator/symbol.go` | `MethodInfo.ParamTypes map[string]string` 추가, `parseGoInterfaces()`에서 Request struct 파싱 |
| `validator/validator.go` | `validateCallInputTypes()` 추가 — @call inputs 타입 비교 |

### 3. 테스트

| 파일 | 변경 |
|---|---|
| `generator/generator_test.go` | @auth 새 코드젠, 미사용 변수 `_`, config.Get 변환 |
| `validator/validator_test.go` | @call 타입 불일치 ERROR |

### 4. 문서

| 파일 | 변경 |
|---|---|
| `artifacts/manual-for-ai.md` | @auth 코드젠, config.Get, 미사용 변수, 테스트 수 |
| `artifacts/manual-for-human.md` | @auth 코드젠 예시, config 사용법, 미사용 변수 |
| `README.md` | @auth 코드젠 설명, config 코드젠, 테스트 수 |
| `CLAUDE.md` | @auth 코드젠, config 코드젠, 테스트 수 |

## 상세 설계

### A. @auth → @call 방식 코드 생성

**현재 템플릿**:
```go
if err := authz.Check(currentUser, "{{.Action}}", "{{.Resource}}", authz.Input{ {{.InputFields}} }); err != nil {
    c.JSON(http.StatusForbidden, gin.H{"error": "{{.Message}}"})
    return
}
```

**변경 템플릿**:
```go
if _, err {{if .FirstErr}}:={{else}}={{end}} authz.Check(authz.CheckRequest{Action: "{{.Action}}", Resource: "{{.Resource}}", {{.InputFields}} }); err != nil {
    c.JSON(http.StatusForbidden, gin.H{"error": "{{.Message}}"})
    return
}
```

Subscribe 버전도 동일 패턴으로 변경 (`return fmt.Errorf(...)` 에러).

**`needsCurrentUser()` 변경**: `seq.Type == parser.SeqAuth` 무조건 true 제거. @auth inputs에 `currentUser.ID` 등이 있으면 기존 Inputs 체크로 감지됨.

**err 추적 변경**: @auth가 `if _, err :=` 패턴을 사용하므로, err 추적은 `@call` no-result와 동일하게 처리 (이미 현재와 동일).

**import 변경**: `authz` import는 그대로 유지 (패키지명 동일).

### B. @call 입력 타입 검증

**MethodInfo 확장**:
```go
type MethodInfo struct {
    Cardinality string
    Params      []string            // 파라미터명
    ParamTypes  map[string]string   // 파라미터명 → Go 타입 (e.g. "amount" → "int")
}
```

**symbol.go 변경**: `parseGoInterfaces()`에서 interface 메서드의 파라미터 타입을 추출하여 `ParamTypes`에 저장.

현재 interface 파싱 로직에서 파라미터 이름만 추출하는 부분을 타입도 함께 추출하도록 확장:
```go
// 현재: Params = append(Params, paramName)
// 변경: Params = append(Params, paramName); ParamTypes[paramName] = paramType
```

**validator.go 변경**: `validateCallInputTypes()` 함수 추가.

SSaC inputs의 각 값에서 타입을 결정:
- `request.Field` → DDL 역추적 (`resolveInputParamType`과 유사 로직)
- `variable.Field` → 해당 result 변수의 모델 테이블에서 Field 컬럼 타입
- `"literal"` → `string`
- `currentUser.Field` → model/*.go의 CurrentUser struct 필드 타입
- `config.Field` → `string` (config.Get 반환 타입)

비교: SSaC에서 결정한 타입 ≠ interface ParamTypes의 타입 → ERROR. 모든 타입은 DDL, model, 또는 고정값에서 확정 가능하므로 스킵 없이 무조건 비교.

### C. 미사용 변수 _ 처리

**사전 스캔**: `generateHTTPFunc()`/`generateSubscribeFunc()`에서 시퀀스 루프 전에 참조 분석:

```go
func collectUsedVars(seqs []parser.Sequence) map[string]bool {
    used := map[string]bool{}
    for _, seq := range seqs {
        // Guard Target
        if seq.Target != "" {
            used[rootVar(seq.Target)] = true
        }
        // Inputs values
        for _, val := range seq.Inputs {
            if !strings.HasPrefix(val, "request.") && !strings.HasPrefix(val, "currentUser.") &&
               !strings.HasPrefix(val, "config.") && !strings.HasPrefix(val, `"`) && val != "query" {
                used[rootVar(val)] = true
            }
        }
        // Response Fields values
        for _, val := range seq.Fields {
            if !strings.HasPrefix(val, `"`) {
                used[rootVar(val)] = true
            }
        }
    }
    return used
}
```

**templateData 확장**:
```go
type templateData struct {
    // ...
    Unused bool // result 변수가 이후 미참조 → _ 사용
}
```

**buildTemplateData()**: result가 있고 `usedVars[result.Var] == false`이면 `d.Unused = true`.

**템플릿 변경** (get, post, call_with_result):
```go
{{if .Unused}}_{{else}}{{.Result.Var}}{{end}}, err := ...
```

### D. config.* → config.Get() 코드젠 변환

**`inputValueToCode()` 변경**:
```go
func inputValueToCode(val string) string {
    if val == "query" { return "opts" }
    if strings.HasPrefix(val, "request.") { return lcFirst(val[len("request."):]) }
    if strings.HasPrefix(val, "config.") {
        key := val[len("config."):]
        return `config.Get("` + toUpperSnake(key) + `")`
    }
    return val
}
```

**`toUpperSnake()` 헬퍼**: PascalCase → UPPER_SNAKE_CASE
```go
func toUpperSnake(s string) string {
    // SMTPHost → SMTP_HOST
    // DatabaseURL → DATABASE_URL
    // AppPort → APP_PORT
    // toSnakeCase(s) 재활용 후 strings.ToUpper
}
```

기존 `toSnakeCase()` (validator에 있음)를 generator에서도 사용하거나 복제 후 `strings.ToUpper()` 적용.

**import 추가**: `collectImports()`에서 config 참조 존재 시 `"config"` 추가:
```go
func needsConfig(seqs []parser.Sequence) bool {
    for _, seq := range seqs {
        for _, val := range seq.Inputs {
            if strings.HasPrefix(val, "config.") {
                return true
            }
        }
    }
    return false
}
```

Subscribe import에도 동일 적용.

## 테스트 계획

### Generator

| 테스트 | 검증 내용 |
|---|---|
| `TestGenerateAuthCallStyle` | `authz.Check(authz.CheckRequest{Action: ..., Resource: ..., ...})`, 403 |
| `TestGenerateAuthNoCurrentUser` | @auth inputs에 currentUser 없으면 currentUser 추출 코드 없음 |
| `TestGenerateUnusedVar` | 미참조 result → `_, err :=` |
| `TestGenerateUsedVar` | 참조되는 result → `varName, err :=` |
| `TestGenerateConfigGet` | `config.SMTPHost` → `config.Get("SMTP_HOST")`, import "config" |
| `TestGenerateSubscribeAuth` | subscribe 함수 내 @auth → `return fmt.Errorf(...)` |

### Validator

| 테스트 | 검증 내용 |
|---|---|
| `TestValidateCallInputTypeMismatch` | @call input 타입 불일치 → ERROR |
| `TestValidateCallInputTypeMatch` | @call input 타입 일치 → OK |

## 의존성

- 기존 generator/validator 인프라
- `toSnakeCase()` (validator에 존재, generator에서 재사용 또는 복제)

## 주의사항

1. **@auth err 추적**: 현재 @auth는 `if err :=` 인라인. 변경 후 `if _, err :=` 인라인으로 변경 — err 추적 로직은 기존 @call no-result와 동일하게 유지.
2. **기존 @auth 테스트**: @auth 관련 기존 generator/validator 테스트가 새 코드젠 형태에 맞게 업데이트 필요.
3. **B 항목 타입 결정**: 모든 소스의 타입이 확정적 — DDL(request, variable), model(currentUser), 고정값(literal→string, config→string). 타입 불일치면 무조건 ERROR.
4. **config import 경로**: `"config"` (shorthand). fullend가 제공하는 `pkg/config` 패키지.
5. **toUpperSnake**: 연속 대문자 그룹은 하나로 유지 (`SMTP` → `SMTP`, `URL` → `URL`). 기존 `toSnakeCase()` + `strings.ToUpper()`로 구현.

## 검증 방법

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
```
