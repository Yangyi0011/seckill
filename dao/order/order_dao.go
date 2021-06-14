package order

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/infra/code"
	"seckill/infra/db"
	"seckill/model"
)

type orderDao struct {
	db *gorm.DB
}

func NewOrderDao() *orderDao {
	return &orderDao{
		db: db.DB,
	}
}

func (d *orderDao) QueryOrderInfoByOrderId(id string) (o model.OrderInfo, e error) {
	if e = d.db.Where("order_id = ?", id).Take(&o).Error; e != nil {
		log.Println(e)
		return
	}
	return
}

func (d *orderDao) QueryByCondition(c model.OrderInfoQueryCondition) (list []model.OrderInfo, e error) {
	sql := c.Model.GetWhereSql()
	if c.OrderId != "" {
		sql.WriteString("and order_id = ")
		sql.WriteString("'")
		sql.WriteString(c.OrderId)
		sql.WriteString("' ")
	}
	if c.UserId != 0 {
		sql.WriteString("and user_id = ")
		sql.WriteString(fmt.Sprintf("%d ", c.UserId))
	}
	if c.GoodsId != 0 {
		sql.WriteString("and goods_id = ")
		sql.WriteString(fmt.Sprintf("%d ", c.GoodsId))
	}
	if c.GoodsName != "" {
		sql.WriteString("and goods_name like ")
		sql.WriteString("'%")
		sql.WriteString(c.GoodsName)
		sql.WriteString("%' ")
	}
	if c.GoodsPrice != 0.0 {
		sql.WriteString("and goods_price = ")
		sql.WriteString(fmt.Sprintf("%f ", c.GoodsPrice))
	}
	if c.PaymentId != 0 {
		sql.WriteString("and payment_id = ")
		sql.WriteString(fmt.Sprintf("%d ", c.PaymentId))
	}
	if c.Status != 0 {
		sql.WriteString("and status = ")
		sql.WriteString(fmt.Sprintf("%d ", c.Status))
	}
	if e = d.db.Limit(c.PageDTO.GetLimit()).Offset(c.PageDTO.GetOffset()).Find(&list, sql.String()).Error; e != nil {
		log.Println(e)
		e = code.DBErr
		return
	}
	return
}

func (d *orderDao) Insert(o model.OrderInfo) (e error) {
	if e = d.db.Create(&o).Error; e != nil {
		return
	}
	return
}

func (d *orderDao) Update(o model.OrderInfo) (e error) {
	if e = d.db.Model(&o).Updates(&o).Error; e != nil {
		return
	}
	return
}

func (d *orderDao) Delete(id string) (e error) {
	o, err := d.QueryOrderInfoByOrderId(id)
	if err != nil {
		return err
	}
	if o.ID == 0 {
		e = code.RecordNotFoundErr
		return
	}
	if e = d.db.Where("order_id = ?", id).Delete(&o).Error; e != nil {
		return
	}
	return
}

// CreateOrder 创建订单
func (d *orderDao) CreateOrder(orderInfo model.OrderInfo) error {
	err := d.db.Transaction(func(tx *gorm.DB) (e error) {
		// 先减少商品库存
		decrStockSql := "update goods set stock = stock - 1 where id = ? and stock > 0"
		tx = tx.Exec(decrStockSql, orderInfo.GoodsId)
		if e = tx.Error; e != nil {
			log.Printf("tx.Exec() failed, err: %v, goodsId: %v", e, orderInfo.GoodsId)
			return
		}
		if tx.RowsAffected == 0 {
			log.Printf("tx.RowsAffected() failed, err: %v, rows: %v", tx.Error, tx.RowsAffected)
			return errors.New("减库存失败")
		}

		// 创建订单
		order := model.Order{
			Model:   orderInfo.Model,
			OrderId: orderInfo.OrderId,
			UserId:  orderInfo.UserId,
			GoodsId: orderInfo.GoodsId,
		}
		tx = tx.Create(&order)
		if e = tx.Error; e != nil {
			log.Printf("tx.Create() failed, err: %v, order: %v", e, order)
			return
		}
		if tx.RowsAffected == 0 {
			log.Printf("tx.Create() failed, err: %v, rows: %v", tx.Error, tx.RowsAffected)
			return errors.New("创建订单失败")
		}

		// 创建订单信息
		tx = tx.Create(&orderInfo)
		if e = tx.Error; e != nil {
			log.Printf("tx.Create() failed, err: %v, order: %v", e, order)
			return
		}
		if tx.RowsAffected == 0 {
			log.Printf("tx.Create() failed, err: %v, rows: %v", tx.Error, tx.RowsAffected)
			return errors.New("创建订单信息失败")
		}
		// 提交事务
		return nil
	})
	return err
}

// CloseOrder 关闭订单
func (d *orderDao) CloseOrder(orderInfo model.OrderInfo) error {
	err := d.db.Transaction(func(tx *gorm.DB) (e error) {
		// 加库存
		incrStockSql := "update goods set stock = stock + 1 where id = ?"
		tx = tx.Exec(incrStockSql, orderInfo.GoodsId)
		if e = tx.Error; e != nil {
			log.Printf("tx.Exec() failed, err: %v, goodsId: %v", e, orderInfo.GoodsId)
			return
		}
		if tx.RowsAffected == 0 {
			log.Printf("tx.RowsAffected() failed, err: %v, rows: %v", tx.Error, tx.RowsAffected)
			e = errors.New("加库存失败")
			return
		}

		// 删除订单
		if e = tx.Delete(model.Order{}, "order_id = ?", orderInfo.OrderId).Error; e != nil {
			log.Printf("tx.Delete() failed, err: %v", e)
			return
		}
		if tx.RowsAffected == 0 {
			log.Printf("rs.RowsAffected() failed, err: %v, rows: %v", tx.Error, tx.RowsAffected)
			e = errors.New("删除订单失败")
			return
		}

		// 修改订单信息
		if e = tx.Model(&orderInfo).Where("order_id = ? and status = ?", orderInfo.OrderId, orderInfo.Status).
			Update("status = ?", model.Unpaid).Error; e != nil {
			log.Printf("tx.Update() failed, err: %v", e)
			return
		}
		if tx.RowsAffected == 0 {
			log.Printf("rs.RowsAffected() failed, err: %v, rows: %v", tx.Error, tx.RowsAffected)
			e = errors.New("修改订单信息失败")
			return
		}
		// 提交事务
		return nil
	})
	return err
}