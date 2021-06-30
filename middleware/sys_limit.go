package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/conf"
	"seckill/infra/code"
	"seckill/infra/utils/limit"
	"seckill/infra/utils/response"
	"seckill/model"
)

var (
	tokenBucket *limit.TokenBucket
)
func init() {
	tokenBucket = limit.NewTokenBucket(int(conf.Config.Total), int(conf.Config.Rate))
}

// SysLimit 使用令牌桶算法针对整个系统进行限流
func SysLimit() gin.HandlerFunc {
	return func(context *gin.Context) {
		var (
			result model.Result
		)
		// 如果请求速率超过了系统的限制，则直接返回
		if !tokenBucket.Limit() {
			context.Abort()
			result.Code = http.StatusTooManyRequests
			result.Message = code.SysBusyErr.Error()
			response.Fail(context, result)
			return
		}
		context.Next()
	}
}
