package dao

import "seckill/model"

type IGoodsDao interface {
	QueryGoodsByID(id int) (g model.Goods, e error)
	QueryByCondition(c model.GoodsQueryCondition) ([]model.Goods, error)
	Insert(g model.Goods) error
	Update(g model.Goods) error
	Delete(id int) error
}
