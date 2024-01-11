package repo

import (
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/core/data"
	"testing"
)

func SetupTestRoomRepo() *RoomRepoInMemory {
	return NewRoomRepoInMemory()
}

func TestAddRoomShouldAddPublicRoomToCorrectMaps(t *testing.T) {
	repo := SetupTestRoomRepo()
	room := data.CreateRoom("TESTS")

	addRoom, err := repo.AddRoom(room)
	assert.NoError(t, err)
	assert.Equal(t, room, addRoom)

	roomGet, err := repo.GetRoomById(room.GetId())
	assert.NoError(t, err)
	assert.Equal(t, room, roomGet)

	publicRoom, err := repo.GetRoomByCode(room.GetCode())
	assert.NoError(t, err)
	assert.Equal(t, room, publicRoom)
}

func TestUpdateRoomShouldNotBeAbleToAlterTheRoomCode(t *testing.T) {
	repo := SetupTestRoomRepo()
	code := "TESTS"
	room := data.CreateRoom(data.RoomCode(code))

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
	room := data.CreateRoom(data.RoomCode(code))

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

func TestDeleteRoomShouldRemovePublicRoom(t *testing.T) {
	repo := SetupTestRoomRepo()
	code := "TESTS"
	room := data.CreateRoom(data.RoomCode(code))

	_, err := repo.AddRoom(room)
	assert.NoError(t, err)

	err = repo.DeleteRoom(room.GetId())
	assert.NoError(t, err)

	_, err = repo.GetRoomById(room.GetId())
	assert.Error(t, data.RoomNotFound.Error(), err)

	_, err = repo.GetRoomByCode(room.GetCode())
	assert.Error(t, data.RoomNotFound.Error(), err)
}
