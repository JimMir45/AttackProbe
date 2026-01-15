package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	CategoryPromptInjection = 1 // 提示词注入
	CategoryJailbreak       = 2 // 越狱攻击
	CategorySensitiveData   = 3 // 敏感信息泄露
	CategoryOther           = 4 // 其他
)

const (
	RiskLevelLow    = 1
	RiskLevelMedium = 2
	RiskLevelHigh   = 3
)

const (
	JudgeByKeyword = 1 // 关键词匹配
	JudgeByRegex   = 2 // 正则表达式
	JudgeByLLM     = 3 // LLM判定
)

type TestCase struct {
	ID               int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name             string         `json:"name" gorm:"size:200;not null"`
	Category         int            `json:"category" gorm:"not null"`
	RiskLevel        int            `json:"risk_level" gorm:"default:2"`
	AttackType       string         `json:"attack_type" gorm:"size:50"`
	Content          string         `json:"content" gorm:"type:text;not null"`
	SystemPrompt     string         `json:"system_prompt" gorm:"type:text"`
	ExpectedBehavior string         `json:"expected_behavior" gorm:"type:text"`
	JudgeMethod      int            `json:"judge_method" gorm:"default:1"`
	JudgeConfig      string         `json:"judge_config" gorm:"type:text"`
	Source           string         `json:"source" gorm:"size:100"`
	Reference        string         `json:"reference" gorm:"size:500"`
	IsBuiltin        int            `json:"is_builtin" gorm:"default:0"`
	Status           int            `json:"status" gorm:"default:1"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

func (TestCase) TableName() string {
	return "testcase"
}

func (t *TestCase) GetCategoryName() string {
	switch t.Category {
	case CategoryPromptInjection:
		return "提示词注入"
	case CategoryJailbreak:
		return "越狱攻击"
	case CategorySensitiveData:
		return "敏感信息泄露"
	case CategoryOther:
		return "其他"
	default:
		return "未知"
	}
}

func (t *TestCase) GetRiskLevelName() string {
	switch t.RiskLevel {
	case RiskLevelLow:
		return "低"
	case RiskLevelMedium:
		return "中"
	case RiskLevelHigh:
		return "高"
	default:
		return "未知"
	}
}
