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
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.AuthResponse	"Ошибка в данных запроса."
// @Failure 500 {object} dto.AuthResponse
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
	user, err := h.service.ValidateAccount(reqStruct)
	if err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	// Генерация токена
	tokens, err := h.SetNewToken(c, user.ID)
	if err != nil {
		//логи
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка генерации токена."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Успешный вход.",
		"user":    user,
		"tokens":  tokens,
	})
}

// @Summary Sign Up
// @Description Регистрация нового пользователя.
// @Tags Auth
// @Accept json
// @Produce json
// @Param account body dto.Account true "Account information"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.AuthResponse	"Ошибка в данных запроса."
// @Failure 500 {object} dto.AuthResponse
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
	reqStruct.ID = id

	tokens, err := h.SetNewToken(c, id)
	if err != nil {
		//логиии
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Аккаунт успешно создан.",
		"user":    reqStruct,
		"tokens":  tokens,
	})

}

// @Summary Refresh Token
// @Description Обновление токенов доступа и обновления.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.AuthTokens
// @Failure 401 {object} dto.AuthTokens
// @Failure 500 {object} dto.AuthTokens
// @Router /auth/refresh [post]
func (h *Handlers) Refresh(c *gin.Context) {

	var tokenStruct dto.AuthTokens
	if err := c.Bind(&tokenStruct); err != nil {
		c.JSON(http.StatusBadRequest, err)
		c.Abort()
		//логи
	}

	fmt.Println(tokenStruct)
	claims, err := h.jwtService.DecodeKey(tokenStruct.RefreshToken)

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

	tokens, err := h.SetNewToken(c, id)
	if err != nil {
		//логиии
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Токены обновлены.",
		"tokens":  tokens,
	})

}

// @Summary Change Password
// @Description Смена пароля пользователя.
// @Tags User
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

// @Summary Change Profile Info
// @Description Смена данных профлия пользователя.
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.ChangeProfileRequest true "Change Profile information"
// @Security BearerAuth
// @Success 200 {object} dto.Response  	 "Информация успешна изменена."
// @Failure 400 {object} dto.Response	"Ошибка в данных запроса."
// @Failure 500 {object} dto.Response	"Ошибка при cмене  данных."
// @Router /auth/change-password [post]
func (h *Handlers) ChangeProfileInfo(c *gin.Context) {
	var reqStruct dto.ChangeProfileRequest

	if err := c.Bind(&reqStruct); err != nil {
		//логи
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}
	err := h.service.ChangeProfileInfo(&reqStruct)
	if err != nil {
		// опять логи
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка при изменение информации."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Информация успешно изменена."})

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
	profile.ID = userId

	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": profile})
}

func (h *Handlers) SetNewToken(c *gin.Context, id string) (*dto.AuthTokens, error) {

	newAccesClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.AccesTimeExpr).Unix()}

	newRefreshClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.RefreshTimeExpr).Unix()}

	var tokens dto.AuthTokens

	access, err := h.jwtService.Encode(newAccesClaims)
	if err != nil {
		// логи

		return nil, err
	}
	tokens.AccesToken = access

	tokens.RefreshToken, err = h.jwtService.Encode(newRefreshClaims)
	if err != nil {
		// логи
		return nil, err
	}

	return &tokens, nil
}
