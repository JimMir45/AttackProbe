package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"llm-security-bas/internal/request"
	"llm-security-bas/internal/response"
	"llm-security-bas/internal/service"
)

type TargetController struct {
	svc *service.TargetService
}

func NewTargetController(db *gorm.DB) *TargetController {
	return &TargetController{
		svc: service.NewTargetService(db),
	}
}

func (ctrl *TargetController) Add(c *gin.Context) {
	var req request.TargetAddRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	id, err := ctrl.svc.Create(&req)
	if err != nil {
		response.FailWithMsg(c, response.CodeTargetNameExists, err.Error())
		return
	}

	response.Success(c, response.AddedID{ID: id})
}

func (ctrl *TargetController) Update(c *gin.Context) {
	var req request.TargetUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.Update(&req)
	if err != nil {
		response.FailWithMsg(c, response.CodeTargetNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *TargetController) Delete(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.Delete(req.ID)
	if err != nil {
		response.FailWithMsg(c, response.CodeTargetNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *TargetController) Detail(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	target, err := ctrl.svc.FindByID(req.ID)
	if err != nil {
		response.Fail(c, response.CodeTargetNotFound)
		return
	}

	// 脱敏API密钥
	data := map[string]interface{}{
		"id":               target.ID,
		"name":             target.Name,
		"type":             target.Type,
		"type_name":        target.GetTypeName(),
		"endpoint":         target.Endpoint,
		"api_key":          service.MaskAPIKey(target.APIKey),
		"model":            target.Model,
		"extra_headers":    target.ExtraHeaders,
		"timeout":          target.Timeout,
		"status":           target.Status,
		"last_test_time":   target.LastTestTime,
		"last_test_status": target.LastTestStatus,
		"created_at":       target.CreatedAt,
		"updated_at":       target.UpdatedAt,
	}

	response.Success(c, data)
}

func (ctrl *TargetController) Page(c *gin.Context) {
	var req request.TargetPageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	targets, total, err := ctrl.svc.FindPage(&req)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	// 构建返回数据（脱敏）
	content := make([]map[string]interface{}, len(targets))
	for i, t := range targets {
		content[i] = map[string]interface{}{
			"id":               t.ID,
			"name":             t.Name,
			"type":             t.Type,
			"type_name":        t.GetTypeName(),
			"endpoint":         t.Endpoint,
			"model":            t.Model,
			"status":           t.Status,
			"last_test_time":   t.LastTestTime,
			"last_test_status": t.LastTestStatus,
			"created_at":       t.CreatedAt,
		}
	}

	response.Success(c, response.PageResult{
		List:  content,
		Total: total,
	})
}

func (ctrl *TargetController) Options(c *gin.Context) {
	options, err := ctrl.svc.FindAllOptions()
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	response.Success(c, options)
}

func (ctrl *TargetController) Test(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	success, msg, err := ctrl.svc.Test(req.ID)
	if err != nil {
		response.FailWithMsg(c, response.CodeTargetNotFound, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"success": success,
		"message": msg,
	})
}
