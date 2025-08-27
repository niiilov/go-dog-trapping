package dto

import "time"

const (
	AccesTimeExpr   = 15 * time.Minute
	RefreshTimeExpr = 48 * time.Hour
)

type Account struct {
	FullName string `json:"full_name"`
	Login    string `json:"login"`
	Role     string `json:"role"`
	Password string `json:"password"`
}
