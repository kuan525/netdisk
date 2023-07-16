package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/kuan525/netdisk/common"
	"github.com/kuan525/netdisk/util"
	"net/http"
)

// IsTokenValid token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}

	// TODO 判断token的时效性，是否过期
	// TODO 从数据库表tbl_user_token查询username对应的token信息
	// TODO 对比两个token是否一致
	return true
}

// Authorize http请求拦截器
func Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Request.FormValue("username")
		token := c.Request.FormValue("token")

		// 验证登陆token是否有效
		if len(username) < 3 || !IsTokenValid(token) {
			// token 校验失败则跳转到登陆页面
			c.Abort()
			resp := util.NewRespMsg(
				int(common.StatusTokenInvalid),
				"token无效",
				nil)
			c.JSON(http.StatusOK, resp)
			return
		}
		c.Next()
	}
}
