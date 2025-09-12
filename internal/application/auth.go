package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

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
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	// Генерация токена
	err = h.SetNewToken(c, id)
	if err != nil {
		//логи
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка генерации токена."})
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

		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	err = h.SetNewToken(c, id)
	if err != nil {
		//логиии
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "я хз какой ответ какой json и ответ тут посылать."})

}

// эендпойнт дял рефреша
func (h *Handlers) Refresh(c *gin.Context) {

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Необходима авторизация."})
		c.Abort()
		return
	}
	claims, err := h.jwtService.DecodeKey(refreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Необходима авторизация."})
		c.Abort()
		return
	}
	if claims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Токен обновления истек."})
		c.Abort()
		return
	}

	id, err := claims.GetSubject()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ошибка."})
		c.Abort()
		return
	}

	h.SetNewToken(c, id)

}

func (h *Handlers) SetNewToken(c *gin.Context, id string) error {

	newAccesClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.AccesTimeExpr).Unix()}

	newRefreshClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.RefreshTimeExpr).Unix()}

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
