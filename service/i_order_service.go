package service

type IOrderService interface {
	SecondKill(userId, goodsId int) (e error)
	GetOrderId(userId, goodsId int) (orderId string, err error)
}
