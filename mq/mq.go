package mq

import (
	"context"
	"seckill/infra/cache"
	"seckill/service"
)

var (
	ctx = context.Background()
	OrderTimeout *orderTimeout
	PrecreateOrder *precreateOrder
)

const (
	// 订单超时延迟队列
	orderTimeoutDelayQueue = "order_timeout_delay_queue"
)

func Init() {
	OrderTimeout = &orderTimeout{
		orderService: service.OrderService,
		redis:        cache.Client,
	}
	PrecreateOrder = &precreateOrder{
		orderService: service.OrderService,
		redis:        cache.Client,
	}
}

// Run 对队列进行消息监听和消费
func Run() {
	go OrderTimeout.Receive()
	go PrecreateOrder.Receive()
}