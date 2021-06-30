package router

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"reflect"
	"seckill/handler"
	"seckill/infra/utils/response"
	"seckill/middleware"
	"seckill/model"
	"seckill/mq"
	"seckill/service/goods"
	"seckill/service/order"
	"seckill/service/user"

	// swagger 文档生成在本项目中的目录，必须导入这个目录文档才能正常显示
	_ "seckill/docs" // docs is generated by Swag CLI, you have to import it.
)

var (
	// 全局路由引擎
	myRouter *gin.Engine

	goodsHandler *handler.GoodsHandler
	userHandler *handler.UserHandler
	orderHandler *handler.OrderHandler
)

func init() {
	myRouter = gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册 model.LocalTime 类型的自定义校验规则
		v.RegisterCustomTypeFunc(ValidateJSONDateType, model.LocalTime{})
	}
}

// ValidateJSONDateType 解决验证器 binding:"required" 无法正常工作的问题
func ValidateJSONDateType(field reflect.Value) interface{} {
	if field.Type() == reflect.TypeOf(model.LocalTime{}) {
		timeStr := field.Interface().(model.LocalTime).String()
		// 0001-01-01 00:00:00 是 go 中 time.Time 类型的空值
		// 这里返回 Nil 则会被 validator 判定为空值，而无法通过 `binding:"required"` 规则
		if timeStr == "0001-01-01 00:00:00" {
			return nil
		}
		return timeStr
	}
	return nil
}

// InitRouter 初始化路由器
func InitRouter() * gin.Engine{
	//myRouter.Use(middleware.CostumerLog())
	myRouter.Use(middleware.Cors())
	myRouter.Use(middleware.SysLimit())
	myRouter.Use(middleware.UserLimit())

	initService()
	initHandler()
	swaggerRouter()
	customRouter()

	mq.Run()
	return myRouter
}

// 初始化 service 层
func initService() {
	defer mq.Init()
	goods.InitService()
	order.InitService()
	user.InitService()
}

// 初始化 handler 层
func initHandler() {
	goodsHandler = handler.NewGoodsHandler()
	userHandler = handler.NewUserHandler()
	orderHandler = handler.NewOrderHandler()
}

// SwaggerRouter swagger 路由
// 文档访问地址：http://localhost:8080/swagger/index.html
// @title Gin swagger
// @version 1.0
// @description Gin swagger 示例项目
// @contact.name 清影
// @contact.email 1024569696@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
func swaggerRouter() {
	// The url pointing to API definition
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	myRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
}

func customRouter() {
	// 首页
	indexGroup := myRouter.Group("/")
	{
		indexGroup.Any("", func(ctx *gin.Context) {
			result := model.Result{
				Code:    http.StatusOK,
				Message: "index 空空如也~",
				Data:    nil,
			}
			response.Success(ctx, result)
			return
		})
	}
	// 业务接口
	api := myRouter.Group("/api")
	{
		userGroup := api.Group("/user")
		{
			userGroup.POST("/register", userHandler.Register)
			userGroup.POST("/login", userHandler.Login)
			userGroup.POST("/logout", userHandler.Logout)
		}

		goodsGroup := api.Group("/goods")
		{
			goodsGroup.GET("/:id", goodsHandler.QueryGoodsVOByID)
			goodsGroup.POST("/list", goodsHandler.QueryGoodsVOByCondition)
			goodsGroup.POST("/", middleware.Auth(), middleware.SellerAuth(), goodsHandler.Insert)
			goodsGroup.PUT("/", middleware.Auth(), middleware.SellerAuth(), goodsHandler.Update)
			goodsGroup.DELETE("/:id", middleware.Auth(), middleware.SellerAuth(), goodsHandler.Delete)
			goodsGroup.POST("/seckillInit", middleware.Auth(), middleware.SellerAuth(), goodsHandler.SecondKillGoodsInit)
		}

		seckill := api.Group("/seckill")
		{
			seckill.POST("/", middleware.Auth(), orderHandler.SecondKill)
			seckill.GET("/:goodsId", middleware.Auth(), orderHandler.GetSecondKillResult)
		}

		orderGroup := api.Group("/order")
		{
			orderGroup.GET("/:id", middleware.Auth(), orderHandler.QueryByID)
			orderGroup.POST("/list", middleware.Auth(), orderHandler.QueryByCondition)
		}
	}
}