package model

import (
	"ginblogtest/routes/errmsg"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	UserId    uint   `json:"user_id"`                                    //用户ID
	ArticleId uint   `json:"article_id"`                                 //文章ID
	Title     string `json:"article_title"`                              //文章名
	Username  string `json:"username"`                                   //用户名
	Content   string `gorm:"type:varchar(500);not null;" json:"content"` //内容
	Status    int8   `gorm:"type:tinyint;default:2" json:"status"`       //状态
}

// AddComment 新增评论
func AddComment(data *Comment) int {
	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// GetComment 查询单个评论
func GetComment(id int) (Comment, int) {
	var comment Comment
	err = db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		return Comment{}, errmsg.ERROR
	}
	return comment, errmsg.SUCCSE
}

// GetCommentCount 获取评论数量
func GetCommentCount(id int) int64 {
	var comment Comment
	var total int64
	db.Model(&comment).Where("article_id = ?", id).Where("status = ?", 1).Count(&total)
	return total
}

// GetCommentList 后台所有获取评论列表
func GetCommentList(pageSize int, pageNum int) ([]Comment, int64, int) {

	var commentList []Comment
	var total int64
	db.Find(&commentList).Count(&total) //所有评论个数
	err = db.Model(&commentList).Limit(pageSize).Offset((pageNum - 1) * pageSize).Order("Created_At DESC").Select("comment.id, article.title,user_id,article_id, user.username, comment.content, comment.status,comment.created_at,comment.deleted_at").Joins("LEFT JOIN article ON comment.article_id = article.id").Joins("LEFT JOIN user ON comment.user_id = user.id").Scan(&commentList).Error
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return commentList, total, errmsg.SUCCSE
}

// GetCommentListFront 展示页面获取评论列表
func GetCommentListFront(id int, pageSize int, pageNum int) ([]Comment, int64, int) {
	var commentList []Comment
	var total int64
	db.Find(&Comment{}).Where("article_id = ?", id).Where("status = ?", 1).Count(&total)
	err = db.Model(&Comment{}).Limit(pageSize).Offset((pageNum-1)*pageSize).Order("Created_At DESC").Select("comment.id, article.title, user_id, article_id, user.username, comment.content, comment.status,comment.created_at,comment.deleted_at").Joins("LEFT JOIN article ON comment.article_id = article.id").Joins("LEFT JOIN user ON comment.user_id = user.id").Where("article_id = ?",
		id).Where("status = ?", 1).Scan(&commentList).Error
	if err != nil {
		return nil, 0, errmsg.ERROR
	}
	return commentList, total, errmsg.SUCCSE
}

// DeleteComment 删除评论
func DeleteComment(id uint) int {
	var comment Comment
	err = db.Where("id = ?", id).Delete(&comment).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// CheckComment 通过评论
func CheckComment(id int, data *Comment) int {
	var comment Comment
	var res Comment
	var article Article
	var maps = make(map[string]interface{})
	maps["status"] = data.Status

	err = db.Model(&comment).Where("id = ?", id).Updates(maps).First(&res).Error
	db.Model(&article).Where("id = ?", res.ArticleId).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

// UncheckComment 撤下评论
func UncheckComment(id int, data *Comment) int {
	var comment Comment
	var res Comment
	var article Article
	var maps = make(map[string]interface{})
	maps["status"] = data.Status

	err = db.Model(&comment).Where("id = ?", id).Updates(maps).First(&res).Error
	db.Model(&article).Where("id = ?", res.ArticleId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
