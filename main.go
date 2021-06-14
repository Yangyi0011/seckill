package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"seckill/router"
)

func main() {
	// 启用发布模式
	gin.SetMode(gin.ReleaseMode)
	// 启动
	err := router.InitRouter().Run()
	if err != nil {
		log.Fatal("项目启动错误：", err)
		return
	}
}
