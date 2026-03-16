package dto

type TerOtdel struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	DistrictID string `json:"district_id"`
}

type CreateTerOtdelDTO struct {
	Name       string `json:"name"`
	DistrictID string `json:"district_id"`
}
