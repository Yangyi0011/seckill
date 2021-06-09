package model

// Result 统一响应模型
type Result struct {
	Code    int         `json:"code" example:"000"`
	Message string      `json:"message" example:"响应信息"`
	Data    interface{} `json:"data" `
}
