package routes

import (
	v1 "ginblogtest/api/v1"
	"ginblogtest/utils"

	"github.com/gin-gonic/gin"
)

//设定路由规则
func InitRouter() {
	gin.SetMode(utils.AppMode) //设置运行模式
	r := gin.Default()         //也可以用gin.New()，区别在于Default自带两个中间件

	//路由组，用于处理逻辑
	routerv1 := r.Group("api/v1")
	{
		routerv1.POST("user/add", v1.AddUser)
		routerv1.GET("users", v1.GetUsers)
	}
	// fmt.Println(utils.HttpPort)
	r.Run(utils.HttpPort)
}
