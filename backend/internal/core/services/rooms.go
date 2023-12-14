package services

import (
	"fmt"
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"log"
	"net"
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

func (r RoomService) CreateLocalRoom(ip *net.IP) (*data.Room, error) {
	if exists, _ := r.repo.GetRoomByIp(ip); exists != nil {
		return exists, nil
	}

	room, err := r.repo.AddRoom(data.CreateLocalRoom(ip))
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r RoomService) CreatePublicRoom() (*data.Room, error) {
	code, err := r.code.CreateRoomCode()
	if err != nil {
		return nil, err
	}

	room := data.CreatePublicRoom(code)
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

func (r RoomService) GetRoomByIp(ip *net.IP) (*data.Room, error) {
	return r.repo.GetRoomByIp(ip)
}

func (r RoomService) GetRoomByCode(code data.RoomCode) (*data.Room, error) {
	return r.repo.GetRoomByCode(code)
}

func (r RoomService) JoinPublicRoom(code data.RoomCode, user *data.User) error {
	room, err := r.GetRoomByCode(code)
	if err != nil {
		log.Println("get room by code", err)
		return err
	}

	room.AddClient(user)
	_, err = r.repo.UpdateRoom(room)
	if err != nil {
		log.Println("Update Room", err)
		return err
	}

	return nil
}

func (r RoomService) LeavePublicRoom(code data.RoomCode, user *data.User) error {
	room, err := r.GetRoomByCode(code)
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

	return nil
}

func (r RoomService) JoinLocalRoom(user *data.User) (*data.Room, error) {
	ip := user.GetIp()
	if user.GetIp().IsPrivate() {
		fmt.Println("PRIVATE IP")
		parsed := net.ParseIP("0.0.0.0")
		ip = &parsed
	}

	room, err := r.CreateLocalRoom(ip)
	if err != nil {
		return nil, err
	}

	room.AddClient(user)
	_, err = r.repo.UpdateRoom(room)
	if err != nil {
		return nil, err
	}
	log.Printf("DBG -->> user %s joined to room with the ip %v\n", user.GetUsername(), ip)

	return room, nil
}

func (r RoomService) LeaveLocalRoom(user *data.User) error {
	room, err := r.GetRoomByIp(user.GetIp())
	if err != nil {
		return err
	}
	room.RemoveClient(user)
	updateRoom, err := r.repo.UpdateRoom(room)
	if err != nil {
		return err
	}

	if len(updateRoom.GetClients()) < 1 {
		err := r.repo.DeleteRoom(room.GetId())
		if err != nil {
			return err
		}
	}
	return nil
}

func (r RoomService) DeleteRoom(id uuid.UUID) error {
	return r.repo.DeleteRoom(id)
}
