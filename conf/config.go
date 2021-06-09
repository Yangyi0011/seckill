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
}

// Datasource 数据源配置信息
type Datasource struct {
	DriverName string `yaml:"driverName"`
	Host string `yaml:"host"`
	BaseName string `yaml:"baseName"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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