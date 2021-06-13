package request

import (
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"seckill/infra/code"
	"seckill/infra/secret"
	"strings"
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

// GetIP 获取当前请求的 IP 地址
func GetIP(r *http.Request) (ip string) {
	// 尝试从 X-Forwarded-For 中获取
	xForwardedFor := r.Header.Get(`X-Forwarded-For`)
	ip = strings.TrimSpace(strings.Split(xForwardedFor, `,`)[0])
	if ip == `` {
		// 尝试从 X-Real-Ip 中获取
		ip = strings.TrimSpace(r.Header.Get(`X-Real-Ip`))
		if ip == `` {
			// 直接从 Remote Addr 中获取
			_ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
			if err != nil {
				panic(err)
			} else {
				ip = _ip
			}
		}
	}
	return
}