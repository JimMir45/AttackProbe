package request

// TargetQuery 目标查询条件
type TargetQuery struct {
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Status int    `json:"status"`
}

// TargetPageRequest 目标分页请求
type TargetPageRequest struct {
	Query TargetQuery `json:"query"`
	Page  PageParam   `json:"page"`
}

// TargetAddRequest 添加目标请求
type TargetAddRequest struct {
	Name         string `json:"name" binding:"required"`
	Type         int    `json:"type" binding:"required"`
	Endpoint     string `json:"endpoint" binding:"required"`
	APIKey       string `json:"api_key"`
	Model        string `json:"model"`
	ExtraHeaders string `json:"extra_headers"`
	Timeout      int    `json:"timeout"`
}

// TargetUpdateRequest 更新目标请求
type TargetUpdateRequest struct {
	ID           int64  `json:"id" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Type         int    `json:"type" binding:"required"`
	Endpoint     string `json:"endpoint" binding:"required"`
	APIKey       string `json:"api_key"`
	Model        string `json:"model"`
	ExtraHeaders string `json:"extra_headers"`
	Timeout      int    `json:"timeout"`
	Status       int    `json:"status"`
}
