package service

import "seckill/model"

type IGoodsService interface {

	// Check 商品信息检查
	Check(g model.Goods) error

	// FindGoodsByID 通过 id 查询一条模型数据
	FindGoodsByID(id int) (model.Goods, error)

	// FindGoodsVOByID 通过 id 查询一条视图数据
	FindGoodsVOByID(id int) (model.GoodsVO, error)

	// ToVO 把 Goods 转为 GoodsVO
	ToVO(g model.Goods) (model.GoodsVO, error)

	// FindByCondition 通过条件查询多条数据
	FindByCondition(c model.GoodsQueryCondition) ([]model.GoodsVO, error)

	// Insert 插入数据
	Insert(dto model.GoodsDTO) error

	// Update 更新数据
	Update(dto model.GoodsDTO) error

	// DeleteWithPhysics 物理删除数据
	DeleteWithPhysics(id int) error

	// DeleteWithLogic 逻辑删除数据
	DeleteWithLogic(id int) error

	// SetGoodsStock 设置商品库存缓存
	SetGoodsStock(goodsId int, stock int) (err error)

	// DecrStock 商品库存缓存原子 -1，并返回减少后的当前库存
	DecrStock(goodsId int) (stock int, err error)

	// IncrStock 商品库存缓存原子 +1，并返回添加后的当前库存
	IncrStock(goodsId int) (err error)

	// InitScekillGoods 初始化参加秒杀的商品
	InitScekillGoods() (e error)
}
