package limit

import (
	"time"
)

// TokenBucket 令牌桶限流算法
// 规定固定容量的桶, token 以固定速度往桶内填充, 当桶满时 token 不会被继续放入,
// 每过来一个请求把 token 从桶中移除, 如果桶中没有 token 不能请求
type TokenBucket struct {
	Cap int				// 桶容量
	NowSize int			// 当前令牌数量
	PreTime int64		// 上一次请求时间
	Rate int			// 令牌更新速率（单位：Rate个/秒）
}

// NewTokenBucket 创建一个令牌桶
func NewTokenBucket(cap, rate int) *TokenBucket {
	return &TokenBucket{
		Cap:     cap,
		NowSize: cap,
		PreTime: time.Now().Unix(),
		Rate:    rate,
	}
}

// Limit 限流，返回 false 说明没有拿到令牌，通不过
func (t *TokenBucket) Limit() bool {
	now := time.Now().Unix()
	// 按速率往桶中添加令牌
	size := t.NowSize + max(0, int(now - t.PreTime) * t.Rate)
	// 当前令牌数不能超出桶容量
	t.NowSize = min(t.Cap, size)
	t.PreTime = now
	if t.NowSize < 1 {
		// 桶中没有令牌，通不过
		return false
	}
	// 拿走一个令牌
	t.NowSize --
	return true
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
