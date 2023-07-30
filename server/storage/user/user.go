package userstorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/redis/go-redis/v9"
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
	rdb *redis.Client
}

func New(pdb *gorm.DB, rdb *redis.Client) *UserStorageImpl {
	return &UserStorageImpl{pdb: pdb, rdb: rdb}
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
	r := storage.pdb.Where("username = ?", username).First(&user)
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
	sessionId, _ := storage.rdb.Get(context.Background(), fmt.Sprint(userId)).Result()
	storage.rdb.Del(context.Background(), sessionId, fmt.Sprint(userId)).Result()

	r := storage.pdb.Delete(&models.User{}, userId)
	if r.Error != nil {
		return r.Error
	} else if r.RowsAffected < 1 {
		return errors.New("user not found")
	}
	return nil
}
