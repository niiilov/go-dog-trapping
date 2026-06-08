package application

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/niiilov/go-dog-trapping/internal/dto"
)

// @Summary Login
// @Description Авторизация пользователя.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.AuthCredentials true "Login information"
// @Success 200 {object} dto.AuthResponse
// @Failure 400 {object} dto.AuthResponse
// @Failure 500 {object} dto.AuthResponse
// @Router /api/auth/login [post]
func (h *Handlers) Login(c *gin.Context) {

	var cred dto.AuthCredentials
	if err := c.Bind(&cred); err != nil {
		fmt.Println(err)
	}
	reqStruct := &dto.AuthCredentials{
		Login:    cred.Login,
		Password: cred.Password,
	}

	user, err := h.service.ValidateAccount(reqStruct)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Ошибка в данных запроса."})
		return
	}

	tokens, err := h.SetNewToken(c, user.ID, user.RoleID, user.TerOtdelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Ошибка генерации токена."})
		return
	}

	c.JSON(http.StatusOK, dto.AuthResponse{
		Status:  "ok",
		Message: "Успешный вход.",
		User:    user,
		Tokens:  tokens,
	})
}

// @Summary Register
// @Description Регистрация нового пользователя.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserDTO true "User information"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/register [post]
func (h *Handlers) Register(c *gin.Context) {

	var reqStruct dto.CreateUserDTO
	if err := c.Bind(&reqStruct); err != nil {
		fmt.Println(err)
	}
	err := h.service.CreateUser(&reqStruct)
	fmt.Println(err)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Аккаунт успешно создан.",
		"user":    reqStruct,
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
	}

	claims, roleID, terOtdelID, err := h.jwtService.DecodeKey(tokenStruct.RefreshToken)

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

	tokens, err := h.SetNewToken(c, id, roleID, terOtdelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Токены обновлены.",
		"tokens":  tokens,
	})

}

func (h *Handlers) SetNewToken(c *gin.Context, id string, role_id string, ter_otdel_id string) (*dto.AuthTokens, error) {

	newAccesClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.AccesTimeExpr).Unix(), "role": role_id, "ter_otdel_id": ter_otdel_id}

	newRefreshClaims := jwt.MapClaims{"sub": id, "exp": time.Now().Add(dto.RefreshTimeExpr).Unix(), "role": role_id, "ter_otdel_id": ter_otdel_id}

	var tokens dto.AuthTokens

	access, err := h.jwtService.Encode(newAccesClaims)
	if err != nil {
		return nil, err
	}
	tokens.AccesToken = access

	tokens.RefreshToken, err = h.jwtService.Encode(newRefreshClaims)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}
