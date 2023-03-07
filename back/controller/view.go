package controller

import (
	"ViewLog/back/common/tools"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type viewHandle struct{}

var ViewHandle = new(viewHandle)

// Install 安装页面
func (*viewHandle) Install(c *gin.Context) {
	c.HTML(http.StatusOK, "install.html", gin.H{})
}

// ViewIndex 首页
func (*viewHandle) ViewIndex(c *gin.Context) {
	resData := make(gin.H, 0)
	resData["time"] = time.Now().Format("2006-01-02 15:04:05")
	resData["title"] = "欢迎使用日志查看器"
	c.HTML(http.StatusOK, "welcome.html", resData)
}

// ViewLog 查看日志
func (*viewHandle) ViewLog(c *gin.Context) {
	resData := make(gin.H, 0)

	//#region 读取文件
	lines, err := tools.ReadFileLine("D:\\1_liuxiaobo\\testlog\\log.txt", 1000, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	//#endregion

	resData["code"] = lines

	c.HTML(http.StatusOK, "log.html", resData)
}

// ViewSimple 单主机
func (*viewHandle) ViewSimple(c *gin.Context) {
	resData := make(gin.H, 0)
	c.HTML(http.StatusOK, "simple.html", resData)
}

// ViewSshAdd 添加ssh
func (*viewHandle) ViewSshAdd(c *gin.Context) {
	resData := make(gin.H, 0)
	c.HTML(http.StatusOK, "add_ssh.html", resData)
}

// ViewFolderAdd 添加文件夹
func (*viewHandle) ViewFolderAdd(c *gin.Context) {
	resData := make(gin.H, 0)
	c.HTML(http.StatusOK, "add_folder.html", resData)
}
