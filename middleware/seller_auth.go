package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/infra/code"
	"seckill/infra/secret"
	"seckill/infra/utils/response"
	"seckill/model"
)

// SellerAuth 商家认证
func SellerAuth() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		result := model.Result{}
		// 获取当前用户数据
		var claims *secret.CustomClaims
		c, ok := ctx.Get("claims")
		if !ok {
			ctx.Abort()
			result.Code = http.StatusUnauthorized
			result.Message = code.TokenInvalidErr.Error()
			response.Fail(ctx, result)
			return
		}
		claims, ok = c.(*secret.CustomClaims)
		// 没有登录信息
		if !ok || claims == nil {
			ctx.Abort()
			result.Code = http.StatusUnauthorized
			result.Message = code.TokenInvalidErr.Error()
			response.Fail(ctx, result)
			return
		}
		// 不是卖家
		if claims.Kind != model.NormalSeller {
			ctx.Abort()
			result.Code = http.StatusForbidden
			result.Message = code.StatusForbiddenErr.Error()
			response.Fail(ctx, result)
			return
		}
		ctx.Next()
	}
}
