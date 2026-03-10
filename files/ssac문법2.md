# SSaC 문법 v2

## 설계 원칙

- 한 시퀀스 = 한 줄
- 파라미터 출처가 문법에 내장 (`request.X`, `variable.X`, `"literal"`)
- `@model` + `@param` + `@result` 태그 분리 폐지 → 함수 호출 문법으로 통합

## 문법

### @get — 조회 (result 필수)

```
@get {Type} {var} = {Model}.{Method}({args...})
```

```go
// @get Course course = Course.FindByID(request.CourseID)
// @get User instructor = User.FindByID(course.InstructorID)
// @get []Reservation reservations = Reservation.ListByRoom(request.RoomID)
```

### @post — 생성 (result 필수)

```
@post {Type} {var} = {Model}.{Method}({args...})
```

```go
// @post Session session = Session.Create(request.ProjectID, request.Command)
```

### @put — 수정 (result 없음)

```
@put {Model}.{Method}({args...})
```

```go
// @put Course.Update(request.Title, course.ID)
```

### @delete — 삭제 (result 없음)

```
@delete {Model}.{Method}({args...})
```

```go
// @delete Reservation.Cancel(reservation.ID)
```

### @empty — guard nil (개체 또는 개체.멤버)

```
@empty {target} "{message}"
```

```go
// @empty course "코스를 찾을 수 없습니다"
// @empty course.InstructorID "강사가 지정되지 않았습니다"
```

### @exists — guard exists (개체 또는 개체.멤버)

```
@exists {target} "{message}"
```

```go
// @exists existing "이미 존재합니다"
```

### @state — 상태 전이 검증

```
@state {diagramID} {inputs} "{transition}" "{message}"
```

- `{diagramID}`: 상태 다이어그램 패키지 식별자
- `{inputs}`: JSON 형식의 입력 매핑 (`{key: variable.Field, ...}`)
- `{transition}`: 시도할 전이 액션
- `{message}`: 실패 시 에러 메시지

```go
// 단순: 상태 필드 1개
// @state reservation {status: reservation.Status} "cancel" "취소할 수 없습니다"

// 복합: 여러 필드 입력
// @state course {status: course.Status, createdAt: course.CreatedAt} "publish" "발행할 수 없습니다"
```

코드젠 결과:
```go
if !coursestate.CanTransition(coursestate.Input{
    Status:    course.Status,
    CreatedAt: course.CreatedAt,
}, "publish") {
    c.JSON(http.StatusConflict, gin.H{"error": "발행할 수 없습니다"})
    return
}
```

### @auth — 권한 검사

```
@auth "{action}" "{resource}" {inputs} "{message}"
```

- `{action}`: OPA 액션 (문자열 리터럴)
- `{resource}`: OPA 리소스 (문자열 리터럴)
- `{inputs}`: JSON 형식의 추가 컨텍스트 매핑 (`{key: variable.Field, ...}`)
- `{message}`: 실패 시 에러 메시지

```go
// 기본: 추가 정보 없음
// @auth "delete" "project" {} "권한 없음"

// 소유권 검사
// @auth "delete" "project" {id: project.ID, owner: project.OwnerID} "권한 없음"

// 조직 기반
// @auth "update" "course" {id: course.ID, orgId: course.OrgID} "권한 없음"
```

코드젠 결과:
```go
authz.Check(currentUser, "delete", "project", authz.Input{
    ID:    project.ID,
    Owner: project.OwnerID,
})
```

### @call — 외부 함수 호출

result 있음:
```
@call {Type} {var} = {package}.{Func}({args...})
```

result 없음:
```
@call {package}.{Func}({args...})
```

```go
// @call Token token = auth.VerifyPassword(user.Email, request.Password)
// @call notification.Send(reservation.ID, "cancelled")
```

### @response — 응답 반환 (필드 매핑)

서비스 레이어의 책임: 모델 결과를 OpenAPI response schema에 맞춰 필드별 매핑.
권한별 응답 차이는 서비스 함수 분리로 해결 (조건 분기 금지).

```
@response {
  {필드명}: {변수},
  {필드명}: {변수}.{멤버},
  {필드명}: "{리터럴}"
}
```

허용되는 값 표현:
- 변수 직접 매핑: `course: course`
- 변수 필드 매핑: `instructor_name: instructor.Name`
- 리터럴: `status: "success"`
- 런타임 함수(`len` 등) 금지 — 집계는 SQL에서 처리

```go
// @response {
//   course: course,
//   instructor_name: instructor.Name,
//   reviews: reviews
// }
```

## 전체 예시

```go
import "myapp/auth"

// @get Course course = Course.FindByID(request.CourseID)
// @empty course "코스를 찾을 수 없습니다"
// @get User instructor = User.FindByID(course.InstructorID)
// @empty instructor "강사를 찾을 수 없습니다"
// @response {
//   course: course,
//   instructor_name: instructor.Name
// }
func GetCourse(c *gin.Context) {}

// @get Course course = Course.FindByID(request.CourseID)
// @empty course "코스를 찾을 수 없습니다"
// @get User instructor = User.FindByID(course.InstructorID)
// @empty instructor "강사를 찾을 수 없습니다"
// @get []Review reviews = Review.ListByCourse(course.ID)
// @response {
//   course: course,
//   instructor: instructor,
//   reviews: reviews,
//   email: instructor.Email
// }
func GetCourseAdmin(c *gin.Context) {}

// @post Session session = Session.Create(request.ProjectID, request.Command)
// @response {
//   session: session
// }
func CreateSession(c *gin.Context) {}

// @get Reservation reservation = Reservation.FindByID(request.ReservationID)
// @empty reservation "예약을 찾을 수 없습니다"
// @state reservation {status: reservation.Status} "cancel" "예약 상태 전이 불가"
// @delete Reservation.Cancel(reservation.ID)
// @call notification.Send(reservation.ID, "cancelled")
// @response {
//   reservation: reservation
// }
func CancelReservation(c *gin.Context) {}

// @get User user = User.FindByEmail(request.Email)
// @empty user "사용자를 찾을 수 없습니다"
// @call Token token = auth.VerifyPassword(user.Email, request.Password)
// @response {
//   token: token
// }
func Login(c *gin.Context) {}
```

## v1 → v2 변경 요약

| v1 | v2 | 비고 |
|---|---|---|
| `@sequence get` | `@get` | 타입 키워드가 태그 |
| `@model Model.Method` | 호출 문법에 통합 | `Model.Method(args)` |
| `@param Name source` | 호출 인자에 통합 | `source.Name` |
| `@result var Type` | 대입 문법에 통합 | `Type var = ...` |
| `@message "msg"` | 각 태그 끝에 통합 | `@empty x "msg"` |
| `@guard nil x` | `@empty x` | 이름 변경 |
| `@guard exists x` | `@exists x` | 이름 변경 |
| `@guard state x` + `@param x.Field` | `@state x {inputs} "action"` | 전이 액션+입력 명시 |
| `@action/@resource/@id` | `@auth "action" "resource" {inputs}` | 한 줄 + JSON 입력 |
| `@var x` + `@sequence response` | `@response { field: var }` | 필드 매핑 블록으로 통합 |
| `@func pkg.Fn` + `@param` | `@call pkg.Fn(args)` | 한 줄로 통합 |
