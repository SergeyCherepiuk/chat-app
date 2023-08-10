package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/jmoiron/sqlx"
)

type GroupChatService struct {
	getChatStmt    *sqlx.NamedStmt
	getHistoryStmt *sqlx.NamedStmt
	createStmt     *sqlx.NamedStmt
	updateColumns  []string
	updateStmts    map[string]*sqlx.NamedStmt
	deleteStmt     *sqlx.NamedStmt
	isAdminStmt    *sqlx.NamedStmt
}

func NewGroupChatService() *GroupChatService {
	service := GroupChatService{
		updateColumns: []string{"name"},
		updateStmts:   make(map[string]*sqlx.NamedStmt),
	}
	utils.MustPrepareNamed(db, &service.getChatStmt, `SELECT * FROM group_chats WHERE id = :chat_id`)
	utils.MustPrepareNamed(db, &service.getHistoryStmt, `SELECT * FROM group_messages WHERE chat_id = :chat_id ORDER BY created_at`)
	utils.MustPrepareNamed(db, &service.createStmt, `INSERT INTO group_chats (name, creator_id) VALUES (:name, :creator_id) RETURNING *`)
	utils.MustPrepareNamedMap(db, service.updateColumns, service.updateStmts, `UPDATE group_chats SET %s = :value WHERE id = :chat_id`)
	utils.MustPrepareNamed(db, &service.deleteStmt, `DELETE FROM group_chats WHERE id = :chat_id`)
	utils.MustPrepareNamed(db, &service.isAdminStmt, `SELECT FROM group_chats WHERE id = :chat_id AND creator_id = :creator_id`)
	return &service
}

func (service GroupChatService) GetChat(chatId uint) (domain.GroupChat, error) {
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	chat := domain.GroupChat{}
	if err := service.getChatStmt.Get(&chat, namedParams); err != nil {
		return domain.GroupChat{}, err
	}

	return chat, nil
}

func (service GroupChatService) GetHistory(chatId uint) ([]domain.GroupMessage, error) {
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	history := []domain.GroupMessage{}
	if err := service.getHistoryStmt.Select(&history, namedParams); err != nil {
		return []domain.GroupMessage{}, err
	}

	return history, nil
}

func (service GroupChatService) Create(chat *domain.GroupChat) error {
	namedParams := map[string]any{
		"name":       chat.Name,
		"creator_id": chat.CreatorID,
	}

	insertedChat := domain.GroupChat{}
	if err := service.createStmt.Get(&insertedChat, namedParams); err != nil {
		return err
	}

	chat.ID = insertedChat.ID
	chat.CreatedAt = insertedChat.CreatedAt
	return nil
}

func (service GroupChatService) Update(chatId uint, updates map[string]any) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// TODO: Potentially n+1 problem
	for _, column := range service.updateColumns {
		if value, ok := updates[column]; ok {
			stmt := tx.NamedStmt(service.updateStmts[column])
			namedParams := map[string]any{
				"value":   value,
				"chat_id": chatId,
			}

			if _, err := stmt.Exec(namedParams); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (service GroupChatService) Delete(chatId uint) error {
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	_, err := service.deleteStmt.Exec(namedParams)
	return err
}

func (service GroupChatService) IsAdmin(chatId, userId uint) (bool, error) {
	namedParams := map[string]any{
		"chat_id":    chatId,
		"creator_id": userId,
	}

	if result, err := service.isAdminStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}
