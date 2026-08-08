package main

import (
	"log"

	"github.com/DevJuloGPHI/go-vue-auth/backend/auth"
	"github.com/DevJuloGPHI/go-vue-auth/backend/config"
	"github.com/DevJuloGPHI/go-vue-auth/backend/database"
	"github.com/DevJuloGPHI/go-vue-auth/backend/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Load environment configuration
	cfg := config.Load()

	// 2. Connect PostgreSQL
	db, err := database.Connect(cfg)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 3. Run migration
	if err := db.AutoMigrate(&user.User{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 4. Create repository
	userRepository := user.NewRepository(db)

	// 5. Create service
	userService := user.NewService(userRepository)

	// 6. Create HTTP handler
	authHandler := auth.NewHandler(userService)

	// 7. Create Gin server
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
	}))

	// 8. Health endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is running",
		})
	})

	// 9. API routes
	api := router.Group("/api/v1")

	auth.RegisterRoutes(api, authHandler)

	// 10. Start server
	log.Printf("Server running on http://localhost:%s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
