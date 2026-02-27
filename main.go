package main

import (
	"fmt"
	"hyperscal-go/config"
	"hyperscal-go/internal/controller"
	"hyperscal-go/internal/repository"
	"hyperscal-go/internal/repository/oracle"
	"hyperscal-go/internal/repository/postgres"
	"hyperscal-go/internal/service"
	"hyperscal-go/pkg/database"
	kafkapkg "hyperscal-go/pkg/kafka"
	"hyperscal-go/pkg/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repository based on database driver
	var countryRepo repository.CountryRepository
	var cityRepo repository.CityRepository
	var userRepo repository.UserRepository

	switch cfg.Database.Driver {
	case "postgres":
		countryRepo = postgres.NewCountryPostgresRepository(db)
		cityRepo = postgres.NewCityPostgresRepository(db)
		userRepo = postgres.NewUserPostgresRepository(db)
		log.Println("Using PostgreSQL repository")
	case "oracle":
		countryRepo = oracle.NewCountryOracleRepository(db)
		log.Println("Using Oracle repository")
	default:
		log.Fatalf("Unsupported database driver: %s", cfg.Database.Driver)
	}

	// Initialize service
	countryService := service.NewCountryService(countryRepo)

	// Initialize controller
	countryController := controller.NewCountryController(countryService)

	// Initialize city service and controller
	cityService := service.NewCityService(cityRepo)
	cityController := controller.NewCityController(cityService)

	// Inizialize user service and controller
	authService := service.NewAuthService(userRepo)
	authController := controller.NewAuthController(authService)

	// Initialize Kafka client and controller
	kafkaClient := kafkapkg.NewKafkaClient(&cfg.Kafka)
	kafkaService, err := service.NewKafkaService(kafkaClient)
	if err != nil {
		log.Printf("Warning: Kafka unavailable, kafka endpoints will return errors: %v", err)
	}
	kafkaController := controller.NewKafkaController(kafkaService)

	// Setup Gin router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes (no authentication required)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
		}

		// Protected routes (require JWT authentication)
		protected := api.Group("")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			countries := protected.Group("/countries")
			{
				countries.POST("", countryController.CreateCountry)
				countries.GET("", countryController.GetAllCountries)
				countries.GET("/:id", countryController.GetCountryByID)
				countries.PUT("/:id", countryController.UpdateCountry)
				countries.DELETE("/:id", countryController.DeleteCountry)
			}

			cities := protected.Group("/cities")
			{
				cities.POST("", cityController.CreateCity)
				cities.GET("", cityController.GetAllCities)
				cities.GET("search", cityController.SearchCities)
			}

			kafka := protected.Group("/kafka")
			{
				kafka.POST("/user-created", kafkaController.PublishUserCreated)
				kafka.POST("/order-placed", kafkaController.PublishOrderPlaced)
			}
		}

	}

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on port %s with %s database...", cfg.Server.Port, cfg.Database.Driver)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
