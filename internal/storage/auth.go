package storage

import (
	"fmt"
	vk_test_task "github.com/fishmanDK/vk-test-task"
	"github.com/jmoiron/sqlx"
)

type Auth interface {
	GetUserInfo(user vk_test_task.User) (UserData, error)
	CreateUser(newUser vk_test_task.CreateUser) error
}

type AuthDb struct {
	db *sqlx.DB
}

func NewAuthDb(db *sqlx.DB) *AuthDb {
	return &AuthDb{
		db: db,
	}
}

type UserData struct {
	ID    int64  `db:"id"`
	Email string `db:"email"`
	Role  string `db:"role"`
}

func (a *AuthDb) GetUserInfo(user vk_test_task.User) (UserData, error) {
	const op = "storage.GetUserInfo"

	query := "SELECT id, email, role FROM users WHERE email=$1 AND hash_password=$2"

	var userData UserData
	err := a.db.Get(&userData, query, user.Email, user.Password)
	if err != nil {
		return userData, fmt.Errorf("%s: %w", op, err)
	}
	return userData, nil
}

func (a *AuthDb) CreateUser(newUser vk_test_task.CreateUser) error {
	const op = "storage.CreateUser"

	query := "INSERT INTO users (email, hash_password, role) VALUES ($1, $2, $3)"

	tx, err := a.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(query, newUser.Email, newUser.Password, newUser.Role)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
