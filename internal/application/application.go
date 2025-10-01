package application

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter(handlers *Handlers) *gin.Engine {
	router := gin.Default()

	// router.Use(cors.Default())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
	router.POST("/api/user/change-profile-info", handlers.ChangeProfileInfo)

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
