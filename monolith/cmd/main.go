package main

import (
	"perfect_api/internal/config"
	"perfect_api/internal/database"
	"perfect_api/internal/logger"

	userHandler "perfect_api/internal/user/handler"
	userRepository "perfect_api/internal/user/repository"
	userService "perfect_api/internal/user/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// initialize the log
	logger.InitLogger()
	logger.Log.Info("Starting Monolith Wallet Application...")

	// 1. load configuration
	cfg := config.LoadConfig()

	// 2. connect to database with retry
	db, err := database.ConnectWithRetry(cfg.DBDSN)
	if err != nil {
		logger.Log.Error("Critical Error: Could not connect to database after retries", "error", err)
	}
	defer db.Close()

	// 1. initiate layer
	uRepo := userRepository.NewMySQLUserRepository(db)
	uSvc := userService.NewUserService(uRepo)
	uHandler := userHandler.NewUserHandler(uSvc)

	// 2. setup gin router
	r := gin.Default()

	// routes
	r.POST("/api/v1/users", uHandler.Register)
	r.GET("/api/v1/users/:id", uHandler.GetProfile)
	r.PUT("/api/v1/users/:id", uHandler.UpdateProfile)

	// start server
	logger.Log.Info("Server running on port 8080....")
	if err := r.Run(":8080"); err != nil {
		logger.Log.Error("Server failed to run", "error", err)
	}
}
