package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go/extra"

	"upgrade-api/src/share/cdn"
	"upgrade-api/src/share/lib/log"
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/rds"
	"upgrade-api/src/share/lib/worker"
	"upgrade-api/src/share/quota"

	"upgrade-api/src/exec/backend/controllers/conf"
	"upgrade-api/src/exec/backend/models"
)

/* 全局对象 */
var BackendCntx = &BackendContext{}

func GetBackendCntx() *BackendContext {
	return BackendCntx
}

// BACKEND上下文信息
type BackendContext struct {
	Conf  *conf.Conf     // 配置信息
	Model *models.Models // 依赖组件资源
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
 **作    者: # Shuangpeng.guo # 2020-06-01 12:51:20 #
 ******************************************************************************/
func Init(cf *conf.Conf) (ctx *BackendContext, err error) {
	ctx = GetBackendCntx()

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
	ctx.Model.Mysql = mysql.GetMysqlPool() // 获取ORM

	/* 3.注册redis */
	ctx.Model.Redis = rds.CreatePool(cf.Redis.Conn, cf.Redis.Pwd, cf.Redis.MaxIdel) // 获取redis

	/* 4.jsoniter 全局设置 */
	// 开启容忍字符串数字互转，容忍空数组作为对象
	extra.RegisterFuzzyDecoders()
	extra.RegisterTimeAsInt64Codec(time.Second)

	/* 5.获取cdn splat句柄 */
	ctx.Model.CdnSplat = cdn.NewSplat(cf.CdnSplat.Uri, ctx.Model.Redis)
	ctx.Model.CdnSplat.RefreshSplatCache()

	/* 6.注册异步worker处理 */
	ctx.Model.Worker = worker.NewWorker(0, cf.ApkChanLen)

	/*  quota */
	quota, err := quota.Init(&quota.Component{
		Mysql: ctx.Model.Mysql,
		Redis: ctx.Model.Redis,
	})
	if err != nil {
		logs.Error("Quota load fail , %s", err.Error())
		return nil, err
	}
	//defer quota.ShutDown()
	ctx.Model.Quota = quota

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
 **作    者: # Shuangpeng.guo # 2020-06-01 13:21:39 #
 ******************************************************************************/
func Launch(ctx *BackendContext) (err error) {

	// 异步处理 非阻塞
	ctx.Model.Worker.Run()

	return nil
}
