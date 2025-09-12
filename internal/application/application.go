package application

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter(handlers *Handlers) *gin.Engine {
	router := gin.Default()

	router.Use(handlers.authMiddleware)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.POST("/api/auth/singup", handlers.singUp)
	router.POST("/api/auth/singin", handlers.singIn)
	router.POST("/api/send_request", handlers.SendRequest)

	router.GET("/api/auth/refresh", handlers.Refresh)
	router.GET("/api/all_requests", handlers.GetAllRequests)
	router.GET("/api/requests_otdel", handlers.GetRequestsByOtdel)

	return router
}

func StartApplication(addr string, handlers *Handlers) {
	router := InitRouter(handlers)
	router.Run(addr)
}
