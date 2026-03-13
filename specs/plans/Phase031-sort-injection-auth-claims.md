# Phase031: sort SQL Injection 수정 + auth Claims 추가 (BUG005, BUG007) ✅ 완료

## 배경

수정지시서002의 4건 중 실제 ssac 범위 2건만 수정.
- BUG006 (password_hash json 노출): ssac는 DDL→Go struct를 생성하지 않음 (model interface만 생성). 범위 밖.
- BUG008 (@call 에러 일괄 500): Phase026에서 이미 구현 완료 (`ErrStatus` 필드, 파서, 코드젠, 테스트 모두 있음).

## BUG005: sort 파라미터 SQL Injection

### 문제

`generateQueryOptsCode` (go_helpers.go:282)가 `c.Query("sort")` 값을 검증 없이 `opts.SortCol`에 대입하는 코드를 생성.

```go
// 현재 생성 코드
opts := QueryOpts{}
if v := c.Query("sort"); v != "" {
    opts.SortCol = v  // ← allowlist 없음
}
```

### 수정

x-sort가 있으면 허용 컬럼 allowlist를 생성하여 검증. allowlist에 없는 값은 무시.

```go
// 기대 생성 코드
opts := QueryOpts{}
if v := c.Query("limit"); v != "" {
    opts.Limit, _ = strconv.Atoi(v)
}
if v := c.Query("offset"); v != "" {
    opts.Offset, _ = strconv.Atoi(v)
}
allowedSort := map[string]bool{"created_at": true, "title": true}
if v := c.Query("sort"); allowedSort[v] {
    opts.SortCol = v
}
if v := c.Query("direction"); v == "asc" || v == "desc" {
    opts.SortDir = v
}
```

### 변경 파일

| 파일 | 변경 |
|---|---|
| `generator/go_helpers.go` | `generateQueryOptsCode` — x-sort allowlist 코드 생성, direction 검증 추가 |
| `generator/go_handler_test.go` | sort allowlist 검증 테스트 추가 |

### 상세

`generateQueryOptsCode(st)` → `generateQueryOptsCode(funcName, st)` 시그니처 변경.
해당 함수의 operationSymbol에서 `XSort.Columns`를 읽어 allowlist map 리터럴 생성.
direction은 `"asc"` / `"desc"` 고정 검증.

`go_handler.go`의 호출부도 `generateQueryOptsCode(sf.Name, st)` 으로 변경.

---

## BUG007: authz.Check에 Claims 미전달

### 문제

`@auth` 시퀀스에서 `authz.Check` 호출 시 Claims가 누락. OPA 소유권 검증 불가.

```go
// 현재 생성 코드
authz.Check(authz.CheckRequest{Action: "read", Resource: "project", ResourceID: project.ID})
```

### 수정

`@auth`에서 `currentUser.*`를 참조하면 `Claims: authz.Claims{UserID: currentUser.ID}` 자동 추가.

```go
// 기대 생성 코드
authz.Check(authz.CheckRequest{Action: "read", Resource: "project", ResourceID: project.ID, Role: currentUser.Role, Claims: authz.Claims{UserID: currentUser.ID}})
```

### 변경 파일

| 파일 | 변경 |
|---|---|
| `generator/go_helpers.go` | `buildTemplateData` — @auth + currentUser 참조 시 `Claims: authz.Claims{UserID: currentUser.ID}` 추가 |
| `generator/go_templates.go` | auth/sub_auth 템플릿에 `{{.ClaimsCode}}` 추가 |
| `generator/go_helpers.go` | `templateData` struct에 `ClaimsCode string` 필드 추가 |
| `generator/go_handler_test.go` | Claims 포함 검증 테스트 추가, 기존 auth 테스트 assertion 수정 |

### 상세

`buildTemplateData`에서 `@auth` + `hasCurrentUserRef(inputs)` 일 때:
- 기존: `Role: currentUser.Role` 자동 추가
- 추가: `d.ClaimsCode = "Claims: authz.Claims{UserID: currentUser.ID}, "` 설정

auth/sub_auth 템플릿:
```
authz.Check(authz.CheckRequest{Action: "{{.Action}}", Resource: "{{.Resource}}", {{.ClaimsCode}}{{.InputFields}} })
```

`currentUser` 참조 없는 `@auth`는 Claims 없이 기존 동작 유지.

---

## 의존성

- 없음 (기존 코드 수정만)

## 검증

```bash
go test ./generator/... -count=1
```

기존 테스트 전체 통과 + 각 버그별 신규 테스트 통과.
