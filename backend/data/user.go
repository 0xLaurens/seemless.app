package data

import (
	"github.com/gofiber/contrib/websocket"
)

type UserID string
type Conn *websocket.Conn

type UserOption func(*User)

type User struct {
	Username   string `json:"username"`
	Device     string `json:"device"`
	Connection Conn   `json:"connection,omitempty"`
}

// WithConnection helper function for create user to pass a connection
func WithConnection(conn Conn) UserOption {
	return func(u *User) {
		u.Connection = conn
	}
}

func CreateUser(username string, device string, options ...UserOption) *User {
	user := &User{
		Username: username,
		Device:   device,
	}

	for _, option := range options {
		option(user)
	}

	return user
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
