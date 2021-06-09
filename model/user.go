package model

import (
	"github.com/jinzhu/gorm"
	"log"
	"seckill/infra/db"
	"time"
)

const (
	// NormalCustomer 买家（顾客）
	NormalCustomer = 0
	// NormalSeller 卖家（商家）
	NormalSeller = 1
)

// LoginUser 用户登录DTO
type LoginUser struct {
	Username string `json:"username" example:"tom" binding:"required"`
	Password string `json:"password" example:"123" binding:"required"`
}

// RegisterUser 用户注册DTO
type RegisterUser struct {
	LoginUser
	Kind int `json:"kind" example:"0"`
}

// User 用户模型DO
type User struct {
	gorm.Model `swaggerignore:"true"`
	Username   string `json:"username" gorm:"type:varchar(20);comment:'用户名';unique_index:idx_users_username"`
	Password   string `json:"password" gorm:"type:varchar(32);comment:'用户密码'"`
	Kind       int    `json:"kind" gorm:"type:tinyint(1);comment:'用户类别（0-客户，1-商家）'"`
}

// TableName 继承接口指定表名
func TableName() string {
	return "users"
}

func init() {
	// 表不存在的时候创建表
	if !db.DB.HasTable(User{}) {
		log.Printf("正在创建 %s 表\n", TableName())
		db.DB.Debug().CreateTable(User{})
	}
}

// Insert 添加用户
func (user User) Insert() error {
	user.CreatedAt = time.Now()
	err := db.DB.Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

// QueryByUsername 通过 username 查询用户信息
func (user User) QueryByUsername() (User, error) {
	err := db.DB.Debug().Where("username = ?", user.Username).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
