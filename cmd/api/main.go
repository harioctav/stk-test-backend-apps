package main

import (
	"log"
	"time"

	"stk-technical-test-api/internal/config"
	"stk-technical-test-api/internal/database"
	"stk-technical-test-api/internal/handler"
	"stk-technical-test-api/internal/repository"
	"stk-technical-test-api/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.NewDatabase(cfg.GetDSN(), cfg.App.Env == "development")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize dependencies (Dependency Injection)
	menuRepo := repository.NewMenuRepository(db.GetDB())
	menuService := service.NewMenuService(menuRepo)
	menuHandler := handler.NewMenuHandler(menuService)

	// Setup Gin router
	router := setupRouter(menuHandler, cfg)

	// Start server
	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func setupRouter(menuHandler *handler.MenuHandler, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS Configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Server Is Running",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Menu routes
		menus := api.Group("/menus")
		{
			menus.GET("/hierarchy", menuHandler.GetMenuHierarchy)
			menus.GET("/uuid/:uuid", menuHandler.GetMenuByUUID)
			menus.GET("", menuHandler.GetAllMenus)
			menus.GET("/:id", menuHandler.GetMenuByID)
			menus.POST("", menuHandler.CreateMenu)
			menus.PUT("/:id", menuHandler.UpdateMenu)
			menus.DELETE("/:id", menuHandler.DeleteMenu)
		}
	}

	return router
}

