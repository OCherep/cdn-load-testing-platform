package load

import "sync"

type StickinessTracker struct {
	mu    sync.Mutex
	edges map[string]map[string]int
}

func NewStickinessTracker() *StickinessTracker {
	return &StickinessTracker{
		edges: make(map[string]map[string]int),
	}
}

func (s *StickinessTracker) Record(client, edge string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.edges[client]; !ok {
		s.edges[client] = make(map[string]int)
	}
	s.edges[client][edge]++
}

func (s *StickinessTracker) Ratio(client string) float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	m := s.edges[client]
	if len(m) == 0 {
		return 0
	}

	var max, sum int
	for _, v := range m {
		sum += v
		if v > max {
			max = v
		}
	}
	return float64(max) / float64(sum)
}
