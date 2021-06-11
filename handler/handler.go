package handler

import (
	"log"
	"seckill/handler/user"
)

var (
	UserHandler *user.UserHandler
)

func init() {
	UserHandler = user.SingleUserHandler()
	if UserHandler == nil {
		log.Fatalln("Err: UserHandler is not find")
	}
}
