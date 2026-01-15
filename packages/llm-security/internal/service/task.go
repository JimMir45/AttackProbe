package service

import (
	"errors"
	"time"

	"llm-security-bas/internal/model"
	"llm-security-bas/internal/repository"
	"llm-security-bas/internal/request"
	"llm-security-bas/internal/service/executor"

	"gorm.io/gorm"
)

type TaskService struct {
	db         *gorm.DB
	repo       *repository.TaskRepository
	resultRepo *repository.ResultRepository
	tcRepo     *repository.TestCaseRepository
	targetRepo *repository.TargetRepository
}

func NewTaskService(db *gorm.DB) *TaskService {
	return &TaskService{
		db:         db,
		repo:       repository.NewTaskRepository(db),
		resultRepo: repository.NewResultRepository(db),
		tcRepo:     repository.NewTestCaseRepository(db),
		targetRepo: repository.NewTargetRepository(db),
	}
}

func (s *TaskService) Create(req *request.TaskAddRequest) (int64, error) {
	// 检查目标是否存在
	_, err := s.targetRepo.FindByID(req.TargetID)
	if err != nil {
		return 0, errors.New("目标不存在")
	}

	// 获取测试用例
	var testcases []model.TestCase
	if len(req.TestCaseIDs) > 0 {
		testcases, err = s.tcRepo.FindByIDs(req.TestCaseIDs)
		if err != nil {
			return 0, err
		}
	} else {
		testcases, err = s.tcRepo.FindAllEnabled()
		if err != nil {
			return 0, err
		}
	}

	if len(testcases) == 0 {
		return 0, errors.New("没有可用的测试用例")
	}

	// 创建任务
	task := &model.Task{
		Name:       req.Name,
		TargetID:   req.TargetID,
		Status:     model.TaskStatusPending,
		TotalCount: len(testcases),
	}

	err = s.repo.Create(task)
	if err != nil {
		return 0, err
	}

	// 关联测试用例
	tcIDs := make([]int64, len(testcases))
	for i, tc := range testcases {
		tcIDs[i] = tc.ID
	}
	err = s.repo.AddTestCases(task.ID, tcIDs)
	if err != nil {
		return 0, err
	}

	// 创建结果记录（待执行状态）
	results := make([]model.TaskResult, len(testcases))
	for i, tc := range testcases {
		results[i] = model.TaskResult{
			TaskID:     task.ID,
			TestCaseID: tc.ID,
			Status:     model.ResultStatusPending,
		}
	}
	err = s.resultRepo.BatchCreate(results)
	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (s *TaskService) Delete(id int64) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("任务不存在")
	}

	if task.Status == model.TaskStatusRunning {
		return errors.New("任务正在执行中，无法删除")
	}

	// 删除关联数据
	s.resultRepo.DeleteByTaskID(id)
	s.db.Where("task_id = ?", id).Delete(&model.TaskTestCase{})

	return s.repo.Delete(id)
}

func (s *TaskService) FindByID(id int64) (*model.Task, error) {
	return s.repo.FindByID(id)
}

func (s *TaskService) FindPage(req *request.TaskPageRequest) ([]model.Task, int64, error) {
	req.Page.SetDefault()
	return s.repo.FindPage(&req.Query, req.Page.GetOffset(), req.Page.Size)
}

func (s *TaskService) GetProgress(id int64) (map[string]interface{}, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	return map[string]interface{}{
		"status":          task.Status,
		"status_name":     task.GetStatusName(),
		"total_count":     task.TotalCount,
		"completed_count": task.CompletedCount,
		"success_count":   task.SuccessCount,
		"failed_count":    task.FailedCount,
		"error_count":     task.ErrorCount,
		"progress":        task.GetProgress(),
	}, nil
}

func (s *TaskService) GetResults(req *request.TaskResultPageRequest) ([]map[string]interface{}, int64, error) {
	req.Page.SetDefault()
	results, total, err := s.resultRepo.FindByTaskID(req.TaskID, req.Page.GetOffset(), req.Page.Size)
	if err != nil {
		return nil, 0, err
	}

	// 获取关联的测试用例信息
	tcIDs := make([]int64, len(results))
	for i, r := range results {
		tcIDs[i] = r.TestCaseID
	}
	testcases, _ := s.tcRepo.FindByIDs(tcIDs)
	tcMap := make(map[int64]*model.TestCase)
	for i := range testcases {
		tcMap[testcases[i].ID] = &testcases[i]
	}

	// 组装返回数据
	data := make([]map[string]interface{}, len(results))
	for i, r := range results {
		item := map[string]interface{}{
			"id":               r.ID,
			"testcase_id":      r.TestCaseID,
			"status":           r.Status,
			"status_name":      r.GetStatusName(),
			"judge_result":     r.JudgeResult,
			"judge_result_name": r.GetJudgeResultName(),
			"judge_reason":     r.JudgeReason,
			"duration":         r.Duration,
			"request_content":  r.RequestContent,
			"response_content": r.ResponseContent,
			"error_message":    r.ErrorMessage,
			"executed_at":      r.ExecutedAt,
		}

		if tc, ok := tcMap[r.TestCaseID]; ok {
			item["testcase_name"] = tc.Name
			item["testcase_category"] = tc.GetCategoryName()
			item["risk_level"] = tc.GetRiskLevelName()
		}

		data[i] = item
	}

	return data, total, nil
}

func (s *TaskService) Start(id int64) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("任务不存在")
	}

	if task.Status == model.TaskStatusRunning {
		return errors.New("任务正在执行中")
	}

	if task.Status == model.TaskStatusCompleted {
		return errors.New("任务已完成")
	}

	// 更新状态为执行中
	now := time.Now()
	task.Status = model.TaskStatusRunning
	task.StartedAt = &now
	if err := s.repo.Update(task); err != nil {
		return err
	}

	// 启动异步执行
	exec := executor.GetExecutor()
	if exec != nil {
		return exec.Execute(id)
	}
	return nil
}

func (s *TaskService) Cancel(id int64) error {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("任务不存在")
	}

	if task.Status != model.TaskStatusRunning {
		return errors.New("任务未在执行")
	}

	// 取消执行中的任务
	exec := executor.GetExecutor()
	if exec != nil {
		exec.Cancel(id)
	}

	task.Status = model.TaskStatusCancelled
	now := time.Now()
	task.FinishedAt = &now
	return s.repo.Update(task)
}
