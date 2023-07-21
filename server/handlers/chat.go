package handlers

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/SergeyCherepiuk/chat-app/storage"
	"github.com/gofiber/contrib/websocket"
)

type ChatHandler struct {
	storage *storage.ChatStorage
}

func NewChatHandler(storage *storage.ChatStorage) *ChatHandler {
	return &ChatHandler{storage: storage}
}

var connections []*websocket.Conn

func (handler ChatHandler) Enter(c *websocket.Conn) {
	defer c.Close()

	userId, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Println(errors.New("Unathorized"))
		return
	}

	chatId, err := strconv.ParseUint(c.Params("chat_id", ""), 10, 64)
	if err != nil {
		log.Println(err)
		return
	}

	if isExists := handler.storage.IsExists(uint(chatId)); !isExists {
		log.Println(err)
		return
	}

	connections = append(connections, c)

	messages, err := handler.storage.GetAllForChat(uint(chatId))
	if err != nil {
		log.Println(err)
		return
	}

	for _, message := range messages {
		if err := c.WriteMessage(websocket.TextMessage, []byte(message.Text)); err != nil {
			log.Println(err)
			return
		}
	}

	for {
		_, text, err := c.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		message := models.Message{
			Text:   string(text),
			UserID: userId,
			ChatID: uint(chatId),
			SentAt: time.Now(),
		}
		if err := handler.storage.Create(message); err != nil {
			log.Println(err)
			return
		}

		for _, ws := range connections {
			if !reflect.DeepEqual(c, ws) {
				if err := ws.WriteMessage(websocket.TextMessage, text); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
