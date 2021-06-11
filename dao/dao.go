package dao

import (
	"log"
	"seckill/dao/goods"
	"seckill/dao/user"
)

var (
	UserDao IUserDao
	GoodsDao IGoodsDao
)
func init() {
	UserDao = user.SingleUserDao()
	if UserDao == nil {
		log.Fatalln("Err: UserDao is not find")
	}
	GoodsDao = goods.SingleGoodsDao()
	if UserDao == nil {
		log.Fatalln("Err: GoodsDao is not find")
	}
}
