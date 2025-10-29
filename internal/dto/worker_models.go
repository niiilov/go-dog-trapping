package dto

// RequestForGeneratingSomething model info
// @Description Структура для генерации несколько заявок. Все поля обязательны к заполнению
type RequestForGeneratingSomething struct {
	Requests []*RequestFull `json:"requests"` // Список заявок
	Number   string         `json:"number"`   // Номер заявки
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
