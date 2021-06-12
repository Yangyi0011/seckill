package model

type Model struct {
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt LocalTime  `json:"createdAt" gorm:"type:datetime;comment:'创建时间'" swaggerignore:"true"`
	UpdatedAt LocalTime  `json:"updatedAt" gorm:"type:datetime;comment:'更新时间'" swaggerignore:"true"`
	DeletedAt *LocalTime `json:"deletedAt" sql:"index" gorm:"type:datetime;comment:'删除时间'" swaggerignore:"true"`
}
