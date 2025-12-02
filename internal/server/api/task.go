package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/topvennie/spotify_organizer/internal/database/model"
	"github.com/topvennie/spotify_organizer/internal/server/dto"
	"github.com/topvennie/spotify_organizer/internal/server/service"
)

type Task struct {
	router fiber.Router

	task service.Task
}

func NewTask(router fiber.Router, service service.Service) *Task {
	api := &Task{
		router: router.Group("/task"),
		task:   *service.NewTask(),
	}

	api.createRoutes()

	return api
}

func (r *Task) createRoutes() {
	r.router.Get("/", r.getTasks)
	r.router.Get("/history", r.getHistory)
	r.router.Post("/start/:uid", r.start)
}

func (r *Task) getTasks(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	tasks, err := r.task.GetTasks(c.Context(), userID)
	if err != nil {
		return err
	}

	return c.JSON(tasks)
}

func (r *Task) getHistory(c *fiber.Ctx) error {
	uid := c.Query("uid")
	resultStr := c.Query("result")

	var result *model.TaskResult
	switch resultStr {
	case string(model.TaskSuccess), string(model.TaskFailed):
		resultTmp := model.TaskResult(resultStr)
		result = &resultTmp
	}

	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 0)
	if limit < 1 || page < 0 {
		return fiber.ErrBadRequest
	}

	tasks, err := r.task.GetHistory(c.Context(), dto.TaskFilter{
		TaskUID: uid,
		Result:  result,
		Limit:   limit,
		Offset:  page * limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(tasks)
}

func (r *Task) start(c *fiber.Ctx) error {
	uid := c.Params("uid")

	userID, ok := c.Locals("userID").(int)
	if !ok {
		return fiber.ErrUnauthorized
	}

	if err := r.task.Start(c.Context(), userID, uid); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusAccepted)
}
