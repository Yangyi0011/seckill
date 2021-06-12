package goods

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/infra/code"
	"seckill/infra/db"
	"seckill/model"
)

type goodsDao struct {
	db *gorm.DB
}

// NewGoodsDao 创建一个 GoodsDao 接口的实例
func NewGoodsDao() *goodsDao {
	return &goodsDao{db: db.DB}
}

func (d *goodsDao) QueryGoodsByID(id int) (g model.Goods, e error) {
	if e = db.DB.Debug().Where("id = ?", id).Take(&g).Error; e != nil {
		log.Println(e)
		return
	}
	return
}

func (d *goodsDao) QueryByCondition(c model.GoodsQueryCondition) (list []model.Goods, e error) {
	timeZero := model.LocalTime{}.ZeroValue()
	var sql bytes.Buffer
	sql.WriteString("1 = 1 ")
	if c.ID != 0 {
		sql.WriteString("and id = ")
		sql.WriteString(fmt.Sprintf("%d ", c.ID))
	}
	if c.CreatedAt != timeZero {
		sql.WriteString("and create_at = ")
		sql.WriteString("'")
		sql.WriteString(c.CreatedAt.String())
		sql.WriteString("' ")
	}
	if c.UpdatedAt != timeZero {
		sql.WriteString("and update_at = ")
		sql.WriteString("'")
		sql.WriteString(c.UpdatedAt.String())
		sql.WriteString("' ")
	}
	if c.DeletedAt != nil && (*c.DeletedAt) != timeZero {
		sql.WriteString("and delete_at = ")
		sql.WriteString("'")
		sql.WriteString(c.DeletedAt.String())
		sql.WriteString("' ")
	}
	if c.Name != "" {
		sql.WriteString("and name like ")
		sql.WriteString("'%")
		sql.WriteString(c.Name)
		sql.WriteString("%' ")
	}
	if c.Price != 0.0 {
		sql.WriteString("and price = ")
		sql.WriteString(fmt.Sprintf("%f ", c.Price))
	}
	if c.Stock != 0 {
		sql.WriteString("and price = ")
		sql.WriteString(fmt.Sprintf("%d ", c.Stock))
	}
	if c.StartTime != timeZero {
		sql.WriteString("and start_time = ")
		sql.WriteString("'")
		sql.WriteString(c.StartTime.String())
		sql.WriteString("' ")
	}
	if c.EndTime != timeZero {
		sql.WriteString("and end_time = ")
		sql.WriteString("'")
		sql.WriteString(c.EndTime.String())
		sql.WriteString("' ")
	}
	offset := 0
	if c.Index != 0 {
		offset = (c.Index-1)*c.Size
	}
	limit := 0
	if c.Size != 0 {
		limit = c.Size
	}
	if e = d.db.Debug().Limit(limit).Offset(offset).Find(&list, sql.String()).Error; e != nil {
		log.Println(e)
		e = code.DBErr
		return
	}
	return
}

func (d *goodsDao) Insert(g model.Goods) error {
	if e := db.DB.Debug().Create(&g).Error; e != nil {
		log.Println(e)
		return e
	}
	return nil
}

// Update 更新数据
func (d *goodsDao) Update(g model.Goods) error {
	if e := db.DB.Debug().Model(&g).Updates(&g).Error; e != nil {
		log.Println(e)
		return e
	}
	return nil
}

// Delete 物理删除数据
func (d *goodsDao) Delete(id int) error {
	var g model.Goods
	g.ID = uint(id)
	if e := db.DB.Debug().Take(&g).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return code.RecordNotFoundErr
		}
		log.Println(e)
		return e
	}
	if g.ID == 0 {
		return code.RecordNotFoundErr
	}
	if e := db.DB.Debug().Delete(&g).Error; e != nil {
		log.Println(e)
		return e
	}
	return nil
}