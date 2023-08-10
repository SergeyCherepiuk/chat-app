package postgres

import (
	"fmt"
	"strings"

	"github.com/SergeyCherepiuk/chat-app/domain"
)

type postgresUserService struct{}

func NewUserService() *postgresUserService {
	return &postgresUserService{}
}

func (storage postgresUserService) GetById(userId uint) (domain.User, error) {
	query := `SELECT * FROM users WHERE id = :user_id`
	namedParams := map[string]any{
		"user_id": userId,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()

	user := domain.User{}
	if err := stmt.Get(&user, namedParams); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (storage postgresUserService) GetByUsername(username string) (domain.User, error) {
	query := `SELECT * FROM users WHERE username = :username`
	namedParams := map[string]any{
		"username": username,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()

	user := domain.User{}
	if err := stmt.Get(&user, namedParams); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (storage postgresUserService) Update(userId uint, updates map[string]any) error {
	query := []byte("UPDATE users SET ")
	namedParams := map[string]any{
		"user_id": userId,
	}

	updatePairs := []string{}
	for k, v := range updates {
		switch v.(type) {
		case string, rune, byte, []byte:
			updatePairs = append(updatePairs, fmt.Sprintf("%s = '%s'", k, v))
		default:
			updatePairs = append(updatePairs, fmt.Sprintf("%s = %s", k, v))
		}
	}
	query = append(query, strings.Join(updatePairs, ", ")...)
	query = append(query, "WHERE id = :user_id"...)

	stmt, err := db.PrepareNamed(string(query))
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}

func (storage postgresUserService) Delete(userId uint) error {
	query := `DELETE FROM users WHERE id = :user_id`
	namedParams := map[string]any{
		"user_id": userId,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(namedParams)
	return err
}
