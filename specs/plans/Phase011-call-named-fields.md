✅ 완료

# Phase 011: @call Request struct named field 생성 + trailing dot 수정

## 목표

`@call` 코드젠의 두 가지 버그를 수정한다:
1. Request struct 초기화를 positional(unkeyed)에서 named field로 변경
2. bare variable arg에서 trailing dot 제거

```go
// 변경 전 (잘못됨)
auth.HashPasswordRequest{ password }        // positional → 컴파일 에러
userModel.Create(email, hashedPassword., role, name)  // trailing dot

// 변경 후
auth.HashPasswordRequest{Password: password}  // named field
userModel.Create(email, hashedPassword, role, name)   // dot 제거
```

## 변경 파일

| 파일 | 내용 |
|---|---|
| `generator/go_target.go` | `argToCode()`: Field 빈 문자열일 때 Source만 반환 |
| `generator/go_target.go` | `buildCallInputFields()`: named field 형식 생성 (`FieldName: value`) |
| `generator/generator_test.go` | `TestGenerateCallWithResult`, `TestGenerateCallWithoutResult` assertion 업데이트 + trailing dot 테스트 추가 |

## 버그 상세

### 버그 1: `buildCallInputFields()` — positional → named field

현재 `argToCode(a)` 값만 나열. Go struct literal은 named field 필수.

필드명 도출 규칙:
- `request.Field` → 필드명 `Field` (a.Field)
- `source.Field` (변수 참조) → 필드명 `Field` (a.Field)
- `currentUser.Field` → 필드명 `Field` (a.Field)
- bare variable `hashedPassword` → 필드명 `HashedPassword` (ucFirst(a.Source))
- `"literal"` → 필드명 `ucFirst(literal)`

```go
func buildCallInputFields(args []parser.Arg) string {
    var fields []string
    for _, a := range args {
        name := callFieldName(a)
        value := argToCode(a)
        fields = append(fields, name+": "+value)
    }
    return strings.Join(fields, ", ")
}

func callFieldName(a parser.Arg) string {
    if a.Literal != "" {
        return ucFirst(a.Literal)
    }
    if a.Field != "" {
        return a.Field  // 이미 PascalCase
    }
    return ucFirst(a.Source)  // bare variable
}
```

### 버그 2: `argToCode()` — trailing dot

`Arg{Source: "hashedPassword", Field: ""}` 일 때 `a.Source + "." + a.Field` = `hashedPassword.`

```go
if a.Source != "" {
    if a.Field == "" {
        return a.Source
    }
    return a.Source + "." + a.Field
}
```

## 검증

```bash
go test ./parser/... ./validator/... ./generator/... -count=1
ssac gen specs/dummy-study/ /tmp/ssac-phase11-check/
```

## 의존성

- 수정지시서v2/005
- Phase 008 (현재 unkeyed 방식) 덮어씀
