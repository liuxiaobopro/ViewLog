package ssh

import (
	"fmt"

	"ViewLog/back/global"
	"ViewLog/back/model"
	"ViewLog/back/tools/constant"
	toolsCrypto "ViewLog/back/tools/crypto"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	User     string
	Password string
	Host     string
	Port     int
}

func (th *Config) Connect() (*ssh.Client, error) {
	// 连接到远程服务器
	config := &ssh.ClientConfig{
		User: th.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(th.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 不验证目标主机的公钥
	}
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", th.Host, th.Port), config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func UpdateGlobalClient() error {
	sess := global.Db

	//#region 查询活跃的ssh
	sshInfo := &model.Ssh{}
	if _, err := sess.Where("is_active=?", constant.SshIsActiveYES).Get(sshInfo); err != nil {
		logrus.Errorf("查询活跃的ssh失败: %v", err)
		return err
	}
	//#endregion

	if sshInfo.Id != 0 {
		//#region 关闭旧的sshClient
		if global.SshClient != nil {
			if err := global.SshClient.Close(); err != nil {
				logrus.Errorf("关闭旧的sshClient失败: %v", err)
				return err
			}
		}
		//#endregion

		//#region 更新全局sshClient
		password, err := toolsCrypto.AesDecrypt(sshInfo.Password)
		if err != nil {
			logrus.Errorf("解密密码失败: %v", err)
			return err
		}
		sshConfig := &Config{
			User:     sshInfo.Username,
			Password: password,
			Host:     sshInfo.Host,
			Port:     sshInfo.Port,
		}
		sshClient, err := sshConfig.Connect()
		if err != nil {
			logrus.Errorf("连接ssh失败: %v", err)
			return err
		}
		global.SshClient = sshClient
		//#endregion
	}
	return nil
}
