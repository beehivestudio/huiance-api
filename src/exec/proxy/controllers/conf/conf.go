package conf

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/lib/mysql"
)

/* 应用升级配置文件 */
type ProxyConf struct {
	Log           string         // 日志配置
	Mysql         mysql.Conf     // mysql 配置
	AkpUpgrade    *AkpUpgrade    // apk升级url
	BusinessInfor *BusinessInfor // 业务配置信息
	BusinessKey   string         // 业务key
	Sign          *Sign          //与客户端签名验证
}

/* 升级接口url */
type AkpUpgrade struct {
	ApkUpgrade    string
	ApkUpAradeAll string
	ApkUninstall  string
}

/* 业务信息 */
type BusinessInfor struct {
	BusinessId int //业务id
	DevTypeId  int //设备类型id
}

/* 业务信息 */
type Sign struct {
	AkSign string
	SkSign string
}

func Load() (conf *ProxyConf, err error) {
	conf = new(ProxyConf)

	// 加载MySQL配置
	conf.Mysql.Conn = beego.AppConfig.String("mysql_conn")
	if 0 == len(conf.Mysql.Conn) {
		return nil, errors.New("Get mysql connection failed!")
	}

	// 加载日志配置
	conf.Log = beego.AppConfig.String("log_conf")
	if 0 == len(conf.Log) {
		return nil, errors.New("Get log configuration failed!")
	}

	// 加载业务key
	conf.BusinessKey = beego.AppConfig.String("business_key")
	if 0 == len(conf.BusinessKey) {
		return nil, errors.New("Get business_key configuration failed!")
	}

	// 加载akp升级接口配置
	conf.AkpUpgrade, err = LoadApkApiConf()
	if nil != err {
		return nil, err
	}

	// 加载业务信息
	conf.BusinessInfor, err = LoadBusinessInfor()
	if nil != err {
		return nil, err
	}

	//加载签名信息
	conf.Sign, err = LoadSign()
	if nil != err {
		return nil, err
	}

	logs.Info("Load configuration completed!")

	return conf, nil
}

func LoadBusinessInfor() (conf *BusinessInfor, err error) {

	conf = new(BusinessInfor)

	conf.BusinessId, _ = beego.AppConfig.Int("business_id")
	if 0 == conf.BusinessId {
		return nil, errors.New("Get business_id failed")
	}

	conf.DevTypeId, err = beego.AppConfig.Int("dev_type_id")
	if 0 == conf.DevTypeId {
		return nil, errors.New("Get dev_type_id failed")
	}

	return conf, nil
}

/* 加载apk升级接口配置 */
func LoadApkApiConf() (conf *AkpUpgrade, err error) {

	conf = new(AkpUpgrade)

	//加载请求升级接口地址
	conf.ApkUpgrade = beego.AppConfig.String("apk_upgrade")
	if 0 == len(conf.ApkUpgrade) {
		return nil, errors.New("Get apk_upgrade failed")
	}

	//加载请求全部升级接口地址
	conf.ApkUpAradeAll = beego.AppConfig.String("apk_upgrade_all")
	if 0 == len(conf.ApkUpAradeAll) {
		return nil, errors.New("Get apk_upgrade_all failed")
	}

	//加载卸载接口地址
	conf.ApkUninstall = beego.AppConfig.String("apk_uninstall")
	if 0 == len(conf.ApkUpAradeAll) {
		return nil, errors.New("Get apk_uninstall failed")
	}

	return conf, nil
}

/* 加载签名配置 */
func LoadSign() (conf *Sign, err error) {

	conf = new(Sign)

	//加载ak
	conf.AkSign = beego.AppConfig.String("my_ak")
	if 0 == len(conf.AkSign) {
		//return nil, errors.New("Get my_ak failed")  //todo 处理去掉注释
	}

	conf.SkSign = beego.AppConfig.String("my_sk")
	if 0 == len(conf.SkSign) {
		//return nil, errors.New("Get my_sk failed") //todo 处理去掉注释
	}

	return conf, nil
}
