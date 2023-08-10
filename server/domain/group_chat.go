package domain

import "time"

type GroupChat struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatorID uint      `json:"creator_id" db:"creator_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GroupChatService interface {
	GetChat(chatId uint) (GroupChat, error)
	GetHistory(chatId uint) ([]GroupMessage, error)
	Create(chat *GroupChat) error
	Update(chatId uint, updates map[string]any) error
	Delete(chatId uint) error
	IsAdmin(chatId, userId uint) (bool, error)
}
