package task

import (
	"context"
	"fmt"
	"slices"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/pkg/config"
	"github.com/topvennie/sortifyr/pkg/utils"
	"go.uber.org/zap"
)

// Manager is the global single task manager instance
var Manager *manager

type job struct {
	task     model.Task
	status   Status
	interval time.Duration

	users []model.User // If it's not empty then an user triggered it and is waiting on it
}

// Manager can be used to schedule  recurring tasks in the background
// It keeps logs inside the database.
// However it does not automatically reshedule tasks after an application reboot
type manager struct {
	scheduler gocron.Scheduler
	repoTask  repository.Task
	repoUser  repository.User

	mu   sync.Mutex
	jobs map[string]job

	isDev bool
}

func newManager(repo repository.Repository) (*manager, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("failed to create new scheduler %w", err)
	}

	scheduler.Start()

	manager := &manager{
		scheduler: scheduler,
		repoTask:  *repo.NewTask(),
		repoUser:  *repo.NewUser(),
		jobs:      make(map[string]job),
		isDev:     config.IsDev(),
	}

	if err := manager.repoTask.SetInactiveAll(context.Background()); err != nil {
		return nil, err
	}

	return manager, nil
}

// Add adds a new task to the manager
// An unique uid is required.
// if you change a task's uid then all it's history will be lost (but still in the DB)
// Recurring tasks (defined by the interval != IntervalOnce) will be schedules according to the interval
// If the interval is production then a recurring task will immediately be run when added.
// Non recurring tasks will immediately be executed.
func (m *manager) Add(ctx context.Context, newTask Task) error {
	zap.S().Infof("Adding task: %s", newTask.Name())

	if _, ok := m.jobs[newTask.UID()]; ok {
		return ErrTaskExists
	}

	isRecurring := newTask.Interval() != IntervalOnce

	task, err := m.repoTask.GetByUID(ctx, newTask.UID())
	if err != nil {
		return err
	}
	if task != nil {
		// Pre-existing task
		// Update it
		task.Name = newTask.Name()
		task.Active = isRecurring
		task.Recurring = isRecurring
		if err := m.repoTask.Update(ctx, *task); err != nil {
			return err
		}
	} else {
		// New task
		// Let's create it
		task = &model.Task{
			UID:       newTask.UID(),
			Name:      newTask.Name(),
			Active:    isRecurring,
			Recurring: isRecurring,
		}
		if err := m.repoTask.Create(ctx, *task); err != nil {
			return err
		}
	}

	// We lock it early so that jobs can't run immediately until we release the lock.
	// We only release is once we add it to the map.
	m.mu.Lock()
	defer m.mu.Unlock()

	options := []gocron.JobOption{
		gocron.WithName(task.UID),
		gocron.WithContext(newTask.Ctx()),
		gocron.WithTags(task.UID),
	}
	if isRecurring && !m.isDev {
		// Only start tasks immediately in production environment
		// This will only impact recurring tasks
		// Tasks run once will always run immediately regardless of the environment
		options = append(options, gocron.WithStartAt(gocron.WithStartImmediately()))
	}

	var def gocron.JobDefinition
	if isRecurring {
		def = gocron.DurationJob(newTask.Interval())
	} else {
		def = gocron.OneTimeJob(gocron.OneTimeJobStartImmediately())
	}

	if _, err := m.scheduler.NewJob(
		def,
		gocron.NewTask(m.wrap(newTask)),
		options...,
	); err != nil {
		return fmt.Errorf("failed to add task %+v | %w", *task, err)
	}

	status := Waiting
	if !isRecurring {
		status = Running
	}

	m.jobs[task.UID] = job{
		task:     *task,
		status:   status,
		interval: newTask.Interval(),
		users:    []model.User{},
	}

	return nil
}

// RunRecurringByUID runs a pre existing recurring task given a task UID.
func (m *manager) RunRecurringByUID(taskUID string, user model.User) error {
	var job gocron.Job
	for _, j := range m.scheduler.Jobs() {
		if taskUID == j.Tags()[0] {
			job = j
			break
		}
	}
	if job == nil {
		return fmt.Errorf("task with uid %s not found", taskUID)
	}

	// Set the user argument
	m.mu.Lock()
	info, ok := m.jobs[taskUID]
	if !ok {
		m.mu.Unlock()
		return fmt.Errorf("task with uid %s not found", taskUID)
	}
	if info.interval == IntervalOnce {
		// It's a one time task
		// Not allowed!
		m.mu.Unlock()
		return fmt.Errorf("task with uid %s is a one time task", taskUID)
	}
	info.users = append(info.users, user)
	m.jobs[taskUID] = info
	m.mu.Unlock()

	if err := job.RunNow(); err != nil {
		return fmt.Errorf("failed to run task with uid %s | %w", taskUID, err)
	}

	return nil
}

// Tasks returns all scheduled tasks
func (m *manager) Tasks() ([]Stat, error) {
	m.mu.Lock()
	jobsGocron := m.scheduler.Jobs()
	jobsLocal := m.jobs
	m.mu.Unlock()

	stats := make([]Stat, 0, len(jobsGocron))

	for _, job := range jobsGocron {
		taskUID := job.Tags()[0]

		nextRun, err := job.NextRun()
		if err != nil {
			return nil, fmt.Errorf("get next run for task %s | %w", job.Name(), err)
		}

		if j, ok := jobsLocal[taskUID]; ok {
			lastRun, err := job.LastRun()
			if err != nil {
				return nil, fmt.Errorf("get last run for task %s | %w", job.Name(), err)
			}

			stats = append(stats, Stat{
				TaskUID:   j.task.UID,
				Name:      j.task.Name,
				Status:    j.status,
				NextRun:   nextRun,
				LastRun:   lastRun,
				Interval:  j.interval,
				Recurring: j.interval != IntervalOnce,
			})
		}
	}

	slices.SortFunc(stats, func(a, b Stat) int { return int(a.NextRun.Sub(b.NextRun).Nanoseconds()) })

	return stats, nil
}

func (m *manager) wrap(task Task) func(context.Context) {
	return func(ctx context.Context) {
		isRecurring := task.Interval() != IntervalOnce

		m.mu.Lock()
		info, ok := m.jobs[task.UID()]
		if !ok {
			// Should not be possible
			m.mu.Unlock()
			zap.S().Errorf("Task %s not found during execution", task.Name())
			return
		}

		var users []model.User
		if len(info.users) > 0 {
			users = info.users
			info.users = []model.User{}
		}
		info.status = Running

		m.jobs[task.UID()] = info
		m.mu.Unlock()

		if len(users) == 0 {
			// It's a generic interval run
			// Add all real users
			usersDB, err := m.repoUser.GetActualAll(ctx)
			if err != nil {
				zap.S().Error(err)
				return
			}
			users = utils.SliceDereference(usersDB)
		}

		// Run task
		start := time.Now()
		results := task.Func()(ctx, users)
		end := time.Now()

		// Save result
		for _, result := range results {
			taskResult := model.TaskSuccess
			if result.Error != nil {
				taskResult = model.TaskFailed
			}

			taskDB := &model.Task{
				UID:      task.UID(),
				UserID:   result.User.ID,
				RunAt:    start,
				Result:   taskResult,
				Message:  result.Message,
				Error:    result.Error,
				Duration: end.Sub(start),
			}

			if errDB := m.repoTask.CreateRun(ctx, taskDB); errDB != nil {
				zap.S().Errorf("Failed to save recurring task result in database %+v | %v", *taskDB, errDB)
			}
		}

		m.mu.Lock()
		defer m.mu.Unlock()

		if isRecurring {
			info = m.jobs[task.UID()]
			info.status = Waiting
			m.jobs[task.UID()] = info
		} else {
			delete(m.jobs, task.UID())
			m.scheduler.RemoveByTags(task.UID())
		}
	}
}
