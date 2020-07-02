package log

import (
	"github.com/astaxie/beego/logs"
)

/******************************************************************************
 **函数名称: SetLogger
 **功    能: 初始化日志信息
 **输入参数:
 **     conf: 日志配置 如: {"filename":"/letv/app.log", "maxsize":1000, "level":7}
 **输出参数: NONE
 **返    回: NONE
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2020-06-29 15:20:26 #
 ******************************************************************************/
func SetLogger(conf string) {
	// 设置日志配置信息
	logs.SetLogger("file", conf)

	// 显示文件名和行号
	logs.EnableFuncCallDepth(true)
}
