package goods

import (
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"seckill/infra/code"
	"seckill/infra/db"
	"seckill/model"
	"sync"
)

var (
	once sync.Once
)

type goodsDao struct {
	db *gorm.DB
}

// NewGoodsDao 创建一个 GoodsDao 接口的实例
func NewGoodsDao() *goodsDao {
	return &goodsDao{db: db.DB}
}

// SingleGoodsDao IUserDao 接口单例模式
func SingleGoodsDao() (d *goodsDao) {
	once.Do(func() {
		d = NewGoodsDao()
	})
	return
}

func (d *goodsDao) QueryGoodsByID(id int) (g model.Goods, e error) {
	if e = db.DB.Debug().Where("id = ?", id).Take(&g).Error; e != nil {
		log.Println(e)
		return
	}
	return
}

func (d *goodsDao) QueryByCondition(c model.GoodsQueryCondition) ([]model.Goods, error) {
	//db.DB.Exec("select * from")
	return []model.Goods{}, nil
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
	if e := db.DB.Debug().Updates(&g).Error; e != nil {
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
			return code.RecordNotFound
		}
		log.Println(e)
		return e
	}
	if g.ID == 0 {
		return code.RecordNotFound
	}
	if e := db.DB.Debug().Delete(&g).Error; e != nil {
		log.Println(e)
		return e
	}
	return nil
}