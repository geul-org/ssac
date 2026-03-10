package service

// @auth "create" "reservation" {id: request.RoomID} "권한이 없습니다"
// @get Room room = Room.FindByID(request.RoomID)
// @empty room "스터디룸이 존재하지 않습니다"
// @get Reservation conflict = Reservation.FindConflict(request.RoomID, request.StartAt, request.EndAt)
// @exists conflict "해당 시간에 이미 예약이 있습니다"
// @post Reservation reservation = Reservation.Create(currentUser.ID, request.RoomID, request.StartAt, request.EndAt)
// @response {
//   reservation: reservation
// }
func CreateReservation() {}
