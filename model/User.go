package model

import "time"

type User struct {
	ID        int64
	UserID    int64
	CreatedAt time.Time
}
