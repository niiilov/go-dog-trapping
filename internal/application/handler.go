package application

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niiilov/go-dog-trapping/internal/dto"
	jw "github.com/niiilov/go-dog-trapping/pkg/jwt"
)

type Service interface {
	CreateAccount(account *dto.Account) (string, error)
	ValidateAccount(account *dto.AuthCredentials) (profile *dto.UserProfile, role_id string, err error)
	SendRequest(request *dto.RequestFull) error
	GetAllRequests() ([]*dto.RequestFull, error)
	GetRequestsByOtdel(otdel_id string) ([]*dto.RequestFull, error)
	ChangePassword(req *dto.ChangePasswordRequest) error
	ChangeProfileInfo(req *dto.ChangeProfileRequest) error
	GetUserProfile(userId string) (*dto.UserProfile, error)
	GetFileURL(objectKey string) string
	ChangeStatusRequest(req *dto.ChangeStatusRequest) error
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
// @Success 200 {object} []dto.RequestFull	"Все заявки"
// @Failure 500 {object} dto.Response 	 "Ошибка при получении запросов."
// @Router /requests [get]
func (h *Handlers) GetAllRequests(c *gin.Context) {
	role_id, exists := c.Get("role_id")

	fmt.Println(role_id, "РОЛЬ ID  В GET AL REQ")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Доступ запрещен."})
		return
	}

	if role_id == "794c900e-f04f-496d-a4e6-5cf5d69fad90" {
		requests, err := h.service.GetAllRequests()
		if err != nil {
			// опять логи
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении запросов."})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "data": requests})
		return
	} else {

		h.GetRequestsByOtdel(c, role_id.(string))
	}

}

// @Summary Get Requests By Otdel
// @Security BearerAuth
// @Description Получение запросов на отлов бродячих собак по отделу. Требуется параметр otdel_id в query, otdel_id находится в справочнике. Доступно только для районных администраторов.
// @Tags Requests
// @Produce json
// @Param otdel_id query string true "Otdel ID"
// @Success 200 {object} []dto.RequestFull
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении запросов."
// @Router /requests_otdel [get]
func (h *Handlers) GetRequestsByOtdel(c *gin.Context, role_id string) {

	requests, err := h.service.GetRequestsByOtdel(role_id)

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
// @Success 200 {object} dto.ResponseUrl
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при получении запросов."
// @Router /requests/download_url [get]
func (h *Handlers) DowloadUrl(c *gin.Context) {
	number := c.Query("number")
	fmt.Println(number)
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}
	filename := "zayavka_" + number + ".xlsx"

	fmt.Println(filename)
	url := h.service.GetFileURL(filename)

	c.JSON(http.StatusOK, gin.H{"status": "ok", "url": url})
}
