package application

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(handlers *Handlers) *gin.Engine {
	router := gin.Default()

	router.Use(cors.Default())

	router.Use(handlers.authMiddleware)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/api/auth/sign-up", handlers.singUp)
	router.POST("/api/auth/sign-in", handlers.singIn)
	router.POST("/api/requests", handlers.SendRequest)
	router.POST("/api/auth/change-password", handlers.ChangePassword)
	router.POST("/api/auth/refresh", handlers.Refresh)

	router.GET("/api/requests", handlers.GetAllRequests)
	router.GET("/api/requests_otdel", handlers.GetRequestsByOtdel)
	router.GET("/api/user/profile", handlers.Profile)
	router.GET("/api/requests/download_url", handlers.DowloadUrl)

	return router
}

func StartApplication(addr string, handlers *Handlers) {
	router := InitRouter(handlers)
	router.Run(addr)
}
