package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"seckill/conf"
	"seckill/infra/cache"
	"seckill/infra/code"
	"seckill/infra/utils/request"
	"seckill/infra/utils/response"
	"seckill/model"
	"time"
)

const (
	// 限流数据在 redis 中的 key
	rateLimitKey = "rate_limit:%s"
)

var (
	ctx = context.Background()
)

// UserLimit 对当个 IP 的请求进行限流
func UserLimit() gin.HandlerFunc {
	return func(context *gin.Context) {
		var (
			result model.Result
			clientIP string
			k string
			e error
			count int64
		)
		// 获取当前请求的 IP
		clientIP = request.GetIP(context.Request)
		k = fmt.Sprintf(rateLimitKey, clientIP)
		// 记录该 IP 的请求次数，使其 + 1
		count, e = cache.Client.Incr(ctx, k).Result()
		if e != nil {
			log.Printf("rdb.Incr() failed, err: %v", e)
			context.Abort()
			result.Code = http.StatusInternalServerError
			result.Message = code.RedisErr.Error()
			response.Fail(context, result)
			return
		}

		// 给该 IP 的请求次数记录设置一个过期时间
		if e = cache.Client.Expire(ctx, k, time.Duration(conf.Config.RateLimit.Time)*time.Second).Err(); e != nil {
			log.Printf("rdb.Expire() failed, err: %v", e)
			context.Abort()
			result.Code = http.StatusInternalServerError
			result.Message = code.RedisErr.Error()
			response.Fail(context, result)
			return
		}

		// 如果在规定时间段内的请求超过了规定的次数上限，则说明该 IP 存在恶意攻击行为，需要对其请求进行限制
		if count > conf.Config.RateLimit.Count {
			context.Abort()
			result.Code = http.StatusTooManyRequests
			result.Message = code.TooManyRequests.Error()
			response.Fail(context, result)
			return
		}
		context.Next()
	}
}
