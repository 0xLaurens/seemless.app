package services

import (
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"log"
)

type RoomService struct {
	repo ports.RoomRepo
	code ports.CodeService
}

var _ ports.RoomService = (*RoomService)(nil)

func NewRoomService(repo ports.RoomRepo, code ports.CodeService) *RoomService {
	return &RoomService{
		repo: repo,
		code: code,
	}
}

func (r RoomService) CreateRoom() (*data.Room, error) {
	code, err := r.code.CreateRoomCode()
	if err != nil {
		return nil, err
	}

	room := data.CreateRoom(code)
	log.Println("DBG ->", room, room.GetCode())

	_, err = r.repo.AddRoom(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r RoomService) JoinRoom(code data.RoomCode, user *data.User) (*data.Room, error) {
	room, err := r.GetRoomByCode(code)
	if err != nil {
		log.Println("get room by code", err)
		return nil, err
	}

	room.AddClient(user)
	_, err = r.repo.UpdateRoom(room)
	if err != nil {
		log.Println("Update Room", err)
		return nil, err
	}

	user.SetRoom(room.GetId())

	return room, nil
}

func (r RoomService) LeaveRoom(id uuid.UUID, user *data.User) error {
	room, err := r.GetRoomById(id)
	if err != nil {
		return err
	}

	room.RemoveClient(user)
	updatedRoom, err := r.repo.UpdateRoom(room)
	if err != nil {
		return err
	}

	if len(updatedRoom.GetClients()) < 1 {
		err = r.repo.DeleteRoom(room.GetId())
		if err != nil {
			return err
		}

		err = r.code.RemoveRoomCode(room.GetCode())
		if err != nil {
			return err
		}
	}
	user.SetRoom(uuid.Nil)

	return nil
}

func (r RoomService) CreatePublicRoom() (*data.Room, error) {
	code, err := r.code.CreateRoomCode()
	if err != nil {
		return nil, err
	}

	room := data.CreateRoom(code)
	log.Println("DBG ->", room, room.GetCode())

	_, err = r.repo.AddRoom(room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r RoomService) GetRoomById(id uuid.UUID) (*data.Room, error) {
	return r.repo.GetRoomById(id)
}

func (r RoomService) GetRoomByCode(code data.RoomCode) (*data.Room, error) {
	return r.repo.GetRoomByCode(code)
}

func (r RoomService) DeleteRoom(id uuid.UUID) error {
	return r.repo.DeleteRoom(id)
}
