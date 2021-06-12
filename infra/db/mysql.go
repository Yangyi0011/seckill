package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/conf"
)

var (
	// DB 数据源连接变量
	DB *gorm.DB
)

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&loc=Local",
		conf.Config.Username, conf.Config.Datasource.Password, conf.Config.Datasource.Host, conf.Config.Datasource.BaseName)
	var err error
	DB, err = gorm.Open(conf.Config.DriverName, dsn)
	if err != nil {
		log.Panicln("数据源连接错误：", err.Error())
	}
}
