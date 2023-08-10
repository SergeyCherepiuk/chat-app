package mocks

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/google/uuid"
)

type SessionManagerService struct{}

func NewSessionManagerService() *SessionManagerService {
	return &SessionManagerService{}
}

func (service SessionManagerService) reset() {
	users = []domain.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Username: "johndoe", Password: "HashedSecret123!"},
		{ID: 2, FirstName: "Mark", LastName: "Watson", Username: "markwatson", Password: "HashedSecret123!"},
	}
}

func (service SessionManagerService) Create(userId uint) (uuid.UUID, error) {
	service.reset()
	for _, user := range users {
		if user.ID == userId {
			return uuid.New(), nil
		}
	}
	return uuid.UUID{}, errors.New("user not found")
}

func (service SessionManagerService) Check(sessionId uuid.UUID) (uint, error) {
	return 1, nil
}

func (service SessionManagerService) Invalidate(sessionId uuid.UUID) error {
	return nil
}
