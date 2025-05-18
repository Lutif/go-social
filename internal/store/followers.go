package store

import (
	"context"
	"database/sql"
)

type Follower struct {
	UserID     int64 `json:"user_id"`
	FollowerID int64 `json:"follower_id"`
	ID         int64 `json:"id"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, userID int64, followerID int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`
	_, err := s.db.QueryContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}
	return nil
}

func (s *FollowerStore) Unfollow(ctx context.Context, userID int64, followerID int64) error {
	query := `
		DELETE From followers 
		WHERE user_id = $1 AND follower_id = $2
	`
	_, err := s.db.QueryContext(ctx, query, userID, followerID)
	if err != nil {
		return err
	}
	return nil
}
