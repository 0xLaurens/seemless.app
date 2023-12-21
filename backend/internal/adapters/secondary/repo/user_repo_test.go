package repo

import (
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"testing"
)

// setup
func SetupUserStore() ports.UserRepo {
	return NewUserRepoInMemory()
}

func SeedData(s ports.UserRepo) (int, error) {
	users := []*data.User{
		data.CreateUser("Android", data.WithUsername("Harry")),
		data.CreateUser("Iphone", data.WithUsername("John")),
		data.CreateUser("Iphone", data.WithUsername("Ingrid")),
		data.CreateUser("Linux", data.WithUsername("Suzi")),
	}

	for i := range users {
		_, err := s.AddUser(users[i])
		if err != nil {
			return -1, err
		}
	}

	return len(users), nil
}

// get all
func TestGetAllUsersZeroUsers(t *testing.T) {
	s := SetupUserStore()
	users, err := s.GetAllUsers()

	if err != nil {
		t.Error(err)
	}

	if len(users) > 0 {
		t.Errorf("Expected zero users returned")
	}
}

func TestGetAllUsersUsersFound(t *testing.T) {
	s := SetupUserStore()
	expected, err := SeedData(s)

	if err != nil {
		t.Errorf("%e", err)
	}
	users, err := s.GetAllUsers()
	if err != nil {
		t.Errorf("%e", err)
	}

	if len(users) != expected {
		t.Errorf("expected a length of %v got a length of %v", expected, len(users))
	}
}

// add user
func TestAddUserShouldAddUserToStore(t *testing.T) {
	s := SetupUserStore()
	user := data.CreateUser("Iphone")
	_, err := s.AddUser(user)
	if err != nil {
		return
	}

	getUser, err := s.GetUserByName(user.Username)
	if err != nil {
		return
	}

	if user != getUser {
		t.Errorf("got \"%+v\\n expected \"%+v\\n", getUser, user)
	}
}

func TestAddUserShouldNotAddDuplicateUsername(t *testing.T) {
	s := SetupUserStore()
	expected, err := SeedData(s)
	if err != nil {
		return
	}

	_, err = s.AddUser(data.CreateUser("Windows", data.WithUsername("Harry")))
	if err != nil {
		return
	}

	users, err := s.GetAllUsers()
	if err != nil {
		return
	}

	if len(users) != expected {
		t.Errorf("expected user not to be added, length was %v expected: %v", len(users), expected)
	}

}

func TestAddUserShouldThrowDuplicateUsernameError(t *testing.T) {
	s := SetupUserStore()
	_, err := SeedData(s)
	if err != nil {
		return
	}
	_, err = s.AddUser(data.CreateUser("Linux", data.WithUsername("Harry")))
	assert.Equal(t, data.UserDuplicateUsername.Error().Error(), err.Error())
}

// update user
func TestUpdateUserShouldThrowNotFoundError(t *testing.T) {
	s := SetupUserStore()
	_, err := s.GetUserByName("1234")
	assert.Equal(t, data.UserNotFound.Error().Error(), err.Error())
}

func TestUpdateUserShouldUpdateUser(t *testing.T) {
	s := SetupUserStore()
	_, err := SeedData(s)
	if err != nil {
		return
	}

	users, err := s.GetAllUsers()
	if err != nil {
		return
	}

	userDTO := data.CreateUser("Android")
	_, err = s.UpdateUser(users[0].Username, userDTO)
	if err != nil {
		return
	}

	user, err := s.GetUserByName(userDTO.Username)
	if err != nil {
		return
	}

	if user.Device != userDTO.Device && user.Username != userDTO.Username {
		t.Errorf("got %s expected %s", user.Device, userDTO.Device)
	}

}

func TestUpdateUserShouldNotAffectUserCount(t *testing.T) {
	s := SetupUserStore()
	_, err := SeedData(s)
	if err != nil {
		return
	}

	preUpdate, err := s.GetAllUsers()
	if err != nil {
		return
	}

	userDTO := data.CreateUser("Android", data.WithUsername("Harry"))
	_, err = s.UpdateUser(preUpdate[0].Username, userDTO)
	if err != nil {
		return
	}

	postUpdate, err := s.GetAllUsers()
	if err != nil {
		return
	}

	if len(preUpdate) != len(postUpdate) {
		t.Errorf("got %v expected %v", len(preUpdate), len(postUpdate))
		t.Errorf("%v : %v", preUpdate, postUpdate)
	}

}

// remove user
func TestRemoveUserShouldNotAffectCountIfNonExistent(t *testing.T) {
	s := SetupUserStore()
	expected, err := SeedData(s)
	if err != nil {
		return
	}

	_, err = s.RemoveUser("123")
	if err != nil {
		return
	}

	users, err := s.GetAllUsers()
	if err != nil {
		return
	}

	if len(users) != expected {
		t.Errorf("got %v expected %v", len(users), expected)
	}
}

func TestRemoveUserShouldDeleteUser(t *testing.T) {
	s := SetupUserStore()
	expected, err := SeedData(s)
	if err != nil {
		return
	}
	preDeleteUsers, err := s.GetAllUsers()

	_, err = s.RemoveUser(preDeleteUsers[0].Username)
	if err != nil {
		return
	}

	users, err := s.GetAllUsers()
	if err != nil {
		return
	}

	if len(users) != expected-1 {
		t.Errorf("got %v expected %v", len(users), expected)
	}
}

// get user
func TestGetUserShouldReturnErrorWhenUserNotFound(t *testing.T) {
	s := SetupUserStore()
	_, err := SeedData(s)
	if err != nil {
		return
	}

	_, err = s.RemoveUser("bsUsername")
	assert.Equal(t, data.UserNotFound.Error().Error(), err.Error())
}
