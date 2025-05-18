package store

import (
	"context"
	"database/sql"

	customerror "github.com/lutif/go-social/internal/custom-error"
)

type Comment struct {
	ID        int64  `json:"id"`
	Content   string `json:"content"`
	AuthorID  int64  `json:"author_id"`
	PostID    int64  `json:"post_id"`
	Likes     int64  `json:"likes"`
	Version   int64  `json:"version"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {

	query := `
		INSERT INTO comments (content, author_id, post_id, likes)
		VALUES ($1, $2, $3, 0)
		RETURNING id, created_at, updated_at, likes
	`

	return s.db.QueryRowContext(
		ctx, query, comment.Content,
		comment.AuthorID, comment.PostID,
	).Scan(
		&comment.ID, &comment.CreatedAt,
		&comment.UpdatedAt, &comment.Likes,
	)
}

func (s *CommentsStore) GetById(ctx context.Context, id int64) (Comment, error) {
	query := `
		SELECT content, author_id, post_id, likes, created_at, updated_at, version
		FROM comments 
		WHERE id=$1
	`
	comment := Comment{
		ID: id,
	}
	err := s.db.QueryRowContext(
		ctx, query, &comment.ID,
	).Scan(
		&comment.Content,
		&comment.AuthorID, &comment.PostID,
		&comment.Likes, &comment.CreatedAt,
		&comment.UpdatedAt, &comment.Version,
	)
	if err != nil {
		return Comment{}, customerror.CheckForCustomErr(err, sql.ErrNoRows, ErrNotFound)
	}
	return comment, nil
}

func (s *CommentsStore) Update(ctx context.Context, comment *Comment) error {
	query := `
		UPDATE comments 
		SET content = $1, likes =$2, version = version + 1
		WHERE id = $3
		RETURNING version
	`

	return s.db.QueryRowContext(
		ctx, query, &comment.Content,
		&comment.Likes, &comment.ID,
	).Scan(&comment.Version)
}

func (s *CommentsStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM comments
		WHERE id = $1
	`
	result, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return customerror.CheckForCustomErr(err, sql.ErrNoRows, ErrNotFound)
	}
	rowEff, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowEff == 0 {
		return customerror.CheckForCustomErr(err, sql.ErrNoRows, ErrNotFound)
	}
	return nil
}
