package conf

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/cdn"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/rds"
)

/* 应用升级配置文件 */
type ApiConf struct {
	Log        string         // 日志配置
	Mysql      mysql.Conf     // mysql 配置
	Redis      rds.Conf       // redis 配置
	CdnSplat   *cdn.SplatConf // CdnSplatConf
	ApkChanLen int            // apk worker channel 长度
}

func Load() (conf *ApiConf, err error) {
	conf = new(ApiConf)

	// 加载MySQL配置
	conf.Mysql.Conn = beego.AppConfig.String("mysql_conn")
	if 0 == len(conf.Mysql.Conn) {
		return nil, errors.New("Get mysql connection failed!")
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

	// 加载Cdn splat配置
	conf.CdnSplat, err = LoadCdnSplatConf()
	if nil != err {
		return nil, err
	}

	// 加载apk channel 长度
	conf.ApkChanLen, err = beego.AppConfig.Int("apk_chan_len")
	if nil != err {
		conf.ApkChanLen = comm.WORKER_APK_CHANNEL_LEN_DEFAULT
	}

	logs.Info("%s", "conf load complete")
	return conf, nil
}

/* 加载Cdn splat配置 */
func LoadCdnSplatConf() (conf *cdn.SplatConf, err error) {
	conf = new(cdn.SplatConf)
	// 加载Mysql参数
	conf.Uri = beego.AppConfig.String("cdn_splat_uri")
	if 0 == len(conf.Uri) {
		return nil, errors.New("Get cdn splat uri failed")
	}
	return conf, nil
}
