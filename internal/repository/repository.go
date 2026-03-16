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
	query := sq.InsertBuilder{}
	if account.SourceID != "" {
		query = sq.Insert("users").
			Columns("full_name", "role", "login", "password_hash", "source_id", "role_name").
			Values(account.FullName, account.Role, account.Login, account.Password, account.SourceID, account.RoleName).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)
	} else {
		query = sq.Insert("users").
			Columns("full_name", "role", "login", "password_hash", "role_name").
			Values(account.FullName, account.Role, account.Login, account.Password, account.RoleName).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)
	}

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

// func (r *Repository) ValidateAccount(account *dto.AuthCredentials) (user *dto.UserProfile, hashPass string, role string, sourceID string, err error) {
// 	fmt.Println("ТУТА")
// 	query := sq.Select("id", "password_hash", "role", "source_id").
// 		From("users").
// 		Where(sq.Eq{"login": account.Login}).
// 		PlaceholderFormat(sq.Dollar)

// 	sql, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, "", "", "", err
// 	}
// 	var id string
// 	var source sd.NullString

// 	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&id, &hashPass, &role, &source)
// 	if err != nil {
// 		return nil, "", "", "", err
// 	}
// 	user, err = r.GetUserProfile(id)

// 	if err != nil {
// 		return nil, "", "", "", err
// 	}

// 	user.ID = id
// 	fmt.Println(user, role, source)
// 	if source.Valid {
// 		sourceID = source.String
// 	} else {
// 		sourceID = ""
// 	}

// 	return user, hashPass, role, sourceID, nil
// }

// func (r *Repository) SendRequest(request *dto.RequestFull) (int, error) {

// 	query := sq.Insert("requests").Columns("source_id", "applicant_id", "address", "dogs_count", "behavior", "urgency", "contact_person").
// 		Values(request.Source.ID, request.Applicant.ID, request.Address, request.DogsCount, request.Behavior, request.Urgency, request.ContactPerson).
// 		Suffix("RETURNING number").
// 		PlaceholderFormat(sq.Dollar)

// 	sql, args, err := query.ToSql()
// 	if err != nil {
// 		fmt.Println(err)
// 		return 0, err
// 	}
// 	fmt.Println("ЗАПРОС", sql)

// 	var number int

// 	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&number)
// 	if err != nil {
// 		fmt.Println(err)
// 		return 0, err
// 	}

// 	selectQuery := sq.Select("applicants.name", "request_sources.name").
// 		From("requests").
// 		Join("applicants ON requests.applicant_id = applicants.id").
// 		Join("request_sources ON requests.source_id = request_sources.id").
// 		Where(sq.Eq{"requests.number": number}).
// 		PlaceholderFormat(sq.Dollar)

// 	sql, selectArgs, err := selectQuery.ToSql()
// 	if err != nil {
// 		return 0, err
// 	}
// 	err = r.pg.QueryRow(context.Background(), sql, selectArgs...).Scan(&request.Applicant.Name, &request.Source.Name)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return number, nil
// }

// func (r *Repository) ChangeStatusRequest(req *dto.ChangeStatusRequest) error {

// 	query := sq.Update("requests").
// 		Set("status", req.Status).
// 		Where(sq.Eq{"id": req.ID}).
// 		PlaceholderFormat(sq.Dollar)
// 	sql, args, err := query.ToSql()
// 	if err != nil {
// 		return err
// 	}

// 	cmdTag, err := r.pg.Exec(context.Background(), sql, args...)
// 	if err != nil {
// 		return err
// 	}

// 	if cmdTag.RowsAffected() == 0 {
// 		return fmt.Errorf("no request found with id: %s", req.ID)
// 	}

// 	return nil
// }

// func (r *Repository) GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error) {
// 	var query sq.SelectBuilder
// 	if year == "" {
// 		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
// 			From("requests").
// 			InnerJoin("request_sources ON requests.source_id = request_sources.id").
// 			InnerJoin("applicants ON requests.applicant_id = applicants.id").
// 			Where(sq.Eq{"requests.source_id": otdel_id}).
// 			OrderBy("requests.created_at DESC").
// 			PlaceholderFormat(sq.Dollar)
// 	} else {
// 		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
// 			From("requests").
// 			InnerJoin("request_sources ON requests.source_id = request_sources.id").
// 			InnerJoin("applicants ON requests.applicant_id = applicants.id").
// 			Where(sq.Eq{"requests.source_id": otdel_id}).
// 			Where(sq.Eq{"requests.year": year}).
// 			OrderBy("requests.created_at DESC").
// 			PlaceholderFormat(sq.Dollar)
// 	}

// 	sql, args, err := query.ToSql()
// 	if err != nil {

// 		return nil, err
// 	}

// 	rows, err := r.pg.Query(context.Background(), sql, args...)
// 	if err != nil {
// 		fmt.Println("Error querying requests:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var requests []*dto.RequestFull
// 	for rows.Next() {
// 		var req dto.RequestFull

// 		err := rows.Scan(
// 			&req.ID,
// 			&req.Number,
// 			&req.Source.ID, &req.Applicant.ID,
// 			&req.Address,
// 			&req.DogsCount,
// 			&req.Behavior, &req.Urgency, &req.ContactPerson, &req.Status, &req.CreatedAt,
// 			&req.Source.Name, &req.Applicant.Name,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(err)
// 		requests = append(requests, &req)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return requests, nil

// }

// func (r *Repository) GetAllRequests(year string) ([]*dto.RequestFull, error) {
// 	var query sq.SelectBuilder
// 	if year == "" {
// 		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
// 			From("requests").
// 			InnerJoin("request_sources ON requests.source_id = request_sources.id").
// 			InnerJoin("applicants ON requests.applicant_id = applicants.id").
// 			OrderBy("requests.created_at DESC").
// 			PlaceholderFormat(sq.Dollar)
// 	} else {
// 		fmt.Println("год передан")
// 		query = sq.Select("requests.id, requests.number, requests.source_id, requests.applicant_id, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, requests.status, requests.created_at, request_sources.name, applicants.name").
// 			From("requests").
// 			InnerJoin("request_sources ON requests.source_id = request_sources.id").
// 			InnerJoin("applicants ON requests.applicant_id = applicants.id").
// 			Where(sq.Eq{"requests.year": year}).
// 			OrderBy("requests.created_at DESC").
// 			PlaceholderFormat(sq.Dollar)
// 	}
// 	sql, args, err := query.ToSql()
// 	if err != nil {

// 		return nil, err
// 	}

// 	rows, err := r.pg.Query(context.Background(), sql, args...)
// 	if err != nil {
// 		fmt.Println("Error querying requests:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var requests []*dto.RequestFull
// 	for rows.Next() {
// 		var req dto.RequestFull

// 		err := rows.Scan(
// 			&req.ID,
// 			&req.Number,
// 			&req.Source.ID, &req.Applicant.ID,
// 			&req.Address,
// 			&req.DogsCount,
// 			&req.Behavior, &req.Urgency, &req.ContactPerson, &req.Status, &req.CreatedAt,
// 			&req.Source.Name, &req.Applicant.Name,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(err)
// 		requests = append(requests, &req)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	return requests, nil
// }

// func (r *Repository) GetRequestsByDate(dateFrom *time.Time, dateTo *time.Time) ([]*dto.RequestForGenerating, error) {
// 	fmt.Println("РЕПОЗИТОРИЙ ТУТАЭ", dateFrom)

// 	query := sq.Select("requests.number, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, request_sources.name, applicants.name").
// 		From("requests").
// 		InnerJoin("request_sources ON requests.source_id = request_sources.id").
// 		InnerJoin("applicants ON requests.applicant_id = applicants.id").
// 		OrderBy("requests.created_at DESC").
// 		Where(sq.And{
// 			sq.GtOrEq{"requests.created_at": dateFrom},
// 			sq.LtOrEq{"requests.created_at": dateTo},
// 		}).
// 		PlaceholderFormat(sq.Dollar)
// 	sql, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows, err := r.pg.Query(context.Background(), sql, args...)
// 	if err != nil {
// 		fmt.Println("Error querying requests:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var requests []*dto.RequestForGenerating
// 	for rows.Next() {
// 		var req dto.RequestForGenerating

// 		err := rows.Scan(

// 			&req.Number,

// 			&req.Address,
// 			&req.DogsCount,
// 			&req.Behavior, &req.Urgency, &req.ContactPerson,
// 			&req.Source, &req.Applicant,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(err)
// 		requests = append(requests, &req)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(requests)

// 	return requests, nil
// }

// func (r *Repository) GetRequestsByIDs(reqIDs []string) ([]*dto.RequestForGenerating, error) {
// 	fmt.Println("РЕПОЗИТОРИЙ ТУТАЭ", reqIDs)

// 	query := sq.Select("requests.number, requests.address, requests.dogs_count, requests.behavior, requests.urgency, requests.contact_person, request_sources.name, applicants.name").
// 		From("requests").
// 		InnerJoin("request_sources ON requests.source_id = request_sources.id").
// 		InnerJoin("applicants ON requests.applicant_id = applicants.id").
// 		OrderBy("requests.created_at DESC").
// 		Where(sq.Eq{"requests.id": reqIDs}).
// 		PlaceholderFormat(sq.Dollar)
// 	sql, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, err
// 	}

// 	rows, err := r.pg.Query(context.Background(), sql, args...)
// 	if err != nil {
// 		fmt.Println("Error querying requests:", err)
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var requests []*dto.RequestForGenerating
// 	for rows.Next() {
// 		var req dto.RequestForGenerating

// 		err := rows.Scan(

// 			&req.Number,

// 			&req.Address,
// 			&req.DogsCount,
// 			&req.Behavior, &req.Urgency, &req.ContactPerson,
// 			&req.Source, &req.Applicant,
// 		)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Println(err)
// 		requests = append(requests, &req)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}

// 	fmt.Println(requests)

// 	return requests, nil
// }

//	func (r *Repository) DeleteRequest(reqID string) error {
//		query := sq.Delete("requests").
//			Where(sq.Eq{"id": reqID}).
//			PlaceholderFormat(sq.Dollar)
//		sql, args, err := query.ToSql()
//		if err != nil {
//			return err
//		}
//		cmdTag, err := r.pg.Exec(context.Background(), sql, args...)
//		if err != nil {
//			return err
//		}
//		if cmdTag.RowsAffected() == 0 {
//			return fmt.Errorf("no request found with id: %s", reqID)
//		}
//		return nil
//	}
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

// func (r *Repository) GetUserProfile(userId string) (*dto.UserProfile, error) {
// 	query := sq.Select("full_name", "login", "role", "role_name").
// 		From("users").
// 		Where(sq.Eq{"id": userId}).
// 		PlaceholderFormat(sq.Dollar)
// 	sql, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var profile dto.UserProfile
// 	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&profile.FullName, &profile.Login, &profile.Role, &profile.RoleName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &profile, nil
// }

func (r *Repository) NewNotPemamentApplicant(name string) (string, error) {

	query := sq.Insert("applicants").
		Columns("name", "is_permanent").
		Values(name, false).
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

func (r *Repository) CreateExternalUser(user *dto.ExternalUser) (int, error) {
	query := sq.Insert("external_users").
		Columns(
			"username",
			"email",
			"first_name",
			"last_name",
			"patronymic",
			"phone_number",
			"city",
			"date_of_birth",
			"bio",
			"password",
			"is_accepted",
		).
		Values(
			user.Username,
			user.Email,
			user.FirstName,
			user.LastName,
			user.Patronymic,
			user.PhoneNumber,
			user.City,
			user.DateOfBirth,
			user.Bio,
			user.Password,
			user.IsAccepted,
		).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, err
	}
	var id int
	err = r.pg.QueryRow(context.Background(), sql, args...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *Repository) GetExternalUserByID(id int) (*dto.ExternalUser, error) {
	query := sq.Select(
		"id",
		"username",
		"email",
		"first_name",
		"last_name",
		"patronymic",
		"phone_number",
		"city",
		"date_of_birth",
		"bio",
		"password",
		"is_accepted",
	).From("external_users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	row := r.pg.QueryRow(context.Background(), sql, args...)
	var user dto.ExternalUser
	if err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Patronymic,
		&user.PhoneNumber,
		&user.City,
		&user.DateOfBirth,
		&user.Bio,
		&user.Password,
		&user.IsAccepted,
	); err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *Repository) GetExternalUsers() ([]*dto.ExternalUser, error) {
	query := sq.Select(
		"id",
		"username",
		"email",
		"first_name",
		"last_name",
		"patronymic",
		"phone_number",
		"city",
		"date_of_birth",
		"bio",
		"password",
		"is_accepted",
	).From("external_users").
		Where(sq.Eq{"is_accepted": true}).
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
	var users []*dto.ExternalUser
	for rows.Next() {
		var user dto.ExternalUser
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Patronymic,
			&user.PhoneNumber,
			&user.City,
			&user.DateOfBirth,
			&user.Bio,
			&user.Password,
			&user.IsAccepted,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}
func (r *Repository) ChangeAcceptedStatusExternalUser(id int, isAccepted bool) error {
	query := sq.Update("external_users").
		Set("is_accepted", isAccepted).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

// Добавление в систему собак нового пользователя
func (r *Repository) AcceptExternalUser(user *dto.AcceptExternalUserRequest) error {

	userDb, err := r.GetExternalUserByID(user.ID)
	if err != nil {
		return err
	}

	account := &dto.Account{
		FullName: user.FullName,
		Login:    userDb.Username,
		Password: userDb.Password,
		Role:     user.Role,
		RoleName: user.RoleName,
		SourceID: user.SourceID,
	}
	_, err = r.CreateAccount(account)
	if err != nil {
		return err
	}
	return nil

}
func (r *Repository) NotAcceptExternalUser(id int) error {

	r.ChangeAcceptedStatusExternalUser(id, false)
	query := sq.Delete("users").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}

// func (r *Repository) GetApplicants() ([]*dto.Applicant, error) {
// 	query := sq.Select("name", "is_permanent").
// 		From("applicants").
// 		PlaceholderFormat(sq.Dollar)
// 	sql, args, err := query.ToSql()
// 	if err != nil {
// 		return nil, err
// 	}
// 	rows, err := r.pg.Query(context.Background(), sql, args...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var applicants []*dto.Applicant
// 	for rows.Next() {
// 		var applicant dto.Applicant
// 		err := rows.Scan(&applicant.Name, &applicant.IsPermanent)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if applicant.IsPermanent {
// 			applicants = append(applicants, &applicant)
// 		}
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return applicants, nil
// }

func (r *Repository) GetTerrOtdels() ([]*dto.Source, error) {
	query := sq.Select("id", "name").
		From("request_sources").
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

	var sources []*dto.Source
	for rows.Next() {
		var source dto.Source
		err := rows.Scan(&source.ID, &source.Name)
		if err != nil {
			return nil, err
		}
		sources = append(sources, &source)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sources, nil
}

func (r *Repository) AddNewTerOtdel(terOtdel *dto.AddNewTerOtdel) error {
	query := sq.Insert("request_sources").
		Columns("name").
		Values(terOtdel.Name).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	if err != nil {
		return err
	}

	return nil
}
