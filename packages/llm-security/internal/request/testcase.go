package request

// TestCasePageRequest 用例分页请求
type TestCasePageRequest struct {
	Name      string    `json:"name"`
	Category  int       `json:"category"`
	RiskLevel int       `json:"risk_level"`
	IsBuiltin *int      `json:"is_builtin"` // 使用指针区分未设置和0
	Status    int       `json:"status"`
	Page      PageParam `json:"page"`
}

// TestCaseAddRequest 添加用例请求
type TestCaseAddRequest struct {
	Name             string `json:"name" binding:"required"`
	Category         int    `json:"category" binding:"required"`
	RiskLevel        int    `json:"risk_level"`
	AttackType       string `json:"attack_type"`
	Content          string `json:"content" binding:"required"`
	SystemPrompt     string `json:"system_prompt"`
	ExpectedBehavior string `json:"expected_behavior"`
	JudgeMethod      int    `json:"judge_method"`
	JudgeConfig      string `json:"judge_config"`
	Source           string `json:"source"`
	Reference        string `json:"reference"`
}

// TestCaseUpdateRequest 更新用例请求
type TestCaseUpdateRequest struct {
	ID               int64  `json:"id" binding:"required"`
	Name             string `json:"name" binding:"required"`
	Category         int    `json:"category" binding:"required"`
	RiskLevel        int    `json:"risk_level"`
	AttackType       string `json:"attack_type"`
	Content          string `json:"content" binding:"required"`
	SystemPrompt     string `json:"system_prompt"`
	ExpectedBehavior string `json:"expected_behavior"`
	JudgeMethod      int    `json:"judge_method"`
	JudgeConfig      string `json:"judge_config"`
	Reference        string `json:"reference"`
	Status           int    `json:"status"`
}

// BatchStatusRequest 批量状态更新
type BatchStatusRequest struct {
	IDs    []int64 `json:"ids" binding:"required"`
	Status int     `json:"status"`
}
