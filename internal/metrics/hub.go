package metrics

import "sync"

type Hub struct {
	mu      sync.Mutex
	clients map[string][]chan LiveMetric
}

func NewHub() *Hub {
	return &Hub{clients: map[string][]chan LiveMetric{}}
}

func (h *Hub) Subscribe(testID string) chan LiveMetric {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan LiveMetric, 10)
	h.clients[testID] = append(h.clients[testID], ch)
	return ch
}

func (h *Hub) Publish(testID string, m LiveMetric) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, ch := range h.clients[testID] {
		ch <- m
	}
}
