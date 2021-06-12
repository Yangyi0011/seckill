package mq

import (
	"github.com/go-redis/redis/v8"
	"log"
	"seckill/conf"
	"seckill/model"
	"seckill/service"
	"strconv"
	"time"
)

// orderTimeout 订单超时模型
type orderTimeout struct {
	orderService service.IOrderService
	redis *redis.Client
}

// Send 往队列中发送数据
func (mq *orderTimeout) Send(orderId string) {
	// 把订单编号发送到延迟队列中，并以 Score 的方式设置超时时间
	if err := mq.redis.ZAdd(ctx, orderTimeoutDelayQueue, &redis.Z{
		Score:  float64(time.Now().Unix() + conf.Config.Expiration),
		Member: orderId,
	}).Err(); err != nil {
		log.Printf("订单【%s】加入延迟队列失败, err: %v", orderId, err)
	} else {
		log.Printf("订单【%s】加入延迟队列", orderId)
	}
	return
}

// Remove 从队列中移除数据
func (mq *orderTimeout) Remove(orderId string) {
	if err := mq.redis.ZRem(ctx, orderTimeoutDelayQueue, orderId).Err(); err != nil {
		log.Printf("订单【%s】移除延迟队列失败, err: %v", orderId, err)
	} else {
		log.Printf("订单【%s】移除延迟队列", orderId)
	}
	return
}

// Receive 消费队列数据
func (mq *orderTimeout) Receive() {
	var (
		list      []string
		err       error
		orderInfo model.OrderInfo
	)
	for {
		// 从延迟队列中的拉取订单数据，以当前时间戳作为最大 Score 来拉取，每次拉取一条数据
		// 即：把到当前时间依旧未支付的订单当做超时订单处理，直接关闭该订单
		if list, err = mq.redis.ZRangeByScore(ctx, orderTimeoutDelayQueue, &redis.ZRangeBy{
			Min:    "0",
			Max:    strconv.FormatInt(time.Now().Unix(), 10),
			Offset: 0,
			Count:  1,
		}).Result(); err != nil {
			log.Printf("redis.ZRangeByScore() failed, err: %v", err)
			continue
		}
		// 没有订单数据时睡眠
		if len(list) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// 用订单编号来查询订单数据
		if orderInfo, err = mq.orderService.GetOrderInfo(list[0]); err != nil {
			log.Printf("orderService.GetOrderInfo() failed, orderId: %s, err: %v", list[0], err)
			continue
		}
		// 关闭该订单
		if err = mq.orderService.CloseOrder(int(orderInfo.UserId), orderInfo.OrderId); err != nil {
			log.Printf("orderService.CloseOrder() failed, orderId: %s, err: %v", orderInfo.OrderId, err)
			continue
		}
		log.Printf("订单【%s】已关闭", orderInfo.OrderId)
	}
}
