package main

import (
	"log"

	authHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/handler"
	authRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/repository"
	authService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/service"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/config"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/database"
	// programHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/handler"
	// programService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/service" 
	// programRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/repository" 
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/email"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration from .env file
	log.Println("Loading configuration...")
	config.LoadConfig()

	// Initialize database connection
	log.Println("Connecting to database...")
	db, err := database.NewMySQLConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	authRepository := authRepo.NewAuthRepository(db)
	// programRepository := programRepo.NewProgramRepository(db) // NONAKTIFKAN

	// Initialize email service
	emailService := email.NewEmailService(
		config.AppConfig.SMTPHost,
		config.AppConfig.SMTPPort,
		config.AppConfig.SMTPUser,
		config.AppConfig.SMTPPassword,
		config.AppConfig.SMTPSenderEmail,
	)

	// Initialize services
	authSvc := authService.NewAuthService(authRepository, config.AppConfig.GoogleClientID, emailService, *config.AppConfig)
	// programSvc := programService.NewProgramService(programRepository) // NONAKTIFKAN

	// Initialize handlers
	authHdl := authHandler.NewAuthHandler(authSvc)
	// programHdl := programHandler.NewProgramHandler(programSvc) // NONAKTIFKAN

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	setupRoutes(router, authHdl) // HAPUS programHdl DARI SINI

	// Get server port from config
	serverPort := config.AppConfig.ServerPort
	if serverPort == "" {
		serverPort = "8080"
	}

	// Start server
	log.Printf("🚀 Server starting on port %s", serverPort)
	log.Printf("📝 API Documentation: http://localhost:%s/api/v1", serverPort)
	if err := router.Run(":" + serverPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// setupRoutes configures all application routes
func setupRoutes(router *gin.Engine, authHandler *authHandler.AuthHandler) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Terra Server is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/login/google", authHandler.LoginWithGoogle)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/forgot-password", authHandler.ForgotPassword)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		// Program routes
		// programs := v1.Group("/programs")
		// {
		// 	programs.POST("", programHandler.CreateProgram)
		// 	programs.GET("", programHandler.GetAllPrograms)
		// 	programs.GET("/:id", programHandler.GetProgramByID)
		// 	programs.PUT("/:id", programHandler.UpdateProgram)
		// 	programs.DELETE("/:id", programHandler.DeleteProgram)
		// }

		// - Events
		// - Articles
		// - Donations
		// - Volunteers
	}
}