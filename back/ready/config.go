package ready

import (
	"flag"
	"fmt"
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

		runmode string
	)

	//#region 获取命令行执行参数
	flag.StringVar(&runmode, "runmode", "", "please input runmode")
	flag.Parse()

	var configPath string
	if runmode == "" {
		configPath = path + "config.yaml"
	} else {
		configPath = fmt.Sprintf("%sconfig_%s.yaml", path, runmode)
	}
	//#endregion

	//#region 读取yaml并映射到conf
	var contentYaml []byte
	if contentYaml, err = os.ReadFile(configPath); err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(contentYaml, &conf); err != nil {
		panic(err)
	}
	//#endregion

	conf.Host = "0.0.0.0"
	conf.Runmode = runmode

	global.Conf = &conf
	logrus.Infof("config: %+v", global.Conf)
}
