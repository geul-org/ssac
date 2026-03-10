✅ 완료

# Phase 004: CLI + dummy-study 스펙 v2 전환

## 목표

1. CLI(`cmd/ssac/main.go`)를 v2 파서/밸리데이터/제너레이터에 연결
2. `specs/dummy-study/service/*.go`를 v2 문법으로 전환
3. 전체 파이프라인 통합 테스트: `ssac validate` → `ssac gen`

## 변경 사항

### 1. cmd/ssac/main.go — v1에서 복사 후 수정

기본 구조 유지 (parse, validate, gen 커맨드). import 경로만 변경.
`runParse()`의 출력 형식을 v2 IR에 맞게 조정.

### 2. specs/dummy-study/service/*.go — v2 문법으로 전환

7개 파일을 v2 문법으로 재작성:

| 파일 | v1 시퀀스 수 | v2 라인 수 (예상) |
|---|---|---|
| `login.go` | 5 sequences (14줄) | 5줄 |
| `create_reservation.go` | 8 sequences (~24줄) | 7줄 |
| `get_reservation.go` | 3 sequences (9줄) | 3줄 |
| `list_my_reservations.go` | 2 sequences (6줄) | 2줄 |
| `cancel_reservation.go` | 8 sequences (~25줄) | 8줄 |
| `update_room.go` | 5 sequences (16줄) | 5줄 |
| `delete_room.go` | 6 sequences (18줄) | 5줄 |

v1 → v2 전환 예시 (cancel_reservation.go):

```go
// v1 (25줄)
// @sequence authorize
// @action cancel
// @resource reservation
// @id ReservationID
//
// @sequence get
// @model Reservation.FindByID
// @param ReservationID request
// @result reservation Reservation
//
// @sequence guard nil reservation
// @message "예약을 찾을 수 없습니다"
//
// @sequence guard state reservation
// @param reservation.Status
// ...

// v2 (8줄)
// @auth "cancel" "reservation" {id: request.ReservationID} "권한 없음"
// @get Reservation reservation = Reservation.FindByID(request.ReservationID)
// @empty reservation "예약을 찾을 수 없습니다"
// @state reservation {status: reservation.Status} "cancel" "취소할 수 없습니다"
// @call Refund refund = billing.CalculateRefund(reservation.ID, reservation.StartAt, reservation.EndAt)
// @put Reservation.UpdateStatus(request.ReservationID, "cancelled")
// @get Reservation reservation = Reservation.FindByID(request.ReservationID)
// @response {
//   reservation: reservation,
//   refund: refund
// }
```

### 3. testdata/ — v2 테스트 fixture 작성

통합 테스트용 fixture 파일 작성.

## 생성/수정 파일

| 파일 | 내용 |
|---|---|
| `cmd/ssac/main.go` | CLI 진입점 (v2 연결) |
| `specs/dummy-study/service/*.go` (7개) | v2 문법으로 전환 |
| `testdata/` | v2 테스트 fixture |

## 통합 테스트

```bash
# 파싱
ssac parse specs/dummy-study

# 검증 (외부 SSOT 교차 검증 포함)
ssac validate specs/dummy-study

# 코드 생성
ssac gen specs/dummy-study /tmp/dummy-study-out
```

## 의존성

- Phase 001 (parser)
- Phase 002 (validator)
- Phase 003 (generator)

## 검증

```bash
go test ./... -count=1
ssac validate specs/dummy-study
```
