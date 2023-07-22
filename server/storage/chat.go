package storage

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/models"
	"gorm.io/gorm"
)

type ChatStorage struct {
	pdb *gorm.DB
}

func NewChatStorage(pdb *gorm.DB) *ChatStorage {
	return &ChatStorage{pdb: pdb}
}

func (storage ChatStorage) GetAllChats() ([]models.Chat, error) {
	chats := []models.Chat{}
	if r := storage.pdb.Find(&chats); r.Error != nil {
		return []models.Chat{}, r.Error
	}

	for i, chat := range chats {
		messages, err := storage.GetAllMessages(chat.ID)
		if err != nil {
			return []models.Chat{}, err
		}
		chats[i].Messages = messages
	}

	return chats, nil
}

func (storage ChatStorage) GetChatById(chatId uint) (models.Chat, error) {
	chat := models.Chat{}
	if r := storage.pdb.First(&chat, chatId); r.Error != nil {
		return models.Chat{}, r.Error
	}

	messages, err := storage.GetAllMessages(chatId)
	if err != nil {
		return models.Chat{}, err
	} 
	chat.Messages = messages
	
	return chat, nil
}

func (storage ChatStorage) CreateChat(chat *models.Chat) error {
	return storage.pdb.Create(chat).Error
}

func (storage ChatStorage) UpdateChat(chatId uint, updates map[string]any) error {
	chat := models.Chat{ID: chatId}
	r := storage.pdb.Model(&chat).Updates(updates)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("chat not found")
	}
	return nil
}

func (storage ChatStorage) DeleteChat(chatId uint) error {
	r := storage.pdb.Delete(&models.Chat{}, chatId)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("chat not found")
	}
	return nil
}

func (storage ChatStorage) IsChatExists(chatId uint) bool {
	r := storage.pdb.First(&models.Chat{}, chatId)
	return r.Error == nil && r.RowsAffected > 0
}

func (storage ChatStorage) GetAllMessages(chatId uint) ([]models.Message, error) {
	messages := []models.Message{}
	r := storage.pdb.Where("chat_id = ?", chatId).Find(&messages)
	if r.Error != nil {
		return []models.Message{}, r.Error
	}
	return messages, nil
}

func (storage ChatStorage) CreateMessage(message *models.Message) error {
	return storage.pdb.Create(message).Error
}

func (storage ChatStorage) UpdateMessage(messageId uint, updatedText string) error {
	message := models.Message{ID: messageId}
	return storage.pdb.Model(&message).Update("message", updatedText).Error
}

func (storage ChatStorage) DeleteMessage(messageId uint) error {
	return storage.pdb.Delete(&models.Message{}, messageId).Error
}
