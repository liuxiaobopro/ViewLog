package service

import (
	"os"
	"time"

	"ViewLog/back/configs"
	"ViewLog/back/global"
	"ViewLog/back/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"xorm.io/core"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
)

func Db() {
	if global.Conf.Mysql == "" {
		logrus.Warningln("mysql连接字符串为空")
		return
	}
	if global.Db != nil {
		if err := global.Db.Ping(); err == nil {
			return
		}
	}
	engine, err := xorm.NewEngine("mysql", global.Conf.Mysql)
	if err != nil {
		panic(err)
	}

	engine.ShowSQL(true)                     // 显示SQL语句
	engine.Logger().SetLevel(log.LOG_DEBUG)  // 设置日志级别
	engine.SetMaxIdleConns(2)                // 设置最大空闲连接数
	engine.SetMaxOpenConns(100)              // 设置最大连接数
	engine.SetConnMaxLifetime(4 * time.Hour) // 设置连接最大存活时间
	engine.SetTableMapper(core.NewPrefixMapper(core.SnakeMapper{}, "vl_"))

	if err != nil {
		panic(err)
	}

	if err := engine.Sync(
		new(model.Demo),
		new(model.Ssh),
	); err != nil {
		panic(err)
	}

	global.Db = engine
}

func Config() {
	var (
		err  error
		path = "back/configs/"
		conf configs.Conf
	)
	defer func(conf *configs.Conf) {
		global.Conf = conf
		logrus.Infof("config: %+v", global.Conf)
	}(&conf)
	// #region 读取通用配置
	var configPath string = path + "config.yaml"
	var configYaml []byte
	if configYaml, err = os.ReadFile(configPath); err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(configYaml, &conf); err != nil {
		panic(err)
	}
	// #endregion

	// #region 读取db
	var dbPath string = path + "db.yaml"
	if _, err = os.Stat(dbPath); err != nil {
		if os.IsNotExist(err) {
			logrus.Warningln("配置文件db.yaml不存在")
			return
		}
	}
	var dbYaml []byte
	if dbYaml, err = os.ReadFile(dbPath); err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(dbYaml, &conf); err != nil {
		panic(err)
	}
	// #endregion
}
