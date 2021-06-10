package model

type Model struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt LocalTime  `gorm:"type:datetime;comment:'创建时间'"`
	UpdatedAt LocalTime  `gorm:"type:datetime;comment:'更新时间'"`
	DeletedAt *LocalTime `sql:"index" gorm:"type:datetime;comment:'删除时间'"`
}
