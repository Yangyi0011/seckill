package main

import (
	"log"
	"seckill/router"
)

func main() {
	// 启动
	err := router.InitRouter().Run()
	if err != nil {
		log.Fatal("项目启动错误：", err)
		return
	}
}
