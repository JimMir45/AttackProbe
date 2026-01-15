package service

import (
	"errors"

	"llm-security-bas/internal/model"
	"llm-security-bas/internal/repository"
	"llm-security-bas/internal/request"
	"llm-security-bas/pkg/llm"

	"gorm.io/gorm"
)

type TargetService struct {
	repo *repository.TargetRepository
}

func NewTargetService(db *gorm.DB) *TargetService {
	return &TargetService{
		repo: repository.NewTargetRepository(db),
	}
}

func (s *TargetService) Create(req *request.TargetAddRequest) (int64, error) {
	// 检查名称是否已存在
	existing, _ := s.repo.FindByName(req.Name)
	if existing != nil {
		return 0, errors.New("目标名称已存在")
	}

	// API传入的是秒，内部存储毫秒
	timeout := req.Timeout * 1000
	if timeout <= 0 {
		timeout = 30000 // 默认30秒
	}

	target := &model.Target{
		Name:         req.Name,
		Type:         req.Type,
		Endpoint:     req.Endpoint,
		APIKey:       req.APIKey,
		Model:        req.Model,
		ExtraHeaders: req.ExtraHeaders,
		Timeout:      timeout,
		Status:       model.TargetStatusEnabled,
	}

	err := s.repo.Create(target)
	if err != nil {
		return 0, err
	}

	return target.ID, nil
}

func (s *TargetService) Update(req *request.TargetUpdateRequest) error {
	target, err := s.repo.FindByID(req.ID)
	if err != nil {
		return errors.New("目标不存在")
	}

	// 检查名称是否重复
	if req.Name != target.Name {
		existing, _ := s.repo.FindByName(req.Name)
		if existing != nil {
			return errors.New("目标名称已存在")
		}
	}

	target.Name = req.Name
	target.Type = req.Type
	target.Endpoint = req.Endpoint
	target.Model = req.Model
	target.ExtraHeaders = req.ExtraHeaders
	target.Status = req.Status

	// API传入的是秒，内部存储毫秒
	if req.Timeout > 0 {
		target.Timeout = req.Timeout * 1000
	}

	// 只有传入了APIKey才更新
	if req.APIKey != "" {
		target.APIKey = req.APIKey
	}

	return s.repo.Update(target)
}

func (s *TargetService) Delete(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("目标不存在")
	}

	return s.repo.Delete(id)
}

func (s *TargetService) FindByID(id int64) (*model.Target, error) {
	return s.repo.FindByID(id)
}

func (s *TargetService) FindPage(req *request.TargetPageRequest) ([]model.Target, int64, error) {
	req.Page.SetDefault()
	return s.repo.FindPage(&req.Query, req.Page.GetOffset(), req.Page.Size)
}

func (s *TargetService) FindAllOptions() ([]map[string]interface{}, error) {
	targets, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	options := make([]map[string]interface{}, len(targets))
	for i, t := range targets {
		options[i] = map[string]interface{}{
			"id":   t.ID,
			"name": t.Name,
			"type": t.Type,
		}
	}
	return options, nil
}

func (s *TargetService) Test(id int64) (bool, string, error) {
	target, err := s.repo.FindByID(id)
	if err != nil {
		return false, "", errors.New("目标不存在")
	}

	client := llm.NewOpenAIClient(&llm.ClientConfig{
		Endpoint:     target.Endpoint,
		APIKey:       target.APIKey,
		Model:        target.Model,
		Timeout:      target.Timeout,
		ExtraHeaders: target.ExtraHeaders,
	})

	err = client.Test()
	status := 1
	msg := "连接成功"
	if err != nil {
		status = 0
		msg = err.Error()
	}

	s.repo.UpdateTestStatus(id, status)

	return status == 1, msg, nil
}

// MaskAPIKey 脱敏显示API密钥
func MaskAPIKey(apiKey string) string {
	if len(apiKey) <= 8 {
		return "****"
	}
	return apiKey[:4] + "****" + apiKey[len(apiKey)-4:]
}
