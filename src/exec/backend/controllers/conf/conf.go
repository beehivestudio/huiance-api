package conf

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/cdn"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/rds"
	"upgrade-api/src/share/sso"
)

/* 应用升级配置文件 */
type Conf struct {
	Log        string           // 日志配置
	ApkChanLen int              // apk worker channel 长度
	Mysql      *mysql.Conf      // mysql 配置
	Redis      *rds.Conf        // redis 配置
	InSso      *sso.InternalSso // 内网用户认证
	CdnSplat   *cdn.SplatConf   // CdnSplatConf
}

/******************************************************************************
 **函数名称: Load
 **功    能: 加载配置信息
 **输入参数: NONE
 **输出参数: NONE
 **返    回:
 **     conf: 配置信息
 **     err: 错误描述
 **实现描述:
 **注意事项: 默认加载配置文件为${BIN_PATH}/conf/app.conf
 **作    者: # Shuangpeng.guo # 2020-06-01 12:48:51 #
 ******************************************************************************/
func Load() (conf *Conf, err error) {
	conf = new(Conf)

	//配置
	conf.Mysql, err = LoadMysqlConf()
	if nil != err {
		return nil, err
	}

	// 加载日志配置
	conf.Log = beego.AppConfig.String("log_conf")
	if 0 == len(conf.Log) {
		return nil, errors.New("Get log configuration failed")
	}

	// 加载apk channel 长度
	conf.ApkChanLen, err = beego.AppConfig.Int("apk_chan_len")
	if nil != err {
		conf.ApkChanLen = comm.WORKER_APK_CHANNEL_LEN_DEFAULT
	}

	// 加载redis配置
	conf.Redis, err = LoadRedisConf()
	if nil != err {
		return nil, err
	}

	// 加载internal sso 配置
	conf.InSso, err = LoadInternalSsoConf()
	if nil != err {
		return nil, err
	}

	// 加载Cdn splat配置
	conf.CdnSplat, err = LoadCdnSplatConf()
	if nil != err {
		return nil, err
	}

	logs.Info("conf load complete.")
	return conf, nil
}

/* 加载Mysql配置 */
func LoadMysqlConf() (conf *mysql.Conf, err error) {
	conf = new(mysql.Conf)

	// 加载Mysql参数
	conf.Conn = beego.AppConfig.String("mysql_conn")
	if 0 == len(conf.Conn) {
		return nil, errors.New("Get mysql connection failed!")
	}
	return conf, nil
}

/* 加载Redis配置 */
func LoadRedisConf() (conf *rds.Conf, err error) {
	conf = new(rds.Conf)
	// 加载Redis参数
	conf.Conn = beego.AppConfig.String("redis_conn")
	if 0 == len(conf.Conn) {
		return nil, errors.New("Get redis conn failed")
	}

	conf.Pwd = beego.AppConfig.String("redis_pwd")
	if 0 == len(conf.Pwd) {
		return nil, errors.New("Get redis password failed")
	}

	conf.MaxIdel, err = beego.AppConfig.Int("redis_max_idel")
	if nil != err {
		return nil, errors.New("Get redis max idel failed")
	}

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
