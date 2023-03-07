package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"ViewLog/back/global"
	"ViewLog/back/model"
	modelReq "ViewLog/back/model/req"

	"xorm.io/xorm"
)

var (
	dbYamlPath = "back/configs/db.yaml"
	lockPath   = "install.lock"
	sess       = global.Db
)

type apiService struct {
}

var ApiService = new(apiService)

// Install 安装
func (th *apiService) Install(req *modelReq.InstallReq) (any, error) {
	mysqlDns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		req.User,
		req.Password,
		req.Host,
		req.Port,
		req.Dbname,
		req.Charset)

	engine, err := xorm.NewEngine("mysql", mysqlDns)
	if err != nil {
		return nil, err
	}

	if err := engine.Ping(); err != nil {
		return nil, errors.New("数据库连接失败, 请检查配置是否正确")
	}

	global.Db = engine

	//#region 写入文件
	dbYaml := fmt.Sprintf("mysql: %s", mysqlDns)
	if err := os.WriteFile(dbYamlPath, []byte(dbYaml), 0666); err != nil {
		return nil, err
	}
	//#endregion

	Db()
	Config()

	//#region 生成文件锁
	if _, err := os.Stat(lockPath); err != nil {
		if os.IsNotExist(err) {
			content := "install: " + time.Now().Format("2006-01-02 15:04:05")
			if err := os.WriteFile(lockPath, []byte(content), 0666); err != nil {
				return nil, err
			}
		}
	}
	//#endregion

	return nil, nil
}

// Reset 重置
func (th *apiService) Reset() error {
	//#region 删除lock和db.yaml
	if err := os.Remove(lockPath); err != nil {
		return err
	}
	if err := os.Remove(dbYamlPath); err != nil {
		return err
	}
	//#endregion
	return nil
}

// AddSsh 添加ssh
func (*apiService) AddSsh(req *modelReq.AddSshReq) error {
	//#region 校验
	if total, err := sess.Where("name = ?", req.Name).Count(&model.Ssh{}); err != nil {
		return err
	} else if total > 0 {
		return errors.New("ssh名称已存在")
	}
	//#endregion

	//#region 添加
	if _, err := sess.Insert(&model.Ssh{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: req.Password,
	}); err != nil {
		return err
	}
	//#endregion
	return nil
}

// DelSsh 删除ssh
func (*apiService) DelSsh(req *modelReq.DelSshReq) error {
	if _, err := sess.ID(req.Id).Delete(&model.Ssh{}); err != nil {
		return err
	}
	return nil
}

// UpdateSsh 更新ssh
func (*apiService) UpdateSsh(req *modelReq.UpdateSshReq) error {
	//#region 校验
	var sshInfo *model.Ssh
	if has, err := sess.ID(req.Id).Get(sshInfo); err != nil {
		return err
	} else if !has {
		return errors.New("ssh不存在")
	}

	if total, err := sess.Where("name = ?", req.Name).Count(&model.Ssh{}); err != nil {
		return err
	} else if total > 0 {
		return errors.New("ssh名称已存在")
	}
	//#endregion

	//#region 构建更新信息
	var (
		updateSSh  *model.Ssh
		updateCols []string
	)
	if req.Name != "" && req.Name != sshInfo.Name {
		updateSSh.Name = req.Name
		updateCols = append(updateCols, "name")
	}
	if req.Host != "" && req.Host != sshInfo.Host {
		updateSSh.Host = req.Host
		updateCols = append(updateCols, "host")
	}
	if req.Port != 0 && req.Port != sshInfo.Port {
		updateSSh.Port = req.Port
		updateCols = append(updateCols, "port")
	}
	if req.Username != "" && req.Username != sshInfo.Username {
		updateSSh.Username = req.Username
		updateCols = append(updateCols, "username")
	}
	if req.Password != "" && req.Password != sshInfo.Password {
		updateSSh.Password = req.Password
		updateCols = append(updateCols, "password")
	}
	//#endregion

	//#region 更新
	if _, err := sess.ID(req.Id).Cols(updateCols...).Update(updateSSh); err != nil {
		return err
	}
	//#endregion
	return nil
}

// DetailSSh 详情ssh
func (*apiService) DetailSsh(req *modelReq.DetailSshReq) (*model.Ssh, error) {
	var sshInfo *model.Ssh
	if has, err := sess.ID(req.Id).Get(sshInfo); err != nil {
		return nil, err
	} else if !has {
		return nil, errors.New("ssh不存在")
	}
	return sshInfo, nil
}

// ListSsh 列表ssh
func (*apiService) ListSsh(req *modelReq.ListSshReq) (any, error) {
	var sshList = make([]*model.Ssh, 0)
	if err := sess.Limit((req.Page-1)*req.Limit, req.Limit).Find(sshList); err != nil {
		return nil, err
	}
	return sshList, nil
}
