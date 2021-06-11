package dao

import (
	"log"
	"seckill/dao/user"
)

var (
	UserDao IUserDao
)
func init() {
	UserDao = user.SingleUserDao()
	if UserDao == nil {
		log.Fatalln("Err: UserDao is not find")
	}
}
