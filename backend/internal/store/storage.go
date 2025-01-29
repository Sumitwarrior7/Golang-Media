package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("record not found")
	ErrConflict          = errors.New("resource already exists")
	QueryTimeoutduration = time.Second * 5
)

// *Composition allows you to keep the Storage struct independent from the specific details of how posts and users are handled, making it more flexible
// and easier to modify or extend.
// *Inheritance would tie everything together in a more rigid way, making it harder to change the implementation of the stores without impacting
// the Storage struct or needing extra complexity.
type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetaData, error)
		GetPostsByUserId(context.Context, int64, PaginatedFeedQuery) ([]PostWithMetaData, error)
	}
	Users interface {
		Create(context.Context, *sql.Tx, *User) error
		GetAllUsers(context.Context, PaginatedFeedQuery) ([]User, error)
		GetById(context.Context, int64) (*User, error)
		GetByEmail(context.Context, string) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		GetById(context.Context, int64) (*Comment, error)
		Create(context.Context, *Comment) error
		Delete(context.Context, int64) error
		Update(context.Context, *Comment) error
		GetByPostId(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(context.Context, int64, int64) error
		Unfollow(context.Context, int64, int64) error
		GetFollowedUsersById(context.Context, int64) ([]FollowedUserDetails, error)
	}
	Roles interface {
		GetByName(context.Context, string) (*Role, error)
	}
}

func NewPostgresStorage(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostsStore{db},
		Users:     &UsersStore{db},
		Comments:  &CommentsStore{db},
		Followers: &FollowersStore{db},
	}
}

// The withTx function manages database transactions in Go. It begins a transaction, executes a provided callback function (fn) with the transaction,
// and ensures proper cleanup: if the callback returns an error, the transaction is rolled back; otherwise, it commits the transaction.
// This simplifies transaction handling and ensures atomicity.
func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
