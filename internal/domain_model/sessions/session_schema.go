package sessions

import "time"

type Session struct {
	SessionId string
	UserId    int32
	ExpiresAt time.Time
}
