package models

import "time"

type GroupChat struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatorID uint      `json:"creator_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
