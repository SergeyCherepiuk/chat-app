package groupchatstorage

import (
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/utils"
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
	pdb            *sqlx.DB
	getChatStmt    *sqlx.NamedStmt
	getHistoryStmt *sqlx.NamedStmt
	createStmt     *sqlx.NamedStmt
	updateColumns  []string
	updateStmts    map[string]*sqlx.NamedStmt
	deleteStmt     *sqlx.NamedStmt
	isAdminStmt    *sqlx.NamedStmt
}

func New(pdb *sqlx.DB) *GroupChatStorageImpl {
	storage := GroupChatStorageImpl{
		pdb:           pdb,
		updateColumns: []string{"name"},
		updateStmts:   make(map[string]*sqlx.NamedStmt),
	}
	utils.MustPrepareNamed(pdb, &storage.getChatStmt, `SELECT * FROM group_chats WHERE id = :chat_id`)
	utils.MustPrepareNamed(pdb, &storage.getHistoryStmt, `SELECT * FROM group_chat_messages WHERE chat_id = :chat_id ORDER BY created_at`)
	utils.MustPrepareNamed(pdb, &storage.createStmt, `INSERT INTO group_chats (name, creator_id) VALUES (:name, :creator_id) RETURNING *`)
	utils.MustPrepareNamedMap(pdb, storage.updateColumns, storage.updateStmts, `UPDATE group_chats SET %s = :value WHERE id = :chat_id`)
	utils.MustPrepareNamed(pdb, &storage.deleteStmt, `DELETE FROM group_chats WHERE id = :chat_id`)
	utils.MustPrepareNamed(pdb, &storage.isAdminStmt, `SELECT FROM group_chats WHERE id = :chat_id AND creator_id = :creator_id`)
	return &storage
}

func (storage GroupChatStorageImpl) GetChat(chatId uint) (models.GroupChat, error) {
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	chat := models.GroupChat{}
	if err := storage.getChatStmt.Get(&chat, namedParams); err != nil {
		return models.GroupChat{}, err
	}

	return chat, nil
}

func (storage GroupChatStorageImpl) GetHistory(chatId uint) ([]models.GroupChatMessage, error) {
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	history := []models.GroupChatMessage{}
	if err := storage.getHistoryStmt.Select(&history, namedParams); err != nil {
		return []models.GroupChatMessage{}, err
	}

	return history, nil
}

func (storage GroupChatStorageImpl) Create(chat *models.GroupChat) error {
	namedParams := map[string]any{
		"name":       chat.Name,
		"creator_id": chat.CreatorID,
	}

	insertedChat := models.GroupChat{}
	if err := storage.createStmt.Get(&insertedChat, namedParams); err != nil {
		return err
	}

	chat.ID = insertedChat.ID
	chat.CreatedAt = insertedChat.CreatedAt
	return nil
}

func (storage GroupChatStorageImpl) Update(chatId uint, updates map[string]any) error {
	tx, err := storage.pdb.Beginx()
	if err != nil {
		return err
	}

	// TODO: Potentially n+1 problem
	for _, column := range storage.updateColumns {
		if value, ok := updates[column]; ok {
			stmt := tx.NamedStmt(storage.updateStmts[column])
			namedParams := map[string]any{
				"value":   value,
				"chat_id": chatId,
			}

			_, err := stmt.Exec(namedParams)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (storage GroupChatStorageImpl) Delete(chatId uint) error {
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	_, err := storage.deleteStmt.Exec(namedParams)
	return err
}

func (storage GroupChatStorageImpl) IsAdmin(chatId, userId uint) (bool, error) {
	namedParams := map[string]any{
		"chat_id":    chatId,
		"creator_id": userId,
	}

	if result, err := storage.isAdminStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}
