package repo

import (
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/core/data"
	"net"
	"testing"
)

func SetupTestRoomRepo() *RoomRepoInMemory {
	return NewRoomRepoInMemory()
}

func TestAddRoomShouldAddPublicRoomToCorrectMaps(t *testing.T) {
	repo := SetupTestRoomRepo()
	room := data.CreatePublicRoom("TESTS")

	addRoom, err := repo.AddRoom(room)
	assert.NoError(t, err)
	assert.Equal(t, room, addRoom)

	roomGet, err := repo.GetRoomById(room.GetId())
	assert.NoError(t, err)
	assert.Equal(t, room, roomGet)

	publicRoom, err := repo.GetRoomByCode(room.GetCode())
	assert.NoError(t, err)
	assert.Equal(t, room, publicRoom)

	privateRoom, err := repo.GetRoomByIp(room.GetIP())
	assert.Error(t, data.RoomNotFound.Error(), err)
	assert.Equal(t, (*data.Room)(nil), privateRoom)
}

func TestAddRoomShouldAddLocalRoomToCorrectMaps(t *testing.T) {
	repo := SetupTestRoomRepo()
	ip := net.ParseIP("127.0.0.1")
	room := data.CreateLocalRoom(&ip)

	addRoom, err := repo.AddRoom(room)
	assert.NoError(t, err)
	assert.Equal(t, room, addRoom)

	roomGet, err := repo.GetRoomById(room.GetId())
	assert.NoError(t, err)
	assert.Equal(t, room, roomGet)

	publicRoom, err := repo.GetRoomByCode(room.GetCode())
	assert.Error(t, data.RoomNotFound.Error(), err)
	assert.Equal(t, (*data.Room)(nil), publicRoom)

	privateRoom, err := repo.GetRoomByIp(room.GetIP())
	assert.NoError(t, err)
	assert.Equal(t, room, privateRoom)
}

func TestUpdateRoomShouldNotBeAbleToAlterTheIpAddress(t *testing.T) {
	repo := SetupTestRoomRepo()
	ip := net.ParseIP("0.0.0.0")
	room := data.CreateLocalRoom(&ip)

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)

	alteredIP := net.ParseIP("192.168.1.1")
	room.SetIP(&alteredIP)

	updateRoom, err := repo.UpdateRoom(room)
	assert.Error(t, data.InvalidRoomUpdate.Error(), err)
	assert.Equal(t, (*data.Room)(nil), updateRoom)
}

func TestUpdateRoomShouldNotBeAbleToAlterTheRoomCode(t *testing.T) {
	repo := SetupTestRoomRepo()
	code := "TESTS"
	room := data.CreatePublicRoom(data.RoomCode(code))

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)

	room.SetCode("CHESS")

	updateRoom, err := repo.UpdateRoom(room)
	assert.Error(t, data.InvalidRoomUpdate.Error(), err)
	assert.Equal(t, (*data.Room)(nil), updateRoom)
}

func TestUpdateRoomShouldChangeTheNumberOfUsersPublicRoom(t *testing.T) {
	repo := SetupTestRoomRepo()
	code := "TESTS"
	room := data.CreatePublicRoom(data.RoomCode(code))

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)
	dummy := data.CreateUser("Android")
	room.AddClient(dummy)

	updatedRoom, err := repo.UpdateRoom(room)
	assert.NoError(t, err)
	assert.Equal(t, room, updatedRoom)

	roomById, err := repo.GetRoomById(room.GetId())
	assert.NoError(t, err)
	assert.Equal(t, room, roomById)

	roomByCode, err := repo.GetRoomByCode(room.GetCode())
	assert.NoError(t, err)
	assert.Equal(t, room, roomByCode)
}

func TestUpdateRoomShouldChangeTheNumberOfUsersLocalRoom(t *testing.T) {
	repo := SetupTestRoomRepo()
	ip := net.ParseIP("0.0.0.0")
	room := data.CreateLocalRoom(&ip)

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)
	dummy := data.CreateUser("Android")
	room.AddClient(dummy)

	updatedRoom, err := repo.UpdateRoom(room)
	assert.NoError(t, err)
	assert.Equal(t, room, updatedRoom)

	roomById, err := repo.GetRoomById(room.GetId())
	assert.NoError(t, err)
	assert.Equal(t, room, roomById)

	roomByCode, err := repo.GetRoomByIp(room.GetIP())
	assert.NoError(t, err)
	assert.Equal(t, room, roomByCode)
}

func TestDeleteRoomShouldNotDeleteWhenRoomIsNotEmpty(t *testing.T) {
	repo := SetupTestRoomRepo()
	ip := net.ParseIP("0.0.0.0")
	room := data.CreateLocalRoom(&ip)

	dummy := data.CreateUser("Android")
	room.AddClient(dummy)

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)

	err = repo.DeleteRoom(room.GetId())
	assert.Error(t, data.RoomNotEmpty.Error(), err)

	roomById, err := repo.GetRoomById(room.GetId())
	assert.NoError(t, err)
	assert.Equal(t, room, roomById)
}

func TestDeleteRoomShouldRemoveLocalRoom(t *testing.T) {
	repo := SetupTestRoomRepo()
	ip := net.ParseIP("0.0.0.0")
	room := data.CreateLocalRoom(&ip)

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)

	err = repo.DeleteRoom(room.GetId())
	assert.NoError(t, err)

	_, err = repo.GetRoomById(room.GetId())
	assert.Error(t, data.RoomNotFound.Error(), err)

	_, err = repo.GetRoomByIp(room.GetIP())
	assert.Error(t, data.RoomNotFound.Error(), err)
}

func TestDeleteRoomShouldRemovePublicRoom(t *testing.T) {
	repo := SetupTestRoomRepo()
	code := "TESTS"
	room := data.CreatePublicRoom(data.RoomCode(code))

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)

	err = repo.DeleteRoom(room.GetId())
	assert.NoError(t, err)

	_, err = repo.GetRoomById(room.GetId())
	assert.Error(t, data.RoomNotFound.Error(), err)

	_, err = repo.GetRoomByCode(room.GetCode())
	assert.Error(t, data.RoomNotFound.Error(), err)
}
