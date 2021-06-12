package dao

import (
	"log"
	"seckill/dao/goods"
	"seckill/dao/order"
	"seckill/dao/user"
)

var (
	UserDao IUserDao
	GoodsDao IGoodsDao
	OrderDao IOrderDao
)
func init() {
	GoodsDao = goods.NewGoodsDao()
	if GoodsDao == nil {
		log.Fatalln("Err: GoodsDao is not find")
	}
	OrderDao = order.NewOrderDao()
	if OrderDao == nil {
		log.Fatalln("Err: OrderDao is not find")
	}
	UserDao = user.NewUserDao()
	if UserDao == nil {
		log.Fatalln("Err: UserDao is not find")
	}
}
