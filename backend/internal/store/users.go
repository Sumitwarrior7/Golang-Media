package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	Id        int64    `json:"id"`
	Username  string   `json:"username"`
	Email     string   `json:"email"`
	Password  Password `json:"-"`
	CreatedAt string   `json:"created_at"`
	IsActive  bool     `json:"is_active"`
	RoleId    int64    `json:"role_id"`
	Role      Role     `json:"role"`
}

type Password struct {
	text *string // Otiginal Password
	hash []byte  // Hashed password
}

func (p *Password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

type UsersStore struct {
	db *sql.DB
}

// Hashes the password with bcrypt
func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash
	return nil
}

func (s *UsersStore) Create(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		INSERT INTO users (username, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE name = $4)) 
		RETURNING id, created_at
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()
	role := user.Role.Name
	if role == "" {
		role = "user"
	}

	err := tx.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password.hash,
		role,
	).Scan(
		&user.Id,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"`:
			return ErrDuplicateUsername
		default:
			return err
		}
	}
	return nil
}

func (s *UsersStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password, u.created_at, r.id, r.name, r.level, r.description
		FROM users AS u
		JOIN roles AS r ON u.role_id = r.id
		WHERE u.id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	user := &User{}

	// Scan object must follow the order in which sql query is being executed
	err := s.db.QueryRowContext(
		ctx,
		query,
		userId,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
		&user.Role.Id,
		&user.Role.Name,
		&user.Role.Level,
		&user.Role.Description,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UsersStore) GetAllUsers(ctx context.Context, fq PaginatedFeedQuery) ([]User, error) {
	log.Println("fq.Search :", fq.Search)
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE (username ILIKE '%' || $3 || '%' OR email ILIKE '%' || $3 || '%')
		LIMIT $1 OFFSET $2;
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, fq.Limit, fq.Offset, fq.Search)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var allUsers []User
	for rows.Next() {
		var u User
		err := rows.Scan(
			&u.Id,
			&u.Username,
			&u.Email,
			&u.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		allUsers = append(allUsers, u)
	}

	return allUsers, nil
}

func (s *UsersStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, password, created_at 
		FROM users WHERE email = $1 AND is_active = true
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	user := &User{}

	// Scan object must follow the order in which sql query is being executed
	err := s.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return user, nil
}

func (s *UsersStore) CreateAndInvite(ctx context.Context, user *User, token string, invitationExp time.Duration) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// Create user
		if err := s.Create(ctx, tx, user); err != nil {
			return err
		}

		// Create user invite
		if err := s.CreateUserInvitation(ctx, tx, token, invitationExp, user.Id); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStore) CreateUserInvitation(ctx context.Context, tx *sql.Tx, token string, exp time.Duration, userId int64) error {
	query := `
		INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1, $2, $3)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, token, userId, time.Now().Add(exp))
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) Activate(ctx context.Context, token string) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		// Find out the user that this token belongs to
		user, err := s.getUserFromInvitation(ctx, tx, token)
		if err != nil {
			return err
		}

		// Update the user
		user.IsActive = true
		if err := s.update(ctx, tx, user); err != nil {
			return err
		}

		// clean the invitations
		if err := s.deleteUserInvitations(ctx, tx, user.Id); err != nil {
			return err
		}
		return nil
	})
}

func (s *UsersStore) getUserFromInvitation(ctx context.Context, tx *sql.Tx, token string) (*User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.created_at, u.is_active
		FROM users AS u 
		JOIN user_invitations AS ui
		ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	user := &User{}
	err := tx.QueryRowContext(ctx, query, hashToken, time.Now()).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.IsActive,
	)
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return user, nil
}

func (s *UsersStore) update(ctx context.Context, tx *sql.Tx, user *User) error {
	query := `
		UPDATE users 
		SET username = $1, email = $2, is_active = $3 
		WHERE id = $4
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, user.Username, user.Email, user.IsActive, user.Id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) deleteUserInvitations(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM user_invitations 
		WHERE user_id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) deleteUser(ctx context.Context, tx *sql.Tx, userID int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	_, err := tx.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStore) Delete(ctx context.Context, userID int64) error {
	return withTx(s.db, ctx, func(tx *sql.Tx) error {
		if err := s.deleteUser(ctx, tx, userID); err != nil {
			return err
		}

		if err := s.deleteUserInvitations(ctx, tx, userID); err != nil {
			return err
		}

		return nil
	})
}
