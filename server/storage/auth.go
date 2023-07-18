package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthStorage struct {
	pdb *gorm.DB
	rdb *redis.Client
}

func NewAuthStorage(pdb *gorm.DB, rdb *redis.Client) *AuthStorage {
	return &AuthStorage{pdb: pdb, rdb: rdb}
}

func (storage AuthStorage) SignUp(user models.User) (uuid.UUID, error) {
	sessionId := uuid.New()

	tx := storage.pdb.Begin()
	pipe := storage.rdb.Pipeline()

	r := tx.Create(&user)
	if r.Error != nil {
		tx.Rollback()
		pipe.Discard()
		return uuid.UUID{}, r.Error
	}

	err := pipe.Set(context.Background(), sessionId.String(), fmt.Sprint(user.ID), 7*24*time.Hour).Err()
	if err != nil {
		tx.Rollback()
		pipe.Discard()
		return uuid.UUID{}, err
	}

	err = pipe.Set(context.Background(), fmt.Sprint(user.ID), sessionId.String(), 7*24*time.Hour).Err()
	if err != nil {
		tx.Rollback()
		pipe.Discard()
		return uuid.UUID{}, err
	}

	tx.Commit()
	pipe.Exec(context.Background())
	return sessionId, nil
}

func (storage AuthStorage) Login(username, password string) (uuid.UUID, error) {
	user := models.User{}
	r := storage.pdb.Where("username = ?", username).First(&user)
	if r.Error != nil {
		return uuid.UUID{}, r.Error
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return uuid.UUID{}, err
	}

	oldSessionId, err := storage.rdb.Get(context.Background(), fmt.Sprint(user.ID)).Result()
	if err == nil {
		storage.rdb.Del(context.Background(), oldSessionId)
	}

	sessionId := uuid.New()
	pipe := storage.rdb.Pipeline()
	pipe.Set(context.Background(), sessionId.String(), fmt.Sprint(user.ID), 7*24*time.Hour)
	pipe.Set(context.Background(), fmt.Sprint(user.ID), sessionId.String(), 7*24*time.Hour)
	_, err = pipe.Exec(context.Background())
	if err != nil {
		return uuid.UUID{}, err
	}

	return sessionId, nil
}

func (storage AuthStorage) Check(sessionId uuid.UUID) (uint, error) {
	userIdStr, err := storage.rdb.Get(context.Background(), sessionId.String()).Result()
	if err != nil {
		return 0, errors.New("session not found")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

func (storage AuthStorage) Logout(sessionId uuid.UUID) error {
	userId, err := storage.rdb.Get(context.Background(), sessionId.String()).Result()
	if err != nil {
		return err
	}

	pipe := storage.rdb.Pipeline()
	pipe.Del(context.Background(), sessionId.String())
	pipe.Del(context.Background(), userId)
	_, err = pipe.Exec(context.Background())
	return err
}
