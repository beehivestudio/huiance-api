package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/lib/log"
	"upgrade-api/src/share/lib/mysql"

	"upgrade-api/src/exec/proxy/controllers/conf"
	"upgrade-api/src/exec/proxy/models"
)

/* 全局上下文对象 */
type ProxyContext struct {
	Conf  *conf.ProxyConf // 配置信息
	Model *models.Models  // 依赖的组件资源
}

var ProxyCntx = &ProxyContext{}

func GetProxyCntx() *ProxyContext {
	return ProxyCntx
}

/******************************************************************************
 **函数名称: Init
 **功    能: 初始化处理
 **输入参数:
 **     cf: 配置信息
 **输出参数: NONE
 **返    回: 错误描述
 **实现描述:
 **注意事项:
 **作    者: # pengjia # 2019-04-15 15:57:51 #
 ******************************************************************************/
func Init(cf *conf.ProxyConf) (ctx *ProxyContext, err error) {
	ctx = GetProxyCntx()

	ctx.Conf = cf

	/* 1.初始化日志 */
	log.SetLogger(ctx.Conf.Log)

	/* 2.注册Mysql */
	if beego.BConfig.RunMode == beego.DEV {
		orm.DebugLog = orm.NewLog(logs.GetBeeLogger())
		orm.Debug = true
	}

	ctx.Model = new(models.Models)

	mysql.RegisterDb(cf.Mysql.Conn)
	//table.RegisterModel()            // 注册定义的Model
	ctx.Model.Mysql = mysql.GetMysqlPool() // 获取ORM
	return ctx, nil
}

/******************************************************************************
 **函数名称: Launch
 **功    能: 启动程序
 **输入参数: NONE
 **输出参数: NONE
 **返    回: 错误描述
 **实现描述: 启动后台工作协程
 **注意事项:
 **作    者: # pengjia # 2019-04-15 15:58:17 #
 ******************************************************************************/
func Launch(ctx *ProxyContext) (err error) {
	return nil
}
