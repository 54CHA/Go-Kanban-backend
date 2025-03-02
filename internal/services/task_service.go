package services

import (
	"fmt"
	"time"

	"mytasks/internal/models"
	"mytasks/internal/repository"
)

type TaskService struct {
	repo *repository.TaskRepository
}

func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetTasks() []models.Task {
	return s.repo.GetTasks()
}

func (s *TaskService) GetTask(id string) (*models.Task, error) {
	return s.repo.GetTask(id)
}

func (s *TaskService) CreateTask(task *models.Task) error {
	// Validate task
	if err := s.validateTask(*task); err != nil {
		return err
	}

	// Generate ID for new task
	task.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	// Set default values if not provided
	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}
	if task.Status == "" {
		task.Status = models.StatusBacklog
	}
	if task.StartDate == "" {
		task.StartDate = time.Now().Format("2006-01-02")
	}
	if task.SubTasks == nil {
		task.SubTasks = make([]models.Task, 0)
	}

	return s.repo.CreateTask(*task)
}

func (s *TaskService) UpdateTask(task models.Task) error {
	// Validate task
	if err := s.validateTask(task); err != nil {
		return err
	}

	return s.repo.UpdateTask(task)
}

func (s *TaskService) DeleteTask(id string) error {
	return s.repo.DeleteTask(id)
}

func (s *TaskService) GetSubtasks(id string) ([]models.Task, error) {
	return s.repo.GetSubtasks(id)
}

func (s *TaskService) CreateSubtask(parentID string, subtask *models.Task) error {
	// Validate subtask
	if err := s.validateTask(*subtask); err != nil {
		return err
	}

	// Generate ID for new subtask
	subtask.ID = fmt.Sprintf("%d", time.Now().UnixNano())

	// Set default values if not provided
	if subtask.Priority == "" {
		subtask.Priority = models.PriorityMedium
	}
	if subtask.Status == "" {
		subtask.Status = models.StatusBacklog
	}
	if subtask.StartDate == "" {
		subtask.StartDate = time.Now().Format("2006-01-02")
	}
	if subtask.SubTasks == nil {
		subtask.SubTasks = make([]models.Task, 0)
	}

	return s.repo.CreateSubtask(parentID, *subtask)
}

func (s *TaskService) validateTask(task models.Task) error {
	if task.Title == "" {
		return fmt.Errorf("task title is required")
	}

	if task.Priority != "" && 
		task.Priority != models.PriorityLow && 
		task.Priority != models.PriorityMedium && 
		task.Priority != models.PriorityHigh {
		return fmt.Errorf("invalid priority value")
	}

	if task.Status != "" && 
		task.Status != models.StatusBacklog && 
		task.Status != models.StatusActive && 
		task.Status != models.StatusFinished {
		return fmt.Errorf("invalid status value")
	}

	if task.StartDate != "" {
		if _, err := time.Parse("2006-01-02", task.StartDate); err != nil {
			return fmt.Errorf("invalid start date format")
		}
	}

	if task.DueDate != nil {
		if _, err := time.Parse("2006-01-02", *task.DueDate); err != nil {
			return fmt.Errorf("invalid due date format")
		}
	}

	if task.DueTime != nil {
		if _, err := time.Parse("15:04", *task.DueTime); err != nil {
			return fmt.Errorf("invalid due time format")
		}
	}

	return nil
} 