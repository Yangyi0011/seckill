package goods

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/infra/code"
	"seckill/infra/secret"
	"seckill/infra/utils/request"
	"seckill/infra/utils/response"
	"seckill/model"
	"seckill/service"
	"strconv"
	"sync"
)

var (
	once sync.Once
)

type GoodsHandler struct {
	goodsService service.IGoodsService
}

// NewGoodsHandler 创建一个 GoodsHandler 实例
func NewGoodsHandler() *GoodsHandler {
	return &GoodsHandler{
		goodsService: service.GoodsService,
	}
}

// SingleUserHandler UserHandler 单例模式
func SingleUserHandler() (h *GoodsHandler) {
	once.Do(func() {
		h = NewGoodsHandler()
	})
	return
}

// QueryGoodsVOByID go doc
// @Summary 查询商品
// @Description 通过 id 查询秒杀商品
// @Tags 商品管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param id path int true "id" "秒杀商品信息"
// @Success 200 object model.Result model.GoodsVO 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 500 object model.Result 操作失败
// @Router /api/goods/{id} [GET]
func (h *GoodsHandler) QueryGoodsVOByID(ctx *gin.Context) {
	result := model.Result{}
	idStr := ctx.Param("id")
	id, e := strconv.Atoi(idStr)
	if  e != nil || id == 0 {
		result.Code = http.StatusBadRequest
		result.Message = code.RequestParamErr.Error()
		response.Fail(ctx, result)
		return
	}

	g, e := h.goodsService.FindGoodsVOByID(id)
	if e != nil {
		result.Code = http.StatusBadRequest
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	result.Data = g
	response.Success(ctx, result)
	return
}

// Insert go doc
// @Summary 添加商品
// @Description 添加秒杀商品进秒杀系统
// @Tags 商品管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param goods body model.GoodsDTO true "秒杀商品传输信息"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 401 object model.Result 需要登录
// @Failure 403 object model.Result 没有操作权限
// @Failure 500 object model.Result 操作失败
// @Router /api/goods [post]
func (h *GoodsHandler) Insert(ctx *gin.Context) {
	result := model.Result{}
	// 获取当前用户数据
	var claims *secret.CustomClaims
	c, ok := ctx.Get("claims")
	if !ok {
		result.Code = http.StatusUnauthorized
		result.Message = code.TokenInvalidErr.Error()
		response.Fail(ctx, result)
		return
	}
	claims, ok = c.(*secret.CustomClaims)
	if !ok || claims == nil {
		result.Code = http.StatusUnauthorized
		result.Message = code.TokenInvalidErr.Error()
		response.Fail(ctx, result)
		return
	}
	// 不是卖家不能添加商品
	if claims.Kind != model.NormalSeller {
		result.Code = http.StatusForbidden
		result.Message = code.StatusForbiddenErr.Error()
		response.Fail(ctx, result)
		return
	}
	dto := model.GoodsDTO{}
	// 数据绑定
	if e := ctx.BindJSON(&dto); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	dto.UserId = claims.UserId
	if e := h.goodsService.Insert(dto); e != nil {
		result.Code = http.StatusInternalServerError
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	response.Success(ctx, result)
	return
}

// Update go doc
// @Summary 更新商品
// @Description 更新秒杀商品信息
// @Tags 商品管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param goods body model.GoodsDTO true "秒杀商品传输信息"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 401 object model.Result 需要登录
// @Failure 403 object model.Result 没有操作权限
// @Failure 500 object model.Result 操作失败
// @Router /api/goods [put]
func (h *GoodsHandler) Update(ctx *gin.Context) {
	result := model.Result{}
	dto := model.GoodsDTO{}
	// 数据绑定
	if e := ctx.BindJSON(&dto); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	// 获取当前用户数据
	claims, err := request.GetCurrentCustomClaims(ctx)
	if err != nil {
		result.Code = http.StatusUnauthorized
		result.Message = err.Error()
		response.Fail(ctx, result)
		return
	}

	// 查询商品信息
	goods, _ := h.goodsService.FindGoodsByID(int(dto.ID))
	// 商家只能更新自己的商品信息
	if goods.UserId != claims.UserId {
		result.Code = http.StatusForbidden
		result.Message = code.StatusForbiddenErr.Error()
		response.Fail(ctx, result)
		return
	}
	if e := h.goodsService.Update(dto); e != nil {
		result.Code = http.StatusInternalServerError
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	response.Success(ctx, result)
	return
}

// Delete go doc
// @Summary 删除商品
// @Description 通过 id 删除商品信息
// @Tags 商品管理
// @version 1.0
// @Accept json
// @Produce  json
// @Param id path int true "id" "秒杀商品信息"
// @Success 200 object model.Result 成功后返回值
// @Failure 400 object model.Result 请求参数有误
// @Failure 401 object model.Result 需要登录
// @Failure 403 object model.Result 没有操作权限
// @Failure 500 object model.Result 操作失败
// @Router /api/goods/{id} [DELETE]
func (h *GoodsHandler) Delete(ctx *gin.Context) {
	result := model.Result{}
	idStr := ctx.Param("id")
	id, e := strconv.Atoi(idStr)
	if  e != nil || id == 0 {
		result.Code = http.StatusBadRequest
		result.Message = code.RequestParamErr.Error()
		response.Fail(ctx, result)
		return
	}
	// 获取当前用户数据
	claims, err := request.GetCurrentCustomClaims(ctx)
	if err != nil {
		result.Code = http.StatusUnauthorized
		result.Message = err.Error()
		response.Fail(ctx, result)
		return
	}
	// 查询商品信息
	goods, _ := h.goodsService.FindGoodsByID(id)
	// 商家只能删除自己的商品信息
	if goods.UserId != claims.UserId {
		result.Code = http.StatusForbidden
		result.Message = code.StatusForbiddenErr.Error()
		response.Fail(ctx, result)
		return
	}
	e = h.goodsService.DeleteWithLogic(id)
	if e != nil {
		result.Code = http.StatusBadRequest
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	response.Success(ctx, result)
	return
}