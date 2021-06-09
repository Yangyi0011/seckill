package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/model"
)

// Success 响应成功
func Success(context *gin.Context, result model.Result) {
	if result.Code == 0 {
		result.Code = http.StatusOK
	}
	if result.Message == "" {
		result.Message = "成功"
	}
	context.JSON(result.Code, gin.H{
		"result": result,
	})
}

// Fail 响应失败
func Fail(context *gin.Context, result model.Result) {
	if result.Code == 0 {
		result.Code = http.StatusBadRequest
	}
	if result.Message == "" {
		result.Message = "失败"
	}
	context.JSON(result.Code, gin.H{
		"result": result,
	})
}

// Response 自定义响应
func Response(context *gin.Context, result model.Result) {
	context.JSON(result.Code, gin.H{
		"result": result,
	})
}