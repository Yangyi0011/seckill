package model

// PageDTO 分页查询 DTO
type PageDTO struct {
	Index int
	Size  int
}

// PageVO 分页视图 VO
type PageVO struct {
	List interface{} `json:"list"` // 页面数据
	Index   int  `json:"index"`  // 当前页
	Total   int  `json:"total"`  // 总页数
	HasPrev bool `json:"has_prev"` // 是否有上一页
	HasNext bool `json:"has_next"` // 是否有下一页
}
