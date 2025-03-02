package models

type Priority string
type Status string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
)

const (
	StatusBacklog  Status = "backlog"
	StatusActive   Status = "active"
	StatusFinished Status = "finished"
)

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    Priority  `json:"priority"`
	Status      Status    `json:"status"` 
	StartDate   string    `json:"startDate"`
	DueDate     *string   `json:"dueDate"`
	DueTime     *string   `json:"dueTime"`
	SubTasks    []Task    `json:"subTasks"`
} 