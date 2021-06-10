package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/infra/code"
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
	Model `swaggerignore:"true"`
	Username   string `json:"username" gorm:"type:varchar(20);comment:'用户名';unique_index:idx_users_username"`
	Password   string `json:"password" gorm:"type:varchar(32);comment:'用户密码'"`
	Kind       int8    `json:"kind" gorm:"type:tinyint(1);comment:'用户类别（0-客户，1-商家）'"`
}

// TableName 继承接口指定表名
func (u User)TableName() string {
	return "users"
}

func init() {
	u := User{}
	// 表不存在的时候创建表
	if !db.DB.HasTable(u) {
		log.Printf("正在创建 %s 表\n", u.TableName())
		db.DB.Debug().CreateTable(u)
	}
}

// Insert 添加用户
func (u User) Insert() error {
	u.CreatedAt = LocalTime(time.Now())
	if e := db.DB.Debug().Create(&u).Error; e != nil {
		return code.DBErr
	}
	return nil
}

// QueryByUsername 通过 username 查询用户信息
func (u User) QueryByUsername() (User, error) {
	var user User
	if e := db.DB.Debug().Where("username = ?", u.Username).Take(&user).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return user, code.RecordNotFound
		}
		return user, code.DBErr
	}
	return user, nil
}
