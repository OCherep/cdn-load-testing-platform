package metrics

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type LiveMetric struct {
	RPS   int     `json:"rps"`
	P95   int     `json:"p95"`
	Error float64 `json:"error"`
}

func Stream(ws *websocket.Conn, ch <-chan LiveMetric) {
	for m := range ch {
		data, _ := json.Marshal(m)
		ws.WriteMessage(websocket.TextMessage, data)
	}
}
