package store

import "laurensdrop/data"

type UserStore interface {
	AddUser(User *data.User) (*data.User, error)
	GetUser(ID data.UserID) (*data.User, error)
	UpdateUser(ID data.UserID, UserDTO *data.User) (*data.User, error)
	RemoveUser(ID data.UserID) ([]*data.User, error)
	GetAllUsers() ([]*data.User, error)
}
