package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/pkg/sqlc"
	"github.com/topvennie/sortifyr/pkg/utils"
)

type Task struct {
	repo Repository
}

func (r *Repository) NewTask() *Task {
	return &Task{
		repo: *r,
	}
}

// GetByUID returns the task given an uid without any runs
func (t *Task) GetByUID(ctx context.Context, taskUID string) (*model.Task, error) {
	task, err := t.repo.queries(ctx).TaskGetByUID(ctx, taskUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get task by uid %s | %w", taskUID, err)
	}

	return model.TaskModel(task, sqlc.TaskRun{}), nil
}

func (t *Task) GetByRunID(ctx context.Context, runID int) (*model.Task, error) {
	task, err := t.repo.queries(ctx).TaskRunGet(ctx, int32(runID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get task run %d | %w", runID, err)
	}

	return model.TaskModel(task.Task, task.TaskRun), err
}

func (t *Task) GetRunFiltered(ctx context.Context, filter model.TaskFilter) ([]*model.Task, error) {
	// A default value is required for sql
	// If no result filter is set then it will ignored anyway
	result := model.TaskSuccess
	if filter.Result != nil {
		result = *filter.Result
	}

	recurring := false
	if filter.Recurring != nil {
		recurring = *filter.Recurring
	}

	params := sqlc.TaskGetFilteredParams{
		Column1:         int32(filter.UserID),
		Uid:             filter.TaskUID,
		FilterTaskUid:   filter.TaskUID != "",
		Result:          sqlc.TaskResult(result),
		FilterResult:    filter.Result != nil,
		Recurring:       recurring,
		FilterRecurring: filter.Recurring != nil,
		Limit:           int32(filter.Limit),
		Offset:          int32(filter.Offset),
	}

	tasks, err := t.repo.queries(ctx).TaskGetFiltered(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get filtered tasks %+v | %w", filter, err)
	}

	return utils.SliceMap(tasks, func(task sqlc.TaskGetFilteredRow) *model.Task {
		return model.TaskModel(task.Task, task.TaskRun)
	}), nil
}

func (t *Task) GetRunLastAllByUser(ctx context.Context, userID int) ([]*model.Task, error) {
	tasks, err := t.repo.queries(ctx).TaskGetLastAllByUser(ctx, int32(userID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get last runs by user %d | %w", userID, err)
	}

	return utils.SliceMap(tasks, func(t sqlc.TaskRun) *model.Task { return model.TaskModel(sqlc.Task{}, t) }), nil
}

func (t *Task) Create(ctx context.Context, task model.Task) error {
	if err := t.repo.queries(ctx).TaskCreate(ctx, sqlc.TaskCreateParams{
		Uid:    task.UID,
		Name:   task.Name,
		Active: task.Active,
	}); err != nil {
		return fmt.Errorf("create task %+v | %w", t, err)
	}

	return nil
}

func (t *Task) CreateRun(ctx context.Context, task *model.Task) error {
	errStr := ""
	if task.Error != nil {
		errStr = task.Error.Error()
	}

	id, err := t.repo.queries(ctx).TaskRunCreate(ctx, sqlc.TaskRunCreateParams{
		TaskUid:  task.UID,
		UserID:   int32(task.UserID),
		RunAt:    toTime(task.RunAt),
		Message:  toString(task.Message),
		Result:   sqlc.TaskResult(task.Result),
		Error:    toString(errStr),
		Duration: task.Duration.Nanoseconds(),
	})
	if err != nil {
		return fmt.Errorf("create task run %+v | %w", *task, err)
	}

	task.ID = int(id)

	return nil
}

func (t *Task) Update(ctx context.Context, task model.Task) error {
	if err := t.repo.queries(ctx).TaskUpdate(ctx, sqlc.TaskUpdateParams{
		Uid:       task.UID,
		Name:      toString(task.Name),
		Active:    toBool(&task.Active),
		Recurring: toBool(&task.Recurring),
	}); err != nil {
		return fmt.Errorf("update task %+v | %w", task, err)
	}

	return nil
}

func (t *Task) SetInactiveAll(ctx context.Context) error {
	if err := t.repo.queries(ctx).TaskSetInactiveAll(ctx); err != nil {
		return fmt.Errorf("set recurring tasks to inactive %w", err)
	}

	return nil
}
