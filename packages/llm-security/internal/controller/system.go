package controller

import (
	"github.com/gin-gonic/gin"

	"llm-security-bas/internal/model"
	"llm-security-bas/internal/response"
)

type SystemController struct{}

func NewSystemController() *SystemController {
	return &SystemController{}
}

func (ctrl *SystemController) Info(c *gin.Context) {
	response.Success(c, map[string]interface{}{
		"name":    "LLM Security BAS",
		"version": "1.0.0",
		"desc":    "大模型安全有效性验证平台",
	})
}

func (ctrl *SystemController) ConfigGet(c *gin.Context) {
	var configs []model.SysConfig
	model.GetDB().Find(&configs)

	data := make(map[string]string)
	for _, cfg := range configs {
		data[cfg.ConfigKey] = cfg.ConfigValue
	}

	response.Success(c, data)
}

func (ctrl *SystemController) ConfigUpdate(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, response.CodeParamError)
		return
	}

	for key, value := range req {
		model.GetDB().Model(&model.SysConfig{}).Where("config_key = ?", key).Update("config_value", value)
	}

	response.Success(c, nil)
}

func (ctrl *SystemController) Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
