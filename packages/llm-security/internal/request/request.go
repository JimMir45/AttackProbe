package request

// PageParam 分页参数
type PageParam struct {
	Current int    `json:"current"`
	Size    int    `json:"size"`
	SortBy  string `json:"sort_by"`
}

// IDRequest ID请求
type IDRequest struct {
	ID int64 `json:"id" binding:"required"`
}

// IDsRequest 多ID请求
type IDsRequest struct {
	IDs []int64 `json:"ids" binding:"required"`
}

// 设置分页默认值
func (p *PageParam) SetDefault() {
	if p.Current <= 0 {
		p.Current = 1
	}
	if p.Size <= 0 {
		p.Size = 20
	}
	if p.Size > 100 {
		p.Size = 100
	}
}

func (p *PageParam) GetOffset() int {
	return (p.Current - 1) * p.Size
}
