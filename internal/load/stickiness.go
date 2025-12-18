package load

import "sync"

type StickinessTracker struct {
	mu     sync.Mutex
	last   map[string]string
	total  int
	sticky int
}

func NewStickinessTracker() *StickinessTracker {
	return &StickinessTracker{
		last: make(map[string]string),
	}
}

func (s *StickinessTracker) Record(clientID, edge string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if prev, ok := s.last[clientID]; ok {
		if prev == edge {
			s.sticky++
		}
	}
	s.last[clientID] = edge
	s.total++
}

func (s *StickinessTracker) Ratio() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.total == 0 {
		return 1
	}
	return float64(s.sticky) / float64(s.total)
}
