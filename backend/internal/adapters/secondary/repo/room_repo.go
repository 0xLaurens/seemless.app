package repo

import (
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"log"
	"sync"
)

type RoomRepoInMemory struct {
	mu          sync.Mutex
	rooms       map[uuid.UUID]*data.Room
	publicRooms map[data.RoomCode]*data.Room
}

var _ ports.RoomRepo = (*RoomRepoInMemory)(nil)

func NewRoomRepoInMemory() *RoomRepoInMemory {
	return &RoomRepoInMemory{
		rooms:       make(map[uuid.UUID]*data.Room),
		publicRooms: make(map[data.RoomCode]*data.Room),
	}
}

func (r *RoomRepoInMemory) roomExistById(id uuid.UUID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	room := r.rooms[id]
	return room != nil
}

func (r *RoomRepoInMemory) roomExistByCode(code data.RoomCode) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	room := r.publicRooms[code]
	return room != nil
}

func (r *RoomRepoInMemory) roomExists(room *data.Room) bool {
	return r.roomExistById(room.GetId()) && r.roomExistByCode(room.GetCode())
}

func (r *RoomRepoInMemory) AddRoom(room *data.Room) (*data.Room, error) {
	log.Println("ADD ROOM", room)
	if r.roomExists(room) {
		return nil, data.DuplicateRoom.Error()
	}

	r.mu.Lock()
	r.rooms[room.GetId()] = room
	r.publicRooms[room.GetCode()] = room

	r.mu.Unlock()

	return room, nil
}

func (r *RoomRepoInMemory) GetRoomById(id uuid.UUID) (*data.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	room := r.rooms[id]
	if room == nil {
		return nil, data.RoomNotFound.Error()
	}

	return room, nil
}

func (r *RoomRepoInMemory) GetRoomByCode(code data.RoomCode) (*data.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	room := r.publicRooms[code]
	if room == nil {
		return nil, data.RoomNotFound.Error()
	}

	return room, nil
}

func (r *RoomRepoInMemory) UpdateRoom(room *data.Room) (*data.Room, error) {
	oldRoom, _ := r.GetRoomById(room.GetId())
	if !r.roomExists(room) {
		return nil, data.RoomNotFound.Error()
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if oldRoom.GetCode() != room.GetCode() {
		return nil, data.InvalidRoomUpdate.Error()
	}

	r.rooms[room.GetId()] = room
	r.publicRooms[room.GetCode()] = room

	return room, nil
}

func (r *RoomRepoInMemory) DeleteRoom(id uuid.UUID) error {
	room, err := r.GetRoomById(id)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if clients := room.GetClients(); len(clients) > 0 {
		return data.RoomNotEmpty.Error()
	}

	delete(r.rooms, room.GetId())
	delete(r.publicRooms, room.GetCode())

	return nil
}
