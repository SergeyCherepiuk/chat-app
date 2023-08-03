package userstorage

import (
	"fmt"
	"strings"

	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type UserStorage interface {
	GetById(userId uint) (models.User, error)
	GetByUsername(username string) (models.User, error)
	Update(userId uint, updates map[string]any) error
	Delete(userId uint) error
}

type UserStorageImpl struct {
	pdb *sqlx.DB
	rdb *redis.Client
}

func New(pdb *sqlx.DB, rdb *redis.Client) *UserStorageImpl {
	return &UserStorageImpl{pdb: pdb, rdb: rdb}
}

func (storage UserStorageImpl) GetById(userId uint) (models.User, error) {
	query := `SELECT * FROM users WHERE id = :user_id`
	namedParams := map[string]any{
		"user_id": userId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	user := models.User{}
	if err := stmt.Get(&user, namedParams); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (storage UserStorageImpl) GetByUsername(username string) (models.User, error) {
	query := `SELECT * FROM users WHERE username = :username`
	namedParams := map[string]any{
		"username": username,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	user := models.User{}
	if err := stmt.Get(&user, namedParams); err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (storage UserStorageImpl) Update(userId uint, updates map[string]any) error {
	query := []byte("UPDATE users SET ")
	namedParams := map[string]any{
		"user_id": userId,
	}

	updatesPairs := []string{}
	for k, v := range updates {
		switch v.(type) {
		case string, rune, byte, []byte:
			updatesPairs = append(updatesPairs, fmt.Sprintf("%s = '%s'", k, v))
		default:
			updatesPairs = append(updatesPairs, fmt.Sprintf("%s = %s", k, v))
		}
	}
	query = append(query, strings.Join(updatesPairs, ", ")...)
	query = append(query, "WHERE id = :user_id"...)

	stmt, err := storage.pdb.PrepareNamed(string(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage UserStorageImpl) Delete(userId uint) error {
	query := `DELETE FROM users WHERE id = :user_id`
	namedParams := map[string]any{
		"user_id": userId,
	}

	stmt, err := storage.pdb.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}
