package dao

import "seckill/model"

type IOrderDao interface {
	QueryOrderInfoByOrderId(id string) (o model.OrderInfo, e error)
	QueryByCondition(c model.OrderInfoQueryCondition) (list []model.OrderInfo, e error)
	Insert(o model.OrderInfo) error
	Update(o model.OrderInfo) error
	Delete(id string) error
	CreateOrder(o model.OrderInfo) error
	CloseOrder(orderInfo model.OrderInfo) error
}
