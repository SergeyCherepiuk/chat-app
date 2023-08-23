package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type SessionManagerService struct{}

func NewSessionManagerService() *SessionManagerService {
	return &SessionManagerService{}
}

func (manager SessionManagerService) Create(userId uint) (uuid.UUID, error) {
	oldSessionId, err := rdb.Get(context.Background(), fmt.Sprint(userId)).Result()
	if err == nil {
		rdb.Del(context.Background(), oldSessionId)
	}

	sessionId := uuid.New()
	pipe := rdb.Pipeline()

	pipe.Set(context.Background(), sessionId.String(), fmt.Sprint(userId), 7*24*time.Hour)
	pipe.Set(context.Background(), fmt.Sprint(userId), sessionId.String(), 7*24*time.Hour)

	_, err = pipe.Exec(context.Background())
	if err != nil {
		return uuid.Nil, err
	}

	return sessionId, nil
}

func (manager SessionManagerService) Check(sessionId uuid.UUID) (uint, error) {
	userIdStr, err := rdb.Get(context.Background(), sessionId.String()).Result()
	if err != nil {
		return 0, errors.New("session not found")
	}

	userId, err := strconv.ParseUint(userIdStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(userId), nil
}

func (manager SessionManagerService) Invalidate(sessionId uuid.UUID) error {
	userId, err := rdb.Get(context.Background(), sessionId.String()).Result()
	if err != nil {
		return err
	}

	pipe := rdb.Pipeline()
	pipe.Del(context.Background(), sessionId.String())
	pipe.Del(context.Background(), userId)
	_, err = pipe.Exec(context.Background())
	return err
}
