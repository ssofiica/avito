package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User interface {
	GetUserIDByUsername(ctx context.Context, name string) (string, error)
	IsResponsible(ctx context.Context, name string) (string, error)
}

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) User {
	return &UserRepo{db: db}
}

func (t *UserRepo) GetUserIDByUsername(ctx context.Context, name string) (string, error) {
	query := `select id from employee where username=$1`
	var id string
	err := t.db.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (t *UserRepo) IsResponsible(ctx context.Context, name string) (string, error) {
	query := `select organization_id from organization_responsible o JOIN employee e
	on o.user_id=e.id where username=$1`
	var id string
	err := t.db.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return id, nil
}
