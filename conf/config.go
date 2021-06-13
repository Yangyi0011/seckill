package conf

import (
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var (
	// Config 配置信息全局变量
	Config *AppConfig
)

// 初始化配置信息
func init() {
	Config = &AppConfig{}
	Config.reloadConfig()
}

// AppConfig yaml 配置信息绑定
type AppConfig struct {
	App `yaml:"app"`
}

// App 系统配置信息
type App struct {
	Datasource `yaml:"datasource"`
	Redis `yaml:"redis"`
	Order `yaml:"order"`
	RateLimit `yaml:"rate_limit"`
}

// Datasource 数据源配置信息
type Datasource struct {
	DriverName string `yaml:"driverName"`
	Host string `yaml:"host"`
	BaseName string `yaml:"baseName"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Redis redis 配置信息
type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
}

// Order 订单配置信息
type Order struct {
	Expiration int64 `yaml:"expiration"`
}

// RateLimit 限流配置信息
type RateLimit struct {
	Time  int64 `yaml:"time"`
	Count int64 `yaml:"count"`
}

// 装载配置信息
func (appConfig *AppConfig) reloadConfig() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Panicln("config.yaml 读取错误：", err)
	}
	err = yaml.Unmarshal(yamlFile, appConfig)
	if err != nil {
		log.Panicln("config.yaml 解析错误：", err)
	}
}