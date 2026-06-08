package dto

type User struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	RoleID       string `json:"role_id"`
	TerOtdelID   string `json:"ter_otdel_id"`
}

type CreateUserDTO struct {
	FullName   string `json:"full_name"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	RoleID     string `json:"role_id"`
	TerOtdelID string `json:"ter_otdel_id"`
}

type GetUserDTO struct {
	ID         string `json:"id"`
	FullName   string `json:"full_name"`
	Login      string `json:"login"`
	RoleID     string `json:"role_id"`
	TerOtdelID string `json:"ter_otdel_id"`
}

type GetUserDTOResponse struct {
	Status string        `json:"status"`
	Users  []*GetUserDTO `json:"users"`
}

type UserProfile struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	Login        string `json:"login"`
	RoleID       string `json:"roleID"`
	RoleName     string `json:"role_name"`
	TerOtdelID   string `json:"ter_otdel_id"`
	TerOtdelName string `json:"ter_otdel_name"`
}
