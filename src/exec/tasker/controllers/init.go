package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/lib/log"
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/rds"

	"upgrade-api/src/exec/tasker/controllers/conf"
	"upgrade-api/src/exec/tasker/models"
)

/* 全局对象 */
type TaskerContext struct {
	Conf  *conf.Conf     // 配置信息
	Model *models.Models // 依赖的组件对象
}

var TaskerCntx = &TaskerContext{}

func GetTaskerCntx() *TaskerContext {
	return TaskerCntx
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
 **作    者: # Linlin.guo # 2019-06-1 15:09:51 #
 ******************************************************************************/
func Init(cf *conf.Conf) (ctx *TaskerContext, err error) {
	ctx = GetTaskerCntx()
	ctx.Conf = cf

	/* 1.初始化日志 */
	log.SetLogger(ctx.Conf.Log)

	/* 2.注册Mysql */
	if beego.BConfig.RunMode == beego.DEV {
		orm.DebugLog = orm.NewLog(logs.GetBeeLogger())
		orm.Debug = true
	}
	mysql.RegisterDb(cf.Mysql.Conn)

	ctx.Model = new(models.Models)

	//table.RegisterModel()            // 注册定义的Model
	ctx.Model.Mysql = mysql.GetMysqlPool() // 获取ORM

	/* 3.注册redis */
	ctx.Model.Redis = rds.CreatePool(cf.Redis.Conn, cf.Redis.Pwd, cf.Redis.MaxIdel) // 获取redis
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
 **作    者: # Linlin.guo # 2019-06-1 15:09:51 #
 ******************************************************************************/
func Launch(ctx *TaskerContext) (err error) {

	/* > 执行定时任务 */
	go LaunchTasker()

	// 启动APK升级超时检查
	go ApkTasker()
	return nil
}
