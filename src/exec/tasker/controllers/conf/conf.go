package conf

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/rds"
)

/* 应用升级配置文件 */
type Conf struct {
	Log   string     // 日志配置
	Mysql mysql.Conf // mysql 配置
	Redis rds.Conf   // redis 配置
}

func Load() (conf *Conf, err error) {
	conf = new(Conf)
	//配置
	conf.Mysql.Conn = beego.AppConfig.String("mysql_conn")
	if 0 == len(conf.Mysql.Conn) {
		return nil, errors.New("Get mysql configuration failed")
	}

	// 加载日志配置
	conf.Log = beego.AppConfig.String("log_conf")
	if 0 == len(conf.Log) {
		return nil, errors.New("Get log configuration failed")
	}

	// 加载redis配置
	conf.Redis.Conn = beego.AppConfig.String("redis_conn")
	conf.Redis.Pwd = beego.AppConfig.String("redis_pwd")
	conf.Redis.MaxIdel, _ = beego.AppConfig.Int("redis_max_idel")
	if 0 == len(conf.Redis.Conn) || 0 == len(conf.Redis.Pwd) || conf.Redis.MaxIdel == 0 {
		return nil, errors.New("Get redis configuration failed")
	}

	logs.Info("%s", "conf load complete")
	return conf, nil
}
