package dto

import "time"

const (
	AccesTimeExpr   = 1 * time.Hour
	RefreshTimeExpr = 60 * time.Hour
)

// Account model info
// @Description Структура для создания учетной записи пользователя. Все поля обязательны к заполнению
type Account struct {
	ID       string `json:"id"`                  // ID пользователя
	FullName string `json:"full_name"`           //Полное имя или название организации
	Login    string `json:"login"`               // Логин для входа в систему, должен быть уникальным
	Role     string `json:"role"`                // Роль пользователя, например, "admin", "user", "otdel" и т.д.
	RoleName string `json:"role_name"`           // Название роли пользователя
	Password string `json:"password"`            // Пароль для входа в систему
	SourceID string `json:"source_id,omitempty"` // id тер отдела
}

// AuthCredentials model info
// @Description Структура для авторизации пользователя. Все поля обязательны к заполнению
type AuthCredentials struct {
	Login    string `json:"login"`    // Логин для входа в систему, должен быть уникальным
	Password string `json:"password"` // Пароль для входа в систему
}

// ChangeProfileRequest model info
// @Description Структура для отправки запроса на изменение данных профиля
type ChangeProfileRequest struct {
	ID       string `json:"id"`        // ID пользователя
	FullName string `json:"full_name"` //Полное имя или название организации
	Login    string `json:"login"`     // Логин для входа в систему, должен быть уникальным
}

// ChangePasswordRequest model info
// @Description Структура для смены пароля пользователя. Все поля обязательны к заполнению
type ChangePasswordRequest struct {
	Login       string `json:"login"`        // Логин для входа в систему, должен быть уникальным
	OldPassword string `json:"old_password"` // Старый пароль для входа в систему
	NewPassword string `json:"new_password"` // Новый пароль для входа в систему
}

// AuthTokens modle info
// @Description Структура c токенами
type AuthTokens struct {
	AccesToken   string `json:"access_token"`
	RefreshToken string `json:"refhresh_token"`
}

type ChangeStatusRequest struct {
	ID     string `json:"id"`     // ID заявки
	Status string `json:"status"` // Новый статус заявки (новая, в работе, выполнена, отменена)
}

// RequestFull model info
// @Description Полная информация о заявке на отлов бродячей собаки. Поля id, number, status и created_at игнорируются при создании заявки.
type RequestFull struct {
	ID            string    `json:"id"`             //уникальный идентификатор заявки при Post запросе не указывается
	Number        string    `json:"number"`         //уникальный номер заявки, генерируется автоматически при создании заявки
	Source        Source    `json:"applicant"`      //"!!!!ВАЖНО название поля json applicant и source поменяты местами!!"Источник информации от кого был получен запрос ID находится в справочнике
	Applicant     Applicant `json:"source"`         //"!!!!ВАЖНО название поля json applicant и source поменяты местами!!"Заявитель ID находится в справочнике
	Address       string    `json:"address"`        //Адрес, где была замечена бродячая собака
	DogsCount     int       `json:"dogs_count"`     //Количество собак
	Behavior      string    `json:"behavior"`       //Поведение собаки
	Urgency       string    `json:"urgency"`        //срочность
	ContactPerson string    `json:"contact_person"` //Контактные данные
	Status        string    `json:"status"`         //Статус заявки (новая, в работе, выполнена, отменена) генерируется после создания заявки
	CreatedAt     time.Time `json:"created_at"`     //Дата и время создания заявки генерируется автоматически

}

// Source model  info
// @Description Источник информации от кого был получен запрос. Поле name игнорируется при создании заявки.
type Source struct {
	ID   string `json:"id"`             //уникальный идентификатор источника при Post запросе на создание заявки ОБЯЗАТЕЛЕН
	Name string `json:"name,omitempty"` //Название источника игнорируется при создании заявки
}

// Applicant model info
// @Description Заявитель. Поле name игнорируется при создании заявки.
// type Applicant struct {
// 	ID          string `json:"id"`             //уникальный идентификатор заявителя при Post запросе на создание заявки ОБЯЗАТЕЛЕН
// 	Name        string `json:"name,omitempty"` //Название заявителя игнорируется при создании заявки
// 	IsPermanent bool   `json:"is_permanent"`   // Является ли заявитель постоянным
// }

// RequestForGeneratingMultipleByDate model info
// @Description генерация документа с несколькими заявками.
type GenerateMultipleRequestByDate struct {
	DateFrom *time.Time `json:"date_from"` // Дата с которой нужно выбрать заявки
	DateTo   *time.Time `json:"date_to"`   // Дата по которую нужно выбрать заявки
}

// GenerateMultipleRequestByID model info
// @Description генерация документа с несколькими заявками по ID
type GenerateMultipleRequestByID struct {
	IDs []string `json:"ids"` // Список ID заявок которые нужно включить в документ
}

// Response model info
// @Description  Структура для ответа с где обычно ок и сообщние
type Response struct {
	Status  string `json:"status"`  // обычно ok  при ошибках error
	Message string `json:"message"` // просто сообщение

}

// Response model info
// @Description  Структура для ответа с данными пользователя
type AuthResponse struct {
	Status  string       `json:"status"`  // обычно ok  при ошибках error
	Message string       `json:"message"` // просто сообщение
	User    *UserProfile `json:"user"`    // данные пользователя
	Tokens  *AuthTokens  `json:"tokens"`  // токены авторизации

}

// ResponseUrl model info
// @Description  Структура для ответа с где обычно ок и url
type ResponseUrl struct {
	Status string `json:"status"` // обычно ok  при ошибках error
	Url    string `json:"url"`    // url на скаччивание файла
}

// UploadActRequests model info
// @Description  Структура для получения данных при загрузке акта
type UploadActRequests struct {
	Status string `json:"status"` // обычно ok  при ошибках error
	ID     string `json:"id"`     // ID заявки к которой прилагается акт
	Number string `json:"number"` // Номер заявки к которой прилагается акт

}

// DownloadActRequest model info
// @Description  Структура для получения данных при скачивании акта
type DownloadActRequest struct {
	Number string `json:"number"` // номер заявки
	Year   string `json:"year"`   // год заявки

}
