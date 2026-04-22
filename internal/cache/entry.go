package cache

import "time"

type Entry struct {
	Value     []byte
	ExpiresAt int64 //unix timestamp
}

func (e *Entry) IsExpired() bool {
	if e.ExpiresAt == 0 {
		return false
	}

	return time.Now().Unix() > e.ExpiresAt
}
