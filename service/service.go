package service

import (
	"log"
	"seckill/service/goods"
	"seckill/service/user"
)

var (
	UserService IUserService
	GoodsService IGoodsService
)

func init() {
	UserService = user.SingleUserService()
	if UserService == nil {
		log.Fatalln("Err: UserService is not find")
	}
	GoodsService = goods.SingleGoodsService()
	if UserService == nil {
		log.Fatalln("Err: GoodsService is not find")
	}
}
