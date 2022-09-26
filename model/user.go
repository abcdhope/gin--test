package model

import (
	"ginblogtest/routes/errmsg"
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model // 使用 ID 作为主键
	//用户名
	Username string `gorm:"type:varchar(20);not null" json:"username" validate:"required,min=4,max=12" label:"用户名"`
	//密码
	Password string `gorm:"type:varchar(200);not null" json:"password" validate:"required,min=6,max=120" label:"密码"`
	//角色码，区分管理员和普通用户，1是管理员，2以上是普通用户
	Role int `gorm:"type:int;default:2" json:"role" validate:"required,gte=2" label:"角色码"`
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

//判断需要更改信息的用户与数据库的用户是否一致
func CheckUpUser(id int, name string) int {
	var user User
	db.Select("id").Where("username=?", name).First(&user)
	//如果更新除名字以外的内容，则判断名字所属的id是否相同
	if user.ID == uint(id) {
		return errmsg.SUCCSE
	}
	if user.ID > 0 {
		return errmsg.ERROR_USERNAME_USED //1001
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
		return User{}, errmsg.ERROR
	}
	return user, errmsg.SUCCSE
}

//查询用户列表
//根据一部分名字来找到一系列相似的用户
func GetUsers(username string, pageSize int, pageNum int) ([]User, int64) {
	var users []User
	var total int64 //用户个数
	if username != "" {
		//查找当前页的所有用户以及统计用户个数
		db.Select("id", "username", "role", "created_at").Where("username LIKE ?", username+"%").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total)
		//统计用户个数
		// db.Model(&users).Where("username LIKE ?", username+"%").Count(&total)
		return users, total
	}
	//查找当前页面的用户数量
	db.Select("id", "username", "role", "created_at").Limit(pageSize).Offset((pageNum - 1) * pageSize).Find(&users).Count(&total)
	// fmt.Println(users)
	// db.Model(&users).Count(&total) //所有用户的数量
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

// ChangePassword 修改密码
func ChangePassword(id int, data *User) int {

	err = db.Select("password").Where("id = ?", id).Updates(&data).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//删除用户
func DeleteUser(id int) int {
	var user User
	//判断数据库是否存在该id
	if db.Where("id = ?", id).First(&user).RowsAffected < 1 {
		return errmsg.ERROR_USER_NOT_EXIST
	}
	err := db.Where("id = ?", id).Delete(&user).Error
	if err != nil {
		return errmsg.ERROR
	}
	return errmsg.SUCCSE
}

//密码加密策略
func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	u.Role = 2
	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.Password = ScryptPw(u.Password)
	return nil
}

//加密
func ScryptPw(password string) string {
	const cost = 10
	HashPw, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		log.Fatal(err)
	}
	return string(HashPw)
}

// CheckLogin 后台登录验证
func CheckLogin(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	//用户不存在
	if user.ID == 0 {
		return User{}, errmsg.ERROR_USER_NOT_EXIST
	}
	//密码错误
	if PasswordErr != nil {
		return User{}, errmsg.ERROR_PASSWORD_WRONG
	}
	//管理员权限
	if user.Role != 1 {
		return User{}, errmsg.ERROR_USER_NO_RIGHT
	}
	return user, errmsg.SUCCSE
}

// CheckLoginFront 前台登录
func CheckLoginFront(username string, password string) (User, int) {
	var user User
	var PasswordErr error

	db.Where("username = ?", username).First(&user)

	PasswordErr = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if user.ID == 0 {
		return user, errmsg.ERROR_USER_NOT_EXIST
	}
	if PasswordErr != nil {
		return user, errmsg.ERROR_PASSWORD_WRONG
	}
	return user, errmsg.SUCCSE
}
