package repository

import (
	"database/sql"
	"fmt"

	"github.com/akyTheDev/currency-bot/internal/domain"
	"github.com/akyTheDev/currency-bot/internal/models"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

type UserRepository interface {
	CreateUser(chatID int64) error
	DeleteUser(chatID int64) error
	GetAllUsers() ([]models.User, error)
}

func (ur *PostgresUserRepository) CreateUser(chatID int64) error {
	query := `
	INSERT INTO users (chat_id) VALUES ($1)
	ON CONFLICT (chat_id) DO NOTHING
	`
	result, err := ur.db.Exec(
		query,
		chatID,
	)

	if err != nil {
		return fmt.Errorf("CreateUser exec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("CreateUser rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserAlreadyExists
	}

	return nil
}

func (ur *PostgresUserRepository) DeleteUser(chatID int64) error {
	query := `
		DELETE FROM users WHERE chat_id = $1
	`
	result, err := ur.db.Exec(
		query,
		chatID,
	)

	if err != nil {
		return fmt.Errorf("DeleteUser exec: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("DeleteUser rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (ur *PostgresUserRepository) GetAllUsers() ([]models.User, error) {
	query := `
	SELECT id, chat_id FROM users
	`

	rows, err := ur.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("GetAllUsers query: %w", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&user.ID,
			&user.ChatID,
		)
		if err != nil {
			return nil, fmt.Errorf("GetAllUsers scan: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
