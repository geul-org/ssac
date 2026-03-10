//go:build ignore

package service

import _ "github.com/geul-org/fullend/pkg/billing"

// @auth "cancel" "reservation" {id: request.ReservationID} "권한이 없습니다"
// @get Reservation reservation = Reservation.FindByID(request.ReservationID)
// @empty reservation "예약을 찾을 수 없습니다"
// @state reservation {status: reservation.Status} "cancel" "취소할 수 없는 상태입니다"
// @call Refund refund = billing.CalculateRefund({ID: reservation.ID, StartAt: reservation.StartAt, EndAt: reservation.EndAt})
// @put Reservation.UpdateStatus(request.ReservationID, "cancelled")
// @get Reservation reservation = Reservation.FindByID(request.ReservationID)
// @response {
//   reservation: reservation,
//   refund: refund
// }
func CancelReservation() {}
