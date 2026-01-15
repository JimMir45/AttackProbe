package repository

import (
	"llm-security-bas/internal/model"
	"llm-security-bas/internal/request"

	"gorm.io/gorm"
)

type TestCaseRepository struct {
	db *gorm.DB
}

func NewTestCaseRepository(db *gorm.DB) *TestCaseRepository {
	return &TestCaseRepository{db: db}
}

func (r *TestCaseRepository) Create(tc *model.TestCase) error {
	return r.db.Create(tc).Error
}

func (r *TestCaseRepository) Update(tc *model.TestCase) error {
	return r.db.Save(tc).Error
}

func (r *TestCaseRepository) Delete(id int64) error {
	return r.db.Delete(&model.TestCase{}, id).Error
}

func (r *TestCaseRepository) FindByID(id int64) (*model.TestCase, error) {
	var tc model.TestCase
	err := r.db.First(&tc, id).Error
	if err != nil {
		return nil, err
	}
	return &tc, nil
}

func (r *TestCaseRepository) FindPage(req *request.TestCasePageRequest, offset, limit int) ([]model.TestCase, int64, error) {
	var testcases []model.TestCase
	var total int64

	db := r.db.Model(&model.TestCase{})

	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Category > 0 {
		db = db.Where("category = ?", req.Category)
	}
	if req.RiskLevel > 0 {
		db = db.Where("risk_level = ?", req.RiskLevel)
	}
	if req.IsBuiltin != nil {
		db = db.Where("is_builtin = ?", *req.IsBuiltin)
	}
	if req.Status > 0 {
		db = db.Where("status = ?", req.Status)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Order("id DESC").Offset(offset).Limit(limit).Find(&testcases).Error
	if err != nil {
		return nil, 0, err
	}

	return testcases, total, nil
}

func (r *TestCaseRepository) FindByIDs(ids []int64) ([]model.TestCase, error) {
	var testcases []model.TestCase
	err := r.db.Where("id IN ?", ids).Find(&testcases).Error
	return testcases, err
}

func (r *TestCaseRepository) FindAllEnabled() ([]model.TestCase, error) {
	var testcases []model.TestCase
	err := r.db.Where("status = 1").Order("id ASC").Find(&testcases).Error
	return testcases, err
}

func (r *TestCaseRepository) BatchUpdateStatus(ids []int64, status int) error {
	return r.db.Model(&model.TestCase{}).Where("id IN ? AND is_builtin = 0", ids).Update("status", status).Error
}

func (r *TestCaseRepository) GetStats() (map[string]interface{}, error) {
	var total, enabled, builtin int64
	var categoryStats []struct {
		Category int   `json:"category"`
		Count    int64 `json:"count"`
	}
	var riskStats []struct {
		RiskLevel int   `json:"risk_level"`
		Count     int64 `json:"count"`
	}

	r.db.Model(&model.TestCase{}).Count(&total)
	r.db.Model(&model.TestCase{}).Where("status = 1").Count(&enabled)
	r.db.Model(&model.TestCase{}).Where("is_builtin = 1").Count(&builtin)
	r.db.Model(&model.TestCase{}).Select("category, count(*) as count").Group("category").Scan(&categoryStats)
	r.db.Model(&model.TestCase{}).Select("risk_level, count(*) as count").Group("risk_level").Scan(&riskStats)

	return map[string]interface{}{
		"total":          total,
		"enabled":        enabled,
		"builtin":        builtin,
		"category_stats": categoryStats,
		"risk_stats":     riskStats,
	}, nil
}
