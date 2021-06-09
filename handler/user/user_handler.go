package user

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"seckill/infra/code"
	"seckill/infra/secret"
	"seckill/infra/utils/bean"
	"seckill/infra/utils/response"
	"seckill/model"
	"strconv"
	"time"
)

// Register go doc
// @Summary 用户注册
// @Description 注册用户并保存到数据库
// @Tags 用户管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param article body model.RegisterUser true "注册用户"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 500 object model.Result 注册失败
// @Router /api/user/register [post]
func Register(ctx *gin.Context) {
	registerUser := model.RegisterUser{}
	result := model.Result{}
	// 数据绑定
	if e := ctx.BindJSON(&registerUser); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	if e := register(registerUser); e != nil {
		result.Code = http.StatusInternalServerError
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	result.Code = http.StatusOK
	result.Message = "注册成功"
	response.Success(ctx, result)
	return
}

func register(registerUser model.RegisterUser) error {
	// 数据转换
	user := model.User{}
	if e := bean.SimpleCopyProperties(&user, registerUser); e != nil {
		log.Println(code.ConvertErr.Error(), e.Error())
		return code.ConvertErr
	}
	// 查重
	u, e := user.QueryByUsername()
	if !errors.Is(e, gorm.ErrRecordNotFound) {
		return code.DBErr
	}
	if u.ID != 0 {
		return code.UsernameExistedErr
	}
	// 数据保存
	if e = user.Insert(); e != nil {
		return e
	}
	return nil
}

// Login go doc
// @Summary 用户登录
// @Description 用户登录签发 JWT
// @Tags 用户管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param article body model.LoginUser true "用户"
// @Success 200 object model.Result 登录成功
// @Failure 400 object model.Result 请求参数有误
// @Failure 500 object model.Result 登录失败
// @Router /api/user/login [post]
func Login(ctx *gin.Context) {
	loginUser := model.LoginUser{}
	result := model.Result{}
	// 数据绑定
	if e := ctx.BindJSON(&loginUser); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	// 数据转换
	var user model.User
	if e := bean.SimpleCopyProperties(&user, loginUser); e != nil {
		log.Println(code.ConvertErr.Error(), e.Error())
		result.Code = http.StatusInternalServerError
		result.Message = code.ConvertErr.Error()
		response.Fail(ctx, result)
		return
	}
	// 通过 username 到数据库查询数据来对比
	u, e := user.QueryByUsername()
	if e != nil || u.Password != user.Password {
		if e != nil {
			log.Println(e)
		}
		result.Code = http.StatusInternalServerError
		result.Message = code.AuthErr.Error()
		response.Fail(ctx, result)
		return
	}
	token, e := generateToken(user)
	if e != nil {
		result.Code = http.StatusInternalServerError
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	result.Message = "登录成功"
	result.Data = map[string]interface{}{
		"token":token,
	}
	result.Code = http.StatusOK
	ctx.Header("Authorization", token)
	response.Success(ctx, result)
	return
}

// 生成 Token
func generateToken(user model.User) (string, error) {
	if user.Username == "" {
		return "", errors.New("请先注册")
	}
	// 查出最新的用户数据
	user, err := user.QueryByUsername()
	if err != nil {
		return "", err
	}
	// 签发 JWT
	j := secret.NewJWT()
	expiresTime := time.Now().Add(secret.ExpiresTime * time.Second).Unix()
	claims := secret.CustomClaims{
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

// Logout go doc
// @Summary 退出登录
// @Description 用户退出登录，清除登录 token
// @Tags 用户管理
// @version 1.0
// @Accept json
// @Produce  json
// @Success 200 object model.Result 登录成功
// @Router /api/user/logout [post]
func Logout(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	if len(auth) > 0 {
		// 退出时使 token 失效
		j := secret.NewJWT()
		_ = j.InvalidToken(auth)
	}
	result := model.Result{}
	result.Code = http.StatusOK
	result.Message = "退出成功"
	ctx.Header("Authorization", "")
	response.Success(ctx, result)
	return
}