package application

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
	jw "github.com/niiilov/go-dog-trapping/pkg/jwt"
)

type Service interface {
	CreateAccount(account *dto.Account) (string, error)
	ValidateAccount(account *dto.AuthCredentials) (profile *dto.UserProfile, role_id string, err error)
	SendRequest(request *dto.RequestFull) error

	UploadAct(req *dto.UploadActRequests, key string, filename string) error

	ChangeStatusRequest(req *dto.ChangeStatusRequest) error
	ChangePassword(req *dto.ChangePasswordRequest) error
	ChangeProfileInfo(req *dto.ChangeProfileRequest) error

	GenerateMultiple(req *dto.GenerateMultipleRequest) (string, error)

	GetAllRequests(year string) ([]*dto.RequestFull, error)
	GetRequestsByOtdel(otdel_id string, year string) ([]*dto.RequestFull, error)
	GetUserProfile(userId string) (*dto.UserProfile, error)
	GetFileURL(objectKey string) string
}
type Handlers struct {
	jwtService *jw.ServiceJWT

	service Service
}

func NewHandlers(service Service, jwtService *jw.ServiceJWT) *Handlers {
	return &Handlers{service: service, jwtService: jwtService}
}

// @Summary Send Request
// @Security BearerAuth
// @Description Отправка запроса на отлов бродячей собаки. Поля source_id и applicant_id заполнять ID из справочников. Поля name в этих полях игнорируются при отправке запроса. Доступно для всех авторизованных пользователей.
// @Tags Requests
// @Accept json
// @Produce json
// @Param request body dto.RequestFull true "Request information"
// @Success 200 {object} dto.Response	"Запрос успешно отправлен"
// @Failure 400 {object} dto.Response  	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при отправке запроса"
// @Router /requests [post]
func (h *Handlers) SendRequest(c *gin.Context) {

	var reqStruct dto.RequestFull
	if err := c.Bind(&reqStruct); err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	err := h.service.SendRequest(&reqStruct)
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при отправке запроса"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Запрос успешно отправлен"})
}

// @Summary Change Status Request
// @Security BearerAuth
// @Description Изменение статуса заявки на отлов бродячей собаки. Поля id и status обязательны к заполнению.
// @Tags Requests
// @Accept json
// @Produce json
// @Param request body dto.ChangeStatusRequest true "Change Status Request"
// @Success 200 {object} dto.Response	"Статус заявки успешно изменен."
// @Failure 400 {object} dto.Response  	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при изменении статуса заявки."
// @Router /requests/change-status [post]
func (h *Handlers) ChangeStatusRequest(c *gin.Context) {
	var reqStruct dto.ChangeStatusRequest
	if err := c.Bind(&reqStruct); err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}
	err := h.service.ChangeStatusRequest(&reqStruct)
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при изменении статуса заявки."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Статус заявки успешно изменен."})
}

// @Summary Get All Requests
// @Security BearerAuth
// @Description Получение всех запросов на отлов бродячих собак. Доступно только для районных администраторов.
// @Tags Requests
// @Produce json
// @Success 200 {object} []dto.RequestFull	"Все заявки"  query param otdel_id - для получения заявок по отделу, доступно только для  админа
// @Failure 400 {object} dto.Response 	 "Ошибка в данных запроса."
// @Failure 500 {object} dto.Response 	 "Ошибка при получении запросов."
// @Router /requests [get]
func (h *Handlers) GetAllRequests(c *gin.Context) {
	role_id, exists := c.Get("role_id")

	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Доступ запрещен."})
		return
	}
	year := c.Query("year")
	fmt.Println("year:", year)

	// Если роль - супер админ, то возвращаем все заявки
	// Иначе возвращаем заявки по отделу
	// 794c900e-f04f-496d-a4e6-5cf5d69fad90 - это роль супер админа

	if role_id == "794c900e-f04f-496d-a4e6-5cf5d69fad90" {

		if otdel_id := c.Query("otdel_id"); otdel_id != "" {
			h.GetRequestsByOtdelFunc(c, otdel_id, year)
			return
		}
		requests, err := h.service.GetAllRequests(year)
		if err != nil {
			// опять логи
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении запросов."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "data": requests})
		return
	} else {

		h.GetRequestsByOtdelFunc(c, role_id.(string), year)
	}

}

// GetRequestsByOtdelFunc - вспомогательная функция для получения заявок по отделу
func (h *Handlers) GetRequestsByOtdelFunc(c *gin.Context, otdel_id string, year string) {

	requests, err := h.service.GetRequestsByOtdel(otdel_id, year)

	if err != nil {

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении запросов."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": requests})
}

// @Summary Download Request File URL
// @Security BearerAuth
// @Description Получение URL для скачивания файла заявки в формате .docx. Требуется параметр number в query, номер заявки. Доступно для всех авторизованных пользователей.
// @Tags Requests
// @Produce json
// @Param number query string true "Request Number"
// @Param year query string true "Request Year"
// @Success 200 {object} dto.ResponseUrl
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении запросов."
// @Router /requests/download_request [get]
func (h *Handlers) DownloadRequest(c *gin.Context) {
	number := c.Query("number")
	year := c.Query("year")
	fmt.Println(number)
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}
	filename := "zayavka_" + number + "_" + year + ".xlsx"

	fmt.Println(filename)
	url := h.service.GetFileURL(filename)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "url": url})
}

func (h *Handlers) GenerateMultiple(c *gin.Context) {
	var request dto.GenerateMultipleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}
	url, err := h.service.GenerateMultiple(&request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при генерации запросов."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "url": url})
}

// @Summary Upload Act File
// @Security BearerAuth
// @Description Загрузка акта выполненого отлова любой формат файла, доступно админу и подрядчику.
// @Tags Requests
// @Accept multipart/form-data
// @Produce json
// @Param number formData string true "Request Number"
// @Param id formData string true "Request ID"
// @Param status formData string true "Request Status"
// @Param file formData file true "Act File"
// @Success 200 {object} dto.Response	"Файл успешно загружен."
// @Failure 400 {object} dto.Response  	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при загрузке акта."
// @Router /requests/upload_act [post]
func (h *Handlers) UploadAct(c *gin.Context) {

	var req dto.UploadActRequests
	req.ID = c.PostForm("id")
	req.Number = c.PostForm("number")
	req.Status = c.PostForm("status")

	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка при получении файла."})
		return
	}

	sharedDir := "/app/shared/"
	filename := "act_" + req.Number + "_" + time.Now().Format("2006") + filepath.Ext(file.Filename)
	fmt.Println(filename)

	if err := c.SaveUploadedFile(file, sharedDir+filename); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при сохранении файла."})
		return
	}

	if err = h.service.UploadAct(&req, filename, sharedDir+filename); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при загрузке акта."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Файл успешно загружен."})
}

// @Summary Download Act File URL
// @Security BearerAuth
// @Description Получение URL для скачивания файла акта выполненого отлова в формате .docx. Требуется параметры number и year в query, номер заявки и год. Доступно для всех авторизованных пользователей.
// @Tags Requests
// @Produce json
// @Param number query string true "Request Number"
// @Param year query string true "Request Year"
// @Success 200 {object} dto.ResponseUrl
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении запросов."
// @Router /requests/download_act [get]
func (h *Handlers) DownloadAct(c *gin.Context) {
	var req dto.DownloadActRequest
	req.Number = c.Query("number")
	req.Year = c.Query("year")

	if req.Number == "" || req.Year == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}
	filename := "act_" + req.Number + "_" + req.Year + ".docx"

	url := h.service.GetFileURL(filename)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "url": url})
}
