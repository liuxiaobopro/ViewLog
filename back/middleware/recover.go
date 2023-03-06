package middleware

import (
	"net/http"

	"ViewLog/back/common/resp"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Recovery middleware
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Infof("middleware recover err: %v", err)
				c.JSON(http.StatusInternalServerError, resp.FailResp(resp.Fail, "服务器错误"))
				return
			}
		}()
		c.Next()
	}
}
