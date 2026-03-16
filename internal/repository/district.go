package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func (r *Repository) CreateDistrict(district *dto.District) error {
	query := sq.Insert("districts").
		Columns("name").
		Values(district.Name).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err

}

func (r *Repository) GetDistricts() ([]*dto.District, error) {
	query := sq.Select("id", "name").From("districts")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var districts []*dto.District
	for rows.Next() {
		var district dto.District
		if err := rows.Scan(&district.ID, &district.Name); err != nil {
			return nil, err
		}
		districts = append(districts, &district)
	}

	return districts, nil
}

func (r *Repository) DeleteDistrict(id string) error {
	query := sq.Delete("districts").Where(sq.Eq{"id": id})
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err
}
