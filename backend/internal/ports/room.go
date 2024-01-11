package ports

import (
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
)

type RoomRepo interface {
	AddRoom(room *data.Room) (*data.Room, error)
	GetRoomById(id uuid.UUID) (*data.Room, error)
	GetRoomByCode(code data.RoomCode) (*data.Room, error)
	UpdateRoom(room *data.Room) (*data.Room, error)
	DeleteRoom(id uuid.UUID) error
}

type RoomService interface {
	GetRoomById(id uuid.UUID) (*data.Room, error)
	GetRoomByCode(code data.RoomCode) (*data.Room, error)

	JoinRoom(code data.RoomCode, user *data.User) (*data.Room, error)
	LeaveRoom(id uuid.UUID, user *data.User) error

	CreateRoom() (*data.Room, error)
	DeleteRoom(id uuid.UUID) error
}
