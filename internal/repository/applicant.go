package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func (r *Repository) CreateApplicant(applicant *dto.CreateApplicantDTO) error {
	// Implementation for creating an applicant in the database

	query := sq.Insert("applicants").Columns("full_name", "position", "district_id").
		Values(applicant.FullName, applicant.Position, applicant.DistrictID).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err
}

func (r *Repository) GetApplicantByDistrictID(districtID string) ([]*dto.Applicant, error) {
	query := sq.Select("id", "full_name", "position", "district_id").
		From("applicants").
		Where(sq.Eq{"district_id": districtID}).
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

	var applicant []*dto.Applicant
	for rows.Next() {
		var a dto.Applicant
		if err := rows.Scan(&a.ID, &a.FullName, &a.Position, &a.DistrictID); err != nil {
			return nil, err
		}
		applicant = append(applicant, &a)
	}

	return applicant, nil
}
