package goods

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/dao"
	"seckill/infra/code"
	"seckill/infra/utils/bean"
	"seckill/model"
	"sync"
	"time"
)

var (
	once sync.Once
)

// 商品常量
const (
	Ended      int8 = -1 // 已结束
	NotStarted int8 = 0  // 未开始
	OnGoing    int8 = 1  // 进行中
	SoldOut    int8 = 2  // 已售罄
)

// service.IGoodsService 接口实现
type goodsService struct {
	dao dao.IGoodsDao
}

// NewGoodsService 创建一个 service.IGoodsService 接口实例
func NewGoodsService() *goodsService {
	return &goodsService{
		dao: dao.GoodsDao,
	}
}

// SingleGoodsService service.IGoodsService 接口单例模式
func SingleGoodsService() (s *goodsService) {
	once.Do(func() {
		s = NewGoodsService()
	})
	return
}

func (s *goodsService) FindGoodsByID(id int) (g model.Goods, e error) {
	g, e = s.dao.QueryGoodsByID(id)
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return g, code.RecordNotFound
		}
		return g, code.DBErr
	}
	return
}

func (s *goodsService) ToVO(g model.Goods) (vo model.GoodsVO, e error) {
	if e = bean.SimpleCopyProperties(&vo, g); e != nil {
		return vo, code.ConvertErr
	}
	startTime := time.Time(g.StartTime).Unix()
	endTime := time.Time(g.EndTime).Unix()
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
	//db.DB.Exec("select * from")
	return []model.GoodsVO{}, nil
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

func (s *goodsService) Update(dto model.GoodsDTO) error {
	var g model.Goods
	// 数据转换
	err := bean.SimpleCopyProperties(&g, dto)
	if err != nil {
		log.Println(err)
		return code.ConvertErr
	}
	g.UpdatedAt = model.LocalTime(time.Now())
	if e := s.dao.Update(g); e != nil {
		log.Println(e)
		return code.DBErr
	}
	return nil
}

func (s *goodsService) DeleteWithPhysics(id int) error {
	if e := s.dao.Delete(id); e != nil {
		log.Println(e)
		return code.DBErr
	}
	return nil
}

func (s *goodsService) DeleteWithLogic(id int) error {
	g, e := s.FindGoodsByID(id)
	if e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return code.RecordNotFound
		}
		log.Println(e)
		return code.DBErr
	}
	now := model.LocalTime(time.Now())
	g.DeletedAt = &now
	e = s.dao.Update(g)
	if e != nil {
		log.Println(e)
		return code.DBErr
	}
	return nil
}
