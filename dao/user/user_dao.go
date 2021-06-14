package user

import (
	"errors"
	"github.com/jinzhu/gorm"
	"seckill/infra/code"
	"seckill/infra/db"
	"seckill/model"
	"time"
)

type userDao struct {
	db *gorm.DB
}

// NewUserDao 创建一个 IUserDao 接口的实例
func NewUserDao () *userDao {
	return &userDao{db: db.DB}
}

// Insert 添加用户
func (d *userDao) Insert(u model.User) error {
	u.CreatedAt = model.LocalTime(time.Now())
	if e := db.DB.Create(&u).Error; e != nil {
		return code.DBErr
	}
	return nil
}

// QueryByUsername 通过 username 查询用户信息
func (d *userDao) QueryByUsername(username string) (model.User, error) {
	var user model.User
	if e := db.DB.Where("username = ?", username).Take(&user).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return user, code.RecordNotFoundErr
		}
		return user, code.DBErr
	}
	return user, nil
}