package global

import (
	"ViewLog/back/configs"

	"golang.org/x/crypto/ssh"
	"xorm.io/xorm"
)

var (
	Db        *xorm.Engine
	Conf      *configs.Conf
	SshClient *ssh.Client
)
