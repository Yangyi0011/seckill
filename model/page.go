package model

// PageDTO 分页查询 DTO
type PageDTO struct {
	Index int `json:"index"`
	Size  int `json:"size"`
}

// PageVO 分页视图 VO
type PageVO struct {
	List interface{} `json:"list"` // 页面数据
	Index   int  `json:"index"`  // 当前页
	Total   int  `json:"total"`  // 总页数
	HasPrev bool `json:"has_prev"` // 是否有上一页
	HasNext bool `json:"has_next"` // 是否有下一页
}

// GetOffset 获取分页偏移量
func (p PageDTO) GetOffset() int {
	offset := 0
	if p.Index != 0 {
		offset = (p.Index-1)*p.Size
	}
	return offset
}

// GetLimit 获取每页数据量
func (p PageDTO) GetLimit() int {
	// 默认每页 10 条数据
	limit := 10
	if p.Size != 0 {
		limit = p.Size
	}
	return limit
}
