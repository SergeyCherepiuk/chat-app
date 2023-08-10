package mocks

import (
	"errors"
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
)

type DirectMessageService struct{}

func NewDirectMessageService() *DirectMessageService {
	return &DirectMessageService{}
}

func (service DirectMessageService) reset() {
	directMessages = []domain.DirectMessage{
		{ID: 1, Message: "First message", From: 1, To: 2, CreatedAt: time.Now()},
		{ID: 2, Message: "Second message", From: 2, To: 1, CreatedAt: time.Now()},
	}
}

func (service DirectMessageService) GetHistory(userId, companionId uint) ([]domain.DirectMessage, error) {
	service.reset()
	history := []domain.DirectMessage{}
	for _, message := range directMessages {
		if (message.From == userId && message.To == companionId) || (message.From == companionId && message.To == userId) {
			history = append(history, message)
		}
	}
	return history, nil
}

func (service DirectMessageService) Create(message *domain.DirectMessage) error {
	service.reset()
	directMessages = append(directMessages, *message)
	return nil
}

func (service DirectMessageService) Update(messageId uint, updatedMessage string) error {
	service.reset()
	for _, message := range directMessages {
		if message.ID == messageId {
			message.Message = updatedMessage
			message.IsEdited = true
			return nil
		}
	}
	return errors.New("message not found")
}

func (service DirectMessageService) Delete(messageId uint) error {
	service.reset()
	for i, message := range directMessages {
		if message.ID == messageId {
			directMessages = append(directMessages[:i], directMessages[i+1:]...)
			return nil
		}
	}
	return errors.New("message not found")
}

func (service DirectMessageService) DeleteAll(userId, companionId uint) error {
	service.reset()
	for i := 0; i < len(directMessages); i++ {
		message := directMessages[i]
		if (message.From == userId && message.To == companionId) || (message.From == companionId && message.To == userId) {
			directMessages = append(directMessages[:i], directMessages[i+1:]...)
			i--
		}
	}
	return nil
}

func (service DirectMessageService) IsBelongsToChat(messageId, userId, companionId uint) (bool, error) {
	service.reset()
	for _, message := range directMessages {
		if message.ID == messageId && ((message.From == userId && message.To == companionId) || (message.From == companionId && message.To == userId)) {
			return true, nil
		}
	}
	return false, errors.New("message not found in chat")
}

func (service DirectMessageService) IsAuthor(messageId, userId uint) (bool, error) {
	service.reset()
	for _, message := range directMessages {
		if message.ID == messageId {
			return message.From == userId, nil
		}
	}
	return false, errors.New("message not found")
}
