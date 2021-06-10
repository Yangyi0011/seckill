package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/infra/code"
	"seckill/infra/db"
	"seckill/infra/utils/bean"
	"time"
)

// 商品常量
const (
	Ended      int8 = -1 // 已结束
	NotStarted int8 = 0  // 未开始
	OnGoing    int8 = 1  // 进行中
	SoldOut    int8 = 2  // 已售罄
)

// Goods 秒杀商品模型
type Goods struct {
	Model
	Name        string    `gorm:"type:varchar(50);comment:'商品名称'"`
	Img         string    `gorm:"type:varchar(255);comment:'商品图片url'"`
	OriginPrice float64   `gorm:"type:decimal(20,2);comment:'原价'"`
	Price       float64   `gorm:"type:decimal(20,2);comment:'现价'"`
	Amount      int       `gorm:"comment:'商品总量'"`
	Stock       int       `gorm:"comment:'商品剩余数量'"`
	StartTime   LocalTime `gorm:"type:datetime;comment:'秒杀开始时间'"`
	EndTime     LocalTime `gorm:"type:datetime;comment:'秒杀结束时间'"`
	UserId      uint      `gorm:"comment:'创建人'"`
}

// GoodsDTO 商品数据传输模型
type GoodsDTO struct {
	Model
	Name        string    `json:"name" binding:"required"`
	Img         string    `json:"img"`
	OriginPrice float64   `json:"originPrice" minimum:"0.0" default:"0.0"`
	Price       float64   `json:"price" binding:"required" minimum:"0.0" default:"0.0"`
	Amount      int       `json:"amount" binding:"required" minimum:"1" default:"0"`
	Stock       int       `json:"stock" binding:"required" minimum:"0" default:"0"`
	StartTime   LocalTime `json:"startTime" binding:"required"`
	EndTime     LocalTime `json:"endTime" binding:"required"`
	UserId      uint      `json:"userId" swaggerignore:"true"`
}

// GoodsVO 商品视图模型
type GoodsVO struct {
	Model
	Name        string    `json:"name"`
	Img         string    `json:"img"`
	OriginPrice float64   `json:"originPrice"`
	Price       float64   `json:"price"`
	Amount      int       `json:"amount"`
	Stock       int       `json:"stock"`
	StartTime   LocalTime `json:"startTime"`
	EndTime     LocalTime `json:"endTime"`
	UserId      uint      `json:"userId"`
	Duration    int64     `json:"duration"`
	Status      int8      `json:"status"`
}

// GoodsQueryCondition 商品查询条件
type GoodsQueryCondition struct {
	Model
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	StartTime LocalTime `json:"startTime"`
	EndTime   LocalTime `json:"endTime"`
}

// TableName 继承接口指定表名
func (g Goods) TableName() string {
	return "goods"
}

func init() {
	g := Goods{}
	// 表不存在的时候创建表
	if !db.DB.HasTable(g) {
		log.Printf("正在创建 %s 表\n", g.TableName())
		db.DB.Debug().CreateTable(g)
	}
}


// ToVO 把 Goods 转为 GoodsVO
func (g Goods) ToVO() (GoodsVO, error) {
	var vo GoodsVO
	if e := bean.SimpleCopyProperties(&vo, g); e != nil {
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

// QueryGoodsByID 通过 id 查询一条模型数据
func (g Goods) QueryGoodsByID() (Goods, error) {
	goods := Goods{}
	if e := db.DB.Debug().Where("id = ?", g.ID).Take(&goods).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return goods, code.RecordNotFound
		}
		log.Println(e)
		return goods, code.DBErr
	}
	return goods, nil
}

//// QueryGoodsVOByID 通过 id 查询一条视图数据
func (g Goods) QueryGoodsVOByID() (GoodsVO, error) {
	goods, e := g.QueryGoodsByID()
	if e != nil {
		return GoodsVO{}, e
	}
	vo, e := goods.ToVO()
	if e != nil {
		return GoodsVO{}, e
	}
	return vo, nil
}

// QueryByCondition 通过条件查询多条数据
//func (g Goods) QueryByCondition(c GoodsQueryCondition) ([]GoodsVO, error) {
//	//db.DB.Exec("select * from")
//	return []GoodsVO{}, nil
//}

// Insert 插入数据
func (g Goods) Insert(dto GoodsDTO) error {
	// 数据转换
	err := bean.SimpleCopyProperties(&g, dto)
	if err != nil {
		return err
	}
	g.CreatedAt = LocalTime(time.Now())
	if e := db.DB.Create(&g).Error; e != nil {
		log.Println(e)
		return code.DBErr
	}
	return nil
}

//
//// Update 更新数据
//func (g Goods) Update() error {
//	var old Goods
//	if e := db.DB.Debug().Take(&old).Error; e != nil {
//		if errors.Is(e, gorm.ErrRecordNotFound) {
//			return code.RecordNotFound
//		}
//		log.Println(e)
//		return code.DBErr
//	}
//	if old.ID == 0 {
//		return code.RecordNotFound
//	}
//	if e := db.DB.Debug().Updates(&g).Error; e != nil {
//		return code.DBErr
//	}
//	return nil
//}
//
//// DeleteByPhysics 物理删除数据
//func (g Goods) DeleteByPhysics() error {
//	var old Goods
//	if e := db.DB.Debug().Take(&old).Error; e != nil {
//		if errors.Is(e, gorm.ErrRecordNotFound) {
//			return code.RecordNotFound
//		}
//		log.Println(e)
//		return code.DBErr
//	}
//	if old.ID == 0 {
//		return code.RecordNotFound
//	}
//	if e := db.DB.Debug().Delete(&g).Error; e != nil {
//		log.Println(e)
//		return code.DBErr
//	}
//	return nil
//}
//
//// DeleteByLogic 逻辑删除数据
//func (g Goods) DeleteByLogic() error {
//	var old Goods
//	if e := db.DB.Debug().Take(&old).Error; e != nil {
//		if errors.Is(e, gorm.ErrRecordNotFound) {
//			return code.RecordNotFound
//		}
//		log.Println(e)
//		return code.DBErr
//	}
//	if old.ID == 0 {
//		return code.RecordNotFound
//	}
//	if e := db.DB.Debug().Model(&g).Update("deleted_at", time.Now()).Error; e != nil {
//		log.Println(e)
//		return code.DBErr
//	}
//	return nil
//}
