package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dettarune/kos-finder/internal/model"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindUserByUsernameOrEmail(ctx context.Context, username, email string) (*model.RegisterRequest, error) {
	query := `
		SELECT id, email, full_name, username, password, phone, role
		FROM users
		WHERE (username = ? OR email = ?) AND deleted_at IS NULL
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, username, email)

	var user model.RegisterRequest
	err := row.Scan(
		&user.Id,
		&user.Email,
		&user.Full_name,
		&user.Username,
		&user.Password,
		&user.Phone,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}


func (r *UserRepo) InsertUser(ctx context.Context, user *model.RegisterRequest) error {
	query := `
		INSERT INTO users (email, full_name, username, password, phone, role)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Full_name,
		user.Username,
		user.Password,
		user.Phone,
	)

	if err != nil {
		return err
	}

	return nil
}
