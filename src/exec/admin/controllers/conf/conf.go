package conf

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/sso"
)

/* 应用升级配置文件 */
type AdminConf struct {
	Log   string           // 日志配置
	Mysql mysql.Conf       // mysql 配置
	InSso *sso.InternalSso // 内网用户认证
}

func Load() (conf *AdminConf, err error) {
	conf = new(AdminConf)
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

	// 加载internal sso 配置
	conf.InSso, err = LoadInternalSsoConf()
	if nil != err {
		return nil, err
	}

	logs.Info("%s", "conf load complete")
	return conf, nil
}

/* 加载InternalSso配置 */
func LoadInternalSsoConf() (conf *sso.InternalSso, err error) {
	conf = new(sso.InternalSso)
	// 加载InternalSso配置
	conf.Host = beego.AppConfig.String("in_sso_host")
	if 0 == len(conf.Host) {
		return nil, errors.New("Get internal sso host failed")
	}

	conf.Site = beego.AppConfig.String("in_sso_site")
	if 0 == len(conf.Site) {
		return nil, errors.New("Get internal sso site failed")
	}

	conf.Key = beego.AppConfig.String("in_sso_key")
	if 0 == len(conf.Key) {
		return nil, errors.New("Get internal sso key failed")
	}

	conf.TransCode = beego.AppConfig.String("in_sso_transcode")
	if 0 == len(conf.TransCode) {
		return nil, errors.New("Get internal sso trans code uri failed")
	}

	conf.UserDetail = beego.AppConfig.String("in_sso_user_detail")
	if 0 == len(conf.UserDetail) {
		return nil, errors.New("Get internal sso user detail uri failed")
	}

	return conf, nil
}
