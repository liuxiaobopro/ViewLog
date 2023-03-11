package configs

type Conf struct {
	Host string `yaml:"host"` // 服务地址
	Port int    `yaml:"post"` // 服务端口

	Mysql string `yaml:"mysql"` // mysql连接字符串

	Aes struct {
		Key string `yaml:"key"` // aes加密key
		IV  string `yaml:"iv"`  // aes加密iv
	} `yaml:"aes"` // aes加密配置

	SshTimeout int `yaml:"sshTimeout"` // ssh超时时间

	FilterStr string `yaml:"filterStr"` // 过滤字符串
}
