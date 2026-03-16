package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func (r *Repository) GetRoles() ([]*dto.Role, error) {
	query := sq.Select("id", "name", "type_role").
		From("roles").PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*dto.Role
	for rows.Next() {
		var role dto.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Type); err != nil {
			return nil, err
		}
		roles = append(roles, &role)
	}

	return roles, nil
}
