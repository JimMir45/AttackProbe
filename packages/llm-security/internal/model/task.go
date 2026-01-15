package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	TaskStatusPending   = 0 // 待执行
	TaskStatusRunning   = 1 // 执行中
	TaskStatusCompleted = 2 // 已完成
	TaskStatusCancelled = 3 // 已取消
	TaskStatusFailed    = 4 // 失败
)

type Task struct {
	ID             int64          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string         `json:"name" gorm:"size:200;not null"`
	TargetID       int64          `json:"target_id" gorm:"not null;index"`
	Status         int            `json:"status" gorm:"default:0"`
	TotalCount     int            `json:"total_count" gorm:"default:0"`
	CompletedCount int            `json:"completed_count" gorm:"default:0"`
	SuccessCount   int            `json:"success_count" gorm:"default:0"`
	FailedCount    int            `json:"failed_count" gorm:"default:0"`
	ErrorCount     int            `json:"error_count" gorm:"default:0"`
	StartedAt      *time.Time     `json:"started_at"`
	FinishedAt     *time.Time     `json:"finished_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Task) TableName() string {
	return "task"
}

func (t *Task) GetStatusName() string {
	switch t.Status {
	case TaskStatusPending:
		return "待执行"
	case TaskStatusRunning:
		return "执行中"
	case TaskStatusCompleted:
		return "已完成"
	case TaskStatusCancelled:
		return "已取消"
	case TaskStatusFailed:
		return "失败"
	default:
		return "未知"
	}
}

func (t *Task) GetProgress() int {
	if t.TotalCount == 0 {
		return 0
	}
	return t.CompletedCount * 100 / t.TotalCount
}

// TaskTestCase 任务-用例关联表
type TaskTestCase struct {
	ID         int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID     int64     `json:"task_id" gorm:"not null;index"`
	TestCaseID int64     `json:"testcase_id" gorm:"not null"`
	CreatedAt  time.Time `json:"created_at"`
}

func (TaskTestCase) TableName() string {
	return "task_testcase"
}
