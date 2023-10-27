package services

import (
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
)

type UserService struct {
	repo ports.UserRepo
}

func NewUserService(repo ports.UserRepo) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (us *UserService) AddUser(User *data.User) (*data.User, error) {
	return us.repo.AddUser(User)
}

func (us *UserService) GetUserByConn(conn data.Conn) (*data.User, error) {
	return us.repo.GetUserByConn(conn)
}

func (us *UserService) GetUserByName(username string) (*data.User, error) {
	return us.repo.GetUserByName(username)
}

func (us *UserService) UpdateUser(username string, UserDTO *data.User) (*data.User, error) {
	return us.repo.UpdateUser(username, UserDTO)
}

func (us *UserService) RemoveUser(username string) ([]*data.User, error) {
	return us.repo.RemoveUser(username)
}

func (us *UserService) GetAllUsers() ([]*data.User, error) {
	return us.repo.GetAllUsers()
}
