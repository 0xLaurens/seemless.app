package repo

import (
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"log"
	"net"
	"sync"
)

type RoomRepoInMemory struct {
	mu          sync.Mutex
	rooms       map[uuid.UUID]*data.Room
	localRooms  map[string]*data.Room //ip
	publicRooms map[data.RoomCode]*data.Room
}

var _ ports.RoomRepo = (*RoomRepoInMemory)(nil)

func NewRoomRepoInMemory() *RoomRepoInMemory {
	return &RoomRepoInMemory{
		rooms:       make(map[uuid.UUID]*data.Room),
		localRooms:  make(map[string]*data.Room),
		publicRooms: make(map[data.RoomCode]*data.Room),
	}
}

func (r *RoomRepoInMemory) roomExistById(id uuid.UUID) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	room := r.rooms[id]
	return room != nil
}

func (r *RoomRepoInMemory) roomExistByIp(ip *net.IP) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if ip == nil {
		return false
	}
	room := r.localRooms[ip.String()]
	return room != nil
}

func (r *RoomRepoInMemory) roomExistByCode(code data.RoomCode) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	room := r.publicRooms[code]
	return room != nil
}

func (r *RoomRepoInMemory) roomExists(room *data.Room) bool {
	return r.roomExistById(room.GetId()) &&
		r.roomExistByIp(room.GetIP()) || r.roomExistByCode(room.GetCode())
}

func (r *RoomRepoInMemory) AddRoom(room *data.Room) (*data.Room, error) {
	log.Println("ADD ROOM", room)
	if r.roomExists(room) {
		return nil, data.DuplicateRoom.Error()
	}

	r.mu.Lock()
	r.rooms[room.GetId()] = room

	switch room.GetType() {
	case data.PublicRoom:
		log.Println("PUBLIC ROOM", room.GetCode())
		r.publicRooms[room.GetCode()] = room
	case data.LocalRoom:
		r.localRooms[room.GetIP().String()] = room
		log.Println("Add local room", r.localRooms[room.GetIP().String()])
	}

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

func (r *RoomRepoInMemory) GetRoomByIp(ip *net.IP) (*data.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if ip == nil {
		return nil, data.RoomNotFound.Error()
	}

	room := r.localRooms[ip.String()]
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

	switch room.GetType() {
	case data.PublicRoom:
		{
			if oldRoom.GetCode() != room.GetCode() {
				return nil, data.InvalidRoomUpdate.Error()
			}

			r.publicRooms[room.GetCode()] = room
		}
	case data.LocalRoom:
		{
			if oldRoom.GetCode() != room.GetCode() {
				return nil, data.InvalidRoomUpdate.Error()
			}

			r.localRooms[room.GetIP().String()] = room
		}
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.rooms[room.GetId()] = room

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
	if room.GetIP() != nil {
		delete(r.localRooms, room.GetIP().String())
	}
	delete(r.publicRooms, room.GetCode())

	return nil
}
