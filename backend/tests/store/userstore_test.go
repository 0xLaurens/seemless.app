package store

import (
	"laurensdrop/data"
	"laurensdrop/store"
	"testing"
)

// setup
func SetupUserStore() store.UserStore {
	return store.NewUserStoreInMemory()
}

func SeedData(s store.UserStore) (int, error) {
	users := []*data.User{
		data.CreateUser("Jane", "Android"),
		data.CreateUser("John", "Iphone"),
		data.CreateUser("Harry", "Iphone"),
		data.CreateUser("Laurens", "Linux"),
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
	user := data.CreateUser("Fritz", "Iphone")
	_, err := s.AddUser(user)
	if err != nil {
		return
	}

	getUser, err := s.GetUser(user.ID)
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

	_, err = s.AddUser(data.CreateUser("Jane", "Windows"))
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
	_, err = s.AddUser(data.CreateUser("Laurens", "Linux"))
	if err.Error() != string(data.UserStoreError.DuplicateUsername) {
		t.Errorf("got %v expected %v", err.Error(), data.UserStoreError.DuplicateUsername)
	}
}

// update user
func TestUpdateUserShouldThrowNotFoundError(t *testing.T) {
	s := SetupUserStore()
	_, err := s.GetUser("1234")

	if err.Error() != string(data.UserStoreError.NotFound) {
		t.Errorf("got %v expected %v", err.Error(), data.UserStoreError.NotFound)
	}
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

	userDTO := data.CreateUser("Jane", "Android")
	_, err = s.UpdateUser(users[0].ID, userDTO)
	if err != nil {
		return
	}

	user, err := s.GetUser(users[0].ID)
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

	userDTO := data.CreateUser("Jane", "Android")
	_, err = s.UpdateUser(preUpdate[0].ID, userDTO)
	if err != nil {
		return
	}

	postUpdate, err := s.GetAllUsers()
	if err != nil {
		return
	}

	if len(preUpdate) != len(postUpdate) {
		t.Errorf("got %v expected %v", len(preUpdate), len(postUpdate))
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

	_, err = s.RemoveUser(preDeleteUsers[0].ID)
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

	_, err = s.RemoveUser("bsID")
	if err.Error() != string(data.UserStoreError.NotFound) {
		t.Errorf("got %s expected %s", err.Error(), string(data.UserStoreError.NotFound))
	}
}
