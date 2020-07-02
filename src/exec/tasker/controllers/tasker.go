package controllers

import (
	"github.com/astaxie/beego/toolbox"
)

//定时任务
func LaunchTasker() {

	//定时更新设备信息
	task := toolbox.NewTask("UpdateDevInfoTask", "0 */10 * * * *", UpdateDevInfoTask)
	toolbox.AddTask("UpdateDevInfoTask", task)

	//定时从redis里处理差分任务
	task = toolbox.NewTask("ApkPatchHandleTask", "*/1 * * * * *", ApkPatchHandleTask)
	toolbox.AddTask("ApkPatchHandleTask", task)

	//定时轮询mysql状态加入待处理队列
	task = toolbox.NewTask("GetPatchMysqlStatusTask", "*/30 * * * * *", GetPatchMysqlStatusTask)
	toolbox.AddTask("GetPatchMysqlStatusTask", task)

	toolbox.StartTask()
}
