package repo

import (
	"github.com/google/uuid"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"net"
	"strings"
	"sync"
)

type UserRepoInMemory struct {
	mu        sync.Mutex
	Users     map[string]*data.User       //username -> user
	UsersById map[string]*data.User       //id -> user
	Conns     map[*net.TCPAddr]*data.User //ws conn -> user
}

var _ ports.UserRepo = (*UserRepoInMemory)(nil)

func NewUserRepoInMemory() *UserRepoInMemory {
	return &UserRepoInMemory{
		Users:     make(map[string]*data.User),
		UsersById: make(map[string]*data.User),
		Conns:     make(map[*net.TCPAddr]*data.User),
	}
}

func (s *UserRepoInMemory) AddUser(u *data.User) (*data.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.Users[strings.ToUpper(u.Username)]
	if exists {
		return nil, data.UserDuplicateUsername.Error()
	}

	s.Users[strings.ToUpper(u.Username)] = u
	s.UsersById[u.Id.String()] = u

	if u.Connection != nil {
		//cast is allowed since this interface is always equal to TCPAddr
		remoteAddr := u.Connection.RemoteAddr().(*net.TCPAddr)
		s.Conns[remoteAddr] = u
	}

	return u, nil
}

func (s *UserRepoInMemory) GetUserByName(username string) (*data.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Users[strings.ToUpper(username)]
	if !exists {
		return nil, data.UserNotFound.Error()
	}

	return user, nil
}

func (s *UserRepoInMemory) GetUserByAddr(addr data.RemoteAddr) (*data.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Conns[addr]
	if !exists {
		return nil, data.UserNotFound.Error()
	}

	return user, nil
}

func (s *UserRepoInMemory) GetUserById(id uuid.UUID) (*data.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.UsersById[id.String()]
	if !exists {
		return nil, data.UserNotFound.Error()
	}

	return user, nil
}

func (s *UserRepoInMemory) UpdateUser(id uuid.UUID, userDTO *data.User) (*data.User, error) {
	s.mu.Lock()

	user, exists := s.UsersById[id.String()]
	if !exists {
		return nil, data.UserNotFound.Error()
	}

	delete(s.Users, user.Username)
	delete(s.UsersById, user.Id.String())

	user.Username = userDTO.Username
	user.Device = userDTO.Device

	s.mu.Unlock()
	_, err := s.AddUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserRepoInMemory) RemoveUser(id uuid.UUID) ([]*data.User, error) {
	s.mu.Lock()

	user, exists := s.UsersById[id.String()]
	if !exists {
		return nil, data.UserNotFound.Error()
	}

	if user.Connection != nil {
		remoteAddr := user.Connection.Conn.RemoteAddr().(*net.TCPAddr)
		delete(s.Conns, remoteAddr)
	}

	delete(s.Users, user.Username)
	delete(s.UsersById, user.Id.String())

	s.mu.Unlock()

	return s.GetAllUsers()
}

func (s *UserRepoInMemory) GetAllUsers() ([]*data.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	users := make([]*data.User, len(s.Users))

	i := 0
	for _, user := range s.Users {
		users[i] = user
		i++
	}

	return users, nil
}
