package main

import (
	"log"

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

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://127.0.0.1:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Accept",
		"X-Requested-With",
	}
	config.ExposeHeaders = []string{"Content-Length"}
	r.Use(cors.New(config))

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

	// Start server
	log.Fatal(r.Run(":8080"))
} 