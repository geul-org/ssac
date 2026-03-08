# 더미 프로젝트: 스터디룸 예약 시스템

> SSaC 외부 검증(Phase 3)을 위한 더미 프로젝트. 인증/권한/비즈니스 규칙을 포함하여 10종 sequence를 자연스럽게 사용한다.

## 도메인 모델

### User

| 컬럼 | 타입 | 설명 |
|---|---|---|
| ID | int64 | PK |
| Email | string | 이메일 (unique) |
| PasswordHash | string | bcrypt 해시 |
| Name | string | 이름 |
| Role | string | "user" / "admin" |
| CreatedAt | time.Time | 가입일 |

### Room

| 컬럼 | 타입 | 설명 |
|---|---|---|
| ID | int64 | PK |
| Name | string | 방 이름 (e.g. "A101") |
| Capacity | int | 수용 인원 |
| Location | string | 위치 (e.g. "3층 동관") |
| CreatedAt | time.Time | 등록일 |

### Reservation

| 컬럼 | 타입 | 설명 |
|---|---|---|
| ID | int64 | PK |
| UserID | int64 | FK → User |
| RoomID | int64 | FK → Room |
| StartAt | time.Time | 시작 시간 |
| EndAt | time.Time | 종료 시간 |
| Status | string | "confirmed" / "cancelled" |
| CreatedAt | time.Time | 생성일 |

## 서비스 함수 (7개)

### 1. Login

이메일로 사용자 조회 → 비밀번호 검증 → 세션 토큰 생성 → 반환

| sequence | 상세 |
|---|---|
| get | User.FindByEmail(Email request) → user User |
| guard nil | user |
| password | user.PasswordHash, Password request |
| post | Session.Create(user.ID) → token Token |
| response json | token |

사용 타입: `get`, `guard nil`, `password`, `post`, `response`

### 2. CreateReservation

권한 확인 → 방 존재 확인 → 시간 충돌 체크 → 예약 생성 → 알림 발송

| sequence | 상세 |
|---|---|
| authorize | action=create, resource=reservation, id=RoomID |
| get | Room.FindByID(RoomID request) → room Room |
| guard nil | room |
| get | Reservation.FindConflict(RoomID request, StartAt request, EndAt request) → conflict Reservation |
| guard exists | conflict → "해당 시간에 이미 예약이 있습니다" |
| post | Reservation.Create(UserID currentUser, RoomID request, StartAt request, EndAt request) → reservation Reservation |
| call @component | notification, reservation, "예약이 확정되었습니다" |
| response json | reservation |

사용 타입: `authorize`, `get`, `guard nil`, `guard exists`, `post`, `call @component`, `response`

### 3. GetReservation

예약 단건 조회

| sequence | 상세 |
|---|---|
| get | Reservation.FindByID(ReservationID request) → reservation Reservation |
| guard nil | reservation |
| response json | reservation |

사용 타입: `get`, `guard nil`, `response`

### 4. ListMyReservations

내 예약 목록 조회

| sequence | 상세 |
|---|---|
| get | Reservation.ListByUserID(UserID currentUser) → reservations []Reservation |
| response json | reservations |

사용 타입: `get`, `response`

### 5. CancelReservation

권한 확인 → 예약 조회 → 환불 계산 → 상태 변경 → 알림 발송

| sequence | 상세 |
|---|---|
| authorize | action=cancel, resource=reservation, id=ReservationID |
| get | Reservation.FindByID(ReservationID request) → reservation Reservation |
| guard nil | reservation |
| call @func | calculateRefund(reservation) → refund Refund |
| put | Reservation.UpdateStatus(ReservationID request, "cancelled") |
| call @component | notification, reservation, "예약이 취소되었습니다" |
| response json | reservation, refund |

사용 타입: `authorize`, `get`, `guard nil`, `call @func`, `put`, `call @component`, `response`

### 6. DeleteRoom

관리자 전용 — 방 삭제 (예약 있으면 불가)

| sequence | 상세 |
|---|---|
| authorize | action=delete, resource=room, id=RoomID |
| get | Room.FindByID(RoomID request) → room Room |
| guard nil | room |
| get | Reservation.CountByRoomID(RoomID request) → reservationCount int |
| guard exists | reservationCount → "예약이 존재하여 삭제할 수 없습니다" |
| delete | Room.Delete(RoomID request) |
| response json | |

사용 타입: `authorize`, `get`, `guard nil`, `guard exists`, `delete`, `response`

### 7. UpdateRoom

관리자 전용 — 방 정보 수정

| sequence | 상세 |
|---|---|
| authorize | action=update, resource=room, id=RoomID |
| get | Room.FindByID(RoomID request) → room Room |
| guard nil | room |
| put | Room.Update(RoomID request, Name request, Capacity request, Location request) |
| response json | room |

사용 타입: `authorize`, `get`, `guard nil`, `put`, `response`

## sequence 타입 커버리지

| 타입 | 사용 함수 |
|---|---|
| authorize | CreateReservation, CancelReservation, DeleteRoom, UpdateRoom |
| get | Login, CreateReservation(x2), GetReservation, ListMyReservations, CancelReservation, DeleteRoom(x2), UpdateRoom |
| guard nil | Login, CreateReservation, GetReservation, CancelReservation, DeleteRoom, UpdateRoom |
| guard exists | CreateReservation, DeleteRoom |
| post | Login, CreateReservation |
| put | CancelReservation, UpdateRoom |
| delete | DeleteRoom |
| password | Login |
| call @component | CreateReservation, CancelReservation |
| call @func | CancelReservation |
| response | 전체 7개 |

## 외부 SSOT 더미 구성

### DDL (specs/dummy-study/db/)

- `users.sql` — users 테이블
- `rooms.sql` — rooms 테이블
- `reservations.sql` — reservations 테이블
- `queries/` — sqlc 쿼리 정의

### OpenAPI (specs/dummy-study/api/)

- `openapi.yaml` — 7개 엔드포인트 정의

### Go interface (specs/dummy-study/model/)

- `notification.go` — notification component interface
- `refund.go` — calculateRefund 함수 시그니처

### Service spec (specs/dummy-study/service/)

- 7개 서비스 함수의 sequence 주석 파일
