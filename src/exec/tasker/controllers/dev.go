package controllers

import (
	"github.com/astaxie/beego/logs"
)

/******************************************************************************
 **函数名称: UpdateDevInfoTask
 **功    能: 定时更新设备信息（平台，机型，分组）
 **输入参数: NONE
 **输出参数: NONE
 **返    回: err 错误信息
 **实现描述:
 **作    者: # zhao.yang # 2019-06-01 15:44:28 #
 ******************************************************************************/
func UpdateDevInfoTask() error {

	logs.Info("Update dev info start.")

	//定时更新设备平台表
	err := UpdateDevPlat()
	if err != nil {
		logs.Error("Update dev plat failed ! errmsg:%s", err.Error())
		return err
	}

	//定时更新设备机型表
	err = UpdateDevModel()
	if err != nil {
		logs.Error("Update dev model failed ! errmsg:%s", err.Error())
		return err
	}

	//定时更新设备分组表
	err = UpdateDevGroup()
	if err != nil {
		logs.Error("Update dev group failed ! errmsg:%s", err.Error())
		return err
	}

	return nil
}
