package main

import (
	"net/http"
	"sync"

	"github.com/web-dashboard-made-by-renz/backend/config"
	"github.com/web-dashboard-made-by-renz/backend/internal/handlers"
	"github.com/web-dashboard-made-by-renz/backend/internal/repository"
	"github.com/web-dashboard-made-by-renz/backend/internal/service"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var (
	router http.Handler
	once   sync.Once
)

func init() {
	functions.HTTP("dashboard", Handler)
}

// Handler dipanggil oleh Google Cloud Function
func Handler(w http.ResponseWriter, r *http.Request) {
	// only initialize once
	once.Do(func() {
		cfg := config.LoadConfig()

		db, err := config.NewDatabase(cfg)
		if err != nil {
			panic("Failed to connect database: " + err.Error())
		}

		colorisRepo := repository.NewColorisRepository(db.DB)
		trainingRepo := repository.NewTrainingRepository(db.DB)
		selloutRepo := repository.NewSelloutRepository(db.DB)

		colorisService := service.NewColorisService(colorisRepo)
		trainingService := service.NewTrainingService(trainingRepo)
		selloutService := service.NewSelloutService(selloutRepo)
		authService := service.NewAuthService(cfg.JWTSecret)

		colorisHandler := handlers.NewColorisHandler(colorisService)
		trainingHandler := handlers.NewTrainingHandler(trainingService)
		selloutHandler := handlers.NewSelloutHandler(selloutService)
		authHandler := handlers.NewAuthHandler(authService)

		router = handlers.SetupRouter(cfg, colorisHandler, trainingHandler, selloutHandler, authHandler)
	})

	router.ServeHTTP(w, r)
}
