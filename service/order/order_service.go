package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/dao"
	"seckill/infra/cache"
	"seckill/infra/code"
	"seckill/infra/utils/key"
	"seckill/model"
	"seckill/mq"
	"seckill/service"
	"sync"
	"time"
)

const (
	// LockKey redis 锁格式
	LockKey = "lock:%d:%d"
	// OrderIdKey 订单缓存 key
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
	dao dao.IOrderDao
	goodsService service.IGoodsService
	redis *redis.Client
}

func newOrderService() *orderService {
	return &orderService{
		dao: dao.OrderDao,
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
	// 异步下单
	msg := mq.PrecreateOrderMsg{
		UserId:  userId,
		GoodsId: goodsId,
	}
	if e = mq.PrecreateOrder.Send(msg); e != nil {
		log.Println(e)
		e = code.RedisErr
		return
	}
	return
}

func (s *orderService) CreateOrder(userId, goodsId int) (err error) {
	var (
		orderId string
		goods   model.Goods
	)
	// 是否重复秒杀
	if orderId, err = s.GetOrderId(userId, goodsId); err != nil {
		log.Printf("s.GetOrderId() failed, err: %v", err)
		return
	}
	if len(orderId) > 0 {
		return
	}
	// 查询商品信息
	if goods, err = s.goodsService.FindGoodsByID(goodsId); err != nil {
		log.Printf("goodsService.GetGoods() failed, err: %v", err)
		return
	}
	// 创建订单
	orderInfo := s.newOrderInfo(userId, goods)
	if err = s.dao.CreateOrder(orderInfo); err != nil {
		log.Printf("orderRepository.CreateOrder() failed, err: %v", err)
		return
	}
	// 创建订单信息缓存
	if err = s.CreateOrderCache(orderInfo); err != nil {
		return
	}
	// 加入订单超时延迟队列
	mq.OrderTimeout.Send(orderInfo.OrderId)
	return
}

func (s *orderService) GetOrderInfo(orderId string) (o model.OrderInfo, e error) {
	if o, e = s.dao.QueryOrderInfoByOrderId(orderId); e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			e = code.RecordNotFoundErr
			return
		}
		println(e)
		e = code.DBErr
		return
	}
	return
}

func (s *orderService) CloseOrder(userId int, orderId string) (e error) {
	var orderInfo model.OrderInfo
	if orderInfo, e = s.GetOrderInfo(orderId); e != nil {
		return
	}
	if orderInfo.ID == 0 || orderInfo.UserId != uint(userId) {
		e = code.OrderNotFoundErr
		return
	}
	if orderInfo.Status != model.Unpaid {
		e = code.OrderStatusErr
		return
	}
	if e = s.dao.CloseOrder(orderInfo); e != nil {
		e = code.OrderCloseErr
		return
	}
	// 加缓存中的库存
	if e = s.goodsService.IncrStock(int(orderInfo.GoodsId)); e != nil {
		e = code.RedisErr
		return
	}
	// 删除订单缓存
	if e = s.DeleteOrderCache(orderInfo); e != nil {
		e = code.RedisErr
		return
	}
	// 移除延迟队列中的订单单号
	mq.OrderTimeout.Remove(orderInfo.OrderId)
	return
}

func (s *orderService) CreateOrderCache(order model.OrderInfo) (err error) {
	k := fmt.Sprintf(OrderIdKey, order.UserId, order.GoodsId)
	if err = s.redis.Set(ctx, k, order.OrderId, -1).Err(); err != nil {
		log.Printf("redis.Set() failed, err: %v, order: %v", err, order)
		err = code.RedisErr
	}
	return
}

func (s *orderService) DeleteOrderCache(order model.OrderInfo) (err error) {
	k := fmt.Sprintf(OrderIdKey, order.UserId, order.GoodsId)
	if err = s.redis.Del(ctx, k).Err(); err != nil {
		log.Printf("redis.Del() falied, err: %v, order: %v", err, order)
		err = code.RedisErr
	}
	return
}

func (s *orderService) GetOrderId(userId, goodsId int) (orderId string, err error) {
	k := fmt.Sprintf(OrderIdKey, userId, goodsId)
	if orderId, err = s.redis.Get(ctx, k).Result(); err != nil {
		if err == redis.Nil {
			err = nil
		} else {
			log.Printf("redis.Get() failed, err: %v", err)
			err = code.RedisErr
		}
	}
	return
}

func (s *orderService)newOrderInfo(userId int, goods model.Goods) model.OrderInfo {
	orderInfo := model.OrderInfo{
		OrderId:    createOrderId(),
		UserId: uint(userId),
		GoodsId:    goods.ID,
		GoodsName:  goods.Name,
		GoodsImg:   goods.Img,
		GoodsPrice: goods.Price,
		Status:     model.Unpaid,		// 默认未支付
	}
	return orderInfo
}

// 创建订单编号
func createOrderId() string {
	return time.Now().Format("20060102150405") + key.CreateKey(key.Number, 6)
}

func (s *orderService) tryLock(userId, goodsId int, lockId int64) (err error) {
	var res bool
	k := fmt.Sprintf(LockKey, userId, goodsId)
	if res, err = s.redis.SetNX(ctx, k, lockId, time.Minute).Result(); err != nil {
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
	k := fmt.Sprintf(LockKey, userId, goodsId)
	if err = s.redis.Del(ctx, k).Err(); err != nil {
		log.Printf("redis.Del() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}