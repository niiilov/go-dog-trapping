package dto

import "time"

const (
	AccesTimeExpr   = 15 * time.Minute
	RefreshTimeExpr = 48 * time.Hour
)

//структура для создания и валидации аккаунта
type Account struct {
	FullName string `json:"full_name"`
	Login    string `json:"login"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

//структура для авторизации
type AuthCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	Login       string `json:"login"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

//Общая структура для заявки
type RequestFull struct {
	ID            string    `json:"id"`
	Source        Source    `json:"source"`
	Applicant     Applicant `json:"applicant,omitempty"`
	Address       string    `json:"address"`
	DogsCount     int       `json:"dogs_count"`
	Behavior      string    `json:"behavior"`
	Urgency       string    `json:"urgency"` //срочность
	ContactPerson string    `json:"contact_person"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type Applicant struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type RequestForGenerating struct {
	Number        string `json:"number"`
	Applicant     string `json:"applicant_name"`
	Source        string `json:"source_name"`
	Address       string `json:"address"`
	DogsCount     int    `json:"dogs_count"`
	Behavior      string `json:"behavior"`
	Urgency       string `json:"urgency"` //срочность
	ContactPerson string `json:"contact_person"`
}
