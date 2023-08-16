package domain

import "time"

type DirectMessage struct {
	ID        uint      `json:"id"`
	Message   string    `json:"message"`
	From      uint      `json:"message_from" db:"message_from"`
	To        uint      `json:"message_to" db:"message_to"`
	IsEdited  bool      `json:"is_edited" db:"is_edited"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type DirectMessageService interface {
	GetHistory(userId, companionId, fromId uint) ([]DirectMessage, error)
	Create(message *DirectMessage) error
	Update(messageId uint, updatedMessage string) error
	Delete(messageId uint) error
	DeleteAll(userId, companionId uint) error
	IsBelongsToChat(messageId, userId, companionId uint) (bool, error)
	IsAuthor(messageId, userId uint) (bool, error)
}
