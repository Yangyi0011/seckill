package model

import (
	"bytes"
	"fmt"
)

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt LocalTime  `json:"createdAt" gorm:"type:datetime;comment:'创建时间'" swaggerignore:"true"`
	UpdatedAt LocalTime  `json:"updatedAt" gorm:"type:datetime;comment:'更新时间'" swaggerignore:"true"`
	DeletedAt *LocalTime `json:"deletedAt" sql:"index" gorm:"type:datetime;comment:'删除时间'" swaggerignore:"true"`
}

// GetWhereSql 获取动态 sql 查询的 where 前缀
func (m Model) GetWhereSql() bytes.Buffer {
	var sql bytes.Buffer
	timeZero := LocalTime{}.ZeroValue()
	sql.WriteString("1 = 1 ")
	if m.ID != 0 {
		sql.WriteString("and id = ")
		sql.WriteString(fmt.Sprintf("%d ", m.ID))
	}
	if m.CreatedAt != timeZero {
		sql.WriteString("and create_at = ")
		sql.WriteString("'")
		sql.WriteString(m.CreatedAt.String())
		sql.WriteString("' ")
	}
	if m.UpdatedAt != timeZero {
		sql.WriteString("and update_at = ")
		sql.WriteString("'")
		sql.WriteString(m.UpdatedAt.String())
		sql.WriteString("' ")
	}
	if m.DeletedAt != nil && (*m.DeletedAt) != timeZero {
		sql.WriteString("and delete_at = ")
		sql.WriteString("'")
		sql.WriteString(m.DeletedAt.String())
		sql.WriteString("' ")
	}
	return sql
}