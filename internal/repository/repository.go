package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

//тут бд

type Repository struct {
	pg *pgxpool.Pool
}

func New(pg *pgxpool.Pool) *Repository {
	return &Repository{
		pg: pg,
	}
}

func (r *Repository) CreateAccount(account *dto.Account) (string, error) {
	query := sq.Insert("users").
		Columns("full_name", "role", "login", "password_hash").
		Values(account.FullName, account.Role, account.Login, account.Password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	var id string
	sql, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) ValidateAccount(account *dto.Account) (id string, hashPass string, err error) {
	query := sq.Select("id", "password_hash").
		From("users").
		Where(sq.Eq{"login": account.Login}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", "", err
	}

	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&id, &hashPass)
	if err != nil {
		return "", "", err
	}

	return id, hashPass, nil
}
