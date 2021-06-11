package service

import "seckill/model"

type IUserService interface {
	// Register 用户注册
	Register(registerUser model.RegisterUser) error
	// Login 用户登录
	Login(loginUser model.LoginUser) (token string, e error)
	// FindByUsername 通过 username 查询用户信息
	FindByUsername(username string) (user model.User, e error)
	// Logout 用户退出登录
	Logout(token string)
}
