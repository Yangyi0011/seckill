package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/infra/code"
	"seckill/infra/utils/request"
	"seckill/infra/utils/response"
	"seckill/model"
	"seckill/service"
	"strconv"
)

type OrderHandler struct {
	orderService service.IOrderService
}

// NewOrderHandler 创建一个 UserHandler 实例
func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderService: service.OrderService,
	}
}

// SecondKill go doc
// @Summary 秒杀
// @Description 对商品进行秒杀
// @Tags 商品秒杀
// @version 1.0
// @Accept json
// @Produce  json
// @Param goods body model.GoodsDTO true "秒杀商品传输信息"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 401 object model.Result 需要登录
// @Failure 403 object model.Result 没有操作权限
// @Failure 500 object model.Result 操作失败
// @Router /api/seckill [post]
func (h *OrderHandler) SecondKill(ctx *gin.Context) {
	result := model.Result{}
	dto := new(
		struct {
			GoodsId int `json:"goodsId" binding:"required"`
		},
	)
	// 数据绑定
	if e := ctx.BindJSON(&dto); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	// 获取当前用户信息
	claims, _ := request.GetCurrentCustomClaims(ctx)
	// 参与秒杀
	if e := h.orderService.SecondKill(int(claims.UserId), dto.GoodsId); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	response.Success(ctx, result)
	return
}

// GetSecondKillResult go doc
// @Summary 查询秒杀结果
// @Description 通过商品id查询对该商品的秒杀结果
// @Tags 商品秒杀
// @version 1.0
// @Accept json
// @Produce  json
// @Param goodsId path int true "goodsId" "商品id"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 401 object model.Result 需要登录
// @Failure 403 object model.Result 没有操作权限
// @Failure 500 object model.Result 操作失败
// @Router /api/seckill/{goodsId} [GET]
func (h *OrderHandler) GetSecondKillResult(ctx *gin.Context) {
	var (
		result model.Result
		e error
		goodsId int
		userId int
		res model.SecondKillResult
	)
	goodsIdStr := ctx.Param("goodsId")
	goodsId, e = strconv.Atoi(goodsIdStr)
	if  e != nil || goodsId == 0 {
		result.Code = http.StatusBadRequest
		result.Message = code.RequestParamErr.Error()
		response.Fail(ctx, result)
		return
	}
	claims, err := request.GetCurrentCustomClaims(ctx)
	if err != nil {
		result.Code = http.StatusUnauthorized
		result.Message = "请先登录"
		response.Fail(ctx, result)
		return
	}
	userId = int(claims.UserId)
	if res, e = h.orderService.GetSecondKillResult(userId, goodsId); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	result.Data = res
	response.Success(ctx, result)
}

// QueryByCondition go doc
// @Summary 查询订单列表
// @Description 通过查询条件查询订单列表数据
// @Tags 订单管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param condition body model.OrderInfoQueryCondition true "订单查询条件"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 401 object model.Result 需要登录
// @Failure 403 object model.Result 没有操作权限
// @Failure 500 object model.Result 操作失败
// @Router /api/order/list [POST]
func (h *OrderHandler) QueryByCondition(ctx *gin.Context) {
	var (
		result    model.Result
		e         error
		condition model.OrderInfoQueryCondition
		list      []model.OrderInfoVO
	)
	// 数据绑定
	if e = ctx.BindJSON(&condition); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	// 默认查询当前登录人的数据
	claims, _ := request.GetCurrentCustomClaims(ctx)
	condition.UserId = claims.UserId
	if list, e = h.orderService.GetOrderInfoVOList(condition); e != nil {
		result.Code = http.StatusInternalServerError
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	result.Data = list
	response.Success(ctx, result)
	return
}

// QueryByID go doc
// @Summary 查询订单
// @Description 通过 id 查询订单数据
// @Tags 订单管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param id path int true "id" "订单编号id"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 500 object model.Result 操作失败
// @Router /api/order/{id} [GET]
func (h *OrderHandler) QueryByID(ctx *gin.Context) {
	result := model.Result{}
	idStr := ctx.Param("id")
	if idStr == "" {
		result.Code = http.StatusBadRequest
		result.Message = code.RequestParamErr.Error()
		response.Fail(ctx, result)
		return
	}
	claims, _ := request.GetCurrentCustomClaims(ctx)
	vo, e := h.orderService.GetOrderInfoVO(idStr, int(claims.UserId))
	if e != nil {
		result.Code = http.StatusBadRequest
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	result.Data = vo
	response.Success(ctx, result)
	return
}
