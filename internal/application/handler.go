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
	ValidateAccount(account *dto.Account) (id string, err error)
	SendRequest(request *dto.RequestFull) error
	GetAllRequests() ([]*dto.RequestFull, error)
	GetRequestsByOtdel(otdel_id string) ([]*dto.RequestFull, error)
}
type Handlers struct {
	jwtService *jw.ServiceJWT

	service Service
}

func NewHandlers(service Service, jwtService *jw.ServiceJWT) *Handlers {
	return &Handlers{service: service, jwtService: jwtService}
}

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

func (h *Handlers) GetAllRequests(c *gin.Context) {
	requests, err := h.service.GetAllRequests()
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении запросов."})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": requests})
}

func (h *Handlers) GetRequestsByOtdel(c *gin.Context) {
	otdel_id := c.Query("otdel_id")
	if otdel_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}
	fmt.Println("Otdel ID:", otdel_id)
	requests, err := h.service.GetRequestsByOtdel(otdel_id)

	if err != nil {

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении запросов."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": requests})
}
