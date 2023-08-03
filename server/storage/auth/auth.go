package authstorage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthStorage interface {
	SignUp(user models.User) (uuid.UUID, uint, error)
	Login(username, password string) (uuid.UUID, uint, error)
	Check(sessionId uuid.UUID) (uint, error)
	Logout(sessionId uuid.UUID) error
}

type AuthStorageImpl struct {
	pdb *sqlx.DB
	rdb *redis.Client
}

func New(pdb *sqlx.DB, rdb *redis.Client) *AuthStorageImpl {
	return &AuthStorageImpl{pdb: pdb, rdb: rdb}
}

func (storage AuthStorageImpl) SignUp(user models.User) (uuid.UUID, uint, error) {
	sessionId := uuid.New()

	tx, err := storage.pdb.Beginx()
	pipe := storage.rdb.Pipeline()

	query := `INSERT INTO users (first_name, last_name, username, password) VALUES (:first_name, :last_name, :username, :password) RETURNING id`
	namedParams := map[string]any{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"username":   user.Username,
		"password":   user.Password,
	}

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	defer stmt.Close()

	var userId uint
	if err := stmt.Get(&userId, namedParams); err != nil {
		return uuid.UUID{}, 0, err
	}

	err = pipe.Set(context.Background(), sessionId.String(), fmt.Sprint(userId), 7*24*time.Hour).Err()
	if err != nil {
		tx.Rollback()
		return uuid.UUID{}, 0, err
	}

	err = pipe.Set(context.Background(), fmt.Sprint(userId), sessionId.String(), 7*24*time.Hour).Err()
	if err != nil {
		tx.Rollback()
		pipe.Discard()
		return uuid.UUID{}, 0, err
	}

	tx.Commit()
	pipe.Exec(context.Background())
	return sessionId, userId, nil
}

func (storage AuthStorageImpl) Login(username, password string) (uuid.UUID, uint, error) {
	query := `SELECT * FROM users WHERE username = :username`
	namedParams := map[string]any{
		"username": username,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	defer stmt.Close()

	user := models.User{}
	if err := stmt.Get(&user, namedParams); err != nil {
		return uuid.UUID{}, 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return uuid.UUID{}, 0, err
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
		return uuid.UUID{}, 0, err
	}

	return sessionId, user.ID, nil
}

func (storage AuthStorageImpl) Check(sessionId uuid.UUID) (uint, error) {
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

func (storage AuthStorageImpl) Logout(sessionId uuid.UUID) error {
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
