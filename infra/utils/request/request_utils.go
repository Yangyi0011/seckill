package request

import (
	"github.com/gin-gonic/gin"
	"seckill/infra/code"
	"seckill/infra/secret"
)

// GetCurrentCustomClaims 获取当前用户的 JWT 信息
func GetCurrentCustomClaims(ctx *gin.Context) (claims *secret.CustomClaims, e error) {
	// 获取当前用户数据
	c, ok := ctx.Get("claims")
	if !ok {
		e = code.TokenInvalidErr
		return nil, e
	}
	claims, ok = c.(*secret.CustomClaims)
	if !ok || claims == nil {
		e = code.TokenInvalidErr
		return nil, e
	}
	return claims, nil
}
