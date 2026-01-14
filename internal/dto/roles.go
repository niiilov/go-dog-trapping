package dto

const (
	RoleAdmin      = "admin"
	RoleContractor = "contractor"
	RoleUser       = "user"
)

// Проверка может ли роль видеть все заявки
func CanSeeAllRequests(role string) bool {
	return role == RoleAdmin || role == RoleContractor
}
