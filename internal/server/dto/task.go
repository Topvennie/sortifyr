package dto

import (
	"time"

	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/task"
)

type TaskHistory struct {
	ID       int              `json:"id"`
	Name     string           `json:"name"`
	Result   model.TaskResult `json:"result"`
	RunAt    time.Time        `json:"run_at"`
	Message  string           `json:"message"`
	Error    string           `json:"error,omitempty"`
	Duration time.Duration    `json:"duration"`
}

func TaskHistoryDTO(task *model.Task) TaskHistory {
	taskError := ""
	if task.Error != nil {
		taskError = task.Error.Error()
	}

	return TaskHistory{
		ID:       task.ID,
		Name:     task.Name,
		Result:   task.Result,
		RunAt:    task.RunAt,
		Message:  task.Message,
		Error:    taskError,
		Duration: task.Duration,
	}
}

type Task struct {
	TaskUID     string           `json:"uid"`
	Name        string           `json:"name"`
	Status      task.Status      `json:"status"`
	NextRun     time.Time        `json:"next_run"`
	LastStatus  model.TaskResult `json:"last_status,omitempty"`
	LastRun     *time.Time       `json:"last_run,omitzero"`
	LastMessage string           `json:"last_message,omitempty"`
	LastError   string           `json:"last_error,omitempty"`
	Interval    *time.Duration   `json:"interval,omitzero"`
}

func TaskDTO(task task.Stat) Task {
	return Task{
		TaskUID:  task.TaskUID,
		Name:     task.Name,
		Status:   task.Status,
		NextRun:  task.NextRun,
		LastRun:  &task.LastRun,
		Interval: &task.Interval,
	}
}

type TaskFilter struct {
	UserID  int
	TaskUID string
	Result  *model.TaskResult
	Limit   int
	Offset  int
}

func (t *TaskFilter) ToModel() *model.TaskFilter {
	taskFilter := model.TaskFilter(*t)

	return &taskFilter
}
