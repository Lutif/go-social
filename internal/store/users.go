package store

import (
	"context"
	"database/sql"
)

type UsersStore struct {
	db *sql.DB
}

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

func (s *UsersStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, password)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	return s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)
}

func (s *UsersStore) GetById(ctx context.Context, userID int64) (user User, err error) {
	query := `
	 SELECT username, email, created_at
	 FROM users 
	 WHERE id=$1
	`
	user = User{
		ID: userID,
	}
	err = s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		err = ErrNotFound
	}
	return
}
