package controller

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"ViewLog/back/common/tools"
	modelReq "ViewLog/back/model/req"
	modelRes "ViewLog/back/model/res"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type apiHandle struct{}

var ApiHandle = new(apiHandle)

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

func (*apiHandle) ShowFolds(c *gin.Context) {
	var (
		menuList      = make([]*modelRes.LogIndexShowFoldsRes, 0)
		menuChildList = make([]*modelRes.LogIndexShowFoldsRes, 0)
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
