package handlers

import (
	"net/http"
	"strconv"

	"github.com/DmitriyGiryntsev/TODO-API/internal/models"
	"github.com/DmitriyGiryntsev/TODO-API/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type TaskHandler struct {
	Repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{Repo: repo}
}

// ErrorResponse структура для ошибок
type ErrorResponse struct {
	Error string `json:"error"`
}

// MessageResponse структура для сообщений
type MessageResponse struct {
	Message string `json:"message"`
}

// GetTasks godoc
// @Summary Получить список задач
// @Description Получает все задачи для текущего пользователя
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Task
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/tasks/ [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	tasks, err := h.Repo.GetAllTasksByUserID(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// GetTask godoc
// @Summary Получить задачу по ID
// @Description Получает задачу по ее ID для текущего пользователя
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid task ID"})
		return
	}

	task, err := h.Repo.GetTaskByID(taskID, userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot get task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// CreateTask godoc
// @Summary Создать новую задачу
// @Description Создает новую задачу для текущего пользователя
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task data"
// @Success 201 {object} models.Task
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/v1/tasks/ [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid task data"})
		return
	}

	if err := validate.Struct(task); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	task.UserID = userID.(int)

	if err := h.Repo.CreateNewTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// UpdateTask godoc
// @Summary Обновить задачу
// @Description Обновляет существующую задачу для текущего пользователя
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task data"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid task ID"})
		return
	}

	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid task data"})
		return
	}

	task.ID = taskID
	task.UserID = userID.(int)

	if err := h.Repo.UpdateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot update task"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "task updated successfully"})
}

// DeleteTask godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по ее ID для текущего пользователя
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "unauthorized"})
		return
	}

	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid task ID"})
		return
	}

	if err := h.Repo.DeleteTask(taskID, userID.(int)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "cannot delete task"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "task deleted successfully"})
}
