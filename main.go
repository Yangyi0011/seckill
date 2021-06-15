package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os"
	"seckill/router"
)

func main() {
	// 收集 http 请求日志
	logfile, err := os.Create("./gin_http.log")
	if err != nil {
		fmt.Println("Could not create log file")
	}
	// 启用发布模式
	gin.SetMode(gin.ReleaseMode)
	// 指定默认的日志输出路径
	gin.DefaultWriter = io.MultiWriter(logfile)

	// 启动
	err = router.InitRouter().Run()
	if err != nil {
		log.Fatal("项目启动错误：", err)
		return
	}
}
