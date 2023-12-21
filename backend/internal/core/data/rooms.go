package data

import (
	"github.com/google/uuid"
	"net"
)

type RoomType string

const (
	LocalRoom  RoomType = "localRoom"
	PublicRoom RoomType = "publicRoom"
)

type RoomCode string

type Room struct {
	id       uuid.UUID
	roomType RoomType
	code     RoomCode
	ip       *net.IP
	clients  map[*User]bool
}

func CreateLocalRoom(ip *net.IP) *Room {
	return &Room{
		id:       uuid.New(),
		roomType: LocalRoom,
		ip:       ip,
		code:     "",
		clients:  map[*User]bool{},
	}
}

func CreatePublicRoom(code RoomCode) *Room {
	return &Room{
		id:       uuid.New(),
		roomType: PublicRoom,
		code:     code,
		ip:       nil,
		clients:  map[*User]bool{},
	}
}

func (r *Room) GetId() uuid.UUID {
	return r.id
}

func (r *Room) GetCode() RoomCode {
	return r.code
}

func (r *Room) SetCode(code RoomCode) {
	r.code = code
}

func (r *Room) GetIP() *net.IP {
	return r.ip
}

func (r *Room) SetIP(ip *net.IP) {
	r.ip = ip
}

func (r *Room) GetClients() []*User {
	users := make([]*User, len(r.clients))

	i := 0
	for user := range r.clients {
		users[i] = user
		i++
	}

	return users
}

func (r *Room) AddClient(user *User) {
	switch r.GetType() {
	case PublicRoom:
		user.PublicRoom = r.GetId()
	case LocalRoom:
		user.LocalRoom = r.GetId()
	}
	r.clients[user] = true
}

func (r *Room) GetType() RoomType {
	return r.roomType
}

func (r *Room) RemoveClient(user *User) {
	delete(r.clients, user)
}

func (r *Room) GetClient(user *User) bool {
	return r.clients[user]
}
