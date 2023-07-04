package models

import (
	"time"
)

type Session struct {
	ID        uint `db:"id"`
	UserID    uint
	CreatedAt time.Time `db:"created_at"`
	ExpiresAt time.Time `db:"expires_at"`
}
