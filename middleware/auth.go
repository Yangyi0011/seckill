package middleware

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"seckill/infra/secret"
	"seckill/infra/utils/response"
	"seckill/model"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result := model.Result{}
		// 首先在请求头获取 token
		auth := ctx.Request.Header.Get("Authorization")
		if len(auth) == 0 {
			ctx.Abort()
			result.Code = http.StatusUnauthorized
			result.Message = "请先登录"
			response.Fail(ctx, result)
			return
		}
		// 解析 token
		j := secret.NewJWT()
		claims, err := j.ParseToken(auth)
		if err != nil {
			ctx.Abort()
			result.Code = http.StatusUnauthorized
			result.Message = err.Error()
			response.Fail(ctx, result)
			return
		}
		log.Println("auth 认证通过")
		// 认证通过则刷新 JWT 过期时间
		_, _ = j.RefreshToken(auth)

		// 继续交由下一个路由处理，并将解析出的信息传递下去
		ctx.Set("claims", claims)
		ctx.Next()
	}
}
