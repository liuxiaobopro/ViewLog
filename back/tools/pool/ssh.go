package pool

import (
	"fmt"
	"time"

	"ViewLog/back/global"
	"ViewLog/back/model"
	"ViewLog/back/tools/constant"
	toolsCrypto "ViewLog/back/tools/crypto"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

type sshSessionPool struct {
	pool      chan *ssh.Session
	sshConfig *ssh.ClientConfig
	sshAddr   string
}

func newSshSessionPool(sshConfig *ssh.ClientConfig, sshAddr string, poolSize int) *sshSessionPool {
	p := &sshSessionPool{
		pool:      make(chan *ssh.Session, poolSize),
		sshConfig: sshConfig,
		sshAddr:   sshAddr,
	}
	return p
}

func (p *sshSessionPool) GetSession() (*ssh.Session, error) {
	select {
	case session := <-p.pool:
		// 如果从池中获取到了一个 session 对象，则将其 stdout 属性设置为 nil 并返回
		session.Stdout = nil
		return session, nil
	default:
		// 如果池中没有 session 对象可用，则创建一个新的 session 对象并返回
		client, err := ssh.Dial("tcp", p.sshAddr, p.sshConfig)
		if err != nil {
			return nil, err
		}
		session, err := client.NewSession()
		if err != nil {
			return nil, err
		}
		return session, nil
	}
}

func (p *sshSessionPool) PutSession(session *ssh.Session) {
	// 将 session 对象还回池中
	p.pool <- session
}

var SessionPool *sshSessionPool

func SshInit() {
	if SessionPool == nil {
		sess := global.Db

		sshInfo := &model.Ssh{}
		if _, err := sess.Where("is_active=?", constant.SshIsActiveYES).Get(sshInfo); err != nil {
			logrus.Errorf("查询活跃的ssh失败: %v", err)
			return
		}

		password, err := toolsCrypto.AesDecrypt(sshInfo.Password)
		if err != nil {
			logrus.Errorf("解密密码失败: %v", err)
			return
		}
		sshConfig := &ssh.ClientConfig{
			User: sshInfo.Username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         time.Duration(global.Conf.SshTimeout) * time.Second,
		}
		sshAddr := fmt.Sprintf("%s:%d", sshInfo.Host, sshInfo.Port)
		poolSize := 10
		SessionPool = newSshSessionPool(sshConfig, sshAddr, poolSize)
	}
}
