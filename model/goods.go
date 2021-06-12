package model

import (
	"log"
	"seckill/infra/db"
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
