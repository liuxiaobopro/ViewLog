package router

import (
	"ViewLog/back/controller"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//#region 页面
	rg1 := r.Group("")
	{
		rg1.GET("/", controller.ViewIndexIndex)
		rg1.GET("/log", controller.ViewLogIndex)
	}
	//#endregion

	//#region api
	rg2 := r.Group("/api")
	{
		rg2.GET("/show_fold", controller.ShowFolds)
		rg2.GET("/open_fold", controller.OpenFold)
		rg2.GET("/read_file", controller.ReadFile)
	}
	//#endregion
}
