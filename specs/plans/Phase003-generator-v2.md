# Phase 003: v2 제너레이터 — gin 코드 생성

## 목표

v2 IR에서 gin 프레임워크 기반 Go 서비스 코드를 생성한다.
Target 인터페이스 구조는 v1에서 유지하되, 템플릿과 코드젠 로직은 v2 IR에 맞게 재작성.

## 변경 사항

### 1. generator/target.go — v1에서 복사

Target 인터페이스와 DefaultTarget() 유지.

### 2. generator/go_target.go — 재작성

v2 IR 기반 코드 생성. 핵심 차이:

| 항목 | v1 | v2 |
|---|---|---|
| 파라미터 추출 | Sequence.Params 순회 | Sequence.Args 순회, Source별 분기 |
| 변수 참조 | `@param Name varName` → 별도 해석 | `Args[i].Source.Field` → 직접 사용 |
| Response | `@var` 나열 → `gin.H{var: var}` | `Fields` 매핑 → `gin.H{key: value}` |
| Auth | action/resource/id 고정 3개 | action/resource + JSON inputs |
| State | `Target` + `@param entity.Field` | `DiagramID` + `Inputs` + `Transition` |

코드 생성 흐름 (함수별):

1. **request 파라미터 수집**: 모든 시퀀스의 Args에서 `Source == "request"` 추출
2. **Path 파라미터 생성**: OpenAPI PathParams 있으면 `c.Param()` + 타입 변환
3. **Query/Body 파라미터 생성**: request 파라미터 → `c.Query()` 또는 `c.ShouldBindJSON()`
4. **currentUser 추출**: auth 시퀀스 있거나 Args에 `Source == "currentUser"` 있으면
5. **시퀀스별 코드 생성**: 타입별 템플릿 적용
6. **import 수집**: 사용된 패키지 자동 감지

시퀀스별 코드젠:

| 타입 | 생성 코드 |
|---|---|
| get/post | `result, err := model.Method(args...)` + 에러 처리 |
| put/delete | `err := model.Method(args...)` + 에러 처리 |
| empty | `if target == nil { c.JSON(404, ...) }` (타입별 zero check) |
| exists | `if target != nil { c.JSON(409, ...) }` |
| state | `diagramstate.CanTransition(Input{...}, "transition")` |
| auth | `authz.Check(currentUser, "action", "resource", Input{...})` |
| call | `pkg.Func(pkg.FuncRequest{...})` (result 유무 분기) |
| response | `c.JSON(http.StatusOK, gin.H{fields...})` |

Args → 코드 변환:
- `Arg{Source: "request", Field: "CourseID"}` → `courseID` (request에서 추출된 변수)
- `Arg{Source: "course", Field: "InstructorID"}` → `course.InstructorID`
- `Arg{Source: "currentUser", Field: "ID"}` → `currentUser.ID`
- `Arg{Literal: "cancelled"}` → `"cancelled"`

### 3. generator/go_templates.go — 재작성

v2에 맞는 템플릿. v1보다 단순 (시퀀스가 자기 완결적이라 컨텍스트 의존 감소).

### 4. generator/generator.go — 재작성

`Generate()`, `GenerateWith()` 래퍼 + `GenerateModelInterfaces()` 유지.

심볼 테이블 연동:
- 타입 변환 코드젠 (DDL 컬럼 타입 → strconv)
- QueryOpts 자동 전달 (x-extensions)
- List 3-tuple 반환 (many + QueryOpts)
- 모델 인터페이스 파생

## 생성 파일

| 파일 | 내용 |
|---|---|
| `generator/target.go` | Target 인터페이스 (v1 복사) |
| `generator/go_target.go` | Go 코드 생성 (재작성) |
| `generator/go_templates.go` | Go 템플릿 (재작성) |
| `generator/generator.go` | 래퍼 + 유틸 (재작성) |
| `generator/generator_test.go` | 테스트 |

## 테스트 케이스

1. get/post: 모델 호출 + result 바인딩
2. put/delete: result 없는 호출
3. empty/exists: nil/exists guard
4. state: JSON 입력 + 전이 코드
5. auth: JSON 입력 + authz.Check 코드
6. call: result 유무 분기
7. response: 필드 매핑 → gin.H
8. request 파라미터: Query, Body(ShouldBindJSON), Path(c.Param)
9. currentUser 자동 추출
10. 심볼 테이블 연동: 타입 변환, QueryOpts
11. 도메인 폴더 출력

## 의존성

- Phase 001 (parser IR)
- Phase 002 (validator, symbol table)

## 검증

```bash
go test ./generator/... -count=1
```
