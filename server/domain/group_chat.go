package domain

import "time"

type GroupChat struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatorID uint      `json:"creator_id" db:"creator_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GroupMessage struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	UserID    uint      `json:"user_id" db:"user_id"`
	ChatID    uint      `json:"chat_id" db:"chat_id"`
	IsEdited  bool      `json:"is_edited" db:"is_edited"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GroupChatService interface {
	GetChat(chatId uint) (GroupChat, error)
	GetHistory(chatId uint) ([]GroupMessage, error)
	CreateChat(chat *GroupChat) error
	CreateMessage(message *GroupMessage) error
	UpdateChat(chatId uint, updates map[string]any) error
	UpdateMessage(messageId uint, updatedMessage string) error
	DeleteChat(chatId uint) error
	DeleteMessage(messageId uint) error
	IsAdminOfChat(userId, chatId uint) (bool, error)
	IsMessageBelongsToChat(messageId, chatId uint) (bool, error)
	IsAuthorOfMessage(messageId, userId uint) (bool, error)
}
