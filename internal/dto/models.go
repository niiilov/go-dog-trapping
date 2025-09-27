package dto

import "time"

const (
	AccesTimeExpr   = 15 * time.Minute
	RefreshTimeExpr = 48 * time.Hour
)

// Account model info
// @Description Структура для создания учетной записи пользователя. Все поля обязательны к заполнению
type Account struct {
	ID       string `json:"id"`        // ID пользователя
	FullName string `json:"full_name"` //Полное имя или название организации
	Login    string `json:"login"`     // Логин для входа в систему, должен быть уникальным
	Role     string `json:"role"`      // Роль пользователя, например, "admin", "user", "otdel" и т.д.
	Password string `json:"password"`  // Пароль для входа в систему
}

// UserProfile model info
// @Description Структура для получения профиля пользователя. Все поля обязательны к заполнению
type UserProfile struct {
	ID       string `json:"id"`        // ID пользователя
	FullName string `json:"full_name"` //Полное имя или название организации
	Login    string `json:"login"`     // Логин для входа в систему, должен быть уникальным
	Role     string `json:"role"`      // Роль пользователя, например, "admin", "user", "otdel" и т.д.
}

// AuthCredentials model info
// @Description Структура для авторизации пользователя. Все поля обязательны к заполнению
type AuthCredentials struct {
	Login    string `json:"login"`    // Логин для входа в систему, должен быть уникальным
	Password string `json:"password"` // Пароль для входа в систему
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

// RequestFull model info
// @Description Полная информация о заявке на отлов бродячей собаки. Поля id, number, status и created_at игнорируются при создании заявки.
type RequestFull struct {
	ID            string    `json:"id"`                  //уникальный идентификатор заявки при Post запросе не указывается
	Number        string    `json:"number"`              //уникальный номер заявки, генерируется автоматически при создании заявки
	Source        Source    `json:"source"`              //Источник информации от кого был получен запрос ID находится в справочнике
	Applicant     Applicant `json:"applicant,omitempty"` //Заявитель ID находится в справочнике
	Address       string    `json:"address"`             //Адрес, где была замечена бродячая собака
	DogsCount     int       `json:"dogs_count"`          //Количество собак
	Behavior      string    `json:"behavior"`            //Поведение собаки
	Urgency       string    `json:"urgency"`             //срочность
	ContactPerson string    `json:"contact_person"`      //Контактные данные
	Status        string    `json:"status"`              //Статус заявки (новая, в работе, выполнена, отменена) генерируется после создания заявки
	CreatedAt     time.Time `json:"created_at"`          //Дата и время создания заявки генерируется автоматически
}

// Source model  info
// @Description Источник информации от кого был получен запрос. Поле name игнорируется при создании заявки.
type Source struct {
	ID   string `json:"id"`             //уникальный идентификатор источника при Post запросе на создание заявки ОБЯЗАТЕЛЕН
	Name string `json:"name,omitempty"` //Название источника игнорируется при создании заявки
}

// Applicant model info
// @Description Заявитель. Поле name игнорируется при создании заявки.
type Applicant struct {
	ID   string `json:"id"`             //уникальный идентификатор заявителя при Post запросе на создание заявки ОБЯЗАТЕЛЕН
	Name string `json:"name,omitempty"` //Название заявителя игнорируется при создании заявки
}

// RequestForGenerating model info
// @Description Структура для генерации заявки. Все поля обязательны к заполнению
type RequestForGenerating struct {
	Number        string `json:"number"`         //уникальный номер заявки ОБЯЗАТЕЛЕН
	Applicant     string `json:"applicant_name"` //Заявитель ID находится в справочнике
	Source        string `json:"source_name"`    //Источник информации от кого был получен запрос ID находится в справочнике
	Address       string `json:"address"`        //Адрес, где была замечена бродячая собака
	DogsCount     int    `json:"dogs_count"`     //Количество собак
	Behavior      string `json:"behavior"`       //Поведение собаки
	Urgency       string `json:"urgency"`        //срочность
	ContactPerson string `json:"contact_person"` //Контактные данные
}

//Response model info
//@Description  Структура для ответа с где обычно ок и сообщние
type Response struct {
	Status  string `json:"status"`  // обычно ok  при ошибках error
	Message string `json:"message"` // просто сообщение

}

//Response model info
//@Description  Структура для ответа с данными пользователя
type AuthResponse struct {
	Status  string      `json:"status"`  // обычно ok  при ошибках error
	Message string      `json:"message"` // просто сообщение
	User    UserProfile `json:"user"`    // данные пользователя
	Tokens  AuthTokens  `json:"tokens"`  // токены авторизации

}

//ResponseUrl model info
//@Description  Структура для ответа с где обычно ок и url
type ResponseUrl struct {
	Status string `json:"status"` // обычно ok  при ошибках error
	Url    string `json:"url"`    // url на скаччивание файла
}
