package repository

import (
	"context"
	"fmt"

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

func (r *Repository) ValidateAccount(account *dto.AuthCredentials) (user *dto.UserProfile, hashPass string, role_id string, err error) {
	fmt.Println("ТУТА")
	query := sq.Select("id", "password_hash", "role_id").
		From("users").
		Where(sq.Eq{"login": account.Login}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, "", "", err
	}
	var id string

	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&id, &hashPass, &role_id)
	if err != nil {
		return nil, "", "", err
	}
	user, err = r.GetUserProfile(id)

	if err != nil {
		return nil, "", "", err
	}

	user.ID = id
	fmt.Println(user, role_id)

	return user, hashPass, role_id, nil
}

func (r *Repository) SendRequest(request *dto.RequestFull) (int, error) {

	query := sq.Insert("requests").Columns("source_id", "applicant_id", "address", "dogs_count", "behavior", "urgency", "contact_person").
		Values(request.Source.ID, request.Applicant.ID, request.Address, request.DogsCount, request.Behavior, request.Urgency, request.ContactPerson).
		Suffix("RETURNING number").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	var number int

	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&number)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	selectQuery := sq.Select("applicants.name", "request_sources.name").
		From("requests").
		Join("applicants ON requests.applicant_id = applicants.id").
		Join("request_sources ON requests.source_id = request_sources.id").
		Where(sq.Eq{"requests.number": number}).
		PlaceholderFormat(sq.Dollar)

	sql, selectArgs, err := selectQuery.ToSql()
	if err != nil {
		return 0, err
	}
	err = r.pg.QueryRow(context.Background(), sql, selectArgs...).Scan(&request.Applicant.Name, &request.Source.Name)
	if err != nil {
		return 0, err
	}

	return number, nil
}

func (r *Repository) ChangeStatusRequest(req *dto.ChangeStatusRequest) error {

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

func (r *Repository) GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error) {
	var query sq.SelectBuilder
	if year == "" {
		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
			From("requests").
			InnerJoin("request_sources ON requests.source_id = request_sources.id").
			InnerJoin("applicants ON requests.applicant_id = applicants.id").
			Where(sq.Eq{"requests.source_id": otdel_id}).
			OrderBy("requests.created_at DESC").
			PlaceholderFormat(sq.Dollar)
	} else {
		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
			From("requests").
			InnerJoin("request_sources ON requests.source_id = request_sources.id").
			InnerJoin("applicants ON requests.applicant_id = applicants.id").
			Where(sq.Eq{"requests.source_id": otdel_id}).
			Where(sq.Eq{"requests.year": year}).
			OrderBy("requests.created_at DESC").
			PlaceholderFormat(sq.Dollar)
	}

	sql, args, err := query.ToSql()
	if err != nil {

		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Println("Error querying requests:", err)
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.RequestFull
	for rows.Next() {
		var req dto.RequestFull

		err := rows.Scan(
			&req.ID,
			&req.Number,
			&req.Source.ID, &req.Applicant.ID,
			&req.Address,
			&req.DogsCount,
			&req.Behavior, &req.Urgency, &req.ContactPerson, &req.Status, &req.CreatedAt,
			&req.Source.Name, &req.Applicant.Name,
		)
		if err != nil {
			return nil, err
		}
		fmt.Println(err)
		requests = append(requests, &req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil

}

func (r *Repository) GetAllRequests(year string) ([]*dto.RequestFull, error) {
	var query sq.SelectBuilder
	if year == "" {
		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
			From("requests").
			InnerJoin("request_sources ON requests.source_id = request_sources.id").
			InnerJoin("applicants ON requests.applicant_id = applicants.id").
			OrderBy("requests.created_at DESC").
			PlaceholderFormat(sq.Dollar)
	} else {
		fmt.Println("год передан")
		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
			From("requests").
			InnerJoin("request_sources ON requests.source_id = request_sources.id").
			InnerJoin("applicants ON requests.applicant_id = applicants.id").
			Where(sq.Eq{"requests.year": year}).
			OrderBy("requests.created_at DESC").
			PlaceholderFormat(sq.Dollar)
	}
	sql, args, err := query.ToSql()
	if err != nil {

		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Println("Error querying requests:", err)
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.RequestFull
	for rows.Next() {
		var req dto.RequestFull

		err := rows.Scan(
			&req.ID,
			&req.Number,
			&req.Source.ID, &req.Applicant.ID,
			&req.Address,
			&req.DogsCount,
			&req.Behavior, &req.Urgency, &req.ContactPerson, &req.Status, &req.CreatedAt,
			&req.Source.Name, &req.Applicant.Name,
		)
		if err != nil {
			return nil, err
		}
		fmt.Println(err)
		requests = append(requests, &req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) GetRequestsByNumber(numbers []string) ([]*dto.RequestFull, error) {
	query := sq.Select(
		"id", "number", "source_id", "applicant_id", "address", "dogs_count", "behavior", "urgency", "contact_person", "status", "created_at",
		"source_name", "applicant_name",
	).
		From("requests").
		Where(sq.Eq{"number": numbers}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.pg.Query(context.Background(), sql, args...)
	if err != nil {
		fmt.Println("Error querying requests:", err)
		return nil, err
	}
	defer rows.Close()

	var requests []*dto.RequestFull
	for rows.Next() {
		var req dto.RequestFull

		err := rows.Scan(
			&req.ID,
			&req.Number,
			&req.Source.ID, &req.Applicant.ID,
			&req.Address,
			&req.DogsCount,
			&req.Behavior, &req.Urgency, &req.ContactPerson, &req.Status, &req.CreatedAt,
			&req.Source.Name, &req.Applicant.Name,
		)
		if err != nil {
			return nil, err
		}
		fmt.Println(err)
		requests = append(requests, &req)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return requests, nil
}

func (r *Repository) ChangePassword(req *dto.ChangePasswordRequest) error {

	query := sq.Update("users").
		Set("password_hash", req.NewPassword).
		Where(sq.Eq{"login": req.Login}).
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
		return fmt.Errorf("no user found with login: %s", req.Login)
	}

	return nil

}

func (r *Repository) ChangeProfileInfo(req *dto.ChangeProfileRequest) error {
	query := sq.Update("users").
		Set("full_name", req.FullName).
		Set("login", req.Login).
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
		return fmt.Errorf("no user found with login: %s", req.Login)
	}

	return nil
}

func (r *Repository) GetUserProfile(userId string) (*dto.UserProfile, error) {
	query := sq.Select("full_name", "login", "role").
		From("users").
		Where(sq.Eq{"id": userId}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	var profile dto.UserProfile
	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&profile.FullName, &profile.Login, &profile.Role)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
