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
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/middleware"
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

// setupRoutes configures all application routes with proper authentication and authorization
func setupRoutes(router *gin.Engine, authHandler *authHandler.AuthHandler, oauth2Handler *authHandler.OAuth2Handler, programHandler *programHandler.ProgramHandler, articleHandler *articleHandler.ArticleHandler, eventHandler *eventHandler.EventHandler, volunteerHandler *volunteerHandler.VolunteerHandler, donationHandler *donationHandler.DonationHandler) {
	// Health check endpoint (Public)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Terra Server is running",
		})
	})

	// OAuth2 routes (Public - outside /api/v1 for cleaner URLs)
	auth := router.Group("/auth")
	{
		auth.GET("/google/login", oauth2Handler.HandleGoogleLogin)
		auth.GET("/google/callback", oauth2Handler.HandleGoogleCallback)
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// ========== PUBLIC ROUTES (No Authentication Required) ==========

		// Public Authentication routes
		authPublic := v1.Group("/auth")
		{
			authPublic.POST("/register", authHandler.Register)
			authPublic.POST("/login", authHandler.Login)
			authPublic.POST("/login/google", authHandler.LoginWithGoogle)
			authPublic.POST("/forgot-password", authHandler.ForgotPassword)
			authPublic.POST("/reset-password", authHandler.ResetPassword)
		}

		// Public Program routes
		programsPublic := v1.Group("/programs")
		{
			programsPublic.GET("", programHandler.GetAllPrograms)
			programsPublic.GET("/:id", programHandler.GetProgramByID)
		}

		// Public Article routes
		articlesPublic := v1.Group("/articles")
		{
			articlesPublic.GET("", articleHandler.GetAllArticles)
			articlesPublic.GET("/:slug", articleHandler.GetArticleBySlug)
		}

		// Public Event routes
		eventsPublic := v1.Group("/events")
		{
			eventsPublic.GET("", eventHandler.GetAllEvents)
			eventsPublic.GET("/:id", eventHandler.GetEventByID)
		}

		// ========== AUTHENTICATED USER ROUTES (Auth Required) ==========

		authRequired := v1.Group("")
		authRequired.Use(middleware.AuthMiddleware())
		{
			// Authenticated Auth routes
			authAuth := authRequired.Group("/auth")
			{
				authAuth.POST("/refresh", authHandler.RefreshToken)
				authAuth.POST("/logout", authHandler.Logout)
			}

			// User-specific routes
			users := authRequired.Group("/users")
			{
				users.GET("/my-events", eventHandler.GetMyEvents)
				users.GET("/my-donations", donationHandler.GetMyDonations)
			}

			// Event participation (User)
			eventsAuth := authRequired.Group("/events")
			{
				eventsAuth.POST("/:id/join", eventHandler.JoinEvent)
			}

			// Volunteer application (User)
			volunteersAuth := authRequired.Group("/volunteers")
			{
				volunteersAuth.POST("/apply", volunteerHandler.SubmitApplication)
			}

			// Donation creation (User)
			donationsAuth := authRequired.Group("/donations")
			{
				donationsAuth.POST("", donationHandler.CreateDonation)
			}
		}

		// ========== ADMIN ROUTES (Auth + Admin Role Required) ==========

		adminRequired := v1.Group("")
		adminRequired.Use(middleware.AuthMiddleware())
		adminRequired.Use(middleware.AdminMiddleware())
		{
			// Admin Program routes
			programsAdmin := adminRequired.Group("/programs")
			{
				programsAdmin.POST("", programHandler.CreateProgram)
				programsAdmin.PUT("/:id", programHandler.UpdateProgram)
				programsAdmin.DELETE("/:id", programHandler.DeleteProgram)
			}

			// Admin Article routes
			articlesAdmin := adminRequired.Group("/articles")
			{
				articlesAdmin.POST("", articleHandler.CreateArticle)
				articlesAdmin.PUT("/:slug", articleHandler.UpdateArticle)
				articlesAdmin.DELETE("/:slug", articleHandler.DeleteArticle)
			}

			// Admin Event routes
			eventsAdmin := adminRequired.Group("/events")
			{
				eventsAdmin.POST("", eventHandler.CreateEvent)
				eventsAdmin.PUT("/:id", eventHandler.UpdateEvent)
				eventsAdmin.DELETE("/:id", eventHandler.DeleteEvent)
			}

			// Admin Volunteer routes
			volunteersAdmin := adminRequired.Group("/volunteers")
			{
				volunteersAdmin.GET("", volunteerHandler.GetAllApplications)
				volunteersAdmin.GET("/:id", volunteerHandler.GetApplicationByID)
				volunteersAdmin.PUT("/:id/status", volunteerHandler.UpdateApplicationStatus)
			}

			// Admin Donation routes
			donationsAdmin := adminRequired.Group("/donations")
			{
				donationsAdmin.GET("", donationHandler.GetAllDonations)
				donationsAdmin.GET("/:id", donationHandler.GetDonationByID)
				donationsAdmin.PUT("/:id/status", donationHandler.UpdateDonationStatus)
			}
		}
	}
}
