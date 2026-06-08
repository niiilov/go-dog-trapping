package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func (r *Repository) CreateTerOtdel(terOtdel *dto.CreateTerOtdelDTO) (string, error) {
	query := sq.Insert("ter_otdels").
		Columns("name").
		Values(terOtdel.Name).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return "", err
	}

	var id string
	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repository) GetTerOtdels() ([]*dto.TerOtdel, error) {
	query := sq.Select("id", "name").
		From("ter_otdels").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var terOtdels []*dto.TerOtdel
	for rows.Next() {
		var terOtdel dto.TerOtdel
		if err := rows.Scan(&terOtdel.ID, &terOtdel.Name); err != nil {
			return nil, err
		}
		terOtdels = append(terOtdels, &terOtdel)
	}

	return terOtdels, nil
}

func (r *Repository) DeleteTerOtdel(id string) error {
	query := sq.Delete("ter_otdels").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err
}
