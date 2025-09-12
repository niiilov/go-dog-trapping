package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func (h *Handlers) singIn(c *gin.Context) {

	var reqStruct dto.Account
	if err := c.Bind(&reqStruct); err != nil {
		//логи
		fmt.Println(err)
	}
	id, err := h.service.ValidateAccount(&reqStruct)
	if err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, `{"status": "error","message": "Ошибка в данных запроса."}`)
		return
	}

	// Генерация токена
	err = h.SetNewToken(c, id)
	if err != nil {
		//логи
		c.JSON(http.StatusInternalServerError, `{"status": "error","message": "Ошибка генерации токена."}`)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Успешный вход."})
}

func (h *Handlers) singUp(c *gin.Context) {

	var reqStruct dto.Account
	if err := c.Bind(&reqStruct); err != nil {
		//логи
	}

	id, err := h.service.CreateAccount(&reqStruct)

	if err != nil {
		// опять логи

		c.JSON(http.StatusBadRequest, `{"status": "error","message": "Ошибка в данных запроса."}`)
		return
	}

	err = h.SetNewToken(c, id)
	if err != nil {
		//логиии
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, `{"status": "ok","message": "я хз какой ответ какой json и ответ тут посылать."}`)

}

func (h *Handlers) SendRequest(c *gin.Context) {

	var reqStruct dto.RequestFull
	if err := c.Bind(&reqStruct); err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, `{"status": "error","message": "Ошибка в данных запроса."}`)
		return
	}

	err := h.service.SendRequest(&reqStruct)
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, `{"status": "error","message": "Ошибка при отправке запроса."}`)
		return
	}

	c.JSON(http.StatusOK, `{"status": "ok","message": "Запрос успешно отправлен."}`)
}

func (h *Handlers) GetAllRequests(c *gin.Context) {
	requests, err := h.service.GetAllRequests()
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, `{"status": "error","message": "Ошибка при получении запросов."}`)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": requests})
}

func (h *Handlers) GetRequestsByOtdel(c *gin.Context) {
	otdel_id := c.Query("otdel_id")
	if otdel_id == "" {
		c.JSON(http.StatusBadRequest, `{"status": "error","message": "Ошибка в данных запроса."}`)
		return
	}
	fmt.Println("Otdel ID:", otdel_id)
	requests, err := h.service.GetRequestsByOtdel(otdel_id)

	if err != nil {

		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, `{"status": "error","message": "Ошибка при получении запросов."}`)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": requests})
}

// // ниже функции для работы с jwt и тп не используются для роутинга
// func (h *Handlers) getIdFromSubject(c *gin.Context) (string, error) {

// 	authHeader := c.Request.Header.Get("Authorization")

// 	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

// 	token, err := h.jwtService.DecodeKey(tokenString)
// 	if err != nil {

// 		return "", err
// 	}

// 	id, err := token.GetSubject()
// 	if err != nil {
// 		return "", err
// 	}
// 	return id, nil
// }

func (h *Handlers) SetNewToken(c *gin.Context, id string) error {

	newAccesClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.AccesTimeExpr)}

	newRefreshClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.RefreshTimeExpr)}

	accessToken, err := h.jwtService.Encode(newAccesClaims)
	if err != nil {
		// логи

		return err
	}

	refreshToken, err := h.jwtService.Encode(newRefreshClaims)
	if err != nil {
		// логи
		return err
	}

	c.SetCookie(
		"refresh_token",
		refreshToken,
		int(dto.RefreshTimeExpr),
		"/",
		"localhost",
		false,
		true,
	)

	c.SetCookie(
		"access_token",
		accessToken,
		int(dto.AccesTimeExpr),
		"/",
		"localhost",
		false,
		true,
	)

	return nil
}
