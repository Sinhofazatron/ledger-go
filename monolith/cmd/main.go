package main

import (
	_ "perfect_api/docs"
	"time"

	"perfect_api/internal/config"
	"perfect_api/internal/database"
	"perfect_api/internal/logger"
	"perfect_api/internal/middleware"

	userHandler "perfect_api/internal/user/handler"
	userRepository "perfect_api/internal/user/repository"
	userService "perfect_api/internal/user/service"

	walletHandler "perfect_api/internal/wallet/handler"
	walletRepository "perfect_api/internal/wallet/repository"
	walletService "perfect_api/internal/wallet/service"

	ledgerRepository "perfect_api/internal/ledger/repository"

	txHandler "perfect_api/internal/transaction/handler"
	txRepository "perfect_api/internal/transaction/repository"
	txService "perfect_api/internal/transaction/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title			GoWallet Monolith API
// @version			1.0
// @description		API Documentation for GoWallet
// @termOfService	http://swagger.io/terms/

// @contact.name	API Support
// @contact.email	bashocode@gmail.com

// @host			localhost:8080
// @basepath		/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer <your_token>" to authenticate.
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

	// connect to redis
	rdb, err := database.ConnectRedis(cfg.RedisAddr)
	if err != nil {
		logger.Log.Error("Critical Error: Could not connect to Redis", "error", err)
	}
	defer rdb.Close()

	// 1. initiate layer
	uRepo := userRepository.NewMySQLUserRepository(db)
	wRepo := walletRepository.NewMySQLWalletRepository(db)
	tRepo := txRepository.NewMySQLTransactionRepository(db)
	lRepo := ledgerRepository.NewMysqlLedgerRepository(db)

	// inject db to user service for transaction
	uSvc := userService.NewUserService(db, rdb, uRepo, wRepo)
	wSvc := walletService.NewWalletService(wRepo, rdb)
	tSvc := txService.NewTransactionService(db, rdb, tRepo, uRepo, wRepo, lRepo)

	uHandler := userHandler.NewUserHandler(uSvc)
	wHandler := walletHandler.NewWalletHandler(wSvc)
	tHandler := txHandler.NewTransactionHandler(tSvc)

	// 2. setup gin router
	r := gin.New()
	r.Use(gin.Recovery())
	// Register global error handling middleware
	r.Use(middleware.ErrorHandler())

	// apply global rate limiter max 60 request per minutes per ip
	r.Use(middleware.RateLimiter(rdb, 60, time.Minute))

	// register the swagger api
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Route grouping
	v1 := r.Group("/api/v1")
	{
		// Public routes
		v1.POST("/users/register", uHandler.Register)
		v1.POST("/users/login", uHandler.Login)

		// Protected routes (requires valid JWT token)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(rdb))
		{
			protected.GET("/users/me", uHandler.GetProfileMe)
			protected.POST("/users/avatar", uHandler.UploadAvatar)
			protected.PUT("/users/:id", uHandler.UpdateProfile)
			protected.GET("/users/:id", uHandler.GetProfile)
			protected.DELETE("/users/me", uHandler.DeleteAccount)
			protected.POST("/users/logout", uHandler.Logout)

			protected.GET("/wallets/me", wHandler.GetMyWallet)

			protected.POST("/transactions/transfer", tHandler.Transfer)
			protected.GET("/transactions/history", tHandler.GetHistory)
		}
	}

	// start server
	logger.Log.Info("Server running on port 8080....")
	if err := r.Run(":8080"); err != nil {
		logger.Log.Error("Server failed to run", "error", err)
	}
}
