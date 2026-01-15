package model

import (
	"time"
)

const (
	ResultStatusPending = 0 // 待执行
	ResultStatusSuccess = 1 // 成功(攻击被阻断)
	ResultStatusFailed  = 2 // 失败(攻击成功)
	ResultStatusError   = 3 // 错误(执行异常)
)

const (
	JudgeResultBlocked = 1 // 攻击被阻断
	JudgeResultSuccess = 0 // 攻击成功
)

type TaskResult struct {
	ID              int64      `json:"id" gorm:"primaryKey;autoIncrement"`
	TaskID          int64      `json:"task_id" gorm:"not null;index"`
	TestCaseID      int64      `json:"testcase_id" gorm:"not null"`
	Status          int        `json:"status" gorm:"default:0"`
	RequestContent  string     `json:"request_content" gorm:"type:text"`
	ResponseContent string     `json:"response_content" gorm:"type:text"`
	JudgeResult     *int       `json:"judge_result"`
	JudgeReason     string     `json:"judge_reason" gorm:"type:text"`
	Duration        int        `json:"duration"` // ms
	ErrorMessage    string     `json:"error_message" gorm:"type:text"`
	ExecutedAt      *time.Time `json:"executed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (TaskResult) TableName() string {
	return "task_result"
}

func (r *TaskResult) GetStatusName() string {
	switch r.Status {
	case ResultStatusPending:
		return "待执行"
	case ResultStatusSuccess:
		return "成功"
	case ResultStatusFailed:
		return "失败"
	case ResultStatusError:
		return "错误"
	default:
		return "未知"
	}
}

func (r *TaskResult) GetJudgeResultName() string {
	if r.JudgeResult == nil {
		return "未判定"
	}
	switch *r.JudgeResult {
	case JudgeResultBlocked:
		return "攻击被阻断"
	case JudgeResultSuccess:
		return "攻击成功"
	default:
		return "未知"
	}
}
