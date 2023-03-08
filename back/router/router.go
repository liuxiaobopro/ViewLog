package router

import (
	"ViewLog/back/controller"
	"ViewLog/back/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	rg := r.Group("")
	rg.Use(middleware.InstallAlready())
	{
		rg.GET("/install", controller.ViewHandle.Install) // 安装页面
		rg.POST("/install", controller.ApiHandle.Install) // 安装
	}

	rg0 := r.Group("")
	{
		rg0.POST("/reset", controller.ApiHandle.Reset) // 重置
	}

	//#region 页面
	rg1 := r.Group("")
	rg1.Use(middleware.Install())
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
	rg2.Use(middleware.Install())
	{
		rg2.GET("/show_fold", controller.ApiHandle.ShowFolds)
		rg2.GET("/open_fold", controller.ApiHandle.OpenFold)
		rg2.GET("/read_file", controller.ApiHandle.ReadFile)

		//#region ssh
		rg2.POST("/ssh", controller.ApiHandle.AddSsh)
		rg2.DELETE("/ssh", controller.ApiHandle.DelSsh)
		rg2.PUT("/ssh", controller.ApiHandle.UpdateSsh)
		rg2.GET("/ssh/:id", controller.ApiHandle.DetailSsh)
		rg2.GET("/ssh", controller.ApiHandle.ListSsh)
		rg2.PUT("/ssh/active", controller.ApiHandle.UpdateActiveSsh)
		//#endregion
	}
	//#endregion
}
