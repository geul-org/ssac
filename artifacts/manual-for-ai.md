# SSaC — AI Compact Reference

## CLI

```
ssac parse [dir]      # 주석 파싱 결과 출력 (기본: specs/backend/service/)
ssac validate [dir]   # 내부 검증 또는 외부 SSOT 교차 검증 (자동 감지)
ssac gen              # validate → codegen → gofmt
```

## 기술 스택

Go 1.24+, `go/ast`(파싱), `text/template`(코드젠), `gopkg.in/yaml.v3`(OpenAPI)

## DSL 문법

```go
// @sequence <type>        — 블록 시작. 10종: authorize|get|guard nil|guard exists|post|put|delete|password|call|response
// @model <Model.Method>   — 리소스 모델.메서드 (get/post/put/delete)
// @param <Name> <source>  — source: request, currentUser, 변수명, "리터럴"
// @result <var> <Type>    — 결과 바인딩 (get/post 필수, call 선택)
// @message "msg"          — 커스텀 에러 메시지 (선택, 기본값 자동생성)
// @var <name>             — response에서 반환할 변수
// @action @resource @id   — authorize 전용 (3개 모두 필수)
// @component | @func      — call 전용 (택일 필수)
```

타입별 필수 태그:

| 타입 | 필수 |
|---|---|
| authorize | @action, @resource, @id |
| get, post | @model, @result |
| put, delete | @model |
| guard nil/exists | target (sequence 라인에 변수명) |
| password | @param 2개 (hash, plain) |
| call | @component 또는 @func (택일) |
| response | (없음, @var는 선택) |

## 디렉토리

```
specs/                           # 선언 (입력, SSOT)
  backend/service/               #   기존 예시 spec
  dummy-study/                   #   스터디룸 예약 더미 프로젝트
    service/  db/queries/  api/  model/
artifacts/                       # 산출 (출력, 코드)
  cmd/ssac/main.go               #   CLI 진입점
  internal/parser/               #   Phase1: 주석 → []ServiceFunc
  internal/generator/            #   Phase2: 타입별 템플릿 → Go 코드
  internal/validator/            #   Phase3: 내부 + 외부 SSOT 검증
  backend/internal/service/      #   생성된 Go 코드
  MANUAL.md                      #   상세 매뉴얼
files/                           # 기초 자료
  SSaC.md                        #   기획서
```

## 외부 검증 프로젝트 구조

`ssac validate <project-root>` 시 자동 감지:
- `<root>/service/*.go` — sequence spec
- `<root>/db/queries/*.sql` — sqlc 쿼리 (파일명→모델, `-- name:`→메서드)
- `<root>/api/openapi.yaml` — OpenAPI 3.0 (operationId=함수명)
- `<root>/model/*.go` — Go interface→component, func→@func

## Coding Conventions

- gofmt 준수, 에러 즉시 처리 (early return)
- 파일명: snake_case, 변수/함수: camelCase, 타입: PascalCase
- 테스트: `go test ./artifacts/internal/... -count=1`
