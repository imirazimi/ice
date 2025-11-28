package http

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"ice/internal/port"
	"ice/internal/todo"
	"ice/pkg/errors"
	"ice/pkg/logger"
	"ice/pkg/validator"

	"go.uber.org/zap"
)

type TodoHandler struct {
	service   port.TodoService
	validator *validator.Validator
}

func NewTodoHandler(s port.TodoService) *TodoHandler {
	return &TodoHandler{
		service:   s,
		validator: validator.New(),
	}
}

// CreateTodo creates a new todo item
// @Summary Create a new todo item
// @Description Create a new todo item and publish it to Redis Stream
// @Tags todos
// @Accept json
// @Produce json
// @Param request body todo.CreateTodoRequest true "Todo creation request"
// @Success 201 {object} todo.CreateTodoResponse
// @Failure 400 {object} errors.AppError
// @Failure 500 {object} errors.AppError
// @Router /todo [post]
func (h *TodoHandler) CreateTodo(c echo.Context) error {
	log := logger.Get()

	var req todo.CreateTodoRequest
	if err := c.Bind(&req); err != nil {
		log.Warn("Invalid request body", zap.Error(err))
		appErr := errors.NewBadRequestError("invalid request body", err)
		return c.JSON(appErr.Code, appErr)
	}

	// Validate request
	if err := h.validator.Validate(&req); err != nil {
		log.Warn("Validation failed", zap.Error(err), zap.Any("request", req))
		appErr := errors.NewValidationError(err.Error())
		return c.JSON(appErr.Code, appErr)
	}

	item := &todo.TodoItem{
		ID:          uuid.New().String(),
		Description: req.Description,
		DueDate:     req.DueDate,
	}

	if err := h.service.CreateTodo(c.Request().Context(), item); err != nil {
		log.Error("Failed to create todo", zap.Error(err), zap.String("todo_id", item.ID))
		appErr := errors.NewInternalError("failed to create todo", err)
		return c.JSON(appErr.Code, appErr)
	}

	log.Info("Todo created successfully", zap.String("todo_id", item.ID))

	return c.JSON(201, todo.CreateTodoResponse{
		TodoItem: *item,
	})
}
