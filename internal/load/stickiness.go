package load

import "sync"

type StickinessTracker struct {
	mu     sync.Mutex
	client map[string]string
	hits   map[string]int
}

func NewStickinessTracker() *StickinessTracker {
	return &StickinessTracker{
		client: make(map[string]string),
		hits:   make(map[string]int),
	}
}

func (s *StickinessTracker) Record(clientID, edge string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.client[clientID]; ok && prev != edge {
		s.hits["switch"]++
	}
	s.client[clientID] = edge
	s.hits["total"]++
}

func (s *StickinessTracker) Ratio() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.hits["total"] == 0 {
		return 1
	}
	return 1 - float64(s.hits["switch"])/float64(s.hits["total"])
}
