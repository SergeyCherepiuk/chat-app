package userstorage

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/models"
)

type UserStorageMock struct{}

func NewMock() *UserStorageMock { 
	return &UserStorageMock{}
}

var users []models.User

func (storage UserStorageMock) reset() {
	users = []models.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Username: "johndoe", Password: "HashedSecret123!"},
		{ID: 2, FirstName: "Mark", LastName: "Watson", Username: "markwatson", Password: "HashedSecret123!"},
	}
}

func (storage UserStorageMock) GetById(userId uint) (models.User, error) {
	storage.reset()
	for _, u := range users {
		if u.ID == userId {
			return u, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (storage UserStorageMock) GetByUsername(username string) (models.User, error) {
	storage.reset()
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

func (storage UserStorageMock) Update(userId uint, updates map[string]any) error {
	storage.reset()
	_, err := storage.GetById(userId)
	return err
}

func (storage UserStorageMock) Delete(userId uint) error {
	storage.reset()
	_, err := storage.GetById(userId)
	return err
}
