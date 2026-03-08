# OpenAPI x- 확장 — SSaC/STML 인프라 파라미터

> OpenAPI 엔드포인트에 페이지네이션, 정렬, 필터, 관계 포함 기능을 선언하는 4개의 x- 확장.

## 개요

목록 조회 엔드포인트의 인프라 기능(페이지네이션, 정렬, 필터, 관계 포함)은 비즈니스 로직이 아니라 공통 패턴이다. OpenAPI `x-` 확장으로 선언하면 SSaC, STML, DDL 세 방향의 교차 검증이 가능하다.

```
STML (프론트)  →  OpenAPI x-  ←  SSaC (백엔드)
                     ↓
                DDL (인덱스)
```

## x-pagination

목록 엔드포인트의 페이지네이션 방식을 선언한다.

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `style` | string | O | `offset` 또는 `cursor` |
| `defaultLimit` | integer | O | 기본 반환 건수 |
| `maxLimit` | integer | O | 최대 반환 건수 |

### 예시

```yaml
/api/reservations:
  get:
    operationId: ListReservations
    x-pagination:
      style: offset
      defaultLimit: 20
      maxLimit: 100
```

### 코드젠 결과 (SSaC → Go)

```go
// offset style
limit := clampLimit(r.URL.Query().Get("limit"), 20, 100)
offset := parseOffset(r.URL.Query().Get("offset"))
reservations, total, err := reservationModel.ListByUserID(userID, QueryOpts{
    Limit:  limit,
    Offset: offset,
})
```

```go
// cursor style
limit := clampLimit(r.URL.Query().Get("limit"), 20, 100)
cursor := r.URL.Query().Get("cursor")
reservations, nextCursor, err := reservationModel.ListByUserID(userID, QueryOpts{
    Limit:  limit,
    Cursor: cursor,
})
```

### STML 사용

```html
<section data-fetch="ListReservations" data-param-limit="20" data-param-offset="0">
  <ul data-each="reservations">
    <li data-bind="title"></li>
  </ul>
</section>
```

### 검증 규칙

- `x-pagination`이 있는 엔드포인트의 response schema에 배열 필드가 존재해야 한다
- `style: offset`이면 response에 `total` 필드 권장
- `style: cursor`이면 response에 `nextCursor` 필드 권장

---

## x-sort

목록 엔드포인트의 정렬 가능 컬럼을 선언한다.

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `allowed` | string[] | O | 정렬 가능한 컬럼 목록 |
| `default` | string | X | 기본 정렬 컬럼 (없으면 allowed[0]) |
| `direction` | string | X | 기본 정렬 방향. `asc` 또는 `desc` (기본값: `asc`) |

### 예시

```yaml
/api/reservations:
  get:
    operationId: ListReservations
    x-sort:
      allowed: [start_at, created_at, title]
      default: start_at
      direction: desc
```

### 코드젠 결과 (SSaC → Go)

```go
sortCol := validateSort(r.URL.Query().Get("sort"), []string{"start_at", "created_at", "title"}, "start_at")
sortDir := validateDirection(r.URL.Query().Get("direction"), "desc")
reservations, total, err := reservationModel.ListByUserID(userID, QueryOpts{
    SortCol: sortCol,
    SortDir: sortDir,
})
```

### STML 사용

```html
<section data-fetch="ListReservations" data-param-sort="start_at" data-param-direction="desc">
  <ul data-each="reservations">
    <li data-bind="title"></li>
  </ul>
</section>
```

### 검증 규칙

- `x-sort.allowed`의 모든 컬럼이 DDL 테이블에 존재해야 한다
- 정렬 대상 컬럼에 인덱스가 없으면 경고
- STML의 `data-param-sort` 값이 `x-sort.allowed`에 포함되어야 한다

---

## x-filter

목록 엔드포인트의 필터 가능 컬럼을 선언한다. 검색은 필터의 특수 형태로 포함된다.

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `allowed` | string[] | O | 필터 가능한 컬럼 목록 |

### 예시

```yaml
/api/reservations:
  get:
    operationId: ListReservations
    x-filter:
      allowed: [status, room_id]
```

### 복합 예시 (검색 포함)

```yaml
/api/rooms:
  get:
    operationId: ListRooms
    x-filter:
      allowed: [name, location, capacity]
```

검색은 `?name=회의`처럼 필터 파라미터로 처리한다. 모델에서 `LIKE` 또는 `tsvector`로 구현하는 것은 DDL/sqlc의 책임이다.

### 코드젠 결과 (SSaC → Go)

```go
filters := parseFilters(r.URL.Query(), []string{"status", "room_id"})
reservations, total, err := reservationModel.ListByUserID(userID, QueryOpts{
    Filters: filters,
})
```

### STML 사용

```html
<section data-fetch="ListReservations" data-param-status="confirmed" data-param-room-id="room-1">
  <ul data-each="reservations">
    <li data-bind="title"></li>
  </ul>
</section>
```

### 검증 규칙

- `x-filter.allowed`의 모든 컬럼이 DDL 테이블에 존재해야 한다
- STML의 `data-param-{name}` 중 필터 대상이 `x-filter.allowed`에 포함되어야 한다

---

## x-include

목록/상세 엔드포인트에서 관계 리소스를 함께 반환할 수 있음을 선언한다.

| 필드 | 타입 | 필수 | 설명 |
|---|---|---|---|
| `allowed` | string[] | O | 포함 가능한 관계 리소스 목록 |

### 예시

```yaml
/api/reservations:
  get:
    operationId: ListReservations
    x-include:
      allowed: [room, user]

/api/projects/{projectId}:
  get:
    operationId: GetProject
    x-include:
      allowed: [sessions, owner]
```

### 코드젠 결과 (SSaC → Go)

```go
includes := parseIncludes(r.URL.Query().Get("include"), []string{"room", "user"})
reservations, total, err := reservationModel.ListByUserID(userID, QueryOpts{
    Includes: includes,
})
```

### STML 사용

```html
<section data-fetch="ListReservations" data-param-include="room,user">
  <ul data-each="reservations">
    <li>
      <span data-bind="title"></span>
      <span data-bind="room.name"></span>
      <span data-bind="user.name"></span>
    </li>
  </ul>
</section>
```

### 검증 규칙

- `x-include.allowed`의 리소스가 DDL에서 FK 관계로 연결되어 있어야 한다
- STML의 `data-param-include` 값이 `x-include.allowed`에 포함되어야 한다
- `data-bind`에서 dot notation(`room.name`)으로 참조하는 리소스가 include에 포함되어야 한다

---

## 복합 예시

```yaml
/api/reservations:
  get:
    operationId: ListReservations
    summary: 예약 목록 조회
    x-pagination:
      style: offset
      defaultLimit: 20
      maxLimit: 100
    x-sort:
      allowed: [start_at, created_at]
      default: start_at
      direction: desc
    x-filter:
      allowed: [status, room_id]
    x-include:
      allowed: [room, user]
    parameters:
      - name: userId
        in: query
        required: true
        schema:
          type: string
    responses:
      "200":
        description: 예약 목록
        content:
          application/json:
            schema:
              type: object
              properties:
                reservations:
                  type: array
                  items:
                    $ref: '#/components/schemas/Reservation'
                total:
                  type: integer
```

### 대응하는 SSaC

```go
// @sequence get
// @model Reservation.ListByUserID
// @param UserID currentUser
// @result reservations []Reservation

// @sequence response json
// @var reservations
func ListReservations(w http.ResponseWriter, r *http.Request) {}
```

SSaC에 페이지네이션, 정렬, 필터, 포함 파라미터를 선언하지 않는다. 이 파라미터들은 OpenAPI `x-`에만 선언되며, 코드젠이 `x-`를 읽어 자동으로 QueryOpts를 구성한다. SSaC는 비즈니스 파라미터(`UserID`)만 선언한다.

### 대응하는 STML

```html
<section data-fetch="ListReservations"
         data-param-sort="start_at"
         data-param-direction="desc"
         data-param-status="confirmed"
         data-param-include="room">
  <ul data-each="reservations">
    <li>
      <span data-bind="title"></span>
      <span data-bind="startAt"></span>
      <span data-bind="room.name"></span>
    </li>
  </ul>
</section>
```

---

## 교차 검증 요약

| 검증 방향 | 검증 내용 |
|---|---|
| STML → OpenAPI x- | sort/filter/include 파라미터가 allowed에 포함되는가 |
| OpenAPI x- → DDL | allowed 컬럼이 테이블에 존재하는가, 인덱스가 있는가 |
| SSaC → OpenAPI x- | x-가 있는 엔드포인트의 모델 호출에 QueryOpts가 포함되는가 |
| OpenAPI x- → response | pagination이면 배열 + total/cursor 필드가 있는가 |
