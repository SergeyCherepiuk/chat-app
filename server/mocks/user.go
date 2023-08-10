package mocks

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/domain"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (service UserService) reset() {
	users = []domain.User{
		{ID: 1, FirstName: "John", LastName: "Doe", Username: "johndoe", Password: "HashedSecret123!"},
		{ID: 2, FirstName: "Mark", LastName: "Watson", Username: "markwatson", Password: "HashedSecret123!"},
	}
}

func (service UserService) GetById(userId uint) (domain.User, error) {
	service.reset()
	for _, u := range users {
		if u.ID == userId {
			return u, nil
		}
	}
	return domain.User{}, errors.New("user not found")
}

func (service UserService) GetByUsername(username string) (domain.User, error) {
	service.reset()
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return domain.User{}, errors.New("user not found")
}

func (service UserService) Update(userId uint, updates map[string]any) error {
	service.reset()
	_, err := service.GetById(userId)
	return err
}

func (service UserService) Delete(userId uint) error {
	service.reset()
	_, err := service.GetById(userId)
	return err
}
