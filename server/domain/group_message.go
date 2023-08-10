package domain

import "time"

type GroupMessage struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	UserID    uint      `json:"user_id" db:"user_id"`
	ChatID    uint      `json:"chat_id" db:"chat_id"`
	IsEdited  bool      `json:"is_edited" db:"is_edited"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GroupMessageService interface {
	Create(message *GroupMessage) error
	Update(messageId uint, updates map[string]any) error
	Delete(messageId uint) error
	IsBelongsToChat(messageId, chatId uint) (bool, error)
	IsAuthor(messageId, userId uint) (bool, error)
}
