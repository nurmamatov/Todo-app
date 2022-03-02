package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "two_services/api/api/handlers/model"
	pbtask "two_services/api/genproto/task"
	l "two_services/api/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateTaskAs godoc
// @Summary Creates new task
// @Description This method create new task
// @Security BearerAuth
// @Tags Task
// @Accept json
// @Produce json
// @Param todo body task.CreateTaskReq true "New Task"
// @Success 201 {object} task.TaskRes
// @Failure 400 {object} task.ErrOrStatus
// @Failure 500 {object} task.ErrOrStatus
// @Router /tasks [POST]
func (h *HandlerV1) CreateTask(c *gin.Context) {
	var (
		body        pbtask.CreateTaskReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TaskService().Create(ctx, &body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create task", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetTask godoc
// @Summary Get Task
// @Description This method Get task task
// @Security BearerAuth
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 201 {object} task.TaskRes
// @Failure 400 {object} task.ErrOrStatus
// @Failure 500 {object} task.ErrOrStatus
// @Router /tasks/{id} [GET]
func (h *HandlerV1) GetTask(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	fmt.Println(guid)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.TaskService().Get(
		ctx, &pbtask.GetAndDeleteTask{
			Id: guid,
		})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get task", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateTask godoc
// @Summary Update Task
// @Description This method Update task
// @Security BearerAuth
// @Tags Task
// @Accept json
// @Produce json
// @Param todo body task.UpdateTaskReq true "Update Task"
// @Success 201 {object} task.TaskRes
// @Failure 400 {object} task.ErrOrStatus
// @Failure 404 {object} task.ErrOrStatus
// @Failure 500 {object} task.ErrOrStatus
// @Router /tasks [PUT]
func (h *HandlerV1) UpdateTask(c *gin.Context) {
	var (
		body        pbtask.UpdateTaskReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TaskService().Update(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListTasks godoc
// @Summary List Tasks
// @Description This method Update task
// @Security BearerAuth
// @Tags Task
// @Accept json
// @Produce json
// @Param page query string true "List Tasks"
// @Param limit query string true "List Tasks"
// @Success 201 {object} task.TasksList
// @Failure 400 {object} task.ErrOrStatus
// @Failure 404 {object} task.ErrOrStatus
// @Failure 500 {object} task.ErrOrStatus
// @Router /tasks [GET]
func (h *HandlerV1) ListTasks(c *gin.Context) {
	var (
		body        pbtask.LimAndPage
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	page, err := strconv.Atoi(c.Query("page"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Page = int64(page)

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Limit = int64(limit)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TaskService().List(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListTasksOverdue godoc
// @Summary List Tasks Overdue
// @Description This method List Overdue
// @Security BearerAuth
// @Tags Task
// @Accept json
// @Produce json
// @Param todo query string true "List Tasks Overdue"
// @Success 201 {object} task.TasksList
// @Failure 400 {object} task.ErrOrStatus
// @Failure 404 {object} task.ErrOrStatus
// @Failure 500 {object} task.ErrOrStatus
// @Router /tasksoverdue [GET]
func (h *HandlerV1) ListTasksOverdue(c *gin.Context) {
	var (
		body        pbtask.LimAndPage
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Page = int64(page)

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Limit = int64(limit)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TaskService().ListOverdue(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteTask godoc
// @Summary Delete Task
// @Description This method Delete task
// @Security BearerAuth
// @Tags Task
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 201 {object} task.ErrOrStatus
// @Failure 400 {object} task.ErrOrStatus
// @Failure 404 {object} task.ErrOrStatus
// @Failure 500 {object} task.ErrOrStatus
// @Router /tasks/{id} [PUT]
func (h *HandlerV1) DeleteTask(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.TaskService().Delete(
		ctx, &pbtask.GetAndDeleteTask{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete task", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
