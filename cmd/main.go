package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"mytasks/internal/config"
	"mytasks/internal/handlers"
	"mytasks/internal/repository"
	"mytasks/internal/services"
)

func main() {
	// Initialize database connection
	config.InitDB()
	defer config.CloseDB()

	// Initialize repository
	repo := repository.NewTaskRepository()

	// Initialize service with repository
	taskService := services.NewTaskService(repo)

	// Initialize handler with service
	taskHandler := handlers.NewTaskHandler(taskService)

	// Initialize Gin router
	r := gin.Default()

	// Add logging middleware
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))

	r.Use(func(c *gin.Context) {
		log.Printf("Incoming request: %s %s", c.Request.Method, c.Request.URL.Path)
		log.Printf("Origin: %s", c.Request.Header.Get("Origin"))
		c.Next()
	})

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"https://mytasks-project.vercel.app",
		"https://kanban-front-nu.vercel.app",
		"https://barsuc.ru",
	}
	log.Printf("Configured CORS with allowed origins: %v", config.AllowOrigins)
	
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"Authorization",
		"X-Requested-With",
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Headers",
	}
	config.ExposeHeaders = []string{"Content-Length", "Content-Type"}
	config.AllowCredentials = false
	r.Use(cors.New(config))

	// Add detailed CORS logging middleware
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		log.Printf("Request from origin: %s", origin)
		log.Printf("Request method: %s", c.Request.Method)
		log.Printf("Request headers: %v", c.Request.Header)
		c.Next()
	})

	// Add health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Routes
	api := r.Group("/api")
	{
		tasks := api.Group("/tasks")
		{
			tasks.GET("", taskHandler.GetTasks)
			tasks.POST("", taskHandler.CreateTask)
			tasks.GET("/:id", taskHandler.GetTask)
			tasks.PUT("/:id", taskHandler.UpdateTask)
			tasks.DELETE("/:id", taskHandler.DeleteTask)
			tasks.GET("/:id/subtasks", taskHandler.GetSubtasks)
			tasks.POST("/:id/subtasks", taskHandler.CreateSubtask)
		}
	}

	// Get port from environment variable for Railway
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Fatal(r.Run(":" + port))
} 