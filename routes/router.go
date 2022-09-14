package routes

import (
	"ginblogtest/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//设定路由规则
func InitRouter() {
	gin.SetMode(utils.AppMode) //设置运行模式
	r := gin.Default()         //也可以用gin.New()，区别在于Default自带两个中间件

	//路由组，用于处理逻辑
	router := r.Group("api/v1")
	{
		router.GET("hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "ok",
			})
		})
	}
	// fmt.Println(utils.HttpPort)
	r.Run(utils.HttpPort)
}
