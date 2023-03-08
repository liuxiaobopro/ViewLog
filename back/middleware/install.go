package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Install() gin.HandlerFunc {
	return func(c *gin.Context) {
		//#region 查询项目根目录install.lock文件是否存在
		if _, err := os.Stat("install.lock"); err != nil {
			if os.IsNotExist(err) {
				c.Redirect(http.StatusFound, "/install")
				c.Abort()
				return
			} else {
				c.Redirect(http.StatusFound, "/")
				c.Abort()
				return
			}
		}
		//#endregion
		c.Next()
	}
}

func InstallAlready() gin.HandlerFunc {
	return func(c *gin.Context) {
		//#region 查询项目根目录install.lock文件是否存在
		if _, err := os.Stat("install.lock"); err != nil {
			if os.IsNotExist(err) {
				c.Next()
				return
			}
		}
		c.Redirect(http.StatusFound, "/")
		c.Abort()
		//#endregion
	}
}
