package services

import (
	"github.com/google/uuid"
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

var _ ports.UserService = (*UserService)(nil)

func (us *UserService) GetUserById(id uuid.UUID) (*data.User, error) {
	return us.repo.GetUserById(id)
}

func (us *UserService) AddUser(User *data.User) (*data.User, error) {
	return us.repo.AddUser(User)
}

func (us *UserService) GetUserByAddr(addr data.RemoteAddr) (*data.User, error) {
	return us.repo.GetUserByAddr(addr)
}

func (us *UserService) GetUserByName(username string) (*data.User, error) {
	return us.repo.GetUserByName(username)
}

func (us *UserService) UpdateUser(id uuid.UUID, UserDTO *data.User) (*data.User, error) {
	return us.repo.UpdateUser(id, UserDTO)
}
func (us *UserService) RemoveUser(id uuid.UUID) ([]*data.User, error) {
	return us.repo.RemoveUser(id)
}

func (us *UserService) GetAllUsers() ([]*data.User, error) {
	return us.repo.GetAllUsers()
}
