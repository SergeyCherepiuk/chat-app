package groupchatstorage

import (
	"fmt"
	"strings"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/jmoiron/sqlx"
)

type GroupChatStorage interface {
	GetChat(chatId uint) (models.GroupChat, error)
	GetHistory(chatId uint) ([]models.GroupChatMessage, error)
	Create(chat *models.GroupChat) error
	Update(chatId uint, updates map[string]any) error
	Delete(chatId uint) error
	IsAdmin(chatId, userId uint) (bool, error)
}

type GroupChatStorageImpl struct {
	pdb *sqlx.DB
}

func New(pdb *sqlx.DB) *GroupChatStorageImpl {
	return &GroupChatStorageImpl{pdb: pdb}
}

func (storage GroupChatStorageImpl) GetChat(chatId uint) (models.GroupChat, error) {
	query := `SELECT * FROM group_chats WHERE id = :chat_id`
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return models.GroupChat{}, err
	}
	defer stmt.Close()

	chat := models.GroupChat{}
	if err := stmt.Get(&chat, namedParams); err != nil {
		return models.GroupChat{}, err
	}

	return chat, nil
}

func (storage GroupChatStorageImpl) GetHistory(chatId uint) ([]models.GroupChatMessage, error) {
	query := `SELECT * FROM group_chat_messages WHERE chat_id = :chat_id ORDER BY created_at`
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return []models.GroupChatMessage{}, err
	}
	defer stmt.Close()

	history := []models.GroupChatMessage{}
	if err := stmt.Select(&history, namedParams); err != nil {
		return []models.GroupChatMessage{}, err
	}

	return history, nil
}

func (storage GroupChatStorageImpl) Create(chat *models.GroupChat) error {
	query := `INSERT INTO group_chats (name) VALUES (:name) RETURNING *`
	namedParams := map[string]any{
		"name": chat.Name,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	insertedChat := models.GroupChat{}
	if err := stmt.Get(&insertedChat, namedParams); err != nil {
		return err
	}
	
	chat.ID = insertedChat.ID
	chat.CreatedAt = insertedChat.CreatedAt
	return nil
}

func (storage GroupChatStorageImpl) Update(chatId uint, updates map[string]any) error {
	query := []byte("UPDATE group_chats SET ")
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	updatePairs := []string{}
	for k, v := range updates {
		switch v.(type) {
		case string, rune, byte, []byte:
			updatePairs = append(updatePairs, fmt.Sprintf("%s = '%s'", k, v))
		default:
			updatePairs = append(updatePairs, fmt.Sprintf("%s = %s", k, v))
		}
	}
	query = append(query, strings.Join(updatePairs, ", ")...)
	query = append(query, "WHERE id = :chat_id"...)

	stmt, err := storage.pdb.PrepareNamed(string(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage GroupChatStorageImpl) Delete(chatId uint) error {
	query := `DELETE FROM group_chats WHERE id = :chat_id`
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	stmt, err := storage.pdb.PrepareNamed(string(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage GroupChatStorageImpl) IsAdmin(chatId, userId uint) (bool, error) {
	query := `SELECT FROM group_chats WHERE id = :chat_id AND creator_id = :creator_id`
	namedParams := map[string]any{
		"chat_id": chatId,
		"creator_id": userId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	if result, err := stmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}