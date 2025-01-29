package store

import (
	"context"
	"database/sql"
	"errors"
)

type Role struct {
	Id          int64
	Name        string
	Level       int64
	Description string
}

type RolesStore struct {
	db *sql.DB
}

func (s *RolesStore) GetByName(ctx context.Context, name string) (*Role, error) {
	query := `
		SELECT id, name, level, description
		FROM roles WHERE name = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutduration)
	defer cancel()

	role := &Role{}

	// Scan objedt must follow the order in which sql query is being executed
	err := s.db.QueryRowContext(
		ctx,
		query,
		name,
	).Scan(
		&role.Id,
		&role.Name,
		&role.Level,
		&role.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}

	return role, nil
}
