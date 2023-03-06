package router

import (
	"ViewLog/back/controller"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//#region 页面
	rg1 := r.Group("")
	{
		rg1.GET("/", controller.ViewHandle.ViewIndex)               // 首页
		rg1.GET("/log", controller.ViewHandle.ViewLog)              // 查看日志
		rg1.GET("/simple", controller.ViewHandle.ViewSimple)        // 单主机
		rg1.GET("/ssh_add", controller.ViewHandle.ViewSshAdd)       // 添加ssh
		rg1.GET("/folder_add", controller.ViewHandle.ViewFolderAdd) // 添加文件夹
	}
	//#endregion

	//#region api
	rg2 := r.Group("/api")
	{
		rg2.GET("/show_fold", controller.ApiHandle.ShowFolds)
		rg2.GET("/open_fold", controller.ApiHandle.OpenFold)
		rg2.GET("/read_file", controller.ApiHandle.ReadFile)
		rg2.POST("/add_ssh", controller.ApiHandle.AddSsh)
	}
	//#endregion
}
