package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/sortifyr/internal/database/model"
	"github.com/topvennie/sortifyr/internal/database/repository"
	"github.com/topvennie/sortifyr/internal/server/dto"
	"github.com/topvennie/sortifyr/internal/task"
	"github.com/topvennie/sortifyr/pkg/utils"
	"go.uber.org/zap"
)

type Task struct {
	service Service

	task repository.Task
	user repository.User
}

func (s *Service) NewTask() *Task {
	return &Task{
		service: *s,
		task:    *s.repo.NewTask(),
		user:    *s.repo.NewUser(),
	}
}

func (t *Task) GetTasks(ctx context.Context, userID int) ([]dto.Task, error) {
	tasks, err := task.Manager.Tasks()
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if tasks == nil {
		return []dto.Task{}, nil
	}

	lastRuns, err := t.task.GetRunLastAllByUser(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	lastRunMap := make(map[string]*model.Task)
	for _, lastRun := range lastRuns {
		lastRunMap[lastRun.UID] = lastRun
	}

	taskDTOs := make([]dto.Task, 0, len(tasks))
	for _, task := range tasks {
		taskDTO := dto.TaskDTO(task)

		if lastRun, ok := lastRunMap[task.TaskUID]; ok {
			lastError := ""
			if lastRun.Error != nil {
				lastError = lastRun.Error.Error()
			}

			taskDTO.LastStatus = lastRun.Result
			taskDTO.LastMessage = lastRun.Message
			taskDTO.LastError = lastError
		}

		taskDTOs = append(taskDTOs, taskDTO)
	}

	return taskDTOs, nil
}

func (t *Task) GetHistory(ctx context.Context, filter dto.TaskFilter) ([]dto.TaskHistory, error) {
	tasks, err := t.task.GetRunFiltered(ctx, *filter.ToModel())
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if tasks == nil {
		return []dto.TaskHistory{}, nil
	}

	return utils.SliceMap(tasks, dto.TaskHistoryDTO), nil
}

func (t *Task) Start(ctx context.Context, userID int, taskUID string) error {
	user, err := t.user.GetByID(ctx, userID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if user == nil {
		return fiber.ErrUnauthorized
	}

	taskModel, err := t.task.GetByUID(ctx, taskUID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if taskModel == nil {
		return fiber.ErrNotFound
	}
	zap.S().Debug(*taskModel)
	if !taskModel.Recurring || !taskModel.Active {
		return fiber.ErrBadRequest
	}

	return task.Manager.RunRecurringByUID(taskUID, *user)
}
