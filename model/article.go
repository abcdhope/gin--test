package model

import (
	"ginblogtest/routes/errmsg"

	"gorm.io/gorm"
)

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
	//文章访问量
	ReadCount int `gorm:"type:int;not null;default:0" json:"read_count"`
	//评论数
	CommentCount int `gorm:"type:int;not null;default:0" json:"comment_count"`
}

// CreateArt 新增文章
func CreateArt(data *Article) int {
	err := db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCSE
}

//查询该分类下的所有文章
func GetCateArt(id int, pagesize int, pagenum int) ([]Article, int64, int) {
	var arts []Article
	var total int64
	//Preload预加载，从Category中再筛选
	err := db.Preload("Category").Limit(pagesize).Offset((pagenum-1)*pagesize).Where("cid = ?", id).Find(&arts).Error
	db.Model(&arts).Where("cid = ?", id).Count(&total)
	if err != nil {
		return nil, 0, errmsg.ERROR_CATE_NOT_EXIST
	}
	return arts, total, errmsg.SUCCSE
}

//查询单个文章信息
func GetArtInfo(id int) (Article, int) {
	var art Article
	//查询目标文章
	err := db.Where("id = ?", id).First(&art).Error
	//更新访问次数
	db.Model(&art).Where("id = ?", id).UpdateColumn("read_count", gorm.Expr("read_count + ?", 1))
	if err != nil {
		return Article{}, errmsg.ERROR_ART_NOT_EXIST
	}
	return art, errmsg.SUCCSE
}

//查询文章列表
func GetArt(pageSize, pageNum int) ([]Article, int64, int) {
	var arts []Article
	var total int64
	err := db.Select("article.id, title, img, created_at, updated_at, `desc`, comment_count, read_count, category.name").Limit(pageSize).Offset((pageNum - 1) * pageSize).Joins("Category").Find(&arts).Error
	db.Model(&arts).Count(&total)
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return arts, total, errmsg.SUCCSE
}

//查询特定名字下的文章列表
func SearchArticle(title string, pageSize int, pageNum int) ([]Article, int64, int) {
	var arts []Article
	var total int64
	err := db.Select("article.id,title, img, created_at, updated_at, `desc`, comment_count, read_count, Category.name").Order("created_at DESC").Joins("Category").Where("title LIKE ?", title+"%").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&arts).Error
	db.Model(&arts).Where("title LIKE ?", title+"%").Count(&total)
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return arts, total, errmsg.SUCCSE
}

//编辑文章
func EditArt(id int, data *Article) int {
	var art Article
	var maps = make(map[string]interface{})
	maps["title"] = data.Title
	maps["cid"] = data.Cid
	maps["desc"] = data.Desc
	maps["content"] = data.Content
	maps["img"] = data.Img

	err = db.Model(&art).Where("id = ? ", id).Updates(&maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// DeleteArt 删除文章
func DeleteArt(id int) int {
	var art Article
	err = db.Where("id = ? ", id).Delete(&art).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
