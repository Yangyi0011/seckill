package mq

import (
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log"
	"seckill/infra/code"
	"seckill/service"
	"time"
)

const (
	// 订单队列在 redis 中的 key
	precreateOrderKey = "precreate_order_queue"
)

// precreateOrder 预创建订单模型
type precreateOrder struct {
	orderService service.IOrderService
	redis *redis.Client
}

// PrecreateOrderMsg 订单消息体
type PrecreateOrderMsg struct {
	UserId  int `json:"user_id"`
	GoodsId int `json:"goods_id"`
}

// Send 把订单消息推送到队列中
func (mq *precreateOrder) Send(msg PrecreateOrderMsg) error {
	data, e := json.Marshal(msg)
	if e != nil {
		return code.SerializeErr
	}
	return mq.redis.LPush(ctx, precreateOrderKey, data).Err()
}

// Receive 从队列中消费消息
func (mq *precreateOrder) Receive() {
	var (
		popStr string	// 获取到的消息字符串
		err    error
	)
	// 循环消费
	for {
		if popStr, err = mq.redis.RPop(ctx, precreateOrderKey).Result(); err != nil {
			if err == redis.Nil {
				err = nil
			} else {
				log.Printf("redis.BRPop() failed, err: %v", err)
				return
			}
		}
		if len(popStr) == 0 {
			// 队列中没有数据时，睡眠 0.5 秒
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var msg PrecreateOrderMsg
		if err = json.Unmarshal([]byte(popStr), &msg); err != nil {
			continue
		}
		if msg.GoodsId == 0 || msg.UserId == 0 {
			continue
		}
		// 创建订单
		if err = mq.orderService.CreateOrder(msg.UserId, msg.GoodsId); err != nil {
			log.Printf("orderService.CreateOrder() failed, err: %v", err)
			continue
		}
		if err = mq.orderService.UnLock(msg.UserId, msg.GoodsId); err != nil {
			log.Printf("orderService.UnLocks() failed, err: %v", err)
			continue
		}
	}
}