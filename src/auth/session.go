package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SergeyCherepiuk/session-auth/src/models"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type SessionManager struct {
	pdb *gorm.DB
	rdb *redis.Client
}

func NewSessionManager(pdb *gorm.DB, rdb *redis.Client) *SessionManager {
	return &SessionManager{pdb: pdb, rdb: rdb}
}

func (manager SessionManager) CreateSession(userId uint) (models.Session, error) {
	session := models.Session{
		UserID:    userId,
		ExpiresAt: time.Now().Add(10 * time.Second),
	}

	r := manager.pdb.Create(&session)
	if r.Error != nil {
		return models.Session{}, r.Error
	}

	sessionJson, err := json.Marshal(session)
	if err == nil {
		manager.rdb.SetEx(
			context.Background(),
			fmt.Sprint(session.ID),
			sessionJson,
			time.Until(session.ExpiresAt).Round(time.Second),
		)
	}

	return session, nil
}

func (manager SessionManager) CheckSession(sessionId uint) (models.Session, error) {
	session := models.Session{}

	data, err := manager.rdb.Get(context.Background(), fmt.Sprint(sessionId)).Result()
	if err != nil {
		r := manager.pdb.First(&session, sessionId)
		if r.Error != nil {
			return models.Session{}, r.Error
		}
		if r.RowsAffected < 1 {
			return models.Session{}, errors.New("session not found")
		}
	} else {
		err = json.Unmarshal([]byte(data), &session)
		if err != nil {
			return models.Session{}, err
		}
	}

	if session.ExpiresAt.Before(time.Now()) {
		return models.Session{}, errors.New("session expired")
	}

	return session, nil
}
