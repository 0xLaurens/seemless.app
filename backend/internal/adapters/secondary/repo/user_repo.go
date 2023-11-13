package repo

import (
	"fmt"
	"laurensdrop/internal/core/data"
	"net"
	"strings"
)

type UserRepoInMemory struct {
	Users map[string]*data.User          //username -> user
	Conns map[data.RemoteAddr]*data.User //ws conn -> user
}

func NewUserRepoInMemory() *UserRepoInMemory {
	return &UserRepoInMemory{
		Users: map[string]*data.User{},
		Conns: map[data.RemoteAddr]*data.User{},
	}
}

func (s *UserRepoInMemory) AddUser(u *data.User) (*data.User, error) {
	_, exists := s.Users[strings.ToUpper(u.Username)]
	if exists {
		return nil, fmt.Errorf(string(data.UserStoreError.DuplicateUsername))
	}

	s.Users[strings.ToUpper(u.Username)] = u

	// cast is allowed since this interface is always equal to TCPAddr
	remoteAddr := u.Connection.Conn.RemoteAddr().(*net.TCPAddr)

	s.Conns[remoteAddr] = u
	return u, nil
}

func (s *UserRepoInMemory) GetUserByName(username string) (*data.User, error) {
	user, exists := s.Users[strings.ToUpper(username)]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	return user, nil
}

func (s *UserRepoInMemory) GetUserByAddr(addr data.RemoteAddr) (*data.User, error) {
	user, exists := s.Conns[addr]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	return user, nil
}

func (s *UserRepoInMemory) UpdateUser(username string, userDTO *data.User) (*data.User, error) {
	user, exists := s.Users[strings.ToUpper(username)]
	if !exists {
		return nil, fmt.Errorf(string(data.UserStoreError.NotFound))
	}

	delete(s.Users, user.Username)

	user.Username = userDTO.Username
	user.Device = userDTO.Device

	_, err := s.AddUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserRepoInMemory) RemoveUser(username string) ([]*data.User, error) {
	user, exists := s.Users[strings.ToUpper(username)]
	if !exists {
		return nil, fmt.Errorf(data.UserStoreErrMessage(data.UserStoreError.NotFound))
	}
	remoteAddr := user.Connection.Conn.RemoteAddr().(*net.TCPAddr)

	delete(s.Users, strings.ToUpper(username))
	delete(s.Conns, remoteAddr)

	return s.GetAllUsers()
}

func (s *UserRepoInMemory) GetAllUsers() ([]*data.User, error) {
	users := make([]*data.User, len(s.Users))

	i := 0
	for _, user := range s.Users {
		users[i] = user
		i++
	}

	return users, nil
}
