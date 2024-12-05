package common

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string][]byte
	log  io.Writer
}

func NewStore(w io.Writer) *Store {
	store := Store{
		data: make(map[string][]byte),
		log:  w,
	}
	return &store
}

func (s *Store) GetKey(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]

	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

func (s *Store) SetKey(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value

	fmt.Fprintf(s.log, "set %s %s\n", key, value)
}

func (s *Store) ListKeys() []string {
	s.mu.RLock()

	keys := make([]string, len(s.data))
	i := 0
	for key := range s.data {
		keys[i] = key
		i++
	}
	s.mu.RUnlock()

	sort.Strings(keys)
	return keys
}
