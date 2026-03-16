package repository

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

func (r *Repository) CreateUser(user *dto.CreateUserDTO, passwordHash string) error {

	ter := sql.NullString{}
	if user.TerOtdelID != "" {
		ter = sql.NullString{String: user.TerOtdelID, Valid: true}
	}

	dis := sql.NullString{}
	if user.DistrictID != "" {
		dis = sql.NullString{String: user.DistrictID, Valid: true}
	}
	query := sq.Insert("users").
		Columns("full_name", "login", "password_hash", "role_id", "district_id", "ter_otdel_id").
		Values(user.FullName, user.Login, passwordHash, user.RoleID, dis, ter).
		PlaceholderFormat(sq.Dollar)
	sql1, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql1, args...)
	return err
}
func (r *Repository) GetUsers() ([]*dto.User, error) {
	query := sq.Select("id", "full_name", "login", "password_hash", "role_id", "district_id", "ter_otdel_id").
		From("users").
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

	var users []*dto.User
	for rows.Next() {
		var dis sql.NullString
		var ter sql.NullString
		var user dto.User
		if err := rows.Scan(&user.ID, &user.FullName, &user.Login, &user.PasswordHash, &user.RoleID, &dis, &ter); err != nil {
			return nil, err
		}
		user.TerOtdelID = ter.String
		user.DistrictID = dis.String
		users = append(users, &user)
	}
	return users, nil
}

func (r *Repository) GetUserByLogin(login string) (*dto.User, error) {
	query := sq.Select("id", "full_name", "login", "password_hash", "role_id", "district_id", "ter_otdel_id").
		From("users").
		Where(sq.Eq{"login": login}).
		PlaceholderFormat(sq.Dollar)
	sql1, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pg.QueryRow(context.Background(), sql1, args...)
	var ter sql.NullString
	var dis sql.NullString
	var user dto.User
	if err := row.Scan(&user.ID, &user.FullName, &user.Login, &user.PasswordHash, &user.RoleID, &dis, &ter); err != nil {
		return nil, err
	}
	user.TerOtdelID = ter.String
	user.DistrictID = dis.String
	return &user, nil
}

func (r *Repository) GetUserProfile(userID string) (*dto.UserProfile, error) {
	query := sq.Select("u.id", "u.full_name", "u.login", "u.role_id", "u.district_id", "u.ter_otdel_id", "d.name", "t.name", "r.name").
		From("users u").
		LeftJoin("districts d ON u.district_id = d.id").
		LeftJoin("ter_otdels t ON u.ter_otdel_id = t.id").
		LeftJoin("roles r ON u.role_id = r.id").
		Where(sq.Eq{"u.id": userID}).
		PlaceholderFormat(sq.Dollar)
	sql1, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.pg.QueryRow(context.Background(), sql1, args...)
	var terID, terName sql.NullString
	var disID, disName sql.NullString
	var user dto.UserProfile
	if err := row.Scan(&user.ID, &user.FullName, &user.Login, &user.RoleID, &disID, &terID, &disName, &terName, &user.RoleName); err != nil {
		return nil, err
	}
	user.TerOtdelID, user.TerOtdelName = terID.String, terName.String
	user.DistrictID, user.DistrictName = disID.String, disName.String
	return &user, nil
}

func (r *Repository) DeleteUser(userID string) error {
	query := sq.Delete("users").
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.pg.Exec(context.Background(), sql, args...)
	return err
}
