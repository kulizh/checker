package state

import (
	"sync"

	"checker/internal/domain"
)

type Store struct {
	mu   sync.RWMutex
	Data map[string]domain.DomainState
}

func New() *Store {
	return &Store{
		Data: make(map[string]domain.DomainState),
	}
}

func (s *Store) Get(domain string) (domain.DomainState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	st, ok := s.Data[domain]
	return st, ok
}

func (s *Store) Set(domain string, st domain.DomainState) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data[domain] = st
}

func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.Data = make(map[string]domain.DomainState)
}
