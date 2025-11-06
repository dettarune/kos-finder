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
		&user.Role,
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

	result, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Full_name,
		user.Username,
		user.Password,
		user.Phone,
		user.Role,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("no rows inserted")
	}

	return nil
}

func (r *UserRepo) UpdateUserVerification(ctx context.Context, username string, verified bool) error {
    query := `UPDATE users SET is_verified = ? WHERE username = ?`
    _, err := r.db.ExecContext(ctx, query, verified, username)
    return err
}