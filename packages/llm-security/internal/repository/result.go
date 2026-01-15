package repository

import (
	"llm-security-bas/internal/model"

	"gorm.io/gorm"
)

type ResultRepository struct {
	db *gorm.DB
}

func NewResultRepository(db *gorm.DB) *ResultRepository {
	return &ResultRepository{db: db}
}

func (r *ResultRepository) Create(result *model.TaskResult) error {
	return r.db.Create(result).Error
}

func (r *ResultRepository) BatchCreate(results []model.TaskResult) error {
	return r.db.CreateInBatches(results, 100).Error
}

func (r *ResultRepository) Update(result *model.TaskResult) error {
	return r.db.Save(result).Error
}

func (r *ResultRepository) FindByID(id int64) (*model.TaskResult, error) {
	var result model.TaskResult
	err := r.db.First(&result, id).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ResultRepository) FindByTaskID(taskID int64, offset, limit int) ([]model.TaskResult, int64, error) {
	var results []model.TaskResult
	var total int64

	db := r.db.Model(&model.TaskResult{}).Where("task_id = ?", taskID)

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id ASC").Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}

func (r *ResultRepository) FindPendingByTaskID(taskID int64) ([]model.TaskResult, error) {
	var results []model.TaskResult
	err := r.db.Where("task_id = ? AND status = ?", taskID, model.ResultStatusPending).Find(&results).Error
	return results, err
}

func (r *ResultRepository) GetStatsByTaskID(taskID int64) (map[string]int64, error) {
	var stats []struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}

	err := r.db.Model(&model.TaskResult{}).
		Select("status, count(*) as count").
		Where("task_id = ?", taskID).
		Group("status").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	result := map[string]int64{
		"total":   0,
		"pending": 0,
		"success": 0,
		"failed":  0,
		"error":   0,
	}

	for _, s := range stats {
		result["total"] += s.Count
		switch s.Status {
		case model.ResultStatusPending:
			result["pending"] = s.Count
		case model.ResultStatusSuccess:
			result["success"] = s.Count
		case model.ResultStatusFailed:
			result["failed"] = s.Count
		case model.ResultStatusError:
			result["error"] = s.Count
		}
	}

	return result, nil
}

func (r *ResultRepository) DeleteByTaskID(taskID int64) error {
	return r.db.Where("task_id = ?", taskID).Delete(&model.TaskResult{}).Error
}
