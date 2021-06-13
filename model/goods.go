package model

import (
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
	PageDTO
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	StartTime LocalTime `json:"startTime"`
	EndTime   LocalTime `json:"endTime"`
	UserId    uint      `json:"userId" swaggerignore:"true"`
}

// TableName 继承接口指定表名
func (g Goods) TableName() string {
	return "goods"
}

// ToVO Goods 转为 GoodsVO
func (g Goods) ToVO() (vo GoodsVO, e error) {
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


func init() {
	g := Goods{}
	// 表不存在的时候创建表
	if !db.DB.HasTable(g) {
		log.Printf("正在创建 %s 表\n", g.TableName())
		db.DB.Debug().CreateTable(g)
	}
}
