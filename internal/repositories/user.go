package repositories

import (
	"database/sql"
	"fmt"

	"github.com/zaher1307/subscription-service/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) Create(user *models.User) error {
	stmt, err := r.DB.Prepare(`
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, created_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	return stmt.QueryRow(user.Name, user.Email).Scan(&user.ID, &user.CreatedAt)
}

func (r *UserRepository) GetByID(id int) (*models.User, error) {
	stmt, err := r.DB.Prepare(`
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRow(id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user %d not found", id)
		}
		return nil, err
	}

	return &user, nil
}
