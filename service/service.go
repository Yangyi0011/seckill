package service

import (
	"log"
	"seckill/service/user"
)

var (
	UserService IUserService
)

func init() {
	UserService = user.SingleUserService()
	if UserService == nil {
		log.Fatalln("Err: UserService is not find")
	}
}
