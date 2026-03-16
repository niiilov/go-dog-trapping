package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func (r *Repository) CreateRequest(request *dto.CreateRequestDTO) error {
	query := sq.Insert("requests").
		Columns("district_id", "ter_otdel_id", "applicant_id", "address", "dogs_count", "behavior", "urgency", "contact_person", "number").
		Values(request.DistrictID, request.TerOtdelID, request.ApplicantID, request.Address, request.DogsCount, request.Behavior, request.Urgency, request.ContactPerson, request.Number).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err
}

// по тер отделу
func (r *Repository) GetRequestsByTerOtdel(id string) ([]*dto.GetRequestsDTO, error) {

	query := sq.Select("r.id", "r.district_id", "r.ter_otdel_id", "r.applicant_id", "r.address", "r.dogs_count", "r.behavior", "r.urgency", "r.contact_person", "r.number", "r.act_file", "r.status", "d.name", "t.name", "a.full_name, a.position").
		From("requests r").
		LeftJoin("districts d ON r.district_id = d.id").
		LeftJoin("ter_otdels t ON r.ter_otdel_id = t.id").
		LeftJoin("applicants a ON r.applicant_id = a.id").
		Where(sq.Eq{"r.ter_otdel_id": id}).
		PlaceholderFormat(sq.Dollar)

	sql1, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql1, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.GetRequestsDTO
	for rows.Next() {
		var actFile sql.NullString
		var request dto.GetRequestsDTO
		if err := rows.Scan(&request.ID, &request.DistrictID, &request.TerOtdelID, &request.ApplicantID, &request.Address, &request.DogsCount, &request.Behavior, &request.Urgency, &request.ContactPerson, &request.Number, &actFile, &request.Status, &request.DistrictName, &request.TerOtdelName, &request.ApplicantName, &request.ApplicantPosition); err != nil {
			return nil, err
		}
		request.ActFile = actFile.String
		requests = append(requests, &request)
	}

	return requests, nil
}

// ПО Муниципальному округу
func (r *Repository) GetRequestsByDistrictID(id string) ([]*dto.GetRequestsDTO, error) {
	query := sq.Select("r.id", "r.district_id", "r.ter_otdel_id", "r.applicant_id", "r.address", "r.dogs_count", "r.behavior", "r.urgency", "r.contact_person", "r.number", "r.act_file", "r.status", "d.name", "t.name", "a.full_name, a.position").
		From("requests r").
		LeftJoin("districts d ON r.district_id = d.id").
		LeftJoin("ter_otdels t ON r.ter_otdel_id = t.id").
		LeftJoin("applicants a ON r.applicant_id = a.id").
		Where(sq.Eq{"r.district_id": id}).
		PlaceholderFormat(sq.Dollar)

	sql1, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql1, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.GetRequestsDTO
	for rows.Next() {
		var actFile sql.NullString
		var request dto.GetRequestsDTO
		if err := rows.Scan(&request.ID, &request.DistrictID, &request.TerOtdelID, &request.ApplicantID, &request.Address, &request.DogsCount, &request.Behavior, &request.Urgency, &request.ContactPerson, &request.Number, &actFile, &request.Status, &request.DistrictName, &request.TerOtdelName, &request.ApplicantName, &request.ApplicantPosition); err != nil {
			return nil, err
		}
		request.ActFile = actFile.String
		requests = append(requests, &request)
	}

	return requests, nil
}

func (r *Repository) GetRequestsByDistrictIDs(district_id string, ids []string) ([]*dto.GetRequestsDTO, error) {
	query := sq.Select("r.id", "r.district_id", "r.ter_otdel_id", "r.applicant_id", "r.address", "r.dogs_count", "r.behavior", "r.urgency", "r.contact_person", "r.number", "r.act_file", "r.status", "d.name", "t.name", "a.full_name, a.position").
		From("requests r").
		LeftJoin("districts d ON r.district_id = d.id").
		LeftJoin("ter_otdels t ON r.ter_otdel_id = t.id").
		LeftJoin("applicants a ON r.applicant_id = a.id").
		Where(sq.Eq{"r.district_id": district_id}, sq.Eq{"r.id": ids}).
		PlaceholderFormat(sq.Dollar)

	sql1, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql1, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.GetRequestsDTO
	for rows.Next() {
		var actFile sql.NullString
		var request dto.GetRequestsDTO
		if err := rows.Scan(&request.ID, &request.DistrictID, &request.TerOtdelID, &request.ApplicantID, &request.Address, &request.DogsCount, &request.Behavior, &request.Urgency, &request.ContactPerson, &request.Number, &actFile, &request.Status, &request.DistrictName, &request.TerOtdelName, &request.ApplicantName, &request.ApplicantPosition); err != nil {
			return nil, err
		}
		request.ActFile = actFile.String
		requests = append(requests, &request)
	}

	return requests, nil
}
func (r *Repository) ChangeStatusRequest(req *dto.ChangeStatusRequestDTO) error {

	query := sq.Update("requests").
		Set("status", req.Status).
		Where(sq.Eq{"id": req.ID}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	cmdTag, err := r.pg.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no request found with id: %s", req.ID)
	}

	return nil
}
func (r *Repository) AddActFile(reqID string, actFile string) error {
	query := sq.Update("requests").
		Set("act_file", actFile).
		Where(sq.Eq{"id": reqID}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	cmdTag, err := r.pg.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no request found with id: %s", reqID)
	}

	return nil
}

func (r *Repository) DeleteRequest(id string) error {
	query := sq.Delete("requests").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err
}
