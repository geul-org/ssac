# Phase 001: v2 파서 — 한 줄 표현식 파싱

## 목표

v2 문법의 Go 주석을 파싱하여 `[]ServiceFunc` IR을 생성한다.
v1의 태그 기반 다중 라인 파싱을 한 줄 표현식 파싱으로 전면 교체.

## IR 구조체 설계

### parser/types.go

```go
package parser

// ServiceFunc는 하나의 서비스 함수 선언이다.
type ServiceFunc struct {
    Name      string     // 함수명 (e.g. "GetCourse")
    FileName  string     // 원본 파일명
    Domain    string     // 도메인 폴더명 (e.g. "auth", 없으면 "")
    Sequences []Sequence // 시퀀스 목록
    Imports   []string   // Go import 경로
}

// Sequence는 하나의 시퀀스 라인이다.
type Sequence struct {
    Type    string   // "get", "post", "put", "delete", "empty", "exists", "state", "auth", "call", "response"

    // get/post/put/delete/call 공통: 함수 호출
    Model   string   // "Course.FindByID" 또는 "auth.VerifyPassword"
    Args    []Arg    // 호출 인자

    // get/post/call: 대입
    Result  *Result  // 결과 바인딩 (nil이면 대입 없음)

    // empty/exists: guard
    Target  string   // "course" 또는 "course.InstructorID"

    // state: 상태 전이
    DiagramID  string          // "reservation"
    Inputs     map[string]string // {status: "reservation.Status"}
    Transition string          // "cancel"

    // auth: 권한 검사
    Action   string            // "delete"
    Resource string            // "project"
    // Inputs 재사용              // {id: "project.ID", owner: "project.OwnerID"}

    // response: 필드 매핑
    Fields  map[string]string  // {course: "course", instructor_name: "instructor.Name"}

    // 공통
    Message string   // 에러 메시지
}

// Arg는 함수 호출 인자다.
type Arg struct {
    Source string // "request", 변수명, 또는 "" (리터럴)
    Field  string // "CourseID", "ID" 등
    Literal string // "cancelled" 등 (Source가 ""일 때)
}

// Result는 결과 바인딩이다.
type Result struct {
    Type  string // "Course", "[]Reservation"
    Var   string // "course", "reservations"
}
```

### 시퀀스 타입 상수

```go
const (
    SeqGet      = "get"
    SeqPost     = "post"
    SeqPut      = "put"
    SeqDelete   = "delete"
    SeqEmpty    = "empty"
    SeqExists   = "exists"
    SeqState    = "state"
    SeqAuth     = "auth"
    SeqCall     = "call"
    SeqResponse = "response"
)
```

## 파싱 규칙

각 주석 라인을 독립적으로 파싱 (v1과 달리 상태 머신 불필요):

| 패턴 | 파싱 |
|---|---|
| `@get Type var = Model.Method(args)` | Type+var → Result, Model.Method → Model, args → Args |
| `@post Type var = Model.Method(args)` | 동일 |
| `@put Model.Method(args)` | Model → Model, args → Args, Result = nil |
| `@delete Model.Method(args)` | 동일 |
| `@empty target "msg"` | target → Target, msg → Message |
| `@exists target "msg"` | target → Target, msg → Message |
| `@state id {inputs} "transition" "msg"` | id → DiagramID, inputs → Inputs, transition → Transition |
| `@auth "action" "resource" {inputs} "msg"` | action → Action, resource → Resource, inputs → Inputs |
| `@call Type var = pkg.Func(args)` | Result 있는 call |
| `@call pkg.Func(args)` | Result 없는 call |
| `@response {` | 멀티라인 시작, 이후 `key: value,` 라인 수집, `}` 에서 종료 |

인자 파싱 (`args`):
- `request.CourseID` → `Arg{Source: "request", Field: "CourseID"}`
- `course.InstructorID` → `Arg{Source: "course", Field: "InstructorID"}`
- `currentUser.ID` → `Arg{Source: "currentUser", Field: "ID"}`
- `"cancelled"` → `Arg{Literal: "cancelled"}`

## parser/parser.go 핵심 함수

- `ParseDir(dir string) ([]ServiceFunc, error)` — 디렉토리 재귀 탐색, `.go` 파일 파싱
- `ParseFile(path string) ([]ServiceFunc, error)` — 단일 파일 파싱
- `parseLine(line string) (*Sequence, error)` — 한 줄 → Sequence (@ 접두사 기준 분기)
- `parseCallExpr(expr string) (model string, args []Arg, err error)` — `Model.Method(args)` 파싱
- `parseArgs(argsStr string) ([]Arg, error)` — 쉼표 분리 인자 파싱
- `parseArg(s string) Arg` — 단일 인자 파싱 (`source.Field` 또는 `"literal"`)
- `parseInputs(s string) (map[string]string, error)` — `{key: value, ...}` JSON 형식 파싱
- `parseResponseBlock(lines []string) map[string]string` — `@response { ... }` 멀티라인 파싱

Go AST로 함수 선언과 import를 추출하는 것은 v1과 동일 (`go/parser` 사용).

## 생성 파일

| 파일 | 내용 |
|---|---|
| `parser/types.go` | IR 구조체 (ServiceFunc, Sequence, Arg, Result) |
| `parser/parser.go` | 파싱 로직 |
| `parser/parser_test.go` | 테스트 |

## 테스트 케이스

1. `@get` 단일/다중 파라미터, 변수 참조 인자
2. `@post` result 바인딩
3. `@put`/`@delete` result 없는 호출
4. `@empty`/`@exists` 개체, 개체.멤버
5. `@state` JSON 입력 + 전이 액션
6. `@auth` action/resource/inputs
7. `@call` result 있음/없음
8. `@response` 멀티라인 필드 매핑
9. 도메인 폴더 재귀 탐색
10. import 수집
11. dummy-study 서비스 파일 전체 파싱 (v2로 변환 후)

## 의존성

- 없음 (순수 파서)

## 검증

```bash
go test ./parser/... -count=1
```
