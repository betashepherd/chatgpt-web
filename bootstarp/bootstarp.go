package bootstrap

import (
	"chatgpt-web/config"
	"chatgpt-web/library/lfs"
	"chatgpt-web/pkg/logger"
	"net/http"
	"strconv"
)

func StartWebServer() {
	// 注册启动所需各类参数
	SetUpRoute()
	SetupDB()
	initTemplateDir()
	initStaticServer()

	lfs.Init("./data", "/data")

	// 启动服务
	port := config.LoadConfig().Port
	portString := strconv.Itoa(port)
	// 自定义监听地址
	listen := config.LoadConfig().Listen
	err := router.Run(listen + ":" + portString)
	if err != nil {
		logger.Danger("run webserver error %s", err)
		return
	}
}

// initTemplate 初始化HTML模板加载路径
func initTemplateDir() {
	router.LoadHTMLGlob("dist/*.html")
}

// initStaticServer 初始化静态文件处理
func initStaticServer() {
	router.StaticFS("/assets", http.Dir("dist/assets"))
	router.StaticFile("avatar.png", "dist/avatar.png")
	router.StaticFile("gitclone.png", "dist/gitclone.png")
	router.StaticFile("user.png", "dist/user.png")
	router.StaticFile("logo192.png", "dist/logo192.png")
	router.StaticFile("logo512.png", "dist/logo512.png")
	router.StaticFile("favicon.ico", "dist/favicon.ico")
	router.StaticFile("manifest.json", "dist/manifest.json")
}
