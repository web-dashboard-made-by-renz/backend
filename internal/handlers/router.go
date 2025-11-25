package handlers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/web-dashboard-made-by-renz/backend/config"
	"github.com/web-dashboard-made-by-renz/backend/internal/middleware"
)

func SetupRouter(cfg *config.Config, colorisHandler *ColorisHandler, trainingHandler *TrainingHandler, selloutHandler *SelloutHandler, authHandler *AuthHandler) *gin.Engine {
	router := gin.Default()
	registerRoutes := func(group *gin.RouterGroup) {
		group.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Backend Dashboard API is running",
			})
		})

		// Auth routes (public)
		auth := group.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.GET("/verify", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.Verify)
		}

		// Protected routes
		protected := group.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			coloris := protected.Group("/coloris")
			{
				coloris.POST("", colorisHandler.CreateColoris)
				coloris.GET("", colorisHandler.GetAllColoris)
				coloris.GET("/:id", colorisHandler.GetColorisById)
				coloris.PUT("/:id", colorisHandler.UpdateColoris)
				coloris.DELETE("/:id", colorisHandler.DeleteColoris)
				coloris.POST("/import", colorisHandler.ImportExcel)
				coloris.GET("/export", colorisHandler.ExportExcel)
			}

			training := protected.Group("/training")
			{
				training.POST("", trainingHandler.CreateTraining)
				training.GET("", trainingHandler.GetAllTraining)
				training.GET("/:id", trainingHandler.GetTrainingById)
				training.PUT("/:id", trainingHandler.UpdateTraining)
				training.DELETE("/:id", trainingHandler.DeleteTraining)
				training.POST("/import", trainingHandler.ImportExcel)
				training.GET("/export", trainingHandler.ExportExcel)
			}

			sellout := protected.Group("/sellout")
			{
				sellout.POST("", selloutHandler.CreateSellout)
				sellout.GET("", selloutHandler.GetAllSellout)
				sellout.GET("/:id", selloutHandler.GetSelloutById)
				sellout.PUT("/:id", selloutHandler.UpdateSellout)
				sellout.DELETE("/:id", selloutHandler.DeleteSellout)
				sellout.POST("/import", selloutHandler.ImportExcel)
				sellout.GET("/export", selloutHandler.ExportExcel)
			}
		}
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.AllowedOrigins},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.MaxMultipartMemory = 32 << 20

	// Default API prefix
	apiV1 := router.Group("/api/v1")
	registerRoutes(apiV1)

	// Cloud Functions adds the function name as path prefix (e.g. /dashboard/*),
	// so expose the same routes under that prefix to avoid 404s.
	funcPrefixed := router.Group("/dashboard")
	registerRoutes(funcPrefixed)
	funcPrefixedV1 := router.Group("/dashboard/api/v1")
	registerRoutes(funcPrefixedV1)

	return router
}
