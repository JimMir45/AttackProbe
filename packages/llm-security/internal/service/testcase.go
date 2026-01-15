package service

import (
	"errors"

	"llm-security-bas/internal/model"
	"llm-security-bas/internal/repository"
	"llm-security-bas/internal/request"

	"gorm.io/gorm"
)

type TestCaseService struct {
	repo *repository.TestCaseRepository
}

func NewTestCaseService(db *gorm.DB) *TestCaseService {
	return &TestCaseService{
		repo: repository.NewTestCaseRepository(db),
	}
}

func (s *TestCaseService) Create(req *request.TestCaseAddRequest) (int64, error) {
	tc := &model.TestCase{
		Name:             req.Name,
		Category:         req.Category,
		RiskLevel:        req.RiskLevel,
		AttackType:       req.AttackType,
		Content:          req.Content,
		SystemPrompt:     req.SystemPrompt,
		ExpectedBehavior: req.ExpectedBehavior,
		JudgeMethod:      req.JudgeMethod,
		JudgeConfig:      req.JudgeConfig,
		Source:           req.Source,
		Reference:        req.Reference,
		IsBuiltin:        0,
		Status:           1,
	}

	if tc.RiskLevel <= 0 {
		tc.RiskLevel = model.RiskLevelMedium
	}
	if tc.JudgeMethod <= 0 {
		tc.JudgeMethod = model.JudgeByKeyword
	}
	if tc.Source == "" {
		tc.Source = "custom"
	}

	err := s.repo.Create(tc)
	if err != nil {
		return 0, err
	}

	return tc.ID, nil
}

func (s *TestCaseService) Update(req *request.TestCaseUpdateRequest) error {
	tc, err := s.repo.FindByID(req.ID)
	if err != nil {
		return errors.New("用例不存在")
	}

	if tc.IsBuiltin == 1 {
		return errors.New("内置用例不允许修改")
	}

	tc.Name = req.Name
	tc.Category = req.Category
	tc.RiskLevel = req.RiskLevel
	tc.AttackType = req.AttackType
	tc.Content = req.Content
	tc.SystemPrompt = req.SystemPrompt
	tc.ExpectedBehavior = req.ExpectedBehavior
	tc.JudgeMethod = req.JudgeMethod
	tc.JudgeConfig = req.JudgeConfig
	tc.Reference = req.Reference
	tc.Status = req.Status

	return s.repo.Update(tc)
}

func (s *TestCaseService) Delete(id int64) error {
	tc, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("用例不存在")
	}

	if tc.IsBuiltin == 1 {
		return errors.New("内置用例不允许删除")
	}

	return s.repo.Delete(id)
}

func (s *TestCaseService) FindByID(id int64) (*model.TestCase, error) {
	return s.repo.FindByID(id)
}

func (s *TestCaseService) FindPage(req *request.TestCasePageRequest) ([]model.TestCase, int64, error) {
	req.Page.SetDefault()
	return s.repo.FindPage(req, req.Page.GetOffset(), req.Page.Size)
}

func (s *TestCaseService) FindByIDs(ids []int64) ([]model.TestCase, error) {
	return s.repo.FindByIDs(ids)
}

func (s *TestCaseService) FindAllEnabled() ([]model.TestCase, error) {
	return s.repo.FindAllEnabled()
}

func (s *TestCaseService) BatchUpdateStatus(req *request.BatchStatusRequest) error {
	return s.repo.BatchUpdateStatus(req.IDs, req.Status)
}

func (s *TestCaseService) GetStats() (map[string]interface{}, error) {
	return s.repo.GetStats()
}
