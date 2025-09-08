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

func (r *Repository) SendRequest(request *dto.RequestFull) error {
	fmt.Println("Repo request:", request, "сдесь")
	query := sq.Insert("requests").Columns("source_id", "applicant_id", "address", "dogs_count", "behavior", "urgency", "contact_person").
		Values(request.Source.ID, request.Applicant.ID, request.Address, request.DogsCount, request.Behavior, request.Urgency, request.ContactPerson).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *Repository) GetRequests() ([]*dto.RequestFull, error) {
	query := sq.Select("requests.id, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, apllicants.name").
		From("requests").
		InnerJoin("request_sources ON requests.source_id = request_sources.id").
		InnerJoin("applicants ON requests.applicant_id = apllicants.id").
		OrderBy("created_at DESC")

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
