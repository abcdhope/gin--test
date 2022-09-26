package v1

import (
	"ginblogtest/model"
	"ginblogtest/routes/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//添加分类
func AddCategory(c *gin.Context) {
	var cate model.Category
	_ = c.ShouldBindJSON(&cate)
	//检查添加的分类名是否存在
	code := model.CheckCategory(cate.Name)
	if code == errmsg.SUCCSE {
		//创建分类
		model.CreateCategory(&cate)
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

//查询分类信息
func GetCateInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	//获取分类信息
	data, code := model.GetCateInfo(id)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

//查询分类列表
func GetCate(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))

	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}

	//获取列表
	data, total := model.GetCate(pageSize, pageNum)
	code := errmsg.SUCCSE
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"total":   total,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

//编辑分类名
func EditCate(c *gin.Context) {
	var data model.Category
	//获取ID、数据
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)
	//判断更改的名字是否存在
	code := model.CheckCategory(data.Name)
	if code == errmsg.SUCCSE {
		//更改信息
		model.EditCate(id, &data)
	} else {
		c.Abort() //出现错误，停止后续函数操作
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

//删除类别
func DeleteCate(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteCate(id)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
