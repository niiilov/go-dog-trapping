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

	router.POST("/api/auth/login", handlers.Login)
	router.POST("/api/auth/register", handlers.Register)

	router.GET("/api/users", handlers.GetUsers)
	router.DELETE("/api/users/:id", handlers.DeleteUser)

	router.POST("/api/ter-otdel", handlers.CreateTerOtdel)
	router.GET("/api/ter-otdel", handlers.GetTerOtdels)
	router.DELETE("/api/ter-otdel/:id", handlers.DeleteTerOtdel)

	router.POST("/api/districts", handlers.CreateDistrict)
	router.GET("/api/districts", handlers.GetDistricts)
	router.DELETE("/api/districts/:id", handlers.DeleteDistrict)

	router.GET("/api/applicants/:districtID", handlers.GetApplicantByDistrictID)
	router.POST("/api/applicants", handlers.CreateApplicant)

	router.GET("/api/requests", handlers.GetRequests)
	router.POST("/api/requests", handlers.CreateRequest)
	router.PUT("/api/requests/status", handlers.ChangeStatusRequest)
	router.DELETE("/api/requests/:id", handlers.DeleteRequest)
	router.POST("/api/requests/act", handlers.AddActFile)
	router.GET("/api/requests/act/:filename", handlers.DownloadAct)

	router.POST("/api/generate", handlers.Generate)

	router.GET("/api/roles", handlers.GetRoles)

	// router.POST("/api/auth/sign-up", handlers.singUp)
	// router.POST("/api/auth/sign-in", handlers.singIn)
	// router.POST("/api/requests", handlers.SendRequest)
	// router.POST("/api/auth/change-password", handlers.ChangePassword)
	// router.POST("/api/auth/refresh", handlers.Refresh)
	// router.POST("/api/user/change-profile-info", handlers.ChangeProfileInfo)
	// router.POST("/api/requests/upload-act", handlers.UploadAct)

	// router.POST("/api/request/download_multiDate", handlers.GenerateMultipleByDate)
	// router.POST("/api/request/download_multiID", handlers.GenerateMultipleByIDs)

	// router.POST("/api/requests/change-status", handlers.ChangeStatusRequest)

	// router.GET("/api/requests", handlers.GetAllRequests)

	// router.GET("/api/user/profile", handlers.Profile)
	// router.GET("/api/requests/download_request", handlers.DownloadRequest)
	// router.GET("/api/requests/download_act", handlers.DownloadAct)

	// router.GET("/api/sources", handlers.GetApplicants)
	// router.GET("/api/applicants", handlers.GetTerrOtdels)

	// router.DELETE("/api/requests/:id", handlers.DeleteRequest)

	// router.POST("/api/sources", handlers.AddNewTerOtdel)

	// router.POST("/api/external/users", handlers.CreateExternalUser)
	// router.PUT("/api/external/users/:id", handlers.ChangeStatusDisAcceptExternalUser)
	// router.PUT("/api/external/users/:id/activate", handlers.ChangeStatusAcceptExternalUser)
	// router.POST("/api/external/:id", handlers.ChangeAcceptedExternalUser)
	// router.GET("/api/external/users", handlers.GetExternalUsers)

	return router
}

func StartApplication(addr string, handlers *Handlers) {
	router := InitRouter(handlers)
	router.Run(addr)
}
