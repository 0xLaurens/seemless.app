package ports

import "laurensdrop/internal/core/data"

type UserRepo interface {
	AddUser(User *data.User) (*data.User, error)
	GetUserByConn(conn data.Conn) (*data.User, error)
	GetUserByName(username string) (*data.User, error)
	UpdateUser(username string, UserDTO *data.User) (*data.User, error)
	RemoveUser(username string) ([]*data.User, error)
	GetAllUsers() ([]*data.User, error)
}

type UserService interface {
	AddUser(User *data.User) (*data.User, error)
	GetUserByConn(conn data.Conn) (*data.User, error)
	GetUserByName(username string) (*data.User, error)
	UpdateUser(username string, UserDTO *data.User) (*data.User, error)
	RemoveUser(username string) ([]*data.User, error)
	GetAllUsers() ([]*data.User, error)
}
