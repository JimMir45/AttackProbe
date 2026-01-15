package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"llm-security-bas/internal/request"
	"llm-security-bas/internal/response"
	"llm-security-bas/internal/service"
)

type TaskController struct {
	svc *service.TaskService
}

func NewTaskController(db *gorm.DB) *TaskController {
	return &TaskController{
		svc: service.NewTaskService(db),
	}
}

func (ctrl *TaskController) Add(c *gin.Context) {
	var req request.TaskAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	id, err := ctrl.svc.Create(&req)
	if err != nil {
		response.FailWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, response.AddedID{ID: id})
}

func (ctrl *TaskController) Delete(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.Delete(req.ID)
	if err != nil {
		if err.Error() == "任务正在执行中，无法删除" {
			response.Fail(c, response.CodeTaskRunning)
			return
		}
		response.FailWithMsg(c, response.CodeTaskNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *TaskController) Detail(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	task, err := ctrl.svc.FindByID(req.ID)
	if err != nil {
		response.Fail(c, response.CodeTaskNotFound)
		return
	}

	data := map[string]interface{}{
		"id":              task.ID,
		"name":            task.Name,
		"target_id":       task.TargetID,
		"status":          task.Status,
		"status_name":     task.GetStatusName(),
		"total_count":     task.TotalCount,
		"completed_count": task.CompletedCount,
		"success_count":   task.SuccessCount,
		"failed_count":    task.FailedCount,
		"error_count":     task.ErrorCount,
		"progress":        task.GetProgress(),
		"started_at":      task.StartedAt,
		"finished_at":     task.FinishedAt,
		"created_at":      task.CreatedAt,
		"updated_at":      task.UpdatedAt,
	}

	response.Success(c, data)
}

func (ctrl *TaskController) Page(c *gin.Context) {
	var req request.TaskPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	tasks, total, err := ctrl.svc.FindPage(&req)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	content := make([]map[string]interface{}, len(tasks))
	for i, t := range tasks {
		content[i] = map[string]interface{}{
			"id":              t.ID,
			"name":            t.Name,
			"target_id":       t.TargetID,
			"status":          t.Status,
			"status_name":     t.GetStatusName(),
			"total_count":     t.TotalCount,
			"completed_count": t.CompletedCount,
			"success_count":   t.SuccessCount,
			"failed_count":    t.FailedCount,
			"error_count":     t.ErrorCount,
			"progress":        t.GetProgress(),
			"started_at":      t.StartedAt,
			"finished_at":     t.FinishedAt,
			"created_at":      t.CreatedAt,
		}
	}

	response.Success(c, response.PageResult{
		List:  content,
		Total: total,
	})
}

func (ctrl *TaskController) Start(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.Start(req.ID)
	if err != nil {
		if err.Error() == "任务正在执行中" {
			response.Fail(c, response.CodeTaskRunning)
			return
		}
		response.FailWithMsg(c, response.CodeTaskNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *TaskController) Cancel(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.Cancel(req.ID)
	if err != nil {
		if err.Error() == "任务未在执行" {
			response.Fail(c, response.CodeTaskNotRunning)
			return
		}
		response.FailWithMsg(c, response.CodeTaskNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *TaskController) Progress(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	progress, err := ctrl.svc.GetProgress(req.ID)
	if err != nil {
		response.Fail(c, response.CodeTaskNotFound)
		return
	}

	response.Success(c, progress)
}

func (ctrl *TaskController) Results(c *gin.Context) {
	var req request.TaskResultPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	results, total, err := ctrl.svc.GetResults(&req)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	response.Success(c, response.PageResult{
		List:  results,
		Total: total,
	})
}
