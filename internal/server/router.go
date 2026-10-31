package server

import (
	"mini-chat/internal/chat"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 接收聊天 hub，创建路由
func NewRouter(hub *chat.Hub) *gin.Engine {
	// 创建带日志和异常恢复中间件的路由
	router := gin.Default()

	// HTTP API 路由
	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// ws 接口
	router.GET("/ws", func(c *gin.Context) {
		chat.ServeWS(hub, c.Writer, c.Request)
	})

	// Vite 编译后的页面
	router.StaticFile("/", "./web/dist/index.html")
	router.Static("/assets", "./web/dist/assets")

	return router
}
