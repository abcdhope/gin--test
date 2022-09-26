package v1

import (
	"ginblogtest/model"
	"ginblogtest/routes/errmsg"
	"ginblogtest/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//添加用户
func AddUser(c *gin.Context) {
	var user model.User
	//对表单进行校验，并传值到user
	_ = c.ShouldBindJSON(&user)
	//检查字段是否正确
	msg, validCode := validator.Validate(&user)
	if validCode != errmsg.SUCCSE {
		c.JSON(
			http.StatusOK, gin.H{
				"status":  validCode,
				"message": msg,
			},
		)
		c.Abort()
		return
	}
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
func GetUserInfo(c *gin.Context) {
	//获取ID
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetUser(id)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    data,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

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
func EditUser(c *gin.Context) {
	var data model.User
	//对表单进行校验，并传值到user
	_ = c.ShouldBindJSON(&data)
	//获取ID
	id, _ := strconv.Atoi(c.Param("id"))
	//判断需要更改的用户是否一致
	code := model.CheckUpUser(id, data.Username)
	if code == errmsg.SUCCSE {
		model.EditUser(id, &data)
	} else {
		c.Abort() //出现错误拒绝执行后面的函数操作
	}
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

// ChangeUserPassword 修改密码
func ChangeUserPassword(c *gin.Context) {
	var data model.User
	id, _ := strconv.Atoi(c.Param("id"))
	_ = c.ShouldBindJSON(&data)

	code := model.ChangePassword(id, &data)

	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}

//删除用户
func DeleteUser(c *gin.Context) {
	//读取ID
	id, _ := strconv.Atoi(c.Param("id"))

	//删除用户
	code := model.DeleteUser(id)
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"message": errmsg.GetErrMsg(code),
		},
	)
}
