package application

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) authMiddleware(c *gin.Context) {
	if c.Request.URL.Path == "/api/auth/sign-up" || c.Request.URL.Path == "/api/auth/sign-in" || c.Request.URL.Path == "/api/auth/refresh" {
		c.Next()
		return
	}

	if c.Request.URL.Path == "/api/requests" && c.Request.Method == http.MethodPost {
		c.Next()
		return
	}

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Необходима авторизация."})
		c.Abort()
		return
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization header format"})
		return
	}

	accessToken := parts[1]

	c.Set("access_token", accessToken)

	claims, role_id, err := h.jwtService.DecodeKey(accessToken)

	c.Set("role_id", role_id)

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Необходима авторизация."})
		c.Abort()
		return
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Срок действия токена истек."})
		c.Abort()
		return
	}

	c.Next()

}
