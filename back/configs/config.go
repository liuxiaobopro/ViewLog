package configs

type Conf struct {
	Runmode string `yaml:"runmode"` // 运行模式

	Host string `yaml:"host"` // 服务地址
	Port int    `yaml:"post"` // 服务端口

	DB struct {
		Host     string `yaml:"host"`     // 数据库地址
		Port     int    `yaml:"port"`     // 数据库端口
		User     string `yaml:"user"`     // 数据库用户名
		Password string `yaml:"password"` // 数据库密码
		Name     string `yaml:"name"`     // 数据库名
		Charset  string `yaml:"charset"`  // 数据库编码
	} `yaml:"db"`
}
