# SSaC — Service Sequences as Code

Go 주석 기반 선언적 서비스 로직을 파싱하여 Go 구현 코드를 생성하는 CLI 도구.

```
specs/service/*.go  →  ssac validate  →  ssac gen  →  artifacts/service/*.go
   (주석 DSL)           (정합성 검증)      (코드 생성)     (gofmt 완료)
```

## 핵심 아이디어

서비스 함수 내부의 실행 흐름을 **10종 sequence 타입**으로 선언하고, 구현 코드는 심볼릭 코드젠이 산출한다. LLM 없이 템플릿 기반으로 동작한다.

```go
// @sequence get
// @model Project.FindByID
// @param ProjectID request
// @result project Project

// @sequence guard nil project
// @message "프로젝트가 존재하지 않습니다"

// @sequence post
// @model Session.Create
// @param ProjectID request
// @param Command request
// @result session Session

// @sequence response json
// @var session
func CreateSession(w http.ResponseWriter, r *http.Request) {}
```

이 10줄 선언에서 아래 코드가 자동 생성된다:

```go
func CreateSession(w http.ResponseWriter, r *http.Request) {
    projectID := r.FormValue("ProjectID")
    command := r.FormValue("Command")

    project, err := projectModel.FindByID(projectID)
    if err != nil {
        http.Error(w, "Project 조회 실패", http.StatusInternalServerError)
        return
    }

    if project == nil {
        http.Error(w, "프로젝트가 존재하지 않습니다", http.StatusNotFound)
        return
    }

    session, err := sessionModel.Create(projectID, command)
    if err != nil {
        http.Error(w, "Session 생성 실패", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "session": session,
    })
}
```

## Sequence 타입 (10종)

| 타입 | 역할 |
|---|---|
| `authorize` | 권한 검증 (OPA 등) |
| `get` | 리소스 조회 |
| `guard nil` | null이면 종료 |
| `guard exists` | 존재하면 종료 |
| `post` | 리소스 생성 |
| `put` | 리소스 수정 |
| `delete` | 리소스 삭제 |
| `password` | 비밀번호 비교 |
| `call` | 외부 호출 (@component / @func) |
| `response` | 응답 반환 (json) |

## 설치 & 실행

```bash
go build -o ssac ./artifacts/cmd/ssac

ssac parse [dir]       # 주석 파싱 결과 출력
ssac validate [dir]    # 내부 + 외부 SSOT 교차 검증
ssac gen               # validate → 코드 생성 → gofmt
```

## 검증

내부 검증 (항상):
- 타입별 필수 태그 누락
- `@model` 형식 (`Model.Method`)
- 변수 흐름 (선언 전 참조)

외부 SSOT 교차 검증 (프로젝트 구조 감지 시):
- 모델/메서드 존재 (sqlc queries, Go interface)
- request/response 필드 존재 (OpenAPI)
- component/func 존재 (Go interface)

```bash
ssac validate specs/dummy-study    # 외부 검증 포함
ssac validate specs/backend/service  # 내부 검증만
```

## 프로젝트 구조

```
specs/                           # 선언 (SSOT)
  backend/service/               #   예시 spec
  dummy-study/                   #   스터디룸 예약 더미 프로젝트
    service/  db/queries/  api/  model/
artifacts/                       # 코드
  cmd/ssac/                      #   CLI
  internal/parser/               #   주석 → []ServiceFunc
  internal/generator/            #   타입별 템플릿 → Go 코드
  internal/validator/            #   내부 + 외부 검증
  MANUAL.md                      #   상세 매뉴얼
```

## 외부 검증 프로젝트 레이아웃

```
<project>/
  service/*.go            # sequence spec
  db/queries/*.sql        # sqlc 쿼리 (-- name: Method :type)
  api/openapi.yaml        # OpenAPI 3.0 (operationId = 함수명)
  model/*.go              # Go interface (component), func
```

## 테스트

```bash
go test ./artifacts/internal/... -v
```

46개 테스트: parser 14 + generator 4 + validator 28

## 라이선스

MIT
