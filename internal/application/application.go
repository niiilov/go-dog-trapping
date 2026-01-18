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
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization", // Явно указываем Authorization
			"X-Requested-With",
			"Access-Control-Allow-Origin",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Authorization",
		},
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
	router.POST("/api/requests/upload-act", handlers.UploadAct)

	router.POST("/api/request/download_multiDate", handlers.GenerateMultipleByDate)
	router.POST("/api/request/download_multiID", handlers.GenerateMultipleByIDs)

	router.POST("/api/requests/change-status", handlers.ChangeStatusRequest)

	router.GET("/api/requests", handlers.GetAllRequests)

	router.GET("/api/user/profile", handlers.Profile)
	router.GET("/api/requests/download_request", handlers.DownloadRequest)
	router.GET("/api/requests/download_act", handlers.DownloadAct)

	router.GET("/api/sources", handlers.GetApplicants)
	router.GET("/api/applicants", handlers.GetTerrOtdels)

	router.DELETE("/api/requests/:id", handlers.DeleteRequest)

	router.POST("/api/sources", handlers.AddNewTerOtdel)

	router.POST("/api/external/users", handlers.CreateExternalUser)
	router.PUT("/api/external/users/:id", handlers.ChangeAcceptedExternalUser)
	router.GET("/api/external/users", handlers.GetExternalUsers)

	return router
}

func StartApplication(addr string, handlers *Handlers) {
	router := InitRouter(handlers)
	router.Run(addr)
}
