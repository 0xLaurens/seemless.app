package services

import (
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"math/rand"
)

type CodeService struct {
	repo ports.CodeRepo
}

var _ ports.CodeService = (*CodeService)(nil)

const (
	Letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func NewCodeService(repo ports.CodeRepo) *CodeService {
	return &CodeService{
		repo: repo,
	}
}

func generateRoomCode() data.RoomCode {
	bytes := make([]byte, 5)
	for i := range bytes {
		bytes[i] = Letters[rand.Intn(len(Letters))]
	}
	return data.RoomCode(bytes)
}

func (cs *CodeService) CreateRoomCode() (data.RoomCode, error) {
	code := generateRoomCode()
	for cs.repo.CodeExists(code) { // prevent collisions
		code = generateRoomCode()
	}

	err := cs.repo.AddCode(code)
	if err != nil {
		return "", err
	}

	return code, nil
}

func (cs *CodeService) RemoveRoomCode(code data.RoomCode) error {
	cs.repo.RemoveCode(code)
	return nil
}
