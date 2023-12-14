package repo

import (
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/core/data"
	"testing"
)

func SetupTestCodeRepo() *CodeRepoInMemory {
	return NewCodeRepoInMemory()
}

func TestAddUniqueCode(t *testing.T) {
	repo := SetupTestCodeRepo()
	err := repo.AddCode("AEFJG")
	assert.NoError(t, err)
}

func TestDuplicateCodeShouldReturnError(t *testing.T) {
	repo := SetupTestCodeRepo()
	err := repo.AddCode("EAFJG")
	assert.NoError(t, err)
	err = repo.AddCode("EAFJG")
	assert.Equal(t, data.DuplicateRoomCode.Error(), err)
}

func TestDeleteNonExistingCode(t *testing.T) {
	repo := SetupTestCodeRepo()
	repo.RemoveCode("EAFJG")
}

func TestDeleteShouldRemoveCode(t *testing.T) {
	repo := SetupTestCodeRepo()
	_ = repo.AddCode("EAFJG")
	assert.Equal(t, 1, len(repo.codes))
	repo.RemoveCode("EAFJG")
	assert.Equal(t, 0, len(repo.codes))
}

func TestCodeShouldExistAfterCreate(t *testing.T) {
	repo := SetupTestCodeRepo()

	exists := repo.CodeExists("EAFJG")
	assert.Equal(t, false, exists)
	_ = repo.AddCode("EAFJG")

	exists = repo.CodeExists("EAFJG")
	assert.Equal(t, true, exists)
}

func TestCodeShouldBeGoneAfterDelete(t *testing.T) {
	repo := SetupTestCodeRepo()

	exists := repo.CodeExists("EAFJG")
	assert.Equal(t, false, exists)
	_ = repo.AddCode("EAFJG")

	exists = repo.CodeExists("EAFJG")
	assert.Equal(t, true, exists)

	repo.RemoveCode("EAFJG")
	exists = repo.CodeExists("EAFJG")
	assert.Equal(t, false, exists)
}

//func SetupTestCodeRepo() *RoomCodeRepoInMemory {
//	return NewRoomCodeRepo()
//}
//
//func TestUniqueRoomCodeShouldBeGenerated(t *testing.T) {
//	rc := SetupTestCodeRepo()
//	amount := 10000
//	for i := 0; i < amount; i++ {
//		code := rc.CreateRoomCode()
//		digitCount := len(strconv.Itoa(code))
//		assert.Equal(t, true, digitCount == 5) // check whether the digit contains 5 chars
//		validCode := rc.ValidRoomCode(code)
//		assert.Equal(t, true, validCode) // check whether the code was added to the memory
//	}
//	assert.Equal(t, amount, len(rc.codes)) // check if it created only unique codes
//}
//
//func TestCodesShouldBeDeleted(t *testing.T) {
//	rc := SetupTestCodeRepo()
//	for i := 0; i < 10000; i++ {
//		code := rc.CreateRoomCode()
//		digitCount := len(strconv.Itoa(code))
//		assert.Equal(t, true, digitCount == 5) // check whether the digit contains 5 chars
//		validCode := rc.ValidRoomCode(code)
//		assert.Equal(t, true, validCode) // check whether the code was added to the memory
//		rc.DeleteRoomCode(code)
//	}
//	assert.Equal(t, 0, len(rc.codes)) // check if it created only unique codes
//}
//
//func TestInvalidRoomCodesShouldNotBeAccepted(t *testing.T) {
//	rc := SetupTestCodeRepo()
//	for i := 0; i <= 99999; i++ {
//		valid := rc.ValidRoomCode(i)
//		assert.Equal(t, false, valid)
//	}
//	assert.Equal(t, 0, len(rc.codes)) // check if it created only unique codes
//}
