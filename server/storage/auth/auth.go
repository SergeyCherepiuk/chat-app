package authstorage

import (
	"github.com/SergeyCherepiuk/chat-app/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthStorage interface {
	SignUp(user models.User) (uuid.UUID, uint, error)
	Login(username, password string) (uuid.UUID, uint, error)
	Check(sessionId uuid.UUID) (uint, error)
	Logout(sessionId uuid.UUID) error
}

type AuthStorageImpl struct {
	pdb            *sqlx.DB
	sessionManager *SessionManager
}

func New(pdb *sqlx.DB, sessionManager *SessionManager) *AuthStorageImpl {
	return &AuthStorageImpl{pdb: pdb, sessionManager: sessionManager}
}

func (storage AuthStorageImpl) SignUp(user models.User) (uuid.UUID, uint, error) {
	tx, err := storage.pdb.Beginx()
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	query := `INSERT INTO users (first_name, last_name, username, password, description, profile_picture) VALUES (:first_name, :last_name, :username, :password, :description, :profile_picture) RETURNING id`
	namedParams := map[string]any{
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"username":        user.Username,
		"password":        user.Password,
		"description":     user.Description,
		"profile_picture": user.ProfilePicture,
	}

	stmt, err := tx.PrepareNamed(query)
	if err != nil {
		tx.Rollback()
		return uuid.UUID{}, 0, err
	}
	defer stmt.Close()

	var userId uint
	if err := stmt.Get(&userId, namedParams); err != nil {
		tx.Rollback()
		return uuid.UUID{}, 0, err
	}

	sessionId, err := storage.sessionManager.Create(userId)
	if err != nil {
		tx.Rollback()
		return uuid.UUID{}, 0, err
	}

	tx.Commit()
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

	sessionId, err := storage.sessionManager.Create(user.ID)
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	return sessionId, user.ID, nil
}

func (storage AuthStorageImpl) Check(sessionId uuid.UUID) (uint, error) {
	return storage.sessionManager.Check(sessionId)
}

func (storage AuthStorageImpl) Logout(sessionId uuid.UUID) error {
	return storage.sessionManager.Invalidate(sessionId)
}
