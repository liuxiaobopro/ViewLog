package router

import (
	"ViewLog/back/controller"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//#region 页面
	rg1 := r.Group("")
	{
		rg1.GET("/", controller.ViewHandle.ViewIndex)
		rg1.GET("/log", controller.ViewHandle.ViewLog)
		rg1.GET("/simple", controller.ViewHandle.ViewSimple)
	}
	//#endregion

	//#region api
	rg2 := r.Group("/api")
	{
		rg2.GET("/show_fold", controller.ApiHandle.ShowFolds)
		rg2.GET("/open_fold", controller.ApiHandle.OpenFold)
		rg2.GET("/read_file", controller.ApiHandle.ReadFile)
	}
	//#endregion
}
