package ready

import (
	"os"

	"ViewLog/back/configs"
	"ViewLog/back/global"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func Config() {
	var (
		err  error
		path = "back/configs/"
		conf configs.Conf
	)
	// #region 读取通用配置
	var configPath string =path + "config.yaml"
	var configYaml []byte
	if configYaml, err = os.ReadFile(configPath); err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(configYaml, &conf); err != nil {
		panic(err)
	}
	// #endregion

	// #region 读取db
	var dbPath string =path + "db.yaml"
	var dbYaml []byte
	if dbYaml, err = os.ReadFile(dbPath); err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(dbYaml, &conf); err != nil {
		panic(err)
	}
	// #endregion

	global.Conf = &conf
	logrus.Infof("config: %+v", global.Conf)
}
