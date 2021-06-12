package service

import "seckill/model"

type IOrderService interface {
	// SecondKill 秒杀
	SecondKill(userId, goodsId int) (e error)

	// GetOrderId 获取订单 id
	GetOrderId(userId, goodsId int) (orderId string, err error)

	// GetOrderInfo 通过订单编号获取订单信息
	GetOrderInfo(orderId string) (o model.OrderInfo, e error)

	// CreateOrder 创建订单
	CreateOrder(userId int, goodsId int) error

	// CreateOrderCache 创建订单缓存信息
	CreateOrderCache(order model.OrderInfo) (err error)

	// DeleteOrderCache 删除订单缓存
	DeleteOrderCache(order model.OrderInfo) (err error)

	// CloseOrder 关闭订单
	CloseOrder(userId int, orderId string) (e error)

	// UnLock 解锁
	UnLock(userId, goodsId int) (err error)
}
