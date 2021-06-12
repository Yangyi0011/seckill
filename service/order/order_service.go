package order

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"seckill/infra/cache"
	"seckill/infra/code"
	"seckill/model"
	"seckill/service"
	"sync"
	"time"
)

const (
	LockKey = "lock:%d:%d"
	OrderIdKey = "order:%d:%d"
)

var (
	once sync.Once
	ctx = context.Background()
)

// InitService 单例模式初始化 IOrderService 接口实例
func InitService() {
	once.Do(func() {
		service.OrderService = newOrderService()
	})
}

type orderService struct {
	goodsService service.IGoodsService
	redis *redis.Client
}

func newOrderService() *orderService {
	return &orderService{
		goodsService: service.GoodsService,
		redis:        cache.Client,
	}
}

func (s *orderService) SecondKill(userId, goodsId int) (e error) {
	var (
		goods model.Goods
		orderId string
		stock   int
	)
	// 获取商品信息
	if goods, e = s.goodsService.FindGoodsByID(goodsId); e != nil {
		return
	}
	// 验证商品信息
	if e = s.goodsService.Check(goods); e != nil {
		return
	}
	// 加锁
	if e = s.tryLock(userId, goodsId, time.Now().Unix()); e != nil {
		return
	}
	// 校验重复秒杀
	if orderId, e = s.GetOrderId(userId, goodsId); e != nil {
		e = code.RedisErr
		return
	}
	if len(orderId) > 0 {
		e = code.RepeateSeckillErr
		return
	}
	// 预减商品库存
	if stock, e = s.goodsService.DecrStock(goodsId); e != nil {
		e = code.RedisErr
		return
	}
	// 没有库存了，商品已售罄
	if stock < 0 {
		e = code.GoodsSaleOut
		return
	}
	// todo 异步下单
	return
}

func (s *orderService) GetOrderId(userId, goodsId int) (orderId string, err error) {
	key := fmt.Sprintf(OrderIdKey, userId, goodsId)
	if orderId, err = s.redis.Get(ctx, key).Result(); err != nil {
		if err == redis.Nil {
			err = nil
		} else {
			log.Printf("redis.Get() failed, err: %v", err)
			err = code.RedisErr
		}
	}
	return
}

func (s *orderService) tryLock(userId, goodsId int, lockId int64) (err error) {
	var res bool
	key := fmt.Sprintf(LockKey, userId, goodsId)
	if res, err = s.redis.SetNX(ctx, key, lockId, time.Minute).Result(); err != nil {
		log.Printf("redis.SetNx() failed, err: %v", err)
		err = code.RedisErr
		return
	}
	if !res {
		err = code.SeckillFailedErr
	}
	return
}

func (s *orderService) UnLock(userId, goodsId int) (err error) {
	if err = s.redis.Del(ctx, fmt.Sprintf("lock:%d:%d", userId, goodsId)).Err(); err != nil {
		log.Printf("redis.Del() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}