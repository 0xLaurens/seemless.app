package ports

import (
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
)

type UserRepo interface {
	AddUser(User *data.User) (*data.User, error)
	GetUserByAddr(addr data.RemoteAddr) (*data.User, error)
	GetUserByName(username string) (*data.User, error)
	GetUserById(id uuid.UUID) (*data.User, error)
	UpdateUser(id uuid.UUID, UserDTO *data.User) (*data.User, error)
	RemoveUser(id uuid.UUID) ([]*data.User, error)
	GetAllUsers() ([]*data.User, error)
}

type UserService interface {
	AddUser(User *data.User) (*data.User, error)
	GetUserByAddr(addr data.RemoteAddr) (*data.User, error)
	GetUserById(id uuid.UUID) (*data.User, error)
	GetUserByName(username string) (*data.User, error)
	UpdateUser(id uuid.UUID, UserDTO *data.User) (*data.User, error)
	RemoveUser(id uuid.UUID) ([]*data.User, error)
	GetAllUsers() ([]*data.User, error)
}
