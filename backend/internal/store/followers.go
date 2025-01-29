package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

type FollowedUserDetails struct {
	UserId    int64
	Email     string
	Username  string
	CreatedAt string
}

type FollowersStore struct {
	db *sql.DB
}

func (s *FollowersStore) Follow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		INSERT INTO followers (user_id, follower_id)
		VALUES ($1, $2)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		userId,
		followerId,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrConflict
		}
	}
	return nil
}

func (s *FollowersStore) Unfollow(ctx context.Context, followerId int64, userId int64) error {
	query := `
		DELETE FROM followers 
		WHERE user_id = $1 AND follower_id = $2
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		userId,
		followerId,
	)

	if err != nil {
		return err
	}
	return nil
}

// Returns all the users that are followed by the user with provided id
func (s *FollowersStore) GetFollowedUsersById(ctx context.Context, userId int64) ([]FollowedUserDetails, error) {
	query := `
		SELECT id, email, username, created_at 
		FROM users 
		WHERE id IN
			(SELECT user_id FROM followers 
			WHERE follower_id = $1)
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var followedUsers []FollowedUserDetails
	for rows.Next() {
		var fu FollowedUserDetails
		err := rows.Scan(
			&fu.UserId,
			&fu.Username,
			&fu.Email,
			&fu.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		followedUsers = append(followedUsers, fu)
	}

	if err != nil {
		return nil, err
	}
	return followedUsers, nil
}
