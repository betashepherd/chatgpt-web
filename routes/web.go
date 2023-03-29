package routes

import (
	. "chatgpt-web/app/http/controllers"
	"chatgpt-web/app/middlewares"
	"github.com/gin-gonic/gin"
)

var chatController = NewChatController()
var authController = NewAuthController()

// RegisterWebRoutes 注册路由
func RegisterWebRoutes(router *gin.Engine) {

	router.Use(middlewares.Cors)
	router.GET("", chatController.Index)
	router.POST("/user/auth", authController.Auth)
	router.POST("/user/delete", authController.DeleteUser)
	chat := router.Group("/chat").Use(middlewares.Jwt)
	{
		chat.POST("/completion", chatController.Completion)
		chat.POST("/completion/stream", chatController.CompletionStream)
	}
	auth := router.Group("/auth").Use(middlewares.Jwt)
	{
		auth.POST("/info", authController.Info)
	}

}
