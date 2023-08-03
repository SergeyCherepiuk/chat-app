package chatstorage

import (
	"fmt"
	"strings"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/jmoiron/sqlx"
)

type ChatStorage interface {
	GetAllChats() ([]models.Chat, error)
	GetChatById(chatId uint) (models.Chat, error)
	CreateChat(chat models.Chat) error
	UpdateChat(chatId uint, updates map[string]any) error
	DeleteChat(chatId uint) error
	IsChatExists(chatId uint) bool
	GetAllMessages(chatId uint) ([]models.Message, error)
	CreateMessage(message *models.Message) error
	UpdateMessage(messageId uint, updatedMessage string) error
	DeleteMessage(messageId uint) error
}

type ChatStorageImpl struct {
	pdb *sqlx.DB
}

func New(pdb *sqlx.DB) *ChatStorageImpl {
	return &ChatStorageImpl{pdb: pdb}
}

func (storage ChatStorageImpl) GetAllChats() ([]models.Chat, error) {
	query := `SELECT * FROM chats`

	stmt, err := storage.pdb.Preparex(query)
	if err != nil {
		return []models.Chat{}, err
	}
	defer stmt.Close()

	chats := []models.Chat{}
	if err := stmt.Select(&chats); err != nil {
		return []models.Chat{}, err
	}

	return chats, nil
}

func (storage ChatStorageImpl) GetChatById(chatId uint) (models.Chat, error) {
	query := `SELECT * FROM chats WHERE id = :chat_id`
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return models.Chat{}, err
	}
	defer stmt.Close()

	chat := models.Chat{}
	if err := stmt.Get(&chat, namedParams); err != nil {
		return models.Chat{}, err
	}

	return chat, nil
}

func (storage ChatStorageImpl) CreateChat(chat models.Chat) error {
	query := `INSERT INTO chats (name) VALUES (:name)`
	namedParams := map[string]any{
		"name": chat.Name,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage ChatStorageImpl) UpdateChat(chatId uint, updates map[string]any) error {
	query := []byte("UPDATE chats SET ")
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	updatesPairs := []string{}
	for k, v := range updates {
		switch v.(type) {
		case string, rune, byte, []byte:
			updatesPairs = append(updatesPairs, fmt.Sprintf("%s = '%s'", k, v))
		default:
			updatesPairs = append(updatesPairs, fmt.Sprintf("%s = %s", k, v))
		}
	}
	query = append(query, strings.Join(updatesPairs, ", ")...)
	query = append(query, "WHERE id = :chat_id"...)

	stmt, err := storage.pdb.PrepareNamed(string(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage ChatStorageImpl) DeleteChat(chatId uint) error {
	query := `DELETE FROM chats WHERE id = :chat_id`
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage ChatStorageImpl) IsChatExists(chatId uint) bool {
	query := `SELECT FROM users WHERE id = :chat_id`
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err == nil
}

func (storage ChatStorageImpl) GetAllMessages(chatId uint) ([]models.Message, error) {
	query := `SELECT * FROM messages WHERE chat_id = :chat_id`
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return []models.Message{}, err
	}
	defer stmt.Close()

	messages := []models.Message{}
	if err := stmt.Select(&messages, namedParams); err != nil {
		return []models.Message{}, err
	}

	return messages, nil
}

func (storage ChatStorageImpl) CreateMessage(message *models.Message) error {
	query := `INSERT INTO messages (message, sent_at, user_id, chat_id) VALUES (:message, :sent_at, :user_id, :chat_id) RETURNING id`
	namedParams := map[string]any{
		"message": message.Message,
		"sent_at": message.SentAt,
		"user_id": message.UserID,
		"chat_id": message.ChatID,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var messageId uint
	if err := stmt.Get(&messageId, namedParams); err != nil {
		return err
	}

	message.ID = messageId
	return nil
}

func (storage ChatStorageImpl) UpdateMessage(messageId uint, updatedMessage string) error {
	query := `UPDATE messages SET message := message`
	namedParams := map[string]any{
		"message_id": messageId,
		"message":    updatedMessage,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage ChatStorageImpl) DeleteMessage(messageId uint) error {
	query := `DELETE FROM message WHERE id = :message_id`
	namedParams := map[string]any{
		"message_id": messageId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}
