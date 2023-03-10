package service

import (
	"ViewLog/back/global"
	"ViewLog/back/model"
	"ViewLog/back/tools/constant"
	toolsSsh "ViewLog/back/tools/ssh"

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

	//#region 活跃sshId
	activeSshInfo := &model.Ssh{}
	if _, err := sess.Where("is_active = ?", constant.SshIsActiveYES).Get(activeSshInfo); err != nil {
		logrus.Errorf("查询活跃ssh失败: %v", err)
		return err
	}
	resData["activeSshId"] = activeSshInfo.Id
	//#endregion

	//#region 更新global.sshClient
	if err := toolsSsh.UpdateGlobalClient(); err != nil {
		return err
	}
	//#endregion

	//#region 获取文件夹
	if activeSshInfo.Id > 0 {
		folderList := make([]*model.Folder, 0)
		if err := sess.Where("ssh_id = ?", activeSshInfo.Id).OrderBy("create_time Desc").Find(&folderList); err != nil {
			logrus.Errorf("查询文件夹列表失败: %v", err)
			return err
		}
		resData["folderList"] = folderList
	}
	//#endregion

	return resData
}

// ViewFolderAdd 添加文件夹
func (*viewService) ViewFolderAdd() any {
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
