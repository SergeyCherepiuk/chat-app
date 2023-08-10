package mocks

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/google/uuid"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (service AuthService) SignUp(user domain.User) (uuid.UUID, uint, error) {
	sessionId := uuid.New()
	var userId uint = 1
	return sessionId, userId, nil
}

func (service AuthService) Login(username, password string) (uuid.UUID, uint, error) {
	sessionId := uuid.New()
	var userId uint = 1
	return sessionId, userId, nil
}

func (service AuthService) Check(sessionId uuid.UUID) (uint, error) {
	var userId uint = 1
	return userId, nil
}

func (service AuthService) Logout(sessionId uuid.UUID) error {
	return nil
}
