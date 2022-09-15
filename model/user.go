package model

import (
	"ginblogtest/routes/errmsg"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model // 使用 ID 作为主键
	//用户名
	Username string `gorm:"type:varchar(20);not null" json:"username"`
	//密码
	Password string `gorm:"type:varchar(20);not null" json:"password"`
	//角色码，区分管理员和普通用户
	Role int `gorm:"type:int;not null" json:"role"`
}

//检查用户是否存在，返回状态码
func CheckUser(name string) int {
	var user User
	//从名字找到对应的ID
	db.Select("id").Where("Username=?", name).First(&user)
	//如果存在该用户则表示已经被抢注了
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED
	}
	return errmsg.SUCCSE
}

//创建用户，成功就返回状态码
func CreateUser(user *User) int {
	err := db.Create(user).Error //返回错误
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//查询单个用户
func GetUser(id int) (User, int) {
	var user User
	//err := db.Limit(1).Where("ID = ?", id).Find(&user).Error
	err := db.First(&user, id).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

//查询用户列表
//根据一部分名字来找到一系列相似的用户
func GetUsers(username string, pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64 //用户个数
	if username != "" {
		//查找当前页的所有用户
		db.Select("id", "username", "role", "created_at").Where("username LIKE ?", username+"%").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
		//统计用户个数
		db.Model(&users).Where("username LIKE ?", username+"%").Count(&total)
		return users, total
	}
	//查找当前页面的用户数量
	db.Select("id", "username", "role", "created_at").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users)
	// fmt.Println(users)
	db.Model(&users).Count(&total) //所有用户的数量
	return users, total
}

//编辑用户,指定目标用户更改信息
func EditUser(id int, data *User) int {
	var user User
	// 存放需要更改的信息
	maps := make(map[string]interface{})
	maps["username"] = data.Username
	maps["role"] = data.Role
	err := db.Model(&user).Where("id = ?", id).Updates(maps).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//删除用户
func DeleteUser(id int) int {
	var user User
	err := db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}
