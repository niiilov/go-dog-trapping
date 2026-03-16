package dto

type User struct {
	ID           string `json:"id"`
	FullName     string `json:"full_name"`
	Login        string `json:"login"`
	PasswordHash string `json:"password_hash"`
	RoleID       string `json:"role_id"`
	DistrictID   string `json:"district_id"`
	TerOtdelID   string `json:"ter_otdel_id"`
}

type CreateUserDTO struct {
	FullName   string `json:"full_name"`
	Login      string `json:"login"`
	Password   string `json:"password"`
	RoleID     string `json:"role_id"`
	DistrictID string `json:"district_id"`
	TerOtdelID string `json:"ter_otdel_id,omitempty"`
}

type GetUserDTO struct {
	ID         string `json:"id"`
	FullName   string `json:"full_name"`
	Login      string `json:"login"`
	RoleID     string `json:"role_id"`
	DistrictID string `json:"district_id"`
	TerOtdelID string `json:"ter_otdel_id"`
}

type GetUserDTOResponse struct {
	Status string        `json:"status"`
	Users  []*GetUserDTO `json:"users"`
}

// UserProfile model info
// @Description Структура для получения профиля пользователя. Все поля обязательны к заполнению
type UserProfile struct {
	ID           string `json:"id"`                       // ID пользователя
	FullName     string `json:"full_name"`                //Полное имя или название организации
	Login        string `json:"login"`                    // Логин для входа в систему, должен быть уникальным
	RoleID       string `json:"roleID"`                   // Роль пользователя, например, "admin", "user", "otdel" и т.д.
	RoleName     string `json:"role_name"`                // Название роли пользователя
	DistrictID   string `json:"district_id"`              // ID района
	DistrictName string `json:"district_name"`            // Название района
	TerOtdelID   string `json:"ter_otdel_id,omitempty"`   // ID территориального отдела
	TerOtdelName string `json:"ter_otdel_name,omitempty"` // Название территориального отдела
}
