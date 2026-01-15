package repository

import (
	"llm-security-bas/internal/model"
	"llm-security-bas/internal/request"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) Update(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) Delete(id int64) error {
	return r.db.Delete(&model.Task{}, id).Error
}

func (r *TaskRepository) FindByID(id int64) (*model.Task, error) {
	var task model.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) FindPage(query *request.TaskQuery, offset, limit int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64

	db := r.db.Model(&model.Task{})

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.TargetID > 0 {
		db = db.Where("target_id = ?", query.TargetID)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id DESC").Offset(offset).Limit(limit).Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *TaskRepository) UpdateStatus(id int64, status int) error {
	return r.db.Model(&model.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TaskRepository) UpdateProgress(id int64, completed, success, failed, errCount int) error {
	return r.db.Model(&model.Task{}).Where("id = ?", id).Updates(map[string]interface{}{
		"completed_count": completed,
		"success_count":   success,
		"failed_count":    failed,
		"error_count":     errCount,
	}).Error
}

func (r *TaskRepository) AddTestCases(taskID int64, testcaseIDs []int64) error {
	for _, tcID := range testcaseIDs {
		ttc := model.TaskTestCase{
			TaskID:     taskID,
			TestCaseID: tcID,
		}
		if err := r.db.Create(&ttc).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskRepository) GetTestCaseIDs(taskID int64) ([]int64, error) {
	var ids []int64
	err := r.db.Model(&model.TaskTestCase{}).Where("task_id = ?", taskID).Pluck("testcase_id", &ids).Error
	return ids, err
}

func (r *TaskRepository) HasRunningTask(targetID int64) (bool, error) {
	var count int64
	err := r.db.Model(&model.Task{}).Where("target_id = ? AND status = ?", targetID, model.TaskStatusRunning).Count(&count).Error
	return count > 0, err
}
