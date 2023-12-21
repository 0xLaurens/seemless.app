package repo

import (
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
	"sync"
)

type CodeRepoInMemory struct {
	codes map[data.RoomCode]bool
	mu    sync.Mutex
}

var _ ports.CodeRepo = (*CodeRepoInMemory)(nil) //interface impl hack

func NewCodeRepoInMemory() *CodeRepoInMemory {
	return &CodeRepoInMemory{
		codes: map[data.RoomCode]bool{},
	}
}

func (rc *CodeRepoInMemory) AddCode(code data.RoomCode) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if exist := rc.codes[code]; exist {
		return data.DuplicateRoomCode.Error()
	}
	rc.codes[code] = true
	return nil
}

func (rc *CodeRepoInMemory) CodeExists(code data.RoomCode) bool {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	return rc.codes[code]
}

func (rc *CodeRepoInMemory) RemoveCode(code data.RoomCode) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	delete(rc.codes, code)
}
