package store

import (
	"fmt"
	"laurensdrop/data"
)

type UserStoreInMemory struct {
	Users map[string]*data.User    //username -> user
	Conns map[data.Conn]*data.User //ws conn -> user
}

func NewUserStoreInMemory() *UserStoreInMemory {
	return &UserStoreInMemory{
		Users: map[string]*data.User{},
		Conns: map[data.Conn]*data.User{},
	}
}

func (s *UserStoreInMemory) AddUser(u *data.User) (*data.User, error) {
	_, exists := s.Users[u.Username]
	if exists {
		return nil, fmt.Errorf(string(data.UserStoreError.DuplicateUsername))
	}

	s.Users[u.Username] = u
	s.Conns[u.Connection] = u
	return u, nil
}

func (s *UserStoreInMemory) GetUserByName(username string) (*data.User, error) {
	user, exists := s.Users[username]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	return user, nil
}

func (s *UserStoreInMemory) GetUserByConn(conn data.Conn) (*data.User, error) {
	user, exists := s.Conns[conn]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	return user, nil
}

func (s *UserStoreInMemory) UpdateUser(username string, userDTO *data.User) (*data.User, error) {
	user, exists := s.Users[username]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	delete(s.Users, user.Username)

	user.Username = userDTO.Username
	user.Device = userDTO.Device

	s.Users[userDTO.Username] = userDTO

	return user, nil
}

func (s *UserStoreInMemory) RemoveUser(username string) ([]*data.User, error) {
	user, exists := s.Users[username]
	if !exists {
		return nil, fmt.Errorf("user does not exist")
	}

	delete(s.Users, username)
	delete(s.Conns, user.Connection)

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
