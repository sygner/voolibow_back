package models

import "sync"

type WebSocketLoginDTO struct {
	Key string `json:"key"`
}

type ThreadSafe struct {
	*sync.Mutex
}

type WebSocketResponse struct {
	Event string      `json:"event"`
	State string      `json:"state"`
	Data  interface{} `json:"data"`
}
