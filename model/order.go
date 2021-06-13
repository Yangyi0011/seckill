package model

import (
	"log"
	"seckill/conf"
	"seckill/infra/db"
	"time"
)

const (
	// Closed 订单已关闭
	Closed int8 = -1
	// Unpaid 订单未支付
	Unpaid int8 = 0
	// Paying 订单支付中
	Paying int8 = 1
	// Paid 订单已支付
	Paid int8 = 2

	// SecondKilling 秒杀排队中
	SecondKilling = 0
	// SecondKillOk 秒杀成功
	SecondKillOk = 1
)

// Order 订单
type Order struct {
	Model
	OrderId string `gorm:"type:varchar(25);comment:'订单id';index:idx_orders_order_id"`
	UserId  uint   `gorm:"type:int;comment:'下单用户id';index:idx_orders_user_id"`
	GoodsId uint   `gorm:"type:int;comment:'商品id';index:idx_orders_goods_id"`
}

// OrderInfo 订单信息
type OrderInfo struct {
	Model
	OrderId    string  `gorm:"type:varchar(25);comment:'订单id';index:idx_order_info_order_id"`
	UserId     uint    `gorm:"type:int;comment:'下单用户id';index:idx_order_info_user_id"`
	GoodsId    uint    `gorm:"type:int;comment:'商品id';index:idx_order_info_goods_id"`
	GoodsName  string  `gorm:"type:varchar(50);comment:'商品名称'"`
	GoodsImg   string  `gorm:"type:varchar(255);comment:'商品图片url'"`
	GoodsPrice float64 `gorm:"type:decimal(20,2);comment:'商品价格'"`
	PaymentId  int     `gorm:"type:int;comment:'支付id';index:idx_order_info_payment_id"`
	Status     int8    `gorm:"type:tinyint(1);comment:'订单状态';index:idx_order_info_status"`
}

// OrderInfoVO 订单信息视图模型
type OrderInfoVO struct {
	Model
	OrderId    string  `json:"orderId"`
	GoodsId    uint    `json:"goodsId"`
	GoodsName  string  `json:"goodsName"`
	GoodsImg   string  `json:"goodsImg"`
	GoodsPrice float64 `json:"goodsPrice"`
	Status     int8    `json:"status"`
	Duration   int64   `json:"duration"`
	Timeout    int64   `json:"timeout"`
}

// OrderInfoQueryCondition 商品查询条件
type OrderInfoQueryCondition struct {
	Model
	PageDTO
	OrderId    string  `json:"orderId"`
	UserId     uint    `json:"userId"`
	GoodsId    uint    `json:"goodsId"`
	GoodsName  string  `json:"goodsName"`
	GoodsPrice float64 `json:"goodsPrice"`
	PaymentId  int     `json:"paymentId"`
	Status     int8    `json:"status"`
}

// SecondKillResult 秒杀结果
type SecondKillResult struct {
	// 秒杀状态，0：排队中，1：成功，2：失败
	Status  int8   `json:"status"`
	// 订单编号
	OrderId string `json:"orderId"`
}

// TableName 继承接口指定表名
func (order Order) TableName() string {
	return "orders"
}

// TableName 继承接口指定表名
func (orderInfo OrderInfo) TableName() string {
	return "order_info"
}

// ToVO 把 OrderInfo 转为 OrderInfoVO
func (orderInfo OrderInfo) ToVO() OrderInfoVO {
	expiration := conf.Config.Order.Expiration
	orderVO := OrderInfoVO{
		Model: orderInfo.Model,
		OrderId:    orderInfo.OrderId,
		GoodsId:    orderInfo.GoodsId,
		GoodsName:  orderInfo.GoodsName,
		GoodsImg:   orderInfo.GoodsImg,
		GoodsPrice: orderInfo.GoodsPrice,
		Status:     orderInfo.Status,
		Duration:   0,
		Timeout:    expiration,
	}
	if orderVO.Status == 0 {
		orderVO.Duration = orderInfo.CreatedAt.Unix() + expiration - time.Now().Unix()
	}
	return orderVO
}

func init() {
	o := Order{}
	// 表不存在的时候创建表
	if !db.DB.HasTable(o) {
		log.Printf("正在创建 %s 表\n", o.TableName())
		db.DB.Debug().CreateTable(o)
	}
	oi := OrderInfo{}
	if !db.DB.HasTable(oi) {
		log.Printf("正在创建 %s 表\n", oi.TableName())
		db.DB.Debug().CreateTable(oi)
	}
}
