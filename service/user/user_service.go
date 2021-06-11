package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/dao"
	"seckill/infra/code"
	"seckill/infra/secret"
	"seckill/infra/utils/bean"
	"seckill/model"
	"strconv"
	"sync"
	"time"
)

var (
	once sync.Once
)

// service.IUserService 接口实现
type userService struct {
	dao dao.IUserDao
}

// NewUserService 创建一个 service.IUserService 接口实例
func NewUserService() *userService {
	return &userService{
		dao: dao.UserDao,
	}
}

// SingleUserService service.IUserService 接口单例模式
func SingleUserService() (s *userService) {
	once.Do(func() {
		s = NewUserService()
	})
	return
}

func (s *userService) Register(registerUser model.RegisterUser) error {
	// 数据转换
	user := model.User{}
	if e := bean.SimpleCopyProperties(&user, registerUser); e != nil {
		log.Println(code.ConvertErr.Error(), e.Error())
		return code.ConvertErr
	}
	// 查重
	u, e := s.FindByUsername(user.Username)
	if e != nil && !errors.Is(e, gorm.ErrRecordNotFound) && !errors.Is(e, code.RecordNotFound){
		return code.DBErr
	}
	if u.ID != 0 {
		return code.UsernameExistedErr
	}
	return s.dao.Insert(user)
}

func (s *userService) Login(loginUser model.LoginUser) (token string, e error) {
	// 数据转换
	var user model.User
	if e = bean.SimpleCopyProperties(&user, loginUser); e != nil {
		log.Println(code.ConvertErr.Error(), e.Error())
		return
	}
	// 通过 username 到数据库查询数据来对比
	oldUser, err := s.FindByUsername(user.Username)
	if err != nil || oldUser.Password != user.Password {
		if err != nil {
			log.Println(e)
		}
		e = code.AuthErr
		return
	}
	token, e = s.generateToken(oldUser)
	return
}

// 生成 Token
func (s *userService) generateToken(user model.User) (string, error) {
	if user.Username == "" {
		return "", errors.New("请先注册")
	}
	// 查出最新的用户数据
	user, err := s.FindByUsername(user.Username)
	if err != nil {
		return "", err
	}
	// 签发 JWT
	j := secret.NewJWT()
	expiresTime := time.Now().Add(secret.ExpiresTime * time.Second).Unix()
	claims := secret.CustomClaims{
		UserId: user.ID,
		Username: user.Username,
		Password: user.Password,
		Kind:     user.Kind,
		StandardClaims: jwt.StandardClaims{
			Audience:  user.Username,              // 受众
			ExpiresAt: expiresTime,                // 失效时间
			Id:        strconv.Itoa(int(user.ID)), // 编号
			IssuedAt:  time.Now().Unix(),          // 签发时间
			Issuer:    secret.Issuer,              // 签发人
			NotBefore: time.Now().Unix(),          // 生效时间
			Subject:   "auth",                     // 主题
		},
	}
	// 获取 token
	token, err := j.CreateToken(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *userService) FindByUsername(username string) (model.User, error) {
	return s.dao.QueryByUsername(username)
}

func (s *userService) Logout(token string) {
	if len(token) > 0 {
		// 退出时使 token 失效
		j := secret.NewJWT()
		_ = j.InvalidToken(token)
	}
}