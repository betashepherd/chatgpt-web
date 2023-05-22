package routes

import (
	. "chatgpt-web/app/http/controllers"
	"chatgpt-web/app/middlewares"

	"github.com/gin-gonic/gin"
)

var chatController = NewChatController()
var authController = NewAuthController()
var paymentController = NewPaymentController()

// RegisterWebRoutes 注册路由
func RegisterWebRoutes(router *gin.Engine) {
	router.Use(middlewares.Recovery)
	router.Use(middlewares.Cors)
	router.GET("", chatController.Index)
	//router.GET("/payment/debug", paymentController.Debug)
	router.POST("/payment/pay", paymentController.Pay)
	router.POST("/payment/notify", paymentController.Notify)
	router.POST("/user/auth", authController.Auth)
	router.POST("/user/delete", authController.DeleteUser)
	router.POST("/user/register", authController.Register)
	chat := router.Group("/chat").Use(middlewares.Jwt)
	{
		chat.POST("/completion", chatController.Completion)
		chat.POST("/completion/wsstream", chatController.CompletionWsStream)
		chat.POST("/question", chatController.Question)
		chat.GET("/reply", chatController.Reply)
	}
	auth := router.Group("/auth").Use(middlewares.Jwt)
	{
		auth.POST("/info", authController.Info)
	}

}
