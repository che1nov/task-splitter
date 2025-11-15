package http

import (
	"net/http"
	"strconv"
	"task-splitter/internal/dto"
	"task-splitter/internal/usecases"
	"task-splitter/pkg/logger"

	"github.com/gin-gonic/gin"
)

// TaskHandler контроллер для работы с задачами
type TaskHandler struct {
	createTaskUC     *usecases.CreateTaskUseCase
	getTaskUC        *usecases.GetTaskUseCase
	getTasksUC       *usecases.GetTasksUseCase
	updateTaskUC     *usecases.UpdateTaskUseCase
	deleteTaskUC     *usecases.DeleteTaskUseCase
	splitTaskUC      *usecases.SplitTaskUseCase
	getSplitStatusUC *usecases.GetSplitStatusUseCase
	log              logger.Logger
}

// NewTaskHandler создает новый контроллер задач
func NewTaskHandler(
	createTaskUC *usecases.CreateTaskUseCase,
	getTaskUC *usecases.GetTaskUseCase,
	getTasksUC *usecases.GetTasksUseCase,
	updateTaskUC *usecases.UpdateTaskUseCase,
	deleteTaskUC *usecases.DeleteTaskUseCase,
	splitTaskUC *usecases.SplitTaskUseCase,
	getSplitStatusUC *usecases.GetSplitStatusUseCase,
	log logger.Logger,
) *TaskHandler {
	return &TaskHandler{
		createTaskUC:     createTaskUC,
		getTaskUC:        getTaskUC,
		getTasksUC:       getTasksUC,
		updateTaskUC:     updateTaskUC,
		deleteTaskUC:     deleteTaskUC,
		splitTaskUC:      splitTaskUC,
		getSplitStatusUC: getSplitStatusUC,
		log:              log,
	}
}

// GetTasks получает список задач пользователя
// @Summary Get user tasks
// @Description Get list of tasks for current user
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} dto.TaskOutput
// @Failure 401 {object} gin.H
// @Router /tasks [get]
func (h *TaskHandler) GetTasks(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, _ := userIDValue.(uint)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	input := dto.GetTasksByUserInput{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}

	output, err := h.getTasksUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// CreateTask создает новую задачу
// @Summary Create task
// @Description Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTaskInput true "Task data"
// @Success 201 {object} dto.TaskOutput
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, _ := userIDValue.(uint)

	var input dto.CreateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = userID

	output, err := h.createTaskUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, output)
}

// GetTask получает задачу по ID
// @Summary Get task by ID
// @Description Get task details by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 200 {object} dto.TaskOutput
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	output, err := h.getTaskUC.Execute(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, output)
}

// UpdateTask обновляет задачу
// @Summary Update task
// @Description Update task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Param request body dto.UpdateTaskInput true "Task update data"
// @Success 200 {object} dto.TaskOutput
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var input dto.UpdateTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.TaskID = uint(id)

	output, err := h.updateTaskUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}

// DeleteTask удаляет задачу
// @Summary Delete task
// @Description Delete task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Task ID"
// @Success 204
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, _ := userIDValue.(uint)

	idStr := c.Param("id")
	taskID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := h.deleteTaskUC.Execute(c.Request.Context(), uint(taskID), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SplitTask отправляет задачу на разбивку
// @Summary Split task
// @Description Send task for AI-powered splitting into subtasks
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SplitTaskInput true "Split task data"
// @Success 202 {object} dto.SplitTaskOutput
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Router /split [post]
func (h *TaskHandler) SplitTask(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
		return
	}

	userID, _ := userIDValue.(uint)

	var input dto.SplitTaskInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.UserID = userID

	output, err := h.splitTaskUC.Execute(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, output)
}

// GetSplitStatus получает статус разбивки задачи
// @Summary Get split status
// @Description Get status of task splitting request
// @Tags tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} dto.SplitStatusOutput
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /split/{id}/status [get]
func (h *TaskHandler) GetSplitStatus(c *gin.Context) {
	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request ID is required"})
		return
	}

	output, err := h.getSplitStatusUC.Execute(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Request not found"})
		return
	}

	c.JSON(http.StatusOK, output)
}

