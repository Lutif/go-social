package store

import (
	"context"
	"database/sql"
	"errors"
)

var (
	ErrNotFound = errors.New("resource not found")
)

type Store struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (Post, error)
		Update(context.Context, *Post) error
		Delete(context.Context, int64) error
		ListByUserId(context.Context, int64, *[]Post) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetById(context.Context, int64) (Comment, error)
		Update(context.Context, *Comment) error
		Delete(context.Context, int64) error
	}
	Users interface {
		Create(context.Context, *User) error
		GetById(context.Context, int64) (User, error)
	}
}

func NewPostgresStorage(db *sql.DB) Store {
	return Store{
		Posts:    &PostsStore{db},
		Users:    &UsersStore{db},
		Comments: &CommentsStore{db},
	}
}
