package dao

import (
	"seckill/model"
)

type IUserDao interface {
	// Insert 添加用户
	Insert(user model.User) error
	// QueryByUsername 通过 username 查询用户信息
	QueryByUsername(username string) (model.User, error)
}
