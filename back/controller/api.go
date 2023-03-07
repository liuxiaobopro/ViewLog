package controller

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"ViewLog/back/common/resp"
	"ViewLog/back/common/tools"
	modelReq "ViewLog/back/model/req"
	modelRes "ViewLog/back/model/res"
	"ViewLog/back/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type apiHandle struct{}

var ApiHandle = new(apiHandle)

// Install 安装
func (*apiHandle) Install(c *gin.Context) {
	req := new(modelReq.InstallReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Errorf("install req error: %v", err)
		c.JSON(http.StatusOK, resp.Param)
		return
	}

	data, err := service.ApiService.Install(req)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailResp(resp.FailCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.SuccResp(data))
}

// Reset 重置
func (*apiHandle) Reset(c *gin.Context) {
	if err := service.ApiService.Reset(); err != nil {
		c.JSON(http.StatusOK, resp.FailResp(resp.FailCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.SuccResp(nil))
}

// OpenFold 打开文件夹
func (*apiHandle) OpenFold(c *gin.Context) {
	//#region 获取参数
	pathParam := c.Query("filepath")
	if pathParam == "" {
		c.JSON(http.StatusBadRequest, "path is empty")
		return
	}
	//#endregion

	//#region 读取目录中的所有文件和子目录
	entries, err := os.ReadDir(pathParam)
	if err != nil {
		logrus.Errorf("read dir error: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	var (
		menuList      = make([]*modelRes.LogIndexShowFoldsRes, 0)
		menuChildList = make([]*modelRes.LogIndexShowFoldsRes, 0)
	)
	//#endregion

	//#region 整理数据
	for k, entry := range entries {
		var (
			title string
			id    = strconv.Itoa(k)
		)

		if entry.IsDir() {
			// title = "<span class=\"layui-badge layui-bg-orange\">F</span>" + entry.Name()
			// menuChildList = append(menuChildList, &modelRes.LogIndexShowFoldsRes{
			// 	Title:  title,
			// 	Id:     id,
			// 	PathId: fmt.Sprintf("%s-%s", pathParam, id),
			// })
		} else {
			title = "<span class=\"layui-badge layui-bg-green\">D</span> " + entry.Name()
			menuChildList = append(menuChildList, &modelRes.LogIndexShowFoldsRes{
				Title: title,
				Id:    id,
				Path:  pathParam,
				Name:  entry.Name(),
			})
		}
	}

	menuList = append(menuList, &modelRes.LogIndexShowFoldsRes{
		Title:    pathParam,
		Spread:   true,
		Children: menuChildList,
	})
	//#endregion

	c.JSON(http.StatusOK, menuList)
}

// ShowFolds 展示文件夹
func (*apiHandle) ShowFolds(c *gin.Context) {
	var (
		menuList      = make([]*modelRes.LogIndexShowFoldsRes, 0)
		menuChildList []*modelRes.LogIndexShowFoldsRes
	)

	//#region 获取参数
	pathParam := c.Query("path")
	//#endregion

	if pathParam == "" {
		//#region 获取所有盘符
		drives := []string{}
		for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
			path := string(drive) + ":\\"
			_, err := os.Stat(path)
			if err == nil {
				drives = append(drives, path)
			}
		}
		//#endregion

		//#region 整理数据
		for k, drive := range drives {
			id := strconv.Itoa(k)
			// 读取目录中的所有文件和子目录
			entries, err := os.ReadDir(drive)
			if err != nil {
				logrus.Errorf("read dir error: %v", err)
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}

			var (
				menuChildList1 = make([]*modelRes.LogIndexShowFoldsRes, 0)
				menuChildList2 = make([]*modelRes.LogIndexShowFoldsRes, 0)
			)

			// 打印每个子目录的名称
			for k1, entry := range entries {
				var (
					title string
					id1   = strconv.Itoa(k1)
				)

				if entry.IsDir() {
					title = "<span class=\"layui-badge layui-bg-orange\">F</span> " + entry.Name()
					menuChildList1 = append(menuChildList1, &modelRes.LogIndexShowFoldsRes{
						Title:  title,
						Id:     id1,
						PathId: fmt.Sprintf("%s-%s", id, id1),
					})
				} else {
					title = "<span class=\"layui-badge layui-bg-green\">D</span> " + entry.Name()
					menuChildList2 = append(menuChildList2, &modelRes.LogIndexShowFoldsRes{
						Title:  title,
						Id:     id1,
						PathId: fmt.Sprintf("%s-%s", id, id1),
					})
				}
			}

			menuChildList = append(menuChildList1, menuChildList2...)

			menuList = append(menuList, &modelRes.LogIndexShowFoldsRes{
				Title:    drive,
				Id:       id,
				PathId:   id,
				Children: menuChildList,
			})
		}
		//#endregion
	}

	c.JSON(http.StatusOK, menuList)
}

// ReadFile 读取文件
func (*apiHandle) ReadFile(c *gin.Context) {
	//#region 获取参数
	var r modelReq.LogIndexReadFileReq
	if err := c.ShouldBind(&r); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	limit := 1000

	n := (r.Page - 1) * limit
	m := limit

	if r.Path == "" || r.Name == "" {
		c.JSON(http.StatusBadRequest, "参数错误")
		return
	}
	//#endregion

	//#region 拼接文件路径
	var filepath string
	if runtime.GOOS == "windows" {
		filepath = fmt.Sprintf("%s\\%s", r.Path, r.Name)
	} else {
		filepath = fmt.Sprintf("%s/%s", r.Path, r.Name)
	}
	logrus.Infof("filepath: %s", filepath)
	//#endregion

	//#region 读取文件
	lines, err := tools.ReadFileLine(filepath, m, n)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	//#endregion

	c.JSON(http.StatusOK, lines)
}

// AddSsh 添加ssh
func (*apiHandle) AddSsh(c *gin.Context) {
	req := new(modelReq.AddSshReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Errorf("AddSsh req error: %v", err)
		c.JSON(http.StatusOK, resp.Param)
		return
	}

	err := service.ApiService.AddSsh(req)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailResp(resp.FailCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.Succ)
}

// DelSsh 删除ssh
func (*apiHandle) DelSsh(c *gin.Context) {
	req := new(modelReq.DelSshReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Errorf("DelSsh req error: %v", err)
		c.JSON(http.StatusOK, resp.Param)
		return
	}

	err := service.ApiService.DelSsh(req)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailResp(resp.FailCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.Succ)
}

// UpdateSsh 更新ssh
func (*apiHandle) UpdateSsh(c *gin.Context) {
	req := new(modelReq.UpdateSshReq)
	if err := c.ShouldBindJSON(req); err != nil {
		logrus.Errorf("UpdateSsh req error: %v", err)
		c.JSON(http.StatusOK, resp.Param)
		return
	}

	err := service.ApiService.UpdateSsh(req)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailResp(resp.FailCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.Succ)
}

// DetailSSh 获取ssh详情
func (*apiHandle) DetailSsh(c *gin.Context) {
	req := new(modelReq.DetailSshReq)
	if err := c.ShouldBind(req); err != nil {
		logrus.Errorf("DetailSSh req error: %v", err)
		c.JSON(http.StatusOK, resp.Param)
		return
	}

	data, err := service.ApiService.DetailSsh(req)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailResp(resp.FailCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.SuccResp(data))
}

// ListSsh 获取ssh列表
func (*apiHandle) ListSsh(c *gin.Context) {
	req := new(modelReq.ListSshReq)
	if err := c.ShouldBind(req); err != nil {
		logrus.Errorf("ListSsh req error: %v", err)
		c.JSON(http.StatusOK, resp.Param)
		return
	}

	data, err := service.ApiService.ListSsh(req)
	if err != nil {
		c.JSON(http.StatusOK, resp.FailResp(resp.FailCode, err.Error()))
		return
	}
	c.JSON(http.StatusOK, resp.SuccResp(data))
}
