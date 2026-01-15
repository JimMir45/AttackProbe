package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"llm-security-bas/internal/request"
	"llm-security-bas/internal/response"
	"llm-security-bas/internal/service"
)

type TestCaseController struct {
	svc *service.TestCaseService
}

func NewTestCaseController(db *gorm.DB) *TestCaseController {
	return &TestCaseController{
		svc: service.NewTestCaseService(db),
	}
}

func (ctrl *TestCaseController) Add(c *gin.Context) {
	var req request.TestCaseAddRequest
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

func (ctrl *TestCaseController) Update(c *gin.Context) {
	var req request.TestCaseUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.Update(&req)
	if err != nil {
		if err.Error() == "内置用例不允许修改" {
			response.Fail(c, response.CodeTestCaseBuiltin)
			return
		}
		response.FailWithMsg(c, response.CodeTestCaseNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *TestCaseController) Delete(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.Delete(req.ID)
	if err != nil {
		if err.Error() == "内置用例不允许删除" {
			response.Fail(c, response.CodeTestCaseBuiltin)
			return
		}
		response.FailWithMsg(c, response.CodeTestCaseNotFound, err.Error())
		return
	}

	response.Success(c, nil)
}

func (ctrl *TestCaseController) Detail(c *gin.Context) {
	var req request.IDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	tc, err := ctrl.svc.FindByID(req.ID)
	if err != nil {
		response.Fail(c, response.CodeTestCaseNotFound)
		return
	}

	data := map[string]interface{}{
		"id":                tc.ID,
		"name":              tc.Name,
		"category":          tc.Category,
		"category_name":     tc.GetCategoryName(),
		"risk_level":        tc.RiskLevel,
		"risk_level_name":   tc.GetRiskLevelName(),
		"attack_type":       tc.AttackType,
		"content":           tc.Content,
		"system_prompt":     tc.SystemPrompt,
		"expected_behavior": tc.ExpectedBehavior,
		"judge_method":      tc.JudgeMethod,
		"judge_config":      tc.JudgeConfig,
		"source":            tc.Source,
		"reference":         tc.Reference,
		"is_builtin":        tc.IsBuiltin,
		"status":            tc.Status,
		"created_at":        tc.CreatedAt,
		"updated_at":        tc.UpdatedAt,
	}

	response.Success(c, data)
}

func (ctrl *TestCaseController) Page(c *gin.Context) {
	var req request.TestCasePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	testcases, total, err := ctrl.svc.FindPage(&req)
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	content := make([]map[string]interface{}, len(testcases))
	for i, tc := range testcases {
		content[i] = map[string]interface{}{
			"id":              tc.ID,
			"name":            tc.Name,
			"category":        tc.Category,
			"category_name":   tc.GetCategoryName(),
			"risk_level":      tc.RiskLevel,
			"risk_level_name": tc.GetRiskLevelName(),
			"attack_type":     tc.AttackType,
			"is_builtin":      tc.IsBuiltin,
			"status":          tc.Status,
			"created_at":      tc.CreatedAt,
		}
	}

	response.Success(c, response.PageResult{
		List:  content,
		Total: total,
	})
}

func (ctrl *TestCaseController) Stats(c *gin.Context) {
	stats, err := ctrl.svc.GetStats()
	if err != nil {
		response.Fail(c, response.CodeInternalError)
		return
	}

	response.Success(c, stats)
}

func (ctrl *TestCaseController) BatchStatus(c *gin.Context) {
	var req request.BatchStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	err := ctrl.svc.BatchUpdateStatus(&req)
	if err != nil {
		response.FailWithMsg(c, response.CodeInternalError, err.Error())
		return
	}

	response.Success(c, nil)
}
