package userhandler

import (
	userstorage "github.com/SergeyCherepiuk/chat-app/storage/user"
)

type UserHandler struct {
	storage userstorage.UserStorage
}

func New(storage userstorage.UserStorage) *UserHandler {
	return &UserHandler{storage: storage}
}
