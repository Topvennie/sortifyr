package model

import (
	"errors"
	"time"

	"github.com/topvennie/spotify_organizer/pkg/sqlc"
)

type TaskResult string

const (
	TaskSuccess TaskResult = "success"
	TaskFailed  TaskResult = "failed"
)

type Task struct {
	// Task result
	ID       int // ID of the task result
	UserID   int // ID of the user that started the task
	RunAt    time.Time
	Result   TaskResult
	Message  string
	Error    error
	Duration time.Duration

	// Task fields
	UID    string // Identifier of the task
	Name   string
	Active bool
}

func TaskModel(task sqlc.Task, taskRun sqlc.TaskRun) *Task {
	message := ""
	if taskRun.Message.Valid {
		message = taskRun.Message.String
	}
	var err error
	if taskRun.Error.Valid {
		err = errors.New(taskRun.Error.String)
	}

	return &Task{
		ID:       int(taskRun.ID),
		UserID:   int(taskRun.UserID),
		RunAt:    taskRun.RunAt.Time,
		Result:   TaskResult(taskRun.Result),
		Message:  message,
		Error:    err,
		Duration: time.Duration(taskRun.Duration),
		UID:      task.Uid,
		Name:     task.Name,
		Active:   task.Active,
	}
}

type TaskFilter struct {
	UserID  int
	TaskUID string
	Result  *TaskResult
	Limit   int
	Offset  int
}
