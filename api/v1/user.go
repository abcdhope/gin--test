package v1

import (
	"ginblogtest/model"
	"ginblogtest/routes/errmsg"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//添加用户
func AddUser(c *gin.Context) {
	var user model.User
	//对表单进行校验，并传值到user
	_ = c.ShouldBindJSON(&user)
	//判断是否有该用户
	code := model.CheckUser(user.Username)
	if code == errmsg.SUCCSE {
		//创建记录
		model.CreateUser(&user)
	}
	//将写入情况返回
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

//查询单个用户
//查询用户列表
func GetUsers(c *gin.Context) {
	//获取分页信息
	//c.Query返回地址上的键值
	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	pageNum, _ := strconv.Atoi(c.Query("pagenum"))
	username := c.Query("username")
	switch {
	case pageSize >= 100:
		pageSize = 100
	case pageSize <= 0:
		pageSize = 10
	}

	if pageNum == 0 {
		pageNum = 1
	}
	//获取用户列表
	data, total := model.GetUsers(username, pageSize, pageNum)
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

//编辑用户
//删除用户
