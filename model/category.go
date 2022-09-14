package model

import "gorm.io/gorm"

//文章类型
type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}
