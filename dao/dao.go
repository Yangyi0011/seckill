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
	UserDao = user.NewUserDao()
	if UserDao == nil {
		log.Fatalln("Err: UserDao is not find")
	}
	GoodsDao = goods.NewGoodsDao()
	if UserDao == nil {
		log.Fatalln("Err: GoodsDao is not find")
	}
}
