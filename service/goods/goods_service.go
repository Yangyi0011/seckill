package goods

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/dao"
	"seckill/infra/cache"
	"seckill/infra/code"
	"seckill/infra/utils/bean"
	"seckill/model"
	"seckill/service"
	"sync"
	"time"
)

var (
	once sync.Once
	ctx  = context.Background()
)

// 商品常量
const (
	Ended      int8 = -1 // 已结束
	NotStarted int8 = 0  // 未开始
	OnGoing    int8 = 1  // 进行中
	SoldOut    int8 = 2  // 已售罄

	CacheKey      = "goods:%d"       // 商品缓存key格式
	CacheExpire   = 12 * time.Hour   // 缓存过期时间
	GoodsStockKey = "goods_stock:%d" // 商品库存key格式
)

// InitService 单例模式初始化 IGoodsService 接口实例
func InitService() {
	once.Do(func() {
		service.GoodsService = newGoodsService()
	})
}

// service.IGoodsService 接口实现
type goodsService struct {
	dao   dao.IGoodsDao
	redis *redis.Client
}

// NewGoodsService 创建一个 service.IGoodsService 接口实例
func newGoodsService() *goodsService {
	return &goodsService{
		dao:   dao.GoodsDao,
		redis: cache.Client,
	}
}

func (s *goodsService) Check(g model.Goods) (e error) {
	now := time.Now().Unix()
	startTime := g.StartTime.Unix()
	endTime := g.EndTime.Unix()
	if now < startTime {
		// 秒杀活动未开始
		e = code.SeckillNotStart
	} else if now > endTime {
		// 秒杀活动已结束
		e = code.SeckillEnded
	} else if g.Stock <= 0 {
		// 商品已售罄
		e = code.GoodsSaleOut
	}
	return
}

func (s *goodsService) FindGoodsByID(id int) (g model.Goods, e error) {
	// 先尝试从缓存里获取
	if g, e = s.getGoodsFromCache(id); e != nil {
		return
	}
	// 成功从缓存中取到值就返回
	if g.ID > 0 {
		return
	}
	// 从数据库中取值
	if g, e = s.dao.QueryGoodsByID(id); e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			e = code.RecordNotFoundErr
			return
		}
		e = code.DBErr
		return
	}
	// 把数据放入缓存中
	e = s.setGoodsCache(g)
	return
}

// 从缓存中获取商品信息
func (s *goodsService) getGoodsFromCache(id int) (g model.Goods, e error) {
	var res string
	key := fmt.Sprintf(CacheKey, id)
	if res, e = s.redis.Get(ctx, key).Result(); e != nil {
		if e == redis.Nil {
			e = nil
		} else {
			log.Printf("Redis Get() faild, err: %v", e)
			e = code.RedisErr
		}
		return
	}
	if e = json.Unmarshal([]byte(res), &g); e != nil {
		log.Printf("json.Unmarshal() failed, err: %v, json: %v", e, res)
		e = code.SerializeErr
	}
	return
}

// 把商品信息添加进缓存中
func (s goodsService) setGoodsCache(g model.Goods) (e error) {
	var data []byte
	if data, e = json.Marshal(g); e != nil {
		log.Printf("json.Marshal() faild, err: %v", e)
		e = code.SerializeErr
		return
	}
	key := fmt.Sprintf(CacheKey, g.ID)
	if e = s.redis.Set(ctx, key, string(data), CacheExpire).Err(); e != nil {
		log.Printf("redis.Set() failed, err: %v", e)
		e = code.RedisErr
		return
	}
	return
}

// 删除商品缓存信息
func (s * goodsService) deleteGoodsCache(id int) (e error) {
	key := fmt.Sprintf(CacheKey, id)
	if e = s.redis.Del(ctx, key).Err(); e != nil {
		log.Printf("redis.Del() failed, err: %v", e)
		e = code.RedisErr
		return
	}
	return
}

func (s *goodsService) ToVO(g model.Goods) (vo model.GoodsVO, e error) {
	if e = bean.SimpleCopyProperties(&vo, g); e != nil {
		return vo, code.ConvertErr
	}
	startTime := g.StartTime.Unix()
	endTime := g.EndTime.Unix()
	now := time.Now().Unix()
	if now < startTime {
		// 当前时间 < 商品秒杀的开始时间，说明秒杀活动未开始，需要计算倒计时
		vo.Status = NotStarted
		vo.Duration = startTime - now
	} else if now >= startTime && now <= endTime {
		// 秒杀开始时间 <= 当前时间 <= 秒杀结束时间，说明活动正在进行中
		if vo.Stock > 0 {
			// 如果商品还有库存，则活动正常进行
			vo.Status = OnGoing
		} else {
			// 商品已售罄
			vo.Status = SoldOut
		}
	} else {
		// 该商品的秒杀活动已结束
		vo.Status = Ended
	}
	return vo, nil
}

func (s *goodsService) FindGoodsVOByID(id int) (vo model.GoodsVO, e error) {
	goods, err := s.FindGoodsByID(id)
	if err != nil {
		return vo, err
	}
	vo, e = s.ToVO(goods)
	if e != nil {
		return vo, e
	}
	return vo, nil
}

func (s *goodsService) FindByCondition(c model.GoodsQueryCondition) ([]model.GoodsVO, error) {
	goodsList, e := s.dao.QueryByCondition(c)
	if e != nil {
		log.Println(e)
		e = code.DBErr
		return nil, e
	}
	list := make([]model.GoodsVO, 0)
	for _, v := range goodsList {
		vo, err := s.ToVO(v)
		if err != nil {
			return nil, err
		}
		list = append(list, vo)
	}
	return list, nil
}

func (s *goodsService) Insert(dto model.GoodsDTO) error {
	var g model.Goods
	// 数据转换
	err := bean.SimpleCopyProperties(&g, dto)
	if err != nil {
		log.Println(err)
		return code.ConvertErr
	}
	g.CreatedAt = model.LocalTime(time.Now())
	if e := s.dao.Insert(g); e != nil {
		log.Println(e)
		return code.DBErr
	}
	return nil
}

func (s *goodsService) Update(dto model.GoodsDTO) (e error) {
	var g model.Goods
	// 数据转换
	if e = bean.SimpleCopyProperties(&g, dto); e != nil {
		log.Println(e)
		e = code.ConvertErr
		return
	}
	// 数据更新
	g.UpdatedAt = model.LocalTime(time.Now())
	if e = s.dao.Update(g); e != nil {
		log.Println(e)
		e = code.DBErr
		return
	}
	// 用更新后的数据去更新缓存
	var goods model.Goods
	if goods, e = s.dao.QueryGoodsByID(int(g.ID)); e !=nil {
		return
	}
	if e = s.setGoodsCache(goods); e != nil {
		return
	}
	return
}

func (s *goodsService) DeleteWithPhysics(id int) (e error) {
	if e = s.dao.Delete(id); e != nil {
		log.Println(e)
		e = code.DBErr
		return
	}
	// 删除缓存信息
	if e = s.deleteGoodsCache(id); e != nil {
		return
	}
	return nil
}

func (s *goodsService) DeleteWithLogic(id int) (e error) {
	var g model.Goods
	if g, e = s.FindGoodsByID(id); e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			e = code.RecordNotFoundErr
			return
		}
		log.Println(e)
		e = code.DBErr
		return
	}
	now := model.LocalTime(time.Now())
	g.DeletedAt = &now
	if e = s.dao.Update(g); e != nil {
		log.Println(e)
		e = code.DBErr
		return
	}
	// 删除缓存信息
	if e = s.deleteGoodsCache(id); e != nil {
		return
	}
	return nil
}

func (s *goodsService) SetGoodsStock(goodsId int, stock int) (err error) {
	key := fmt.Sprintf(GoodsStockKey, goodsId)
	if err = s.redis.Set(ctx, key, stock, -1).Err(); err != nil {
		log.Printf("redis.Set() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}

func (s *goodsService) DecrStock(goodsId int) (stock int, err error) {
	var res int64
	key := fmt.Sprintf(GoodsStockKey, goodsId)
	if res, err = s.redis.Decr(ctx, key).Result(); err != nil {
		log.Printf("redis.Decr() failed, err: %v", err)
		err = code.RedisErr
	}
	stock = int(res)
	return
}

func (s *goodsService) IncrStock(goodsId int) (err error) {
	key := fmt.Sprintf(GoodsStockKey, goodsId)
	if err = s.redis.Incr(ctx, key).Err(); err != nil {
		log.Printf("redis.Incr() failed, err: %v", err)
		err = code.RedisErr
	}
	return
}

func (s *goodsService) InitScekillGoods(userId int) (e error){
	var list []model.Goods
	if list, e = s.dao.QueryByCondition(model.GoodsQueryCondition{UserId: uint(userId)}); e != nil {
		log.Println(e)
		e = code.DBErr
		return
	}
	for _, v := range list {
		if e = s.SetGoodsStock(int(v.ID), v.Stock); e != nil {
			return
		}
	}
	return
}