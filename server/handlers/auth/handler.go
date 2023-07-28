package authhandler

import "github.com/SergeyCherepiuk/chat-app/storage"

type AuthHandler struct {
	storage *storage.AuthStorage
}

func NewAuthHandler(storage *storage.AuthStorage) *AuthHandler {
	return &AuthHandler{storage: storage}
}
