package mocks

import (
	"errors"
	"time"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/pkg/settings"
)

type GroupChatService struct{}

func NewGroupChatService() *GroupChatService {
	return &GroupChatService{}
}

func (service GroupChatService) reset() {
	groupChats = []domain.GroupChat{
		{ID: 1, Name: "First group chat", CreatedAt: time.Now()},
	}
	groupMessages = []domain.GroupMessage{
		{ID: 1, Message: "First message", UserID: 1, ChatID: 1, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Second message", UserID: 1, ChatID: 2, IsEdited: false, CreatedAt: time.Now()},
		{ID: 1, Message: "Third message", UserID: 1, ChatID: 3, IsEdited: true, CreatedAt: time.Now()},
	}
}

func (service GroupChatService) GetChat(chatId uint) (domain.GroupChat, error) {
	service.reset()
	for _, chat := range groupChats {
		if chat.ID == chatId {
			return chat, nil
		}
	}
	return domain.GroupChat{}, errors.New("chat not found")
}

func (service GroupChatService) GetHistory(chatId, fromId uint) ([]domain.GroupMessage, error) {
	service.reset()
	history := []domain.GroupMessage{}
	for _, message := range groupMessages {
		if len(history) >= settings.CHAT_HISTORY_BLOCK_SIZE {
			break
		}
		if message.ChatID == chatId && message.ID <= fromId {
			history = append(history, message)
		}
	}
	return history, nil
}

func (service GroupChatService) CreateChat(chat *domain.GroupChat) error {
	service.reset()
	groupChats = append(groupChats, *chat)
	return nil
}

func (service GroupChatService) CreateMessage(message *domain.GroupMessage) error {
	service.reset()
	groupMessages = append(groupMessages, *message)
	return nil
}

func (service GroupChatService) UpdateChat(chatId uint, updates map[string]any) error {
	service.reset()
	_, err := service.GetChat(chatId)
	return err
}

func (service GroupChatService) UpdateMessage(messageId uint, updatedMessage string) error {
	service.reset()
	for _, message := range groupMessages {
		if message.ID == messageId {
			return nil
		}
	}
	return errors.New("message not found")
}

func (service GroupChatService) DeleteChat(chatId uint) error {
	service.reset()
	for i, chat := range groupChats {
		if chat.ID == chatId {
			groupChats = append(groupChats[:i], groupChats[i+1:]...)
			return nil
		}
	}
	return errors.New("chat not found")
}

func (service GroupChatService) DeleteMessage(messageId uint) error {
	service.reset()
	for i, message := range groupMessages {
		if message.ID == messageId {
			groupMessages = append(groupMessages[:i], groupMessages[i+1:]...)
			return nil
		}
	}
	return errors.New("message not found")
}

func (service GroupChatService) IsAdminOfChat(chatId, userId uint) (bool, error) {
	service.reset()
	for _, chat := range groupChats {
		if chat.ID == chatId {
			return chat.CreatorID == userId, nil
		}
	}
	return false, errors.New("chat not found")
}

func (service GroupChatService) IsMessageBelongsToChat(messageId, chatId uint) (bool, error) {
	service.reset()
	for _, message := range groupMessages {
		if message.ID == messageId {
			return message.ChatID == chatId, nil
		}
	}
	return false, errors.New("message not found")
}

func (service GroupChatService) IsAuthorOfMessage(messageId, userId uint) (bool, error) {
	service.reset()
	for _, message := range groupMessages {
		if message.ID == messageId {
			return message.UserID == userId, nil
		}
	}
	return false, errors.New("message not found")
}