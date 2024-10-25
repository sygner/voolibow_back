package global

import "time"

var TOKENS = map[string]*TokenValue{}

type TokenValue struct {
	UserId     int32
	Username   string
	ExpireTime time.Time
}
