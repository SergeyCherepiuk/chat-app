package authhandler

import (
	authstorage "github.com/SergeyCherepiuk/chat-app/storage/auth"
)

type AuthHandler struct {
	storage authstorage.AuthStorage
}

func New(storage authstorage.AuthStorage) *AuthHandler {
	return &AuthHandler{storage: storage}
}
