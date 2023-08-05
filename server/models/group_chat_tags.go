package models

import "time"

type GroupChatTags struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Color     uint      `json:"color"`
	ChatID    uint      `json:"chat_id" db:"chat_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
