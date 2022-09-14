package model

import "gorm.io/gorm"

//文章内容
type Article struct {
	Category Category `gorm:"foreignkey:Cid"`
	gorm.Model
	//标题
	Title string `gorm:"type:varchar(100);not null" json:"title"`
	//文章类型的序号
	Cid int `gorm:"type:int;not null" json:"cid"`
	//文章描述
	Desc string `gorm:"type:varchar(200)" json:"desc"`
	//文章内容
	Content string `gorm:"type:longtext" json:"content"`
	//文章图片
	Img string `gorm:"type:varchar(100)" json:"img"`
}
