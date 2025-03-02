package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mytasks/internal/models"
	"mytasks/internal/services"
	"log"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks := h.service.GetTasks()
	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := h.service.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	log.Printf("Starting task creation request")
	
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received task creation request: %+v", task)

	if err := h.service.CreateTask(&task); err != nil {
		log.Printf("Error creating task: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Successfully created task with ID: %s", task.ID)
	
	response := gin.H{
		"status": "success",
		"task": task,
	}
	log.Printf("Sending response: %+v", response)
	c.JSON(http.StatusCreated, response)
	log.Printf("Response sent")
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = id
	if err := h.service.UpdateTask(task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTask(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) GetSubtasks(c *gin.Context) {
	id := c.Param("id")
	subtasks, err := h.service.GetSubtasks(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, subtasks)
}

func (h *TaskHandler) CreateSubtask(c *gin.Context) {
	parentID := c.Param("id")
	var subtask models.Task
	if err := c.ShouldBindJSON(&subtask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateSubtask(parentID, &subtask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, subtask)
} 