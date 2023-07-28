package userhandler

import "github.com/SergeyCherepiuk/chat-app/storage"

type UserHandler struct {
	storage *storage.UserStorage
}

func NewUserHandler(storage *storage.UserStorage) *UserHandler {
	return &UserHandler{storage: storage}
}
