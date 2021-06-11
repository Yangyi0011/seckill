package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/infra/utils/response"
	"seckill/model"
	"seckill/service"
	"sync"
)

var (
	once sync.Once
)

type UserHandler struct {
	userService service.IUserService
}

// NewUserHandler 创建一个 UserHandler 实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.UserService,
	}
}

// SingleUserHandler UserHandler 单例模式
func SingleUserHandler() (h *UserHandler) {
	once.Do(func() {
		h = NewUserHandler()
	})
	return
}

// Register go doc
// @Summary 用户注册
// @Description 注册用户并保存到数据库
// @Tags 用户管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param registerUser body model.RegisterUser true "注册用户"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 500 object model.Result 注册失败
// @Router /api/user/register [post]
func (h *UserHandler) Register(ctx *gin.Context) {
	registerUser := model.RegisterUser{}
	result := model.Result{}
	// 数据绑定
	if e := ctx.BindJSON(&registerUser); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	if e := h.userService.Register(registerUser); e != nil {
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

// Login go doc
// @Summary 用户登录
// @Description 用户登录签发 JWT
// @Tags 用户管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param loginUser body model.LoginUser true "用户"
// @Success 200 object model.Result 登录成功
// @Failure 400 object model.Result 请求参数有误
// @Failure 500 object model.Result 登录失败
// @Router /api/user/login [post]
func (h *UserHandler) Login(ctx *gin.Context) {
	loginUser := model.LoginUser{}
	result := model.Result{}
	// 数据绑定
	if e := ctx.BindJSON(&loginUser); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	token, e := h.userService.Login(loginUser)
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

// Logout go doc
// @Summary 退出登录
// @Description 用户退出登录，清除登录 token
// @Tags 用户管理
// @version 1.0
// @Accept json
// @Produce  json
// @Success 200 object model.Result 登录成功
// @Router /api/user/logout [post]
func (h *UserHandler) Logout(ctx *gin.Context) {
	auth := ctx.Request.Header.Get("Authorization")
	h.userService.Logout(auth)
	result := model.Result{}
	result.Code = http.StatusOK
	result.Message = "退出成功"
	ctx.Header("Authorization", "")
	response.Success(ctx, result)
	return
}