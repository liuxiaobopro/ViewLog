package ssh

import (
	"fmt"

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
