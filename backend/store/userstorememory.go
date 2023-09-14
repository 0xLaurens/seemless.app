package store

import (
	"fmt"
	"laurensdrop/data"
)

type UserStoreInMemory struct {
	Users     map[data.UserID]*data.User
	Usernames map[string]bool
}

func NewUserStoreInMemory() *UserStoreInMemory {
	return &UserStoreInMemory{
		Users:     map[data.UserID]*data.User{},
		Usernames: map[string]bool{},
	}
}

func (s *UserStoreInMemory) AddUser(u *data.User) (*data.User, error) {
	_, exists := s.Usernames[u.Username]
	if exists {
		return nil, fmt.Errorf(string(data.UserStoreError.DuplicateUsername))
	}

	s.Users[u.ID] = u
	s.Usernames[u.Username] = true
	return u, nil
}

func (s *UserStoreInMemory) GetUser(ID data.UserID) (*data.User, error) {
	user, exists := s.Users[ID]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	return user, nil
}

func (s *UserStoreInMemory) UpdateUser(ID data.UserID, UserDTO *data.User) (*data.User, error) {
	user, exists := s.Users[ID]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	delete(s.Usernames, user.Username)

	user.Username = UserDTO.Username
	user.Device = UserDTO.Device

	s.Usernames[user.Username] = true
	s.Users[ID] = user

	return user, nil
}

func (s *UserStoreInMemory) RemoveUser(ID data.UserID) ([]*data.User, error) {
	user, exists := s.Users[ID]
	if !exists {
		return nil, fmt.Errorf("user does not exist")
	}

	delete(s.Users, ID)
	delete(s.Usernames, user.Username)

	return s.GetAllUsers()
}

func (s *UserStoreInMemory) GetAllUsers() ([]*data.User, error) {
	users := make([]*data.User, len(s.Users))

	i := 0
	for _, user := range s.Users {
		users[i] = user
		i++
	}

	return users, nil
}
