package repo

import (
	"laurensdrop/internal/core/data"
	"laurensdrop/internal/ports"
)

type CodeRepoInMemory struct {
	codes map[data.RoomCode]bool
}

var _ ports.CodeRepo = (*CodeRepoInMemory)(nil) //interface impl hack

func NewCodeRepoInMemory() *CodeRepoInMemory {
	return &CodeRepoInMemory{
		codes: map[data.RoomCode]bool{},
	}
}

func (rc *CodeRepoInMemory) AddCode(code data.RoomCode) error {
	if exist := rc.codes[code]; exist {
		return data.DuplicateRoomCode.Error()
	}
	rc.codes[code] = true
	return nil
}

func (rc *CodeRepoInMemory) CodeExists(code data.RoomCode) bool {
	return rc.codes[code]
}

func (rc *CodeRepoInMemory) RemoveCode(code data.RoomCode) {
	delete(rc.codes, code)
}
