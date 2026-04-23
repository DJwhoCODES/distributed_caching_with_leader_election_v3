package cache

import (
	"sync"
	"time"
)

type Store struct {
	mu    sync.RWMutex
	items map[string]*Entry
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]*Entry),
	}
}

func (s *Store) Set(key string, value []byte, ttl uint32) {
	var exp int64 = 0

	if ttl > 0 {
		exp = time.Now().Add(time.Duration(ttl) * time.Second).Unix()
	}

	s.mu.Lock()
	s.items[key] = &Entry{
		Value:     value,
		ExpiresAt: exp,
	}

	s.mu.Unlock()
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.Lock()
	entry, ok := s.items[key]
	s.mu.RUnlock()

	if !ok {
		return nil, false
	}

	if entry.IsExpired() {
		s.Delete(key)
		return nil, false
	}

	return entry.Value, true
}

func (s *Store) Delete(key string) {
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
}

func (s *Store) StartJanitor(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now().Unix()

			s.mu.Lock()

			for k, v := range s.items {
				if v.ExpiresAt != 0 && v.ExpiresAt < now {
					delete(s.items, k)
				}
			}

			s.mu.Unlock()
		}
	}()
}
