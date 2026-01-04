// Package task provides an interface to schedule background recurring tasks
package task

import (
	"context"
	"errors"
	"time"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"go.uber.org/zap"
)

var (
	IntervalOnce     = time.Duration(0)
	ErrTaskExists    = errors.New("task already exists")
	ErrTaskNotExists = errors.New("task doesn't exist")
)

// Init intializes the global task manager instance
func Init(repo repository.Repository) error {
	manager, err := newManager(repo)
	if err != nil {
		return err
	}

	Manager = manager

	return nil
}

// Task is the interface to which a task should adhere to
// You can manually implement all methods or make use of the `NewTask` function
// which will automatically add some logging
type Task interface {
	// UID is an unique identifier for a check
	// History is kept by linking the UID's of tasks
	// Changing the UID will make you lose all the task history
	// Changing the frontend name can be done with the Name() function
	UID() string
	// Name is an user friendly task name
	// You can change this as much as you like
	Name() string
	// Interval returns the time between executions.
	// An interval == IntervalOnce means it will only run once.
	Interval() time.Duration
	// Hidden determines if the task is returned when Tasks is
	// called on the manager.
	Hidden() bool
	// The function that actually gets executed when it's time
	// The user slice contains all users for who the task needs to executed
	// In reality this will either be a single user (if the user started the task from the api)
	// or contain all users if it's a regular interval run
	// It's up to the function to decide how to handle it
	// If the returned task result does not contain one of the users that was given as argument
	// then the task result is not saved for that user
	Func() func(context.Context, []model.User) []TaskResult
	Ctx() context.Context
}

// TaskResult is the expected return from the actual task function
type TaskResult struct {
	User    model.User
	Message string
	Error   error
}

type Status string

const (
	Waiting Status = "waiting"
	Running Status = "running"
)

// Stat contains the information about a current running or scheduled task
type Stat struct {
	TaskUID   string
	Name      string
	Status    Status
	NextRun   time.Time
	LastRun   time.Time
	Interval  time.Duration
	Recurring bool
}

type internalTask struct {
	uid      string
	name     string
	interval time.Duration
	hidden   bool
	fn       func(context.Context, []model.User) []TaskResult
	ctx      context.Context
}

// Interface compliance
var _ Task = (*internalTask)(nil)

// NewTask creates a new task
// It supports an optional context, if none is given the background context is used
// Logs (info level) when a task starts and ends
// Logs (error level) any error that occurs during the task execution
func NewTask(uid, name string, interval time.Duration, hidden bool, fn func(context.Context, []model.User) []TaskResult, ctx ...context.Context) Task {
	c := context.Background()
	if len(ctx) > 0 {
		c = ctx[0]
	}

	return &internalTask{
		uid:      uid,
		name:     name,
		interval: interval,
		hidden:   hidden,
		fn:       fn,
		ctx:      c,
	}
}

func (t *internalTask) UID() string {
	return t.uid
}

func (t *internalTask) Name() string {
	return t.name
}

func (t *internalTask) Interval() time.Duration {
	return t.interval
}

func (t *internalTask) Hidden() bool {
	return t.hidden
}

func (t *internalTask) Func() func(context.Context, []model.User) []TaskResult {
	return func(ctx context.Context, users []model.User) []TaskResult {
		zap.S().Infof("Task running %s", t.name)

		results := t.fn(ctx, users)
		for _, result := range results {
			if result.Error != nil {
				zap.S().Errorf("Task %s failed for user %+v | %+v", t.name, result.User, result)
			}
		}

		zap.S().Infof("Task finished %s", t.name)

		return results
	}
}

func (t *internalTask) Ctx() context.Context {
	return t.ctx
}
