package store

import (
	"context"
	"database/sql"
	"errors"
)

type Comment struct {
	Id        int64
	PostId    int64
	UserId    int64
	Content   string
	CreatedAt string
	User      User
}

type CommentsStore struct {
	db *sql.DB
}

func (s *CommentsStore) GetById(ctx context.Context, commentId int64) (*Comment, error) {
	query := `
		SELECT id, user_id, post_id, content, created_at
		FROM comments
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	// Scan objedt must follow the order in which sql query is being executed
	comment := &Comment{}
	err := s.db.QueryRowContext(
		ctx,
		query,
		commentId,
	).Scan(
		&comment.Id,
		&comment.UserId,
		&comment.PostId,
		&comment.Content,
		&comment.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return comment, nil
}

func (s *CommentsStore) Create(ctx context.Context, comment *Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3) RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		comment.PostId,
		comment.UserId,
		comment.Content,
	).Scan(
		&comment.Id,
		&comment.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *CommentsStore) Delete(ctx context.Context, commentId int64) error {
	query := `
		DELETE FROM comments WHERE id = $1	
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
		commentId,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *CommentsStore) Update(ctx context.Context, comment *Comment) error {
	query := `
		UPDATE comments 
		SET content = $1
		WHERE id = $2;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	_, err := s.db.ExecContext(
		ctx,
		query,
		comment.Content,
		comment.Id,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *CommentsStore) GetByPostId(ctx context.Context, postId int64) ([]Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username
		FROM comments AS c
		JOIN users AS u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	// Scan objedt must follow the order in which sql query is being executed
	rows, err := s.db.QueryContext(
		ctx,
		query,
		postId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}
	for rows.Next() {
		var c Comment
		c.User = User{}
		err := rows.Scan(
			&c.Id,
			&c.PostId,
			&c.UserId,
			&c.Content,
			&c.CreatedAt,
			&c.User.Username,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	return comments, nil
}
