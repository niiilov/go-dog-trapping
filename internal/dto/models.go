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
	Name string `json:"name"`
}

type Applicant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
