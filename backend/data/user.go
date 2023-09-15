package data

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
)

type UserID string
type Conn *websocket.Conn

type User struct {
	ID         UserID `json:"client_id"`
	Username   string `json:"username"`
	Device     string `json:"device"`
	Connection Conn   `json:"connection"`
}

func CreateUser(Username string, Device string, conn Conn) *User {
	return &User{
		ID:         UserID(uuid.NewString()),
		Username:   Username,
		Device:     Device,
		Connection: conn,
	}
}

func (u *User) GetUsername() string {
	return u.Username
}

func (u *User) GetDevice() string {
	return u.Device
}

func (u *User) UpdateUsername(Username string) {
	u.Username = Username
}

func (u *User) GetConnection() Conn {
	return u.Connection
}
