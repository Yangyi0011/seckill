package goods

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"seckill/infra/code"
	"seckill/infra/secret"
	"seckill/infra/utils/response"
	"seckill/model"
	"strconv"
)

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
// @Failure 500 object model.Result 添加失败
// @Router /api/goods [post]
func Insert(ctx *gin.Context) {
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

	dto := model.GoodsDTO{}
	// 数据绑定
	if e := ctx.BindJSON(&dto); e != nil {
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	dto.UserId = claims.UserId
	g := model.Goods{}
	if e := g.Insert(dto); e != nil {
		result.Code = http.StatusInternalServerError
		result.Message = e.Error()
		response.Fail(ctx, result)
		return
	}
	response.Success(ctx, result)
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
// @Failure 500 object model.Result 添加失败
// @Router /api/goods/{id} [GET]
func QueryGoodsVOByID(ctx *gin.Context) {
	result := model.Result{}
	idStr := ctx.Param("id")
	id, e := strconv.Atoi(idStr)
	if  e != nil || id == 0 {
		result.Code = http.StatusBadRequest
		result.Message = code.RequestParamErr.Error()
		response.Fail(ctx, result)
		return
	}
	goods := model.Goods{
		Model:       model.Model{
			ID:        uint(id),
		},
	}
	g, e := goods.QueryGoodsVOByID()
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
