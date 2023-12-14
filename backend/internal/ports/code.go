package ports

import "laurensdrop/internal/core/data"

type CodeRepo interface {
	AddCode(code data.RoomCode) error
	RemoveCode(code data.RoomCode)
	CodeExists(code data.RoomCode) bool
}

type CodeService interface {
	CreateRoomCode() (data.RoomCode, error)
	RemoveRoomCode(code data.RoomCode) error
}
