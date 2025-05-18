package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	customerror "github.com/lutif/go-social/internal/custom-error"
)

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	Tags      []string  `json:"tags"`
	Version   int64     `json:"version"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Comments  []Comment `json:"comments"`
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query :=
		`
			INSERT INTO posts (content, title, user_id, tags)
			VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
		`

	return s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserID,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
}

func (s *PostsStore) GetById(ctx context.Context, id int64) (Post, error) {
	query := `
	 SELECT id, content, title, user_id, tags, version
	 FROM posts
	 WHERE id=$1 
	`
	var post = Post{}
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Content,
		&post.Title,
		&post.UserID,
		pq.Array(&post.Tags),
		&post.Version,
	)
	if err != nil {
		return Post{}, customerror.CheckForCustomErr(err, sql.ErrNoRows, ErrNotFound)
	}
	return post, nil
}

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts
		SET content = $1, title = $2, tags = $3, version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING id, created_at, updated_at, version
	`
	err := s.db.QueryRowContext(ctx, query, post.Content, post.Title, pq.Array(post.Tags), post.ID, post.Version).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)
	if err != nil {
		return customerror.CheckForCustomErr(err, sql.ErrNoRows, ErrNotFound)
	}
	return nil
}

func (s *PostsStore) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM posts WHERE id=$1	
	`
	result, err := s.db.ExecContext(ctx, query, id)

	if err != nil {
		return customerror.CheckForCustomErr(err, sql.ErrNoRows, ErrNotFound)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *PostsStore) ListByUserId(ctx context.Context, userID int64, posts *[]Post) error {
	query := `
	SELECT title, content, tags, posts.created_at, username 
	FROM posts 
	LEFT JOIN users
	ON posts.user_id = users.id
	WHERE posts.user_id = $1
	`

	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		var username string
		if err := rows.Scan(&post.Title, &post.Content, pq.Array(&post.Tags), &post.CreatedAt, &username); err != nil {
			return err
		}
		*posts = append(*posts, post)
	}
	return rows.Err()
}
