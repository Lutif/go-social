package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type FeedAuthor struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type FeedPost struct {
	Post
	Comments_count int64     `json:"comments_count"`
	Comments       []Comment `json:"comments"`
	Username       string    `json:"username"`
}

type FeedStore struct {
	db *sql.DB
}

func (s *FeedStore) GetUserFeed(ctx context.Context, userId int64, pagination Paginated) ([]FeedPost, error) {
	// Validate and sanitize sort direction
	sortDirection := "DESC"
	if pagination.SORT == "ASC" {
		sortDirection = "ASC"
	}

	query := `
		SELECT
			p.id,
			p.content,
			p.title,
			p.tags,
			p.version,
			p.updated_at,
			p.user_id,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON c.post_id = p.id
		JOIN followers f ON p.user_id = f.user_id
		WHERE f.follower_id = $1 OR p.user_id = $1
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + sortDirection + `
		LIMIT $2
		OFFSET $3
	`

	var feed []FeedPost

	res, err := s.db.QueryContext(
		ctx, query, userId, pagination.LIMIT, pagination.OFFSET,
	)

	if err != nil {
		return feed, err
	}

	defer res.Close()

	for res.Next() {
		var feedPost FeedPost

		res.Scan(
			&feedPost.ID,
			&feedPost.Content,
			&feedPost.Title,
			pq.Array(&feedPost.Tags),
			&feedPost.Version,
			&feedPost.UpdatedAt,
			&feedPost.UserID,
			&feedPost.Username,
			&feedPost.Comments_count,
		)
		feed = append(feed, feedPost)
	}

	return feed, nil
}
