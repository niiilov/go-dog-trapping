package dto

type Applicant struct {
	ID         string `json:"id"`
	FullName   string `json:"full_name"`
	Position   string `json:"position"`
	DistrictID string `json:"district_id"`
}

type CreateApplicantDTO struct {
	FullName   string `json:"full_name"`
	Position   string `json:"position"`
	DistrictID string `json:"district_id"`
}
