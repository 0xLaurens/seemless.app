package store

import "laurensdrop/data"

type UserStore interface {
	AddUser(User *data.User) (*data.User, error)
	GetUserByConn(conn data.Conn) (*data.User, error)
	GetUserByName(username string) (*data.User, error)
	UpdateUser(username string, UserDTO *data.User) (*data.User, error)
	RemoveUser(username string) ([]*data.User, error)
	GetAllUsers() ([]*data.User, error)
}
