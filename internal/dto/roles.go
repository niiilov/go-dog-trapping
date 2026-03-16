package dto

const (
	RoleTerOtdel      = "e57ef349-176a-4b06-9116-fb12c0e21f58"
	RoleDistrictAdmin = "c8b92434-db65-42fc-97aa-b516848ee1e9"
	RoleRegionalAdmin = "f2fdedcc-6957-4674-999b-d613531226bb"
)

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
