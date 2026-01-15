package executor

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"llm-security-bas/internal/model"
	"llm-security-bas/pkg/judge"
	"llm-security-bas/pkg/llm"

	"gorm.io/gorm"
)

// Executor 任务执行器
type Executor struct {
	db          *gorm.DB
	concurrency int
	timeout     time.Duration
	mu          sync.Mutex
	running     map[int64]context.CancelFunc // taskID -> cancel func
}

// NewExecutor 创建执行器
func NewExecutor(db *gorm.DB, concurrency int, timeout time.Duration) *Executor {
	if concurrency <= 0 {
		concurrency = 5
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Executor{
		db:          db,
		concurrency: concurrency,
		timeout:     timeout,
		running:     make(map[int64]context.CancelFunc),
	}
}

// Execute 执行任务（异步）
func (e *Executor) Execute(taskID int64) error {
	// 获取任务信息
	var task model.Task
	if err := e.db.First(&task, taskID).Error; err != nil {
		return err
	}

	// 获取目标信息
	var target model.Target
	if err := e.db.First(&target, task.TargetID).Error; err != nil {
		return err
	}

	// 获取待执行的结果记录
	var results []model.TaskResult
	if err := e.db.Where("task_id = ? AND status = ?", taskID, model.ResultStatusPending).Find(&results).Error; err != nil {
		return err
	}

	if len(results) == 0 {
		log.Printf("Task %d: no pending results", taskID)
		return nil
	}

	// 获取测试用例
	testcaseIDs := make([]int64, len(results))
	for i, r := range results {
		testcaseIDs[i] = r.TestCaseID
	}
	var testcases []model.TestCase
	e.db.Where("id IN ?", testcaseIDs).Find(&testcases)
	tcMap := make(map[int64]*model.TestCase)
	for i := range testcases {
		tcMap[testcases[i].ID] = &testcases[i]
	}

	// 创建LLM客户端
	client := llm.NewOpenAIClient(&llm.ClientConfig{
		Endpoint:     target.Endpoint,
		APIKey:       target.APIKey,
		Model:        target.Model,
		Timeout:      target.Timeout,
		ExtraHeaders: target.ExtraHeaders,
	})

	// 创建取消上下文
	ctx, cancel := context.WithCancel(context.Background())
	e.mu.Lock()
	e.running[taskID] = cancel
	e.mu.Unlock()

	// 异步执行
	go func() {
		defer func() {
			e.mu.Lock()
			delete(e.running, taskID)
			e.mu.Unlock()
		}()

		e.executeTask(ctx, &task, client, results, tcMap)
	}()

	return nil
}

// executeTask 执行任务（内部方法）
func (e *Executor) executeTask(ctx context.Context, task *model.Task, client *llm.OpenAIClient, results []model.TaskResult, tcMap map[int64]*model.TestCase) {
	log.Printf("Task %d: starting execution, %d cases", task.ID, len(results))

	// 创建工作通道
	jobs := make(chan *model.TaskResult, len(results))
	done := make(chan struct{})

	// 启动worker
	var wg sync.WaitGroup
	for i := 0; i < e.concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for result := range jobs {
				select {
				case <-ctx.Done():
					return
				default:
					tc := tcMap[result.TestCaseID]
					if tc != nil {
						e.executeOne(ctx, client, result, tc)
						e.updateProgress(task.ID)
					}
				}
			}
		}(i)
	}

	// 分发任务
	go func() {
		for i := range results {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			default:
				jobs <- &results[i]
			}
		}
		close(jobs)
	}()

	// 等待完成
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		log.Printf("Task %d: cancelled", task.ID)
		e.db.Model(&model.Task{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
			"status":      model.TaskStatusCancelled,
			"finished_at": time.Now(),
		})
	case <-done:
		log.Printf("Task %d: completed", task.ID)
		e.db.Model(&model.Task{}).Where("id = ?", task.ID).Updates(map[string]interface{}{
			"status":      model.TaskStatusCompleted,
			"finished_at": time.Now(),
		})
	}
}

// executeOne 执行单个用例
func (e *Executor) executeOne(ctx context.Context, client *llm.OpenAIClient, result *model.TaskResult, tc *model.TestCase) {
	startTime := time.Now()

	// 构建请求
	messages := []llm.Message{}
	if tc.SystemPrompt != "" {
		messages = append(messages, llm.Message{Role: "system", Content: tc.SystemPrompt})
	}
	messages = append(messages, llm.Message{Role: "user", Content: tc.Content})

	reqJSON, _ := json.Marshal(messages)
	result.RequestContent = string(reqJSON)

	// 设置超时上下文
	reqCtx, cancel := context.WithTimeout(ctx, e.timeout)
	defer cancel()

	// 发送请求
	resp, err := client.Chat(reqCtx, &llm.ChatRequest{Messages: messages})

	now := time.Now()
	result.ExecutedAt = &now
	result.Duration = int(time.Since(startTime).Milliseconds())

	if err != nil {
		result.Status = model.ResultStatusError
		result.ErrorMessage = err.Error()
		e.db.Save(result)
		log.Printf("Task result %d: error - %s", result.ID, err.Error())
		return
	}

	result.ResponseContent = resp.Content

	// 判定结果
	j := judge.GetJudge(judge.JudgeMethod(tc.JudgeMethod))
	blocked, reason, err := j.Judge(resp.Content, tc.JudgeConfig)

	if err != nil {
		result.Status = model.ResultStatusError
		result.ErrorMessage = "判定失败: " + err.Error()
		e.db.Save(result)
		return
	}

	result.JudgeReason = reason
	judgeResult := 0
	if blocked {
		judgeResult = model.JudgeResultBlocked
		result.Status = model.ResultStatusSuccess
	} else {
		judgeResult = model.JudgeResultSuccess
		result.Status = model.ResultStatusFailed
	}
	result.JudgeResult = &judgeResult

	e.db.Save(result)
	log.Printf("Task result %d: %s (blocked=%v)", result.ID, reason, blocked)
}

// updateProgress 更新任务进度
func (e *Executor) updateProgress(taskID int64) {
	var stats []struct {
		Status int   `json:"status"`
		Count  int64 `json:"count"`
	}

	e.db.Model(&model.TaskResult{}).
		Select("status, count(*) as count").
		Where("task_id = ?", taskID).
		Group("status").
		Scan(&stats)

	var completed, success, failed, errCount int
	for _, s := range stats {
		switch s.Status {
		case model.ResultStatusSuccess:
			success = int(s.Count)
			completed += success
		case model.ResultStatusFailed:
			failed = int(s.Count)
			completed += failed
		case model.ResultStatusError:
			errCount = int(s.Count)
			completed += errCount
		}
	}

	e.db.Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"completed_count": completed,
		"success_count":   success,
		"failed_count":    failed,
		"error_count":     errCount,
	})
}

// Cancel 取消任务
func (e *Executor) Cancel(taskID int64) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if cancel, ok := e.running[taskID]; ok {
		cancel()
		return true
	}
	return false
}

// IsRunning 检查任务是否在运行
func (e *Executor) IsRunning(taskID int64) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, ok := e.running[taskID]
	return ok
}

// 全局执行器实例
var defaultExecutor *Executor

// Init 初始化全局执行器
func Init(db *gorm.DB, concurrency int, timeout time.Duration) {
	defaultExecutor = NewExecutor(db, concurrency, timeout)
}

// GetExecutor 获取全局执行器
func GetExecutor() *Executor {
	return defaultExecutor
}
