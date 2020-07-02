package controllers

import (
	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/dms"
	"upgrade-api/src/share/strategy"
)

/******************************************************************************
 **函数名称: UpdateDevPlat
 **功    能: 定时更新设备平台表
 **输入参数:NONE
 **输出参数: NONE
 **返    回: NONE
 **实现描述:定时从 http://internal.scloud.letv.com/basedata/get?code=MODEL_Defines 接口更新设设备平台表
 **作    者: # zhao.yang # 2019-06-01 15:44:28 #
 ******************************************************************************/
func UpdateDevPlat() (err error) {
	ctx := GetTaskerCntx()

	o := ctx.Model.Mysql.O

	data, err := dms.GetTvModelList()
	if nil != err {
		logs.Error("Get tv model list failed! errmsg:%s", err.Error())
		return err
	}

	plats := map[string]string{}

	//遍历设备平台接口中设备平台信息
	for _, v := range data {
		platform := v.Platform
		description := v.Description

		//设备平台名称去重处理
		plats[platform] = description
	}

	devTypeId := comm.DEFAULT_DEV_TYPE_ID

	logs.Info("Get plats len msg :%v", len(plats))

	//更新数据库
	for platform, description := range plats {
		err := upgrade.UpdatePlat(o, platform, description, int64(devTypeId))
		if nil != err {
			logs.Error("Update dev plat failed! platform:%s errmsg:%s",
				platform, err.Error())
			return err
		}
	}

	//缓存平台
	err = strategy.CachePlat(o, ctx.Model.Redis)
	if nil != err {
		logs.Error("Cache plat failed! errmsg:%s", err.Error())
		//Don't return err
	}
	return nil
}
