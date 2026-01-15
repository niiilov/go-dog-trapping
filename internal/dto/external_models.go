package dto

import "time"

// ExternalUser model info
// @Description Структура для создания внешнего пользователя. Все поля обязательны к заполнению кроме Patronymic, PhoneNumber, City, DateOfBirth, Bio
type ExternalUser struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Patronymic  string     `json:"patronymic,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	City        string     `json:"city,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Bio         string     `json:"bio,omitempty"`
	Password    string     `json:"password"`
	IsAccepted  bool       `json:"is_accepted"`
}

//@Description Структура для пполучения внешних пользователей
type GetExternalUser struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Patronymic  string     `json:"patronymic,omitempty"`
	PhoneNumber string     `json:"phone_number,omitempty"`
	City        string     `json:"city,omitempty"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Bio         string     `json:"bio,omitempty"`
	IsAccepted  bool       `json:"is_accepted"`
}

type ChangeAcceptedExternalUserRequest struct {
	IsAccepted bool `json:"is_accepted"`
}

// AcceptExternalUserRequest model info
// @Description Структура для принятия внешнего пользователя. Все поля обязательны к заполнению ID берётся из param. Для роли admin, contractor необязательно заполнять SourceID
type AcceptExternalUserRequest struct {
	ID       string `json:"id,omitempty"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
	RoleName string `json:"role_name"`
	SourceID string `json:"source_id,omitempty"`
}

type AddNewTerOtdel struct {
	Name string `json:"name"`
}
