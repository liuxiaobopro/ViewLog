package configs

type Conf struct {
	Runmode string `yaml:"runmode"` // 运行模式

	Host string `yaml:"host"` // 服务地址
	Port int    `yaml:"post"` // 服务端口

	Mysql string `yaml:"mysql"` // mysql连接字符串
}
