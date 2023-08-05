package models

import "time"

type GroupChatMessage struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	UserID    uint      `json:"user_id" db:"user_id"`
	ChatID    uint      `json:"chat_id" db:"chat_id"`
	IsEdited  bool      `json:"is_edited" db:"is_edited"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
