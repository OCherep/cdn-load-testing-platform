package state

import "sync"

type Store struct {
	mu    sync.RWMutex
	tests map[string]TestState
}

func NewStore(_ string) *Store {
	return &Store{
		tests: map[string]TestState{},
	}
}

func (s *Store) MarkSLABreached(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.tests[id]
	t.SLABreached = true
	s.tests[id] = t
}

func (s *Store) UpdateNodes(id string, nodes int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := s.tests[id]
	t.Nodes = nodes
	s.tests[id] = t
}

func (s *Store) GetNodes(id string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.tests[id].Nodes
}
