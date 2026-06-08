package dto

type TerOtdel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateTerOtdelDTO struct {
	Name string `json:"name"`
}
