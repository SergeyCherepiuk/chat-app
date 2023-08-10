package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	sessionManager domain.SessionManagerService
	signUpStmt     *sqlx.NamedStmt
	loginStmt      *sqlx.NamedStmt
}

func NewAuthService(sessionManager domain.SessionManagerService) *AuthService {
	service := AuthService{sessionManager: sessionManager}
	utils.MustPrepareNamed(db, &service.signUpStmt, `INSERT INTO users (first_name, last_name, username, password, description, profile_picture) VALUES (:first_name, :last_name, :username, :password, :description, :profile_picture) RETURNING id`)
	utils.MustPrepareNamed(db, &service.loginStmt, `SELECT * FROM users WHERE username = :username`)
	return &service
}

func (service AuthService) SignUp(user domain.User) (uuid.UUID, uint, error) {
	tx, err := db.Beginx()
	if err != nil {
		return uuid.UUID{}, 0, err
	}

	namedParams := map[string]any{
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
		"username":        user.Username,
		"password":        user.Password,
		"description":     user.Description,
		"profile_picture": user.ProfilePicture,
	}

	var userId uint
	if err := service.signUpStmt.Get(&userId, namedParams); err != nil {
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
	namedParams := map[string]any{
		"username": username,
	}

	user := domain.User{}
	if err := service.loginStmt.Get(&user, namedParams); err != nil {
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
