package routers

import (
	"net/http"
	"web_app/logger"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	// 新建一个没有任何默认中间件的路由
	r := gin.New()

	// 注册2个全局中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 返回r
	return r
}
