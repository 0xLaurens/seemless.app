package ports

import (
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
	"net"
)

type RoomRepo interface {
	AddRoom(room *data.Room) (*data.Room, error)
	GetRoomById(id uuid.UUID) (*data.Room, error)
	GetRoomByCode(code data.RoomCode) (*data.Room, error)
	GetRoomByIp(ip *net.IP) (*data.Room, error)
	UpdateRoom(room *data.Room) (*data.Room, error)
	DeleteRoom(id uuid.UUID) error
}

type RoomService interface {
	CreateLocalRoom(ip *net.IP) (*data.Room, error)
	CreatePublicRoom() (*data.Room, error)
	GetRoomById(id uuid.UUID) (*data.Room, error)
	GetRoomByIp(ip *net.IP) (*data.Room, error)
	GetRoomByCode(code data.RoomCode) (*data.Room, error)
	JoinPublicRoom(code data.RoomCode, user *data.User) error
	LeavePublicRoom(code data.RoomCode, user *data.User) error
	JoinLocalRoom(user *data.User) (*data.Room, error)
	LeaveLocalRoom(user *data.User) error
	DeleteRoom(id uuid.UUID) error
}
