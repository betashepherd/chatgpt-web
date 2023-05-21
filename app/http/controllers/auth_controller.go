package controllers

import (
	"chatgpt-web/library/util"
	"chatgpt-web/pkg/auth"
	"chatgpt-web/pkg/logger"
	"chatgpt-web/pkg/model/user"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// AuthController 认证控制器
type AuthController struct {
	BaseController
}

func NewAuthController() *AuthController {
	return &AuthController{}
}

// authRequest 认证请求
type authRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
	Realname string `json:"realname"`
}

// Auth 认证
func (c *AuthController) Auth(ctx *gin.Context) {
	var req authRequest
	err := ctx.BindJSON(&req)
	if err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if req.Name == "" || req.Password == "" {
		c.ResponseJson(ctx, http.StatusUnauthorized, "请输入用户名密码", nil)
		return
	}

	authUser, err := user.GetByName(req.Name)
	if err != nil && err == gorm.ErrRecordNotFound {
		c.ResponseJson(ctx, http.StatusUnauthorized, "请求认证的用户不存在", nil)
		return
	}
	if !authUser.ComparePassword(req.Password) {
		c.ResponseJson(ctx, http.StatusUnauthorized, "密码错误", nil)
		return
	}
	if authUser.Stat != 0 {
		c.ResponseJson(ctx, http.StatusUnauthorized, "账号暂时无法使用，请联系管理员! QQ Mail: 1299587848#qq.com", nil)
		return
	}
	if authUser.ExpireTimestamp > 0 && authUser.ExpireTimestamp < util.GetCurrentTime().Unix() {
		c.ResponseJson(ctx, http.StatusUnauthorized, "VIP会员已到期, 请续费!", nil)
		return
	}
	token, err := auth.Encode(authUser)
	if err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.ResponseJson(ctx, http.StatusOK, "", gin.H{
		"token": token,
	})
}

// Info 登录用户信息
func (c *AuthController) Info(ctx *gin.Context) {
	authUser, ok := ctx.Get("authUser")
	if !ok {
		c.ResponseJson(ctx, http.StatusInternalServerError, "获取登录用户信息失败", nil)
		return
	}

	_, ok = authUser.(*user.User)
	if !ok {
		c.ResponseJson(ctx, http.StatusInternalServerError, "断言登录用户信息失败", nil)
		return
	}
	// 未实现权限系统，写死
	c.ResponseJson(ctx, http.StatusOK, "", gin.H{
		//"info":             userInfo,
		"permissionRoutes": []string{"chat", "chat/completion", "user/auth/info", "user/auth"},
	})
}

func (c *AuthController) DeleteUser(ctx *gin.Context) {
	var req authRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if req.Name == "" || req.Password == "" {
		c.ResponseJson(ctx, http.StatusUnauthorized, "请输入用户名", nil)
		return
	}

	if req.Password != "#EScaz#^W5JQ8j" {
		c.ResponseJson(ctx, http.StatusUnauthorized, "授权码错误", nil)
		return
	}

	user.DeleteUser(req.Name)

	c.ResponseJson(ctx, http.StatusOK, "", nil)
}

func (c *AuthController) Register(ctx *gin.Context) {
	var req authRequest

	if err := ctx.BindJSON(&req); err != nil {
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	if req.Name == "" || req.Password == "" || req.Realname == "" {
		c.ResponseJson(ctx, http.StatusUnauthorized, "请输入用户名密码", nil)
		return
	}

	if _, err := user.CreateUser(req.Name, req.Password, req.Realname, util.GetCurrentTime().Unix()+24*3600); err != nil {
		logger.Info("create user error:", err)
		c.ResponseJson(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	c.ResponseJson(ctx, http.StatusOK, "", nil)
}
