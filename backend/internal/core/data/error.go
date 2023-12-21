package data

import "errors"

type UserError string

const (
	UserDuplicateUsername UserError = "UserDuplicateUsername"
	UserNotFound          UserError = "UserNotFound"
)

func (ue UserError) Error() error {
	switch ue {
	case UserDuplicateUsername:
		return errors.New("username is not unique")
	case UserNotFound:
		return errors.New("user not found")
	}

	return nil
}

type WebsocketError string

const (
	InvalidRequestBody WebsocketError = "InvalidRequestBody"
)

func (we WebsocketError) Error() error {
	switch we {
	case InvalidRequestBody:
		return errors.New("the message has a invalid request body")
	}

	return nil
}

type RoomError int

const (
	DuplicateRoomCode RoomError = iota
	InvalidRoomCode
	RoomNotEmpty
	RoomNotFound
	DuplicateRoom
	InvalidRoomUpdate
)

func (re RoomError) Error() error {
	switch re {
	case DuplicateRoomCode:
		return errors.New("duplicate room code")
	case InvalidRoomCode:
		return errors.New("invalid room code")
	case RoomNotEmpty:
		return errors.New("room is not empty cannot delete")
	case RoomNotFound:
		return errors.New("room does not exist")
	case DuplicateRoom:
		return errors.New("room already exists")
	case InvalidRoomUpdate:
		return errors.New("you are not allowed to edit this property")
	}
	return nil
}
