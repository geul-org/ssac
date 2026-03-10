package service

// @auth "update" "room" {id: request.RoomID} "권한이 없습니다"
// @get Room room = Room.FindByID(request.RoomID)
// @empty room "스터디룸이 존재하지 않습니다"
// @put Room.Update(request.RoomID, request.Name, request.Capacity, request.Location)
// @get Room room = Room.FindByID(request.RoomID)
// @response {
//   room: room
// }
func UpdateRoom() {}
