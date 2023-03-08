package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"time"

	"ViewLog/back/common/resp"
	"ViewLog/back/global"
	"ViewLog/back/model"
	modelReq "ViewLog/back/model/req"
	"ViewLog/back/tools/constant"
	toolsCrypto "ViewLog/back/tools/crypto"

	"github.com/sirupsen/logrus"
	"xorm.io/xorm"
)

var (
	dbYamlPath = "back/configs/db.yaml"
	lockPath   = "install.lock"
)

type apiService struct{}

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
		logrus.Errorf("数据库连接失败, 请检查配置是否正确: %v", err)
		return nil, err
	}

	if err := engine.Ping(); err != nil {
		logrus.Errorf("数据库连接失败, 请检查配置是否正确: %v", err)
		return nil, errors.New("数据库连接失败, 请检查配置是否正确")
	}

	global.Db = engine

	//#region 写入文件
	dbYaml := fmt.Sprintf("mysql: %s", mysqlDns)
	if err := os.WriteFile(dbYamlPath, []byte(dbYaml), 0666); err != nil {
		logrus.Errorf("写入文件失败: %v", err)
		return nil, err
	}
	//#endregion

	Db()
	Config()

	//#region 生成文件锁
	if _, err := os.Stat(lockPath); err != nil {
		logrus.Errorf("读取文件失败: %v", err)
		if os.IsNotExist(err) {
			content := "install: " + time.Now().Format("2006-01-02 15:04:05")
			if err := os.WriteFile(lockPath, []byte(content), 0666); err != nil {
				logrus.Errorf("写入文件失败: %v", err)
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
		logrus.Errorf("删除文件失败1: %v", err)
		return err
	}
	if err := os.Remove(dbYamlPath); err != nil {
		logrus.Errorf("删除文件失败2: %v", err)
		return err
	}
	//#endregion
	return nil
}

// AddSsh 添加ssh
func (*apiService) AddSsh(req *modelReq.AddSshReq) error {
	sess := global.Db

	//#region 校验
	if total, err := sess.Where("name = ?", req.Name).Count(&model.Ssh{}); err != nil {
		logrus.Errorf("查询ssh名称是否存在失败: %v", err)
		return err
	} else if total > 0 {
		return errors.New("ssh名称已存在")
	}
	//#endregion

	//#region 加密密码
	password, err := toolsCrypto.AesEncrypt(req.Password)
	if err != nil {
		logrus.Errorf("加密密码失败: %v", err)
		return err
	}
	//#endregion

	// //#region debug aes解密
	// p, err := toolsCrypto.AesDecrypt(password)
	// if err != nil {
	// 	logrus.Errorf("解密密码失败: %v", err)
	// 	return err
	// }
	// logrus.Infof("解密密码: %s \n", p)
	// //#endregion

	//#region 添加
	if _, err := sess.Insert(&model.Ssh{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		Username: req.Username,
		Password: password,
	}); err != nil {
		logrus.Errorf("添加ssh失败: %v", err)
		return err
	}
	//#endregion
	return nil
}

// DelSsh 删除ssh
func (*apiService) DelSsh(req *modelReq.DelSshReq) error {
	sess := global.Db

	if _, err := sess.ID(req.Id).Delete(&model.Ssh{}); err != nil {
		logrus.Errorf("删除ssh失败: %v", err)
		return err
	}
	return nil
}

// UpdateSsh 更新ssh
func (*apiService) UpdateSsh(req *modelReq.UpdateSshReq) error {
	sess := global.Db

	//#region 校验
	sshInfo := &model.Ssh{}
	if has, err := sess.ID(req.Id).Get(sshInfo); err != nil {
		logrus.Errorf("查询ssh失败: %v", err)
		return err
	} else if !has {
		return errors.New("ssh不存在")
	}

	if total, err := sess.Where("name = ?", req.Name).Count(&model.Ssh{}); err != nil {
		logrus.Errorf("查询ssh名称是否存在失败: %v", err)
		return err
	} else if total > 0 {
		return errors.New("ssh名称已存在")
	}
	//#endregion

	//#region 构建更新信息
	var (
		updateSsh  = &model.Ssh{}
		updateCols = make([]string, 0)
	)
	if req.Name != "" && req.Name != sshInfo.Name {
		updateSsh.Name = req.Name
		updateCols = append(updateCols, "name")
	}
	if req.Host != "" && req.Host != sshInfo.Host {
		updateSsh.Host = req.Host
		updateCols = append(updateCols, "host")
	}
	if req.Port != 0 && req.Port != sshInfo.Port {
		updateSsh.Port = req.Port
		updateCols = append(updateCols, "port")
	}
	if req.Username != "" && req.Username != sshInfo.Username {
		updateSsh.Username = req.Username
		updateCols = append(updateCols, "username")
	}
	if req.Password != "" && req.Password != fmt.Sprintf("%x", md5.Sum([]byte(sshInfo.Password))) {
		password, err := toolsCrypto.AesEncrypt(req.Password)
		if err != nil {
			logrus.Errorf("加密密码失败: %v", err)
			return err
		}
		updateSsh.Password = password
		updateCols = append(updateCols, "password")
	}
	if req.IsActive != 0 && req.IsActive != sshInfo.IsActive {
		updateSsh.IsActive = req.IsActive
		updateCols = append(updateCols, "is_active")
	}
	//#endregion

	//#region 更新
	if _, err := sess.ID(req.Id).Cols(updateCols...).Update(updateSsh); err != nil {
		logrus.Errorf("更新ssh失败: %v", err)
		return err
	}
	//#endregion
	return nil
}

// DetailSSh 详情ssh
func (*apiService) DetailSsh(req *modelReq.DetailSshReq) (*model.Ssh, error) {
	sess := global.Db

	sshInfo := &model.Ssh{}
	if has, err := sess.ID(req.Id).Get(sshInfo); err != nil {
		logrus.Errorf("查询ssh失败: %v", err)
		return nil, err
	} else if !has {
		return nil, errors.New("ssh不存在")
	}
	return sshInfo, nil
}

// ListSsh 列表ssh
func (*apiService) ListSsh(req *modelReq.ListSshReq) (*resp.ListResult, error) {
	sess := global.Db

	var sshList = make([]*model.Ssh, 0)
	total, err := sess.Limit(req.Limit, (req.Page-1)*req.Limit).OrderBy("create_time Desc").FindAndCount(&sshList)
	if err != nil {
		logrus.Errorf("查询ssh列表失败: %v", err)
		return nil, err
	}

	listResult := &resp.ListResult{
		Total: total,
		List:  sshList,
	}
	return listResult, nil
}

// UpdateActiveSsh 更新ssh激活状态
func (*apiService) UpdateActiveSsh(req *modelReq.UpdateActiveSshReq) error {
	sess := global.Db.NewSession()

	//#region 事务开始
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	//#endregion

	//#region 校验
	sshInfo := &model.Ssh{}
	if has, err := sess.ID(req.Id).Get(sshInfo); err != nil {
		logrus.Errorf("查询ssh失败: %v", err)
		return err
	} else if !has {
		return errors.New("ssh不存在")
	}

	if req.IsActive == sshInfo.IsActive {
		return errors.New("ssh激活状态未改变")
	}
	//#endregion

	//#region 更新
	if req.IsActive == constant.SshIsActiveYES {
		if _, err := sess.Where("is_active = ?", constant.SshIsActiveYES).Cols("is_active").Update(&model.Ssh{IsActive: constant.SshIsActiveNO}); err != nil {
			logrus.Errorf("重置所有ssh激活状态失败: %v", err)
			return err
		}
	}
	if _, err := sess.ID(req.Id).Cols("is_active").Update(&model.Ssh{IsActive: req.IsActive}); err != nil {
		logrus.Errorf("更新ssh激活状态失败: %v", err)
		return err
	}
	//#endregion

	//#region 提交事务
	if err := sess.Commit(); err != nil {
		_ = sess.Rollback()
		return err
	}
	//#endregion

	return nil
}
