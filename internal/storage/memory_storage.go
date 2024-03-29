package storage

import (
	vk_test_task "github.com/fishmanDK/vk-test-task"
	"sync"
)

type MemoryStorage interface {
	SaveTokens(userID int64, tokens vk_test_task.Tokens)
}

type MyStorage struct {
	Db    map[int64]vk_test_task.Tokens
	mutex sync.Mutex
}

func MustMyStorage() *MyStorage {
	return &MyStorage{
		Db:    make(map[int64]vk_test_task.Tokens),
		mutex: sync.Mutex{},
	}
}

func (m *MyStorage) SaveTokens(userID int64, tokens vk_test_task.Tokens) {
	m.mutex.Lock()
	m.Db[userID] = tokens
	m.mutex.Unlock()
}
