package common

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

type Logger struct {
	mu     sync.Mutex
	writer io.Writer
}

type Store struct {
	data   map[string][]byte
	logger Logger
}

func NewStore(w io.Writer) *Store {
	store := Store{
		logger: Logger{
			writer: w,
		},
	}

	store.data = make(map[string][]byte)

	return &store
}

func (s *Store) AppendLog(message string) {
	s.logger.mu.Lock()
	defer s.logger.mu.Unlock()
	fmt.Fprintln(s.logger.writer, message)
}

func (s *Store) GetKey(key string) ([]byte, error) {
	value, ok := s.data[key]

	if !ok {
		return nil, fmt.Errorf("key not found: %s", key)
	}

	return value, nil
}

func (s *Store) SetKey(key string, value []byte) error {
	s.AppendLog(fmt.Sprintf("set %s %s", key, value))

	s.data[key] = value

	return nil
}

func (s *Store) ListKeys() []string {
	keys := []string{}
	for key := range s.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
