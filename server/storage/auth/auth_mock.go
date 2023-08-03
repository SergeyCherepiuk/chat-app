package authstorage

import (
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/google/uuid"
)

type AuthStorageMock struct{}

func NewMock() *AuthStorageMock {
	return &AuthStorageMock{}
}

func (storage AuthStorageMock) SignUp(user models.User) (uuid.UUID, uint, error) {
	sessionId := uuid.New()
	var userId uint = 1
	return sessionId, userId, nil
}

func (storage AuthStorageMock) Login(username, password string) (uuid.UUID, uint, error) {
	sessionId := uuid.New()
	var userId uint = 1
	return sessionId, userId, nil
}

func (storage AuthStorageMock) Check(sessionId uuid.UUID) (uint, error) {
	var userId uint = 1
	return userId, nil
}

func (storage AuthStorageMock) Logout(sessionId uuid.UUID) error {
	return nil
}
