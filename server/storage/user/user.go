package userstorage

import (
	"errors"

	"github.com/SergeyCherepiuk/chat-app/models"
	"gorm.io/gorm"
)

type UserStorage interface {
	GetById(userId uint) (models.User, error)
	GetByUsername(username string) (models.User, error)
	Update(userId uint, updates map[string]any) error
	Delete(userId uint) error
}

type UserStorageImpl struct {
	pdb *gorm.DB
}

func New(pdb *gorm.DB) *UserStorageImpl {
	return &UserStorageImpl{pdb: pdb}
}

func (storage UserStorageImpl) GetById(userId uint) (models.User, error) {
	user := models.User{}
	r := storage.pdb.First(&user, userId)
	if r.Error != nil {
		return models.User{}, r.Error
	}
	return user, nil
}

func (storage UserStorageImpl) GetByUsername(username string) (models.User, error) {
	user := models.User{}
	r := storage.pdb.First(&user).Where("username = ?", username)
	if r.Error != nil {
		return models.User{}, r.Error
	} else if r.RowsAffected < 1 {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (storage UserStorageImpl) Update(userId uint, updates map[string]any) error {
	user := models.User{ID: userId}
	r := storage.pdb.Model(&user).Updates(updates)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("user not found")
	}
	return nil
}

func (storage UserStorageImpl) Delete(userId uint) error {
	r := storage.pdb.Delete(&models.User{}, userId)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("user not found")
	}
	return nil
}
