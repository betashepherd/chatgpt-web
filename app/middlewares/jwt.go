package middlewares

import (
	"chatgpt-web/app/http/controllers"
	"chatgpt-web/pkg/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

var base = controllers.BaseController{}

// Jwt jwt认证
func Jwt(c *gin.Context) {
	claims, err := auth.EncodeByCtx(c)
	if err != nil {
		base.ResponseJson(c, http.StatusUnauthorized, err.Error(), nil)
		return
	}

	if claims.User.ID == 0 {
		base.ResponseJson(c, http.StatusUnauthorized, "用户信息错误，未知的token", nil)
		return
	}
	c.Set("authUser", claims.User)
	c.Next()
}
