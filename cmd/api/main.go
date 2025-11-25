package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/web-dashboard-made-by-renz/backend/config"
	"github.com/web-dashboard-made-by-renz/backend/internal/handlers"
	"github.com/web-dashboard-made-by-renz/backend/internal/repository"
	"github.com/web-dashboard-made-by-renz/backend/internal/service"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

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

	router := handlers.SetupRouter(cfg, colorisHandler, trainingHandler, selloutHandler, authHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Printf("Server is running on port %s", cfg.ServerPort)
		log.Println("Open on http://localhost:" + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
