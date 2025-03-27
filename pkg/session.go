package pkg

import (
	"sync"
)

// SessionStore 定义了 session 存储的接口
type SessionStore interface {
	// Store 存储一个 session
	Store(key string, value interface{})
	// Load 加载一个 session
	Load(key string) (interface{}, bool)
	// Delete 删除一个 session
	Delete(key string)
	// Range 遍历所有 session
	Range(f func(key string, value interface{}) bool)
}

// MemorySessionStore 使用内存实现的 session 存储
type MemorySessionStore struct {
	store sync.Map
}

// NewMemorySessionStore 创建一个新的内存 session 存储
func NewMemorySessionStore() *MemorySessionStore {
	return &MemorySessionStore{}
}

func (s *MemorySessionStore) Store(key string, value interface{}) {
	s.store.Store(key, value)
}

func (s *MemorySessionStore) Load(key string) (interface{}, bool) {
	return s.store.Load(key)
}

func (s *MemorySessionStore) Delete(key string) {
	s.store.Delete(key)
}

func (s *MemorySessionStore) Range(f func(key string, value interface{}) bool) {
	s.store.Range(func(key, value interface{}) bool {
		return f(key.(string), value)
	})
}
