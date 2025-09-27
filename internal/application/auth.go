package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Sign In
// @Description Авторизация пользователя.
// @Tags Auth
// @Accept json
// @Produce json
// @Param account body dto.AuthCredentials true "Account credentials"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response
// @Router /auth/sign-in [post]
func (h *Handlers) singIn(c *gin.Context) {

	var cred dto.AuthCredentials
	if err := c.Bind(&cred); err != nil {
		//логи
		fmt.Println(err)
	}
	reqStruct := &dto.AuthCredentials{
		Login:    cred.Login,
		Password: cred.Password,
	}
	id, err := h.service.ValidateAccount(reqStruct)
	if err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	// Генерация токена
	access, refhresh, err := h.SetNewToken(c, id)
	if err != nil {
		//логи
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка генерации токена."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "ok",
		"message":       "Успешный вход.",
		"access_token":  access,
		"refresh_token": refhresh,
	})
}

// @Summary Sign Up
// @Description Регистрация нового пользователя.
// @Tags Auth
// @Accept json
// @Produce json
// @Param account body dto.Account true "Account information"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response
// @Router /auth/sign-up [post]
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

	access, refhresh, err := h.SetNewToken(c, id)
	if err != nil {
		//логиии
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "ok",
		"message":       "Аккаунт успешно создан.",
		"access_token":  access,
		"refresh_token": refhresh,
	})

}

// @Summary Refresh Token
// @Description Обновление токенов доступа и обновления. Если access токен истёк то делаешь get запрос без body с куками и получаешь новые токены в куках.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /auth/refresh [get]
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

	access, refhresh, err := h.SetNewToken(c, id)

	c.JSON(http.StatusOK, gin.H{
		"status":        "ok",
		"message":       "Токены обновлены.",
		"access_token":  access,
		"refresh_token": refhresh,
	})

}

// @Summary Change Password
// @Description Смена пароля пользователя.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Change Password information"
// @Security BearerAuth
// @Success 200 {object} dto.Response  	 "Пароль успешно изменен."
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при cмене  пароля."
// @Router /auth/change-password [post]
func (h *Handlers) ChangePassword(c *gin.Context) {
	var reqStruct dto.ChangePasswordRequest
	if err := c.Bind(&reqStruct); err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	err := h.service.ChangePassword(&reqStruct)
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при смене пароля."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Пароль успешно изменен."})

}

// @Summary Get User Profile
// @Description Получение профиля пользователя.
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserProfile
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /user/profile [get]
func (h *Handlers) Profile(c *gin.Context) {

	token := c.Value("access_token").(string)
	claims, err := h.jwtService.DecodeKey(token)

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

	userId, err := claims.GetSubject()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "ошибка."})
		c.Abort()
		return
	}

	profile, err := h.service.GetUserProfile(userId)
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при получении профиля."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": profile})
}

func (h *Handlers) SetNewToken(c *gin.Context, id string) (accessToken string, refhreshToken string, err error) {

	newAccesClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.AccesTimeExpr).Unix()}

	newRefreshClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.RefreshTimeExpr).Unix()}

	accessToken, err = h.jwtService.Encode(newAccesClaims)
	if err != nil {
		// логи

		return "", "", err
	}

	refreshToken, err := h.jwtService.Encode(newRefreshClaims)
	if err != nil {
		// логи
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
