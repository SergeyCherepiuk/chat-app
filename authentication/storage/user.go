package storage

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/authentication/models"
	"gorm.io/gorm"
)

type UserStorage struct {
	pdb *gorm.DB
}

func NewUserStorage(pdb *gorm.DB) *UserStorage {
	return &UserStorage{pdb: pdb}
}

func (storage UserStorage) GetById(userId uint) (models.User, error) {
	user := models.User{}
	r := storage.pdb.First(&user, userId)
	if r.Error != nil {
		return models.User{}, r.Error
	} else if r.RowsAffected < 1 {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (storage UserStorage) GetByUsername(username string) (models.User, error) {
	user := models.User{}
	r := storage.pdb.First(&user).Where("username = ?", username)
	if r.Error != nil {
		return models.User{}, r.Error
	} else if r.RowsAffected < 1 {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (storage UserStorage) Update(userId uint, updates map[string]any) error {
	user := models.User{ID: userId}
	r := storage.pdb.Model(&user).Updates(updates)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("user not found")
	}
	return nil
}

func (storage UserStorage) Delete(userId uint) error {
	r := storage.pdb.Delete(&models.User{}, userId)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("user not found")
	}
	return nil
}
