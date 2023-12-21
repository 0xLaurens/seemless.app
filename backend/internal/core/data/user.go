package data

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"laurensdrop/internal/core/utils"
	"net"
)

type RemoteAddr *net.TCPAddr
type UserOption func(*User)

type User struct {
	Id         uuid.UUID       `json:"id"`
	Username   string          `json:"username"`
	Device     string          `json:"device"`
	Connection *websocket.Conn `json:"-"`
	RemoteAddr RemoteAddr      `json:"-"`
	PublicRoom uuid.UUID       `json:"-"`
	LocalRoom  uuid.UUID       `json:"-"`
}

// WithConnection helper function for create user to pass a connection
func WithConnection(conn *websocket.Conn) UserOption {
	return func(u *User) {
		u.Connection = conn
		u.RemoteAddr = conn.RemoteAddr().(*net.TCPAddr)
	}
}

func WithUsername(username string) UserOption {
	return func(u *User) {
		u.Username = username
	}
}

func CreateUser(device string, options ...UserOption) *User {
	user := &User{
		Id:         uuid.New(),
		Username:   utils.GenerateRandomDisplayName(),
		Device:     device,
		PublicRoom: uuid.Nil,
		LocalRoom:  uuid.Nil,
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

func (u *User) GetConnection() *websocket.Conn {
	return u.Connection
}

func (u *User) GetIp() *net.IP {
	return &u.RemoteAddr.IP
}
