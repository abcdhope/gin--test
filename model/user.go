package model

import "gorm.io/gorm"

type User struct {
	gorm.Model // 使用 ID 作为主键
	//用户名
	Username string `gorm:"type:varchar(20);not null" json:"username"`
	//密码
	Password string `gorm:"type:varchar(20);not null" json:"password"`
	//角色码，区分管理员和普通用户
	Role int `gorm:"type:int;not null" json:"role"`
}
