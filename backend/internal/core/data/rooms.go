package data

import (
	"github.com/google/uuid"
)

type RoomType string

const (
	LocalRoom  RoomType = "localRoom"
	PublicRoom RoomType = "publicRoom"
)

type RoomCode string

type Room struct {
	id      uuid.UUID
	code    RoomCode
	clients map[*User]bool
}

func CreateRoom(code RoomCode) *Room {
	return &Room{
		id:      uuid.New(),
		code:    code,
		clients: map[*User]bool{},
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
	r.clients[user] = true
}

func (r *Room) RemoveClient(user *User) {
	delete(r.clients, user)
}

func (r *Room) GetClient(user *User) bool {
	return r.clients[user]
}
