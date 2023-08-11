package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/jmoiron/sqlx"
)

type GroupChatService struct {
	getChatStmt                *sqlx.NamedStmt
	getHistoryStmt             *sqlx.NamedStmt
	createChatStmt             *sqlx.NamedStmt
	createMessageStmt          *sqlx.NamedStmt
	updateChatColumns          []string
	updateChatStmts            map[string]*sqlx.NamedStmt
	updateMessageStmt          *sqlx.NamedStmt
	deleteChatStmt             *sqlx.NamedStmt
	deleteMessageStmt          *sqlx.NamedStmt
	isAdminOfChatStmt          *sqlx.NamedStmt
	isMessageBelongsToChatStmt *sqlx.NamedStmt
	isAuthorOfMessageStmt      *sqlx.NamedStmt
}

func NewGroupChatService() *GroupChatService {
	service := GroupChatService{
		updateChatColumns: []string{"name"},
		updateChatStmts:   make(map[string]*sqlx.NamedStmt),
	}
	utils.MustPrepareNamed(db, &service.getChatStmt, `SELECT * FROM group_chats WHERE id = :chat_id`)
	utils.MustPrepareNamed(db, &service.getHistoryStmt, `SELECT * FROM group_messages WHERE chat_id = :chat_id ORDER BY created_at`)
	utils.MustPrepareNamed(db, &service.createChatStmt, `INSERT INTO group_chats (name, creator_id) VALUES (:name, :creator_id) RETURNING *`)
	utils.MustPrepareNamed(db, &service.createMessageStmt, `INSERT INTO group_messages (message, user_id, chat_id, is_edited) VALUES (:message, :user_id, :chat_id, :is_edited) RETURNING *`)
	utils.MustPrepareNamedMap(db, service.updateChatColumns, service.updateChatStmts, `UPDATE group_chats SET %s = :value WHERE id = :chat_id`)
	utils.MustPrepareNamed(db, &service.updateMessageStmt, `UPDATE group_messages SET message = :message WHERE id = :message_id`)
	utils.MustPrepareNamed(db, &service.deleteChatStmt, `DELETE FROM group_chats WHERE id = :chat_id`)
	utils.MustPrepareNamed(db, &service.deleteMessageStmt, `DELETE FROM group_messages WHERE id = :message_id`)
	utils.MustPrepareNamed(db, &service.isAdminOfChatStmt, `SELECT FROM group_chats WHERE id = :chat_id AND creator_id = :creator_id`)
	utils.MustPrepareNamed(db, &service.isMessageBelongsToChatStmt, `SELECT FROM group_messages WHERE id = :message_id AND chat_id = :chat_id`)
	utils.MustPrepareNamed(db, &service.isAuthorOfMessageStmt, `SELECT FROM group_messages WHERE id = :message_id AND user_id = :user_id`)
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

func (service GroupChatService) CreateChat(chat *domain.GroupChat) error {
	namedParams := map[string]any{
		"name":       chat.Name,
		"creator_id": chat.CreatorID,
	}

	insertedChat := domain.GroupChat{}
	if err := service.createChatStmt.Get(&insertedChat, namedParams); err != nil {
		return err
	}

	chat.ID = insertedChat.ID
	chat.CreatedAt = insertedChat.CreatedAt
	return nil
}

func (service GroupChatService) CreateMessage(message *domain.GroupMessage) error {
	namedParams := map[string]any{
		"message":   message.Message,
		"user_id":   message.UserID,
		"chat_id":   message.ChatID,
		"is_edited": message.IsEdited,
	}

	insertedMessage := domain.GroupMessage{}
	if err := service.createMessageStmt.Get(&insertedMessage, namedParams); err != nil {
		return err
	}

	message.ID = insertedMessage.ID
	message.CreatedAt = insertedMessage.CreatedAt
	return nil
}

func (service GroupChatService) UpdateChat(chatId uint, updates map[string]any) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// TODO: Potentially n+1 problem
	for _, column := range service.updateChatColumns {
		if value, ok := updates[column]; ok {
			stmt := tx.NamedStmt(service.updateChatStmts[column])
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

func (service GroupChatService) UpdateMessage(messageId uint, updatedMessage string) error {
	namedParams := map[string]any{
		"message":    updatedMessage,
		"message_id": messageId,
	}

	_, err := service.updateMessageStmt.Exec(namedParams)
	return err
}

func (service GroupChatService) DeleteChat(chatId uint) error {
	namedParams := map[string]any{
		"chat_id": chatId,
	}

	_, err := service.deleteChatStmt.Exec(namedParams)
	return err
}

func (service GroupChatService) DeleteMessage(messageId uint) error {
	namedParams := map[string]any{
		"message_id": messageId,
	}

	_, err := service.deleteMessageStmt.Exec(namedParams)
	return err
}

func (service GroupChatService) IsAdminOfChat(userId, chatId uint) (bool, error) {
	namedParams := map[string]any{
		"chat_id":    chatId,
		"creator_id": userId,
	}

	if result, err := service.isAdminOfChatStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}

func (service GroupChatService) IsMessageBelongsToChat(messageId, chatId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id": messageId,
		"chat_id":    chatId,
	}

	if result, err := service.isMessageBelongsToChatStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}

func (service GroupChatService) IsAuthorOfMessage(messageId, userId uint) (bool, error) {
	namedParams := map[string]any{
		"message_id": messageId,
		"user_id":    userId,
	}

	if result, err := service.isAuthorOfMessageStmt.Exec(namedParams); err != nil {
		return false, err
	} else if rowsAffected, err := result.RowsAffected(); err != nil {
		return false, err
	} else {
		return rowsAffected > 0, nil
	}
}
