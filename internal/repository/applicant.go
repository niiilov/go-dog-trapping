package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func (r *Repository) CreateApplicant(applicant *dto.CreateApplicantDTO) error {
	query := sq.Insert("applicants").Columns("full_name", "position", "ter_otdel_id").
		Values(applicant.FullName, applicant.Position, applicant.TerOtdelID).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err
}

func (r *Repository) GetApplicantByTerOtdelID(terOtdelID string) ([]*dto.Applicant, error) {
	query := sq.Select("id", "full_name", "position", "ter_otdel_id").
		From("applicants").
		Where(sq.Eq{"ter_otdel_id": terOtdelID}).
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

	var applicants []*dto.Applicant
	for rows.Next() {
		var a dto.Applicant
		if err := rows.Scan(&a.ID, &a.FullName, &a.Position, &a.TerOtdelID); err != nil {
			return nil, err
		}
		applicants = append(applicants, &a)
	}

	return applicants, nil
}
