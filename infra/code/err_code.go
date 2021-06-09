package code

import "strconv"

var (
	errCode = make(map[int]string, 0)
)

// 构建自定义 Code 类型
func buildCode(code int, msg string) Code {
	if _, ok := errCode[code]; ok {
		panic("Code already exist: " + strconv.Itoa(code))
	}
	// 保存自定义 Code 类型到字典中
	errCode[code] = msg
	return Code(code)
}

// CodeMsg 自定义接口，规定错误码的信息格式
type CodeMsg interface {
	Code() int
	// Error builtin.error.Error()
	Error() string
}

// Code 自定义 Code 类型
type Code int

// Error 实现 CodeMsg.Error() 接口来自定义 err
func (e Code) Error() string {
	return errCode[e.Code()]
}

// Code 实现 CodeMsg.Code()
func (e Code) Code() int {
	return int(e)
}