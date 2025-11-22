package main

import (
	"log"

	articleHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/article/handler"
	articleRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/article/repository"
	articleService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/article/service"
	authHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/handler"
	authRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/repository"
	authService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/service"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/config"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/database"
	donationHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/donation/handler"
	donationRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/donation/repository"
	donationService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/donation/service"
	eventHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/event/handler"
	eventRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/event/repository"
	eventService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/event/service"
	programHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/handler"
	programRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/repository"
	programService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/service"
	volunteerHandler "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/volunteer/handler"
	volunteerRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/volunteer/repository"
	volunteerService "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/volunteer/service"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/email"
	googleOAuth "github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/google"
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
	programRepository := programRepo.NewProgramRepository(db)
	articleRepository := articleRepo.NewArticleRepository(db)
	eventRepository := eventRepo.NewEventRepository(db)
	volunteerRepository := volunteerRepo.NewVolunteerRepository(db)
	donationRepository := donationRepo.NewDonationRepository(db)

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
	programSvc := programService.NewProgramService(programRepository)
	articleSvc := articleService.NewArticleService(articleRepository)
	eventSvc := eventService.NewEventService(eventRepository)
	volunteerSvc := volunteerService.NewVolunteerService(volunteerRepository)
	donationSvc := donationService.NewDonationService(donationRepository, programRepository)

	// Initialize OAuth2 configuration
	oauth2Config := googleOAuth.NewOAuth2Config(
		config.AppConfig.GoogleClientID,
		config.AppConfig.GoogleClientSecret,
		config.AppConfig.GoogleRedirectURL,
	)

	// Initialize handlers
	authHdl := authHandler.NewAuthHandler(authSvc)
	oauth2Hdl := authHandler.NewOAuth2Handler(authSvc, oauth2Config)
	programHdl := programHandler.NewProgramHandler(programSvc)
	articleHdl := articleHandler.NewArticleHandler(articleSvc)
	eventHdl := eventHandler.NewEventHandler(eventSvc)
	volunteerHdl := volunteerHandler.NewVolunteerHandler(volunteerSvc)
	donationHdl := donationHandler.NewDonationHandler(donationSvc)

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	setupRoutes(router, authHdl, oauth2Hdl, programHdl, articleHdl, eventHdl, volunteerHdl, donationHdl)

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
func setupRoutes(router *gin.Engine, authHandler *authHandler.AuthHandler, oauth2Handler *authHandler.OAuth2Handler, programHandler *programHandler.ProgramHandler, articleHandler *articleHandler.ArticleHandler, eventHandler *eventHandler.EventHandler, volunteerHandler *volunteerHandler.VolunteerHandler, donationHandler *donationHandler.DonationHandler) {
	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Terra Server is running",
		})
	})

	// OAuth2 routes (outside /api/v1 for cleaner URLs)
	auth := router.Group("/auth")
	{
		auth.GET("/google/login", oauth2Handler.HandleGoogleLogin)
		auth.GET("/google/callback", oauth2Handler.HandleGoogleCallback)
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Authentication routes
		authV1 := v1.Group("/auth")
		{
			authV1.POST("/register", authHandler.Register)
			authV1.POST("/login", authHandler.Login)
			authV1.POST("/login/google", authHandler.LoginWithGoogle)
			authV1.POST("/refresh", authHandler.RefreshToken)
			authV1.POST("/logout", authHandler.Logout)
			authV1.POST("/forgot-password", authHandler.ForgotPassword)
			authV1.POST("/reset-password", authHandler.ResetPassword)
		}

		// Program routes
		programs := v1.Group("/programs")
		{
			programs.POST("", programHandler.CreateProgram)
			programs.GET("", programHandler.GetAllPrograms)
			programs.GET("/:id", programHandler.GetProgramByID)
			programs.PUT("/:id", programHandler.UpdateProgram)
			programs.DELETE("/:id", programHandler.DeleteProgram)
		}

		// Article routes
		articles := v1.Group("/articles")
		{
			articles.POST("", articleHandler.CreateArticle)
			articles.GET("", articleHandler.GetAllArticles)
			articles.GET("/:slug", articleHandler.GetArticleBySlug)
			articles.PUT("/:slug", articleHandler.UpdateArticle)
			articles.DELETE("/:slug", articleHandler.DeleteArticle)
		}

		// Event routes
		events := v1.Group("/events")
		{
			events.POST("", eventHandler.CreateEvent)
			events.GET("", eventHandler.GetAllEvents)
			events.GET("/:id", eventHandler.GetEventByID)
			events.PUT("/:id", eventHandler.UpdateEvent)
			events.DELETE("/:id", eventHandler.DeleteEvent)
			events.POST("/:id/join", eventHandler.JoinEvent)
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/my-events", eventHandler.GetMyEvents)
			users.GET("/my-donations", donationHandler.GetMyDonations)
		}

		// Volunteer routes
		volunteers := v1.Group("/volunteers")
		{
			volunteers.POST("/apply", volunteerHandler.SubmitApplication)
			volunteers.GET("", volunteerHandler.GetAllApplications)
			volunteers.GET("/:id", volunteerHandler.GetApplicationByID)
			volunteers.PUT("/:id/status", volunteerHandler.UpdateApplicationStatus)
		}

		// Donation routes
		donations := v1.Group("/donations")
		{
			donations.POST("", donationHandler.CreateDonation)
			donations.GET("", donationHandler.GetAllDonations)
			donations.GET("/:id", donationHandler.GetDonationByID)
			donations.PUT("/:id/status", donationHandler.UpdateDonationStatus)
		}
	}
}
