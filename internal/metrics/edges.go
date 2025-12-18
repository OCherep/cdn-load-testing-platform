package metrics

import "sync"

type EdgeMetric struct {
	Edge    string
	IP      string
	Latency int64
}

var (
	mu    sync.Mutex
	edges = map[string][]int64{}
)

func RecordEdge(edge string, ip string, lat int64) {
	key := edge + "|" + ip
	mu.Lock()
	defer mu.Unlock()
	edges[key] = append(edges[key], lat)
}

func Snapshot() map[string][]int64 {
	mu.Lock()
	defer mu.Unlock()
	out := edges
	edges = map[string][]int64{}
	return out
}
