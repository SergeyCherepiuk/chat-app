package postgres

import (
	"github.com/SergeyCherepiuk/chat-app/domain"
	"github.com/SergeyCherepiuk/chat-app/utils"
	"github.com/jmoiron/sqlx"
)

type UserService struct {
	getByIdStmt       *sqlx.NamedStmt
	getByUsernameStmt *sqlx.NamedStmt
	updateColumns     []string
	updateStmts       map[string]*sqlx.NamedStmt
	deleteStmt        *sqlx.NamedStmt
}

func NewUserService() *UserService {
	service := UserService{
		updateColumns: []string{"first_name", "last_name", "username", "description"},
		updateStmts: make(map[string]*sqlx.NamedStmt),
	}
	utils.MustPrepareNamed(db, &service.getByIdStmt, `SELECT * FROM users WHERE id = :user_id`)
	utils.MustPrepareNamed(db, &service.getByUsernameStmt, `SELECT * FROM users WHERE username = :username`)
	utils.MustPrepareNamedMap(db, service.updateColumns, service.updateStmts, `UPDATE users SET %s = :value WHERE id = :user_id`)
	utils.MustPrepareNamed(db, &service.deleteStmt, `DELETE FROM users WHERE id = :user_id`)
	return &service
}

func (service UserService) GetById(userId uint) (domain.User, error) {
	namedParams := map[string]any{
		"user_id": userId,
	}

	user := domain.User{}
	if err := service.getByIdStmt.Get(&user, namedParams); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (service UserService) GetByUsername(username string) (domain.User, error) {
	namedParams := map[string]any{
		"username": username,
	}

	user := domain.User{}
	if err := service.getByUsernameStmt.Get(&user, namedParams); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (service UserService) Update(userId uint, updates map[string]any) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	// TODO: Potentially n+1 problem
	for _, column := range service.updateColumns {
		if value, ok := updates[column]; ok {
			stmt := tx.NamedStmt(service.updateStmts[column])
			namedParams := map[string]any{
				"value":   value,
				"user_id": userId,
			}

			if _, err := stmt.Exec(namedParams); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	return nil
}

func (service UserService) Delete(userId uint) error {
	namedParams := map[string]any{
		"user_id": userId,
	}

	_, err := service.deleteStmt.Exec(namedParams)
	return err
}
