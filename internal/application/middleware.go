package application

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) authMiddleware(c *gin.Context) {
	if c.Request.URL.Path == "/api/singup" || c.Request.URL.Path == "/api/singin" || c.Request.URL.Path == "/api/send_request" || c.Request.URL.Path == "/api/requests" {
		c.Next()
		return
	}
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, `{"status": "error","message": "Необходима авторизация."}`)
		c.Abort()
		return
	}
	claims, err := h.jwtService.DecodeKey(accessToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, `{"status": "error","message": "Необходима авторизация."}`)
		c.Abort()
		return
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, `{"status": "error","message": "Необходима авторизация."}`)
			c.Abort()
			return
		}
		claims, err := h.jwtService.DecodeKey(refreshToken)

		if err != nil {
			c.JSON(http.StatusUnauthorized, `{"status": "error","message": "Необходима авторизация."}`)
			c.Abort()
			return
		}
		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, `{"status": "error","message": "Токен обновления истек."}`)
			c.Abort()
			return
		}
	}
	c.Next()

}
