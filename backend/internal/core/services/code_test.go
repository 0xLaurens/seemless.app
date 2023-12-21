package services

import (
	"github.com/stretchr/testify/assert"
	"laurensdrop/internal/adapters/secondary/repo"
	"laurensdrop/internal/core/data"
	"slices"
	"testing"
)

func SetupCodeService() *CodeService {
	codeRepo := repo.NewCodeRepoInMemory()
	return NewCodeService(codeRepo)
}

func TestCreateRoomCreatesUniqueCodes(t *testing.T) {
	service := SetupCodeService()
	amountOfCodes := 10000
	seen := make(map[data.RoomCode]bool)
	for i := 0; i < amountOfCodes; i++ {
		code, err := service.CreateRoomCode()
		assert.NoError(t, err)
		seen[code] = true
	}
	assert.Equal(t, amountOfCodes, len(seen))
}

func TestCreateRoomCodeCreatesCodesWithALengthOf5(t *testing.T) {
	service := SetupCodeService()
	amountOfCodes := 10000
	for i := 0; i < amountOfCodes; i++ {
		code, err := service.CreateRoomCode()
		assert.NoError(t, err)
		assert.Equal(t, 5, len(code))
	}
}

func TestCreateRoomCodeCreatesCodesThatOnlyContainAllowedLetters(t *testing.T) {
	service := SetupCodeService()
	allowedChars := make([]int, len(Letters))
	for char := range Letters {
		allowedChars = append(allowedChars, char)
	}

	amountOfCodes := 10000
	for i := 0; i < amountOfCodes; i++ {
		code, err := service.CreateRoomCode()
		assert.NoError(t, err)
		for char := range code {
			assert.Equal(t, true, slices.Contains(allowedChars, char))
		}
	}
}
