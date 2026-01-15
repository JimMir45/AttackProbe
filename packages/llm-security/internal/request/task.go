package request

// TaskQuery 任务查询条件
type TaskQuery struct {
	Name     string `json:"name"`
	TargetID int64  `json:"target_id"`
	Status   *int   `json:"status"` // 使用指针区分未设置和0
}

// TaskPageRequest 任务分页请求
type TaskPageRequest struct {
	Query TaskQuery `json:"query"`
	Page  PageParam `json:"page"`
}

// TaskAddRequest 创建任务请求
type TaskAddRequest struct {
	Name        string  `json:"name" binding:"required"`
	TargetID    int64   `json:"target_id" binding:"required"`
	TestCaseIDs []int64 `json:"testcase_ids"` // 为空则使用所有启用的用例
}

// TaskResultPageRequest 任务结果分页请求
type TaskResultPageRequest struct {
	TaskID int64     `json:"task_id" binding:"required"`
	Page   PageParam `json:"page"`
}
