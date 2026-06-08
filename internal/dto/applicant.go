package dto

type Applicant struct {
	ID         string `json:"id"`
	FullName   string `json:"full_name"`
	Position   string `json:"position"`
	TerOtdelID string `json:"ter_otdel_id"`
}

type CreateApplicantDTO struct {
	FullName   string `json:"full_name"`
	Position   string `json:"position"`
	TerOtdelID string `json:"ter_otdel_id"`
}
