package handler

import (
	"log"
	"seckill/handler/goods"
	"seckill/handler/user"
)

var (
	UserHandler *user.UserHandler
	GoodsHandler *goods.GoodsHandler
)

func init() {
	UserHandler = user.SingleUserHandler()
	if UserHandler == nil {
		log.Fatalln("Err: UserHandler is not find")
	}
	GoodsHandler = goods.SingleUserHandler()
	if UserHandler == nil {
		log.Fatalln("Err: GoodsHandler is not find")
	}
}
