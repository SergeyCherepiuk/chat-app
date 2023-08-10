package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	sessionManager domain.SessionManagerService
}

func NewAuthService(sessionManager domain.SessionManagerService) *AuthService {
	return &AuthService{sessionManager: sessionManager}
}

func (service AuthService) SignUp(user domain.User) (uuid.UUID, uint, error) {
	tx, err := db.Beginx()
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

	sessionId, err := service.sessionManager.Create(userId)
	if err != nil {
		tx.Rollback()
		return uuid.UUID{}, 0, err
	}

	tx.Commit()
	return sessionId, userId, nil
}

func (service AuthService) Login(username, password string) (uuid.UUID, uint, error) {
	query := `SELECT * FROM users WHERE username = :username`
	namedParams := map[string]any{
		"username": username,
	}

	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return uuid.UUID{}, 0, err
	}
	defer stmt.Close()

	user := domain.User{}
	if err := stmt.Get(&user, namedParams); err != nil {
		return uuid.UUID{}, 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return uuid.UUID{}, 0, err
	}

	sessionId, err := service.sessionManager.Create(user.ID)
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	return sessionId, user.ID, nil
}

func (service AuthService) Check(sessionId uuid.UUID) (uint, error) {
	return service.sessionManager.Check(sessionId)
}

func (service AuthService) Logout(sessionId uuid.UUID) error {
	return service.sessionManager.Invalidate(sessionId)
}
