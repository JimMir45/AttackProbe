package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	TargetTypeOpenAI = 1 // OpenAI兼容API
	TargetTypeRAG    = 2 // RAG应用
	TargetTypeAgent  = 3 // Agent系统
)

const (
	TargetStatusEnabled  = 1
	TargetStatusDisabled = 0
)

type Target struct {
	ID             int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string         `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Type           int            `json:"type" gorm:"not null;default:1"`
	Endpoint       string         `json:"endpoint" gorm:"size:500;not null"`
	APIKey         string         `json:"-" gorm:"column:api_key;size:500"`
	Model          string         `json:"model" gorm:"size:100"`
	ExtraHeaders   string         `json:"extra_headers" gorm:"type:text"`
	Timeout        int            `json:"timeout" gorm:"default:30000"`
	Status         int            `json:"status" gorm:"default:1"`
	LastTestTime   *time.Time     `json:"last_test_time"`
	LastTestStatus *int           `json:"last_test_status"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Target) TableName() string {
	return "target"
}

func (t *Target) GetTypeName() string {
	switch t.Type {
	case TargetTypeOpenAI:
		return "OpenAI兼容"
	case TargetTypeRAG:
		return "RAG应用"
	case TargetTypeAgent:
		return "Agent系统"
	default:
		return "未知"
	}
}
