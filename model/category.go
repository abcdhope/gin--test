package model

import (
	"ginblogtest/routes/errmsg"
)

//文章类型
type Category struct {
	ID   uint   `gorm:"primary_key;auto_increment" json:"id"`
	Name string `gorm:"type:varchar(20);not null" json:"name"`
}

//查询分类是否存在
func CheckCategory(name string) int {
	var cate Category
	db.Select("id").Where("name = ?", name).First(&cate)
	if cate.ID > 0 {
		return errmsg.ERROR_CATENAME_USED
	}
	return errmsg.SUCCSE
}

//创建分类
func CreateCategory(cate *Category) int {
	err := db.Create(cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//查询单个分类信息
func GetCateInfo(id int) (Category, int) {
	var cate Category
	db.First(&cate, id)
	return cate, errmsg.SUCCSE
}

//查询分类列表
func GetCate(pageSize, pageNum int) ([]Category, int64) {
	var cates []Category
	var total int64 //类别个数
	db.Select("name").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&cates)
	db.Model(&cates).Count(&total)
	return cates, total
}

//编辑分类名
func EditCate(id int, data *Category) int {
	var cate Category
	maps := make(map[string]interface{})
	maps["name"] = data.Name
	err := db.Model(&cate).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//删除用户对应的类别
func DeleteCate(id int) int {
	var cate Category
	err := db.Where("id = ?", id).Delete(&cate).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
