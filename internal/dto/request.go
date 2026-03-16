package dto

import "time"

type Request struct {
	ID            string    `json:"id"`
	DistrictID    string    `json:"district_id"`
	TerOtdelID    string    `json:"ter_otdel_id"`
	ApplicantID   string    `json:"applicant_id"`
	Address       string    `json:"address"`
	DogsCount     int       `json:"dogs_count"`
	Behavior      string    `json:"behavior"`
	Urgency       string    `json:"urgency"`
	ContactPerson string    `json:"contact_person"`
	Status        string    `json:"status"`
	Number        int       `json:"number"`
	ActFile       string    `json:"act_file,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}
type GetRequestsDTO struct {
	ID                string    `json:"id"`
	DistrictID        string    `json:"district_id"`
	DistrictName      string    `json:"district_name"`
	TerOtdelID        string    `json:"ter_otdel_id"`
	TerOtdelName      string    `json:"ter_otdel_name"`
	ApplicantID       string    `json:"applicant_id"`
	ApplicantName     string    `json:"applicant_name"`
	ApplicantPosition string    `json:"applicant_position"`
	Address           string    `json:"address"`
	DogsCount         int       `json:"dogs_count"`
	Behavior          string    `json:"behavior"`
	Urgency           string    `json:"urgency"`
	ContactPerson     string    `json:"contact_person"`
	Status            string    `json:"status"`
	Number            int       `json:"number"`
	ActFile           string    `json:"act_file,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
}

type CreateRequestDTO struct {
	DistrictID    string `json:"district_id"`
	TerOtdelID    string `json:"ter_otdel_id"`
	ApplicantID   string `json:"applicant_id"`
	Address       string `json:"address"`
	DogsCount     int    `json:"dogs_count"`
	Behavior      string `json:"behavior"`
	Urgency       string `json:"urgency"`
	ContactPerson string `json:"contact_person"`
	Number        int    `json:"number"`
}

type ChangeStatusRequestDTO struct {
	ID     string `json:"id"`     // ID заявки
	Status string `json:"status"` // Новый статус заявки (новая, в работе, выполнена, отменена)
}
