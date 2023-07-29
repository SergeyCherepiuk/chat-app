package chatstorage

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/models"
	"gorm.io/gorm"
)

type ChatStorage interface {
	GetAllChats() ([]models.Chat, error)
	GetChatById(chatId uint) (models.Chat, error)
	CreateChat(chat *models.Chat) error
	UpdateChat(chatId uint, updates map[string]any) error
	DeleteChat(chatId uint) error
	IsChatExists(chatId uint) bool 
	GetAllMessages(chatId uint) ([]models.Message, error)
	CreateMessage(message *models.Message) error
	UpdateMessage(messageId uint, updatedText string) error
	DeleteMessage(messageId uint) error
}

type ChatStorageImpl struct {
	pdb *gorm.DB
}

func New(pdb *gorm.DB) *ChatStorageImpl {
	return &ChatStorageImpl{pdb: pdb}
}

func (storage ChatStorageImpl) GetAllChats() ([]models.Chat, error) {
	chats := []models.Chat{}
	if r := storage.pdb.Find(&chats); r.Error != nil {
		return []models.Chat{}, r.Error
	}

	return chats, nil
}

func (storage ChatStorageImpl) GetChatById(chatId uint) (models.Chat, error) {
	chat := models.Chat{}
	if r := storage.pdb.First(&chat, chatId); r.Error != nil {
		return models.Chat{}, r.Error
	}

	return chat, nil
}

func (storage ChatStorageImpl) CreateChat(chat *models.Chat) error {
	return storage.pdb.Create(chat).Error
}

func (storage ChatStorageImpl) UpdateChat(chatId uint, updates map[string]any) error {
	chat := models.Chat{ID: chatId}
	r := storage.pdb.Model(&chat).Updates(updates)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("chat not found")
	}
	return nil
}

func (storage ChatStorageImpl) DeleteChat(chatId uint) error {
	r := storage.pdb.Delete(&models.Chat{}, chatId)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("chat not found")
	}
	return nil
}

func (storage ChatStorageImpl) IsChatExists(chatId uint) bool {
	r := storage.pdb.First(&models.Chat{}, chatId)
	return r.Error == nil && r.RowsAffected > 0
}

func (storage ChatStorageImpl) GetAllMessages(chatId uint) ([]models.Message, error) {
	messages := []models.Message{}
	r := storage.pdb.Where("chat_id = ?", chatId).Find(&messages)
	if r.Error != nil {
		return []models.Message{}, r.Error
	}
	return messages, nil
}

func (storage ChatStorageImpl) CreateMessage(message *models.Message) error {
	return storage.pdb.Create(message).Error
}

func (storage ChatStorageImpl) UpdateMessage(messageId uint, updatedText string) error {
	message := models.Message{ID: messageId}
	return storage.pdb.Model(&message).Update("message", updatedText).Error
}

func (storage ChatStorageImpl) DeleteMessage(messageId uint) error {
	return storage.pdb.Delete(&models.Message{}, messageId).Error
}
