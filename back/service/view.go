package service

import (
	"ViewLog/back/global"
	"ViewLog/back/model"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type viewService struct{}

var ViewService = new(viewService)

// ViewSimple
func (*viewService) ViewSimple() any {
	var (
		sess    = global.Db
		resData = make(gin.H, 0)
	)

	//#region ssh列表
	sshList := make([]*model.Ssh, 0)
	if err := sess.OrderBy("create_time Desc").Find(&sshList); err != nil {
		logrus.Errorf("查询ssh列表失败: %v", err)
		return err
	}
	resData["sshList"] = sshList
	//#endregion

	return resData
}
