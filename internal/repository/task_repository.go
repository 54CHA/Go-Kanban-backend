package repository

import (
	"context"
	"fmt"
	"log"

	"mytasks/internal/config"
	"mytasks/internal/models"
)

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) GetTasks() []models.Task {
	query := `
		SELECT id, title, description, priority, status, start_date, due_date, due_time 
		FROM tasks 
		WHERE parent_id IS NULL
	`
	rows, err := config.DB.Query(context.Background(), query)
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		return []models.Task{}
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate, dueTime *string
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.StartDate,
			&dueDate,
			&dueTime,
		)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}

		if dueDate != nil {
			task.DueDate = dueDate
		}
		if dueTime != nil {
			task.DueTime = dueTime
		}

		// Get subtasks
		task.SubTasks = r.getSubtasksForTask(task.ID)
		tasks = append(tasks, task)
	}

	return tasks
}

func (r *TaskRepository) getSubtasksForTask(parentID string) []models.Task {
	query := `
		SELECT id, title, description, priority, status, start_date, due_date, due_time 
		FROM tasks 
		WHERE parent_id = $1
	`
	rows, err := config.DB.Query(context.Background(), query, parentID)
	if err != nil {
		log.Printf("Error querying subtasks: %v", err)
		return []models.Task{}
	}
	defer rows.Close()

	var subtasks []models.Task
	for rows.Next() {
		var task models.Task
		var dueDate, dueTime *string
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&task.StartDate,
			&dueDate,
			&dueTime,
		)
		if err != nil {
			log.Printf("Error scanning subtask: %v", err)
			continue
		}

		if dueDate != nil {
			task.DueDate = dueDate
		}
		if dueTime != nil {
			task.DueTime = dueTime
		}

		// Recursively get nested subtasks
		task.SubTasks = r.getSubtasksForTask(task.ID)
		subtasks = append(subtasks, task)
	}

	return subtasks
}

func (r *TaskRepository) GetTask(id string) (*models.Task, error) {
	query := `
		SELECT id, title, description, priority, status, start_date, due_date, due_time 
		FROM tasks 
		WHERE id = $1
	`
	var task models.Task
	var dueDate, dueTime *string
	err := config.DB.QueryRow(context.Background(), query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&task.StartDate,
		&dueDate,
		&dueTime,
	)
	if err != nil {
		return nil, fmt.Errorf("task not found: %v", err)
	}

	if dueDate != nil {
		task.DueDate = dueDate
	}
	if dueTime != nil {
		task.DueTime = dueTime
	}

	task.SubTasks = r.getSubtasksForTask(task.ID)
	return &task, nil
}

func (r *TaskRepository) CreateTask(task models.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, priority, status, start_date, due_date, due_time, parent_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULL)
	`
	_, err := config.DB.Exec(context.Background(), query,
		task.ID,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.StartDate,
		task.DueDate,
		task.DueTime,
	)
	if err != nil {
		return fmt.Errorf("error creating task: %v", err)
	}
	return nil
}

func (r *TaskRepository) UpdateTask(task models.Task) error {
	query := `
		UPDATE tasks 
		SET title = $1, description = $2, priority = $3, status = $4, 
			start_date = $5, due_date = $6, due_time = $7
		WHERE id = $8
	`
	result, err := config.DB.Exec(context.Background(), query,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.StartDate,
		task.DueDate,
		task.DueTime,
		task.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating task: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) DeleteTask(id string) error {
	query := "DELETE FROM tasks WHERE id = $1"
	result, err := config.DB.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("error deleting task: %v", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("task not found")
	}

	return nil
}

func (r *TaskRepository) CreateSubtask(parentID string, subtask models.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, priority, status, start_date, due_date, due_time, parent_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := config.DB.Exec(context.Background(), query,
		subtask.ID,
		subtask.Title,
		subtask.Description,
		subtask.Priority,
		subtask.Status,
		subtask.StartDate,
		subtask.DueDate,
		subtask.DueTime,
		parentID,
	)
	if err != nil {
		return fmt.Errorf("error creating subtask: %v", err)
	}
	return nil
}

func (r *TaskRepository) GetSubtasks(id string) ([]models.Task, error) {
	subtasks := r.getSubtasksForTask(id)
	return subtasks, nil
} 