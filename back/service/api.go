package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"ViewLog/back/common/resp"
	"ViewLog/back/global"
	"ViewLog/back/model"
	modelReq "ViewLog/back/model/req"
	"ViewLog/back/tools/constant"
	toolsCrypto "ViewLog/back/tools/crypto"
	toolsPool "ViewLog/back/tools/pool"
	toolsSsh "ViewLog/back/tools/ssh"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"xorm.io/xorm"
)

var (
	dbYamlPath = "back/configs/db.yaml"
	lockPath   = "install.lock"
)

type apiService struct {
	DetailFileLock sync.Mutex
}

var ApiService = new(apiService)

// Install 安装
func (th *apiService) Install(req *modelReq.InstallReq) (interface{}, error) {
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

	//#region 写入文件
	dbYaml := fmt.Sprintf("mysql: %s", mysqlDns)
	if err := os.WriteFile(dbYamlPath, []byte(dbYaml), 0666); err != nil {
		logrus.Errorf("写入文件失败: %v", err)
		return nil, err
	} else {
		logrus.Info("写入db.yaml文件成功")
		global.Conf.Mysql = mysqlDns
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
			} else {
				logrus.Info("写入install.lock文件成功")
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
	//#region 校验
	if req.Host == "127.0.0.1" || req.Host == "localhost" {
		return errors.New("不允许添加本机ssh(研发中)")
	}
	//#endregion

	sess := global.Db.NewSession()

	//#region 事务开始
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	//#endregion

	//#region 校验
	//#region 名称
	if total, err := sess.Where("name = ?", req.Name).Count(&model.Ssh{}); err != nil {
		logrus.Errorf("查询ssh名称是否存在失败: %v", err)
		return err
	} else if total > 0 {
		return errors.New("ssh名称已存在")
	}
	//#endregion

	//#region ssh真实性
	sshConf := toolsSsh.Config{
		Host:     req.Host,
		Port:     req.Port,
		User:     req.Username,
		Password: req.Password,
	}
	sshClient, err := sshConf.Connect()
	if err != nil {
		logrus.Errorf("ssh连接失败: %v", err)
		return errors.New("ssh连接失败")
	}
	if global.SshClient != nil {
		_ = global.SshClient.Close()
	}
	global.SshClient = sshClient
	//#endregion
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

	//#region 重置活跃ssh状态
	if _, err := sess.Where("is_active = ?", constant.SshIsActiveYES).Cols("is_active").Update(&model.Ssh{IsActive: constant.SshIsActiveNO}); err != nil {
		logrus.Errorf("重置活跃ssh状态失败: %v", err)
		return err
	}
	//#endregion

	//#region 添加
	if _, err := sess.Insert(&model.Ssh{
		Name:     req.Name,
		Host:     req.Host,
		Port:     req.Port,
		IsActive: constant.SshIsActiveYES,
		Username: req.Username,
		Password: password,
	}); err != nil {
		logrus.Errorf("添加ssh失败: %v", err)
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

// DelSsh 删除ssh
func (*apiService) DelSsh(req *modelReq.DelSshReq) error {
	sess := global.Db.NewSession()

	//#region 事务开始
	defer sess.Close()
	if err := sess.Begin(); err != nil {
		return err
	}
	//#endregion

	//#region 删除ssh
	if _, err := sess.ID(req.Id).Delete(&model.Ssh{}); err != nil {
		logrus.Errorf("删除ssh失败: %v", err)
		return err
	}
	//#endregion

	//#region 删除ssh关联的folder
	if _, err := sess.Where("ssh_id = ?", req.Id).Delete(&model.Folder{}); err != nil {
		logrus.Errorf("删除ssh关联的folder失败: %v", err)
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
	password, err := toolsCrypto.AesDecrypt(sshInfo.Password)
	if err != nil {
		logrus.Errorf("解密密码失败: %v", err)
		return err
	}
	sshConf := &toolsSsh.Config{
		Host:     sshInfo.Host,
		Port:     sshInfo.Port,
		User:     sshInfo.Username,
		Password: password,
	}
	sshClient, err := sshConf.Connect()
	if err != nil {
		logrus.Errorf("连接ssh失败: %v", err)
		return errors.New("连接ssh失败")
	} else {
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
		global.SshClient = sshClient
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

// ListSshFolder 列表ssh文件夹
func (*apiService) ListSshFolder(req *modelReq.ListSshFolderReq) (*resp.ListResult, error) {
	sess := global.Db

	var folderList = make([]*model.Folder, 0)
	if req.Limit > 0 {
		sess.Limit(req.Limit, (req.Page-1)*req.Limit)
	}
	total, err := sess.Where("ssh_id = ?", req.SshId).OrderBy("create_time Desc").FindAndCount(&folderList)
	if err != nil {
		logrus.Errorf("查询ssh文件夹列表失败: %v", err)
		return nil, err
	}

	listResult := &resp.ListResult{
		Total: total,
		List:  folderList,
	}
	return listResult, nil
}

// AddFolder 添加文件夹
func (*apiService) AddFolder(req *modelReq.AddFolderReq) error {
	sess := global.Db

	//#region 校验
	//#region ssh是否存在
	if total, err := sess.ID(req.SshId).Count(&model.Ssh{}); err != nil {
		logrus.Errorf("查询ssh是否存在失败: %v", err)
		return err
	} else if total <= 0 {
		return errors.New("ssh不存在")
	}
	//#endregion

	//#region 文件夹名称是否存在
	if total, err := sess.Where("name = ?", req.Name).Count(&model.Folder{}); err != nil {
		logrus.Errorf("查询文件夹名称是否存在失败: %v", err)
		return err
	} else if total > 0 {
		return errors.New("文件夹名称已存在")
	}
	//#endregion

	//#region 文件夹路径是否存在
	sshClient := global.SshClient
	if sshClient == nil {
		return errors.New("ssh未激活")
	}
	sshSess, err := sshClient.NewSession()
	if err != nil {
		logrus.Errorf("创建ssh会话失败: %v", err)
		return errors.New("创建ssh会话失败")
	}
	defer sshSess.Close()
	cmd := "ls " + req.Path
	output, err := sshSess.Output(cmd)
	if err != nil {
		logrus.Errorf("文件夹路径不存在: %v", err)
		return errors.New("文件夹路径不存在")
	}
	fmt.Println(string(output))
	//#endregion
	//#endregion

	//#region 构建添加信息
	folder := &model.Folder{
		SshId: req.SshId,
		Name:  req.Name,
		Path:  req.Path,
	}
	if _, err := sess.Insert(folder); err != nil {
		logrus.Errorf("添加文件夹失败: %v", err)
		return err
	}
	//#endregion

	return nil
}

// DelFolder 删除文件夹
func (*apiService) DelFolder(req *modelReq.DelFolderReq) error {
	sess := global.Db
	if _, err := sess.ID(req.Id).Delete(&model.Folder{}); err != nil {
		logrus.Errorf("删除文件夹失败: %v", err)
		return err
	}
	return nil
}

// ListFolderChild 列表文件夹子文件夹
func (th *apiService) ListFolderChild(req *modelReq.ListFolderChildReq) (interface{}, error) {
	var (
		sess = global.Db
	)

	//#region 查询文件夹详
	folderInfo := &model.Folder{}
	if has, err := sess.ID(req.FolderId).Get(folderInfo); err != nil {
		logrus.Errorf("查询文件夹失败: %v", err)
		return nil, err
	} else if !has {
		return nil, errors.New("文件夹不存在")
	}
	//#endregion

	//#region 查询子文件夹
	//#region 穿件ssh会话
	sshSess, err := global.SshClient.NewSession()
	if err != nil {
		logrus.Errorf("创建ssh会话失败: %v", err)
		return nil, errors.New("创建ssh会话失败")
	}
	defer sshSess.Close()
	//#endregion

	//#region 获取文件
	defaultCmd := " && find . -type f -name \"*.log\" -not -path \"*/.*/*\""
	defaultCmd += " && find . -type f -name \"*.txt\" -not -path \"*/.*/*\""
	cmd := fmt.Sprintf("cd %s %s", folderInfo.Path, defaultCmd)
	output, err := sshSess.Output(cmd)
	if err != nil {
		logrus.Errorf("查询文件失败: %v", err)
		return nil, errors.New("查询文件失败")
	}
	fileList := strings.Split(string(output), "\n")
	//#endregion

	//#region 过滤空目录
	var filterFileList = make([]string, 0)
	for _, file := range fileList {
		if file != "" {
			filterFileList = append(filterFileList, file)
		}
	}
	//#endregion
	//#endregion

	return filterFileList, nil
}

func (th *apiService) ListFolderChild1(req *modelReq.ListFolderChildReq) (interface{}, error) {
	var (
		sess = global.Db
	)

	//#region 查询文件夹详
	folderInfo := &model.Folder{}
	if has, err := sess.ID(req.FolderId).Get(folderInfo); err != nil {
		logrus.Errorf("查询文件夹失败: %v", err)
		return nil, err
	} else if !has {
		return nil, errors.New("文件夹不存在")
	}
	//#endregion

	//#region 查询子文件夹
	wd, err := GetWalkDir(folderInfo.Path)
	if err != nil {
		logrus.Errorf("查询子文件夹失败: %v", err)
		return nil, errors.New("查询子文件夹失败")
	}
	//#endregion

	return wd, nil
}

type WalkDir struct {
	Title    string     `json:"title"`    // 文件夹名称
	Children []*WalkDir `json:"children"` // 子文件夹
}

// walkDir 返回指定路径下所有文件和文件夹的层次结构
func walkDir(sshSession *ssh.Session, path string, wd *WalkDir) error {
	// 使用 SSH Session 打开远程目录
	remoteDir, err := sshSession.CombinedOutput("ls -p " + path)
	if err != nil {
		logrus.Errorf("打开远程目录失败: %v", err)
		return err
	}

	// 分割 ls 命令输出的每一行
	remoteFiles := strings.Split(string(remoteDir), "\n")

	// 遍历 ls 命令输出中的每个文件/文件夹
	for _, file := range remoteFiles {
		// // 去除目录名末尾的斜杠
		// file = strings.TrimRight(file, "/")

		// 跳过以点号开头的目录（隐藏目录）
		if strings.HasPrefix(file, ".") {
			continue
		}

		// 检查文件是否是一个目录
		isDir := false
		if strings.HasSuffix(file, "/") {
			isDir = true
			file = strings.TrimRight(file, "/")
		}

		// 如果是一个目录，则递归调用 walkDir 函数
		if isDir {
			newDir := &WalkDir{
				Title:    file,
				Children: make([]*WalkDir, 0),
			}
			wd.Children = append(wd.Children, newDir)
			if err := walkDir(sshSession, filepath.Join(path, file), newDir); err != nil {
				return err
			}
		} else {
			// 如果是一个文件，则将其添加到当前目录的文件列表中
			wd.Children = append(wd.Children, &WalkDir{
				Title: file,
			})
		}
	}

	return nil
}

// GetWalkDir 返回指定路径下所有文件和文件夹的层次结构
func GetWalkDir(path string) (*WalkDir, error) {
	toolsPool.SshInit()
	sessionPool := toolsPool.SessionPool
	sshSession, err := sessionPool.GetSession()
	if err != nil {
		logrus.Errorf("获取ssh会话失败: %v", err)
		return nil, errors.New("获取ssh会话失败")
	}
	defer sessionPool.PutSession(sshSession)

	root := &WalkDir{
		Title:    path,
		Children: make([]*WalkDir, 0),
	}

	if err := walkDir(sshSession, path, root); err != nil {
		return nil, err
	}

	// 过滤空目录
	filterEmptyDir(root)

	return root, nil
}

// filterEmptyDir 过滤空目录
func filterEmptyDir(wd *WalkDir) {
	for i := 0; i < len(wd.Children); {
		if len(wd.Children[i].Children) == 0 {
			// 如果是空目录，则将其从当前目录的子目录列表中移除
			wd.Children = append(wd.Children[:i], wd.Children[i+1:]...)
		} else {
			// 如果不是空目录，则递归调用 filterEmptyDir 函数
			filterEmptyDir(wd.Children[i])
			i++
		}
	}
}

// DetailFile 查询文件
func (th *apiService) DetailFile(req *modelReq.DetailFileReq) (interface{}, error) {
	th.DetailFileLock.Lock()
	defer th.DetailFileLock.Unlock()
	var (
		sess     = global.Db
		filePath string // 文件路径
	)

	//#region 获取文件夹详情
	folderInfo := &model.Folder{}
	if has, err := sess.ID(req.FolderId).Get(folderInfo); err != nil {
		logrus.Errorf("查询文件夹失败: %v", err)
		return nil, err
	} else if !has {
		return nil, errors.New("文件夹不存在")
	}
	filePath = folderInfo.Path + req.Path[1:]
	//#endregion

	//#region 字符串过滤

	for _, v := range global.Conf.FilterStr {
		fStr := string(v)
		filePath = strings.ReplaceAll(filePath, fStr, fmt.Sprintf("\\%s", fStr))
	}
	// filePath = strings.ReplaceAll(filePath, ")", "\\)")
	//#endregion

	//#region 查询文件
	//#region 穿件ssh会话
	sshSess, err := global.SshClient.NewSession()
	if err != nil {
		logrus.Errorf("创建ssh会话失败: %v", err)
		return nil, errors.New("创建ssh会话失败")
	}
	defer sshSess.Close()
	//#endregion

	//#region 查询最后1000行
	logrus.Infof("查询文件: %s", filePath)
	cmd := fmt.Sprintf("tail -n 1000 %s", filePath)
	output, err := sshSess.Output(cmd)
	if err != nil {
		logrus.Errorf("查询文件失败: %v", err)
		logrus.Errorf("查询文件命令: %s", cmd)
		return nil, errors.New("查询文件失败")
	}
	//#endregion
	//#endregion

	return string(output), nil
}
