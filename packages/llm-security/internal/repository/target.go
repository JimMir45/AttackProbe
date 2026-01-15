package repository

import (
	"llm-security-bas/internal/model"
	"llm-security-bas/internal/request"

	"gorm.io/gorm"
)

type TargetRepository struct {
	db *gorm.DB
}

func NewTargetRepository(db *gorm.DB) *TargetRepository {
	return &TargetRepository{db: db}
}

func (r *TargetRepository) Create(target *model.Target) error {
	return r.db.Create(target).Error
}

func (r *TargetRepository) Update(target *model.Target) error {
	return r.db.Save(target).Error
}

func (r *TargetRepository) Delete(id int64) error {
	return r.db.Delete(&model.Target{}, id).Error
}

func (r *TargetRepository) FindByID(id int64) (*model.Target, error) {
	var target model.Target
	err := r.db.First(&target, id).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *TargetRepository) FindByName(name string) (*model.Target, error) {
	var target model.Target
	err := r.db.Where("name = ?", name).First(&target).Error
	if err != nil {
		return nil, err
	}
	return &target, nil
}

func (r *TargetRepository) FindPage(query *request.TargetQuery, offset, limit int) ([]model.Target, int64, error) {
	var targets []model.Target
	var total int64

	db := r.db.Model(&model.Target{})

	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Type > 0 {
		db = db.Where("type = ?", query.Type)
	}
	if query.Status > 0 {
		db = db.Where("status = ?", query.Status)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id DESC").Offset(offset).Limit(limit).Find(&targets).Error
	if err != nil {
		return nil, 0, err
	}

	return targets, total, nil
}

func (r *TargetRepository) FindAll() ([]model.Target, error) {
	var targets []model.Target
	err := r.db.Where("status = 1").Order("id DESC").Find(&targets).Error
	return targets, err
}

func (r *TargetRepository) UpdateTestStatus(id int64, status int) error {
	return r.db.Model(&model.Target{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_test_time":   gorm.Expr("datetime('now')"),
		"last_test_status": status,
	}).Error
}
