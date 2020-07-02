package controllers

import (
	"time"

	"github.com/astaxie/beego/logs"
)

/******************************************************************************
 **函数名称: ApkTasker
 **功    能: 扫描APK表
 **输入参数:
 **输出参数: NONE
 **返    回: 创建Id
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-09 14:14:35 #
 ******************************************************************************/
func ApkTasker() {
	ctx := GetTaskerCntx()

	for {
		for {
			// 获取超时通知
			list, err := ctx.Model.GetExpireApk(10)
			if nil != err {
				logs.Error("Get apk time out list failed! errmsg:%s", err.Error())
				break
			} else if 0 == len(list) {
				break
			}

			// 处理超时超时任务
			for _, item := range list {
				ctx.Model.ApkExpireHandler(item)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
