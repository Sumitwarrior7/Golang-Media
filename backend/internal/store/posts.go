package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Post struct {
	Id        int64
	Content   string
	Title     string
	UserId    int64
	Tags      []string
	CreatedAt string
	UpdatedAt string
	Version   int64 // Getting added through add_version migrations[It is mainly used for optimistic concurrency]
	Comments  []Comment
	User      User
}

type PostWithMetaData struct {
	Post
	CommentCount int
}

type PostsStore struct {
	db *sql.DB
}

func (s *PostsStore) Create(ctx context.Context, post *Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&post.Id,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostsStore) GetById(ctx context.Context, postId int64) (*Post, error) {
	query := `
		SELECT id, title, user_id, content, created_at, updated_at, tags, version 
		FROM posts WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	post := &Post{}

	// Scan objedt must follow the order in which sql query is being executed
	err := s.db.QueryRowContext(
		ctx,
		query,
		postId,
	).Scan(
		&post.Id,
		&post.Title,
		&post.UserId,
		&post.Content,
		&post.CreatedAt,
		&post.UpdatedAt,
		pq.Array(&post.Tags),
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return post, nil
}

func (s *PostsStore) Delete(ctx context.Context, postId int64) error {
	query := `
		DELETE FROM posts WHERE id = $1	
	`

	res, err := s.db.ExecContext(
		ctx,
		query,
		postId,
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

func (s *PostsStore) Update(ctx context.Context, post *Post) error {
	query := `
		UPDATE posts 
		SET title = $1, content = $2, version = version+1
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.Id,
		post.Version,
	).Scan(
		&post.Version,
	)

	if err != nil {
		return err
	}
	return nil
}

// Shows the posts of the user and the other users that he followed
func (s *PostsStore) GetUserFeed(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE 
			(
				p.user_id = $1 
				AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
			)
			OR  
			(
				p.user_id IN (
					SELECT user_id 
					FROM followers 
					WHERE follower_id = $1
				)
				AND (p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%')
			)	
		GROUP BY p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2 OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset, fq.Search)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var feed []PostWithMetaData
	for rows.Next() {
		var p PostWithMetaData
		err := rows.Scan(
			&p.Id,
			&p.UserId,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		feed = append(feed, p)
	}

	return feed, nil
}

func (s *PostsStore) GetPostsByUserId(ctx context.Context, userID int64, fq PaginatedFeedQuery) ([]PostWithMetaData, error) {
	query := `
		SELECT 
			p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		JOIN users u ON p.user_id = u.id
		LEFT JOIN comments c ON p.id = c.post_id
		WHERE p.user_id = $1 
		GROUP BY p.id, p.user_id, p.title, p.content, p.created_at, p.version, p.tags, u.username
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3;
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, userID, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []PostWithMetaData
	for rows.Next() {
		var p PostWithMetaData
		err := rows.Scan(
			&p.Id,
			&p.UserId,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
			&p.Version,
			pq.Array(&p.Tags),
			&p.User.Username,
			&p.CommentCount,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}
