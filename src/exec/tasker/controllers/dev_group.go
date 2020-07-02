package controllers

import (
	"github.com/astaxie/beego/logs"
	"upgrade-api/src/share/strategy"

	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/dms"
)

/******************************************************************************
 **函数名称: UpdateDevGroup
 **功    能: 定时更新设备分组表
 **输入参数:NONE
 **输出参数: NONE
 **返    回: NONE
 **实现描述:定时从 http://internal.scloud.letv.com/device/manager/api/v1/tvgroup 接口更新设备分组表
 **作    者: # zhao.yang # 2019-06-01 15:44:28 #
 ******************************************************************************/
func UpdateDevGroup() (err error) {
	ctx := GetTaskerCntx()

	o := ctx.Model.Mysql.O

	list, err := dms.GetTvGroupList()
	if nil != err {
		logs.Error("Request tv group failed! errmsg:%s", err.Error())
		return err
	}

	logs.Info("Get group len msg :%v", len(list))

	for _, v := range list {

		//通过分组id获取设备机型id
		modelId, err := upgrade.GetModelByModelId(o, v.Model)
		if nil != err {
			logs.Error("Get model id failed! errmsg:%s", err.Error())
			return err
		}

		//更新或插入设备分组信息
		err = upgrade.UpdateOrInsertDevGroup(o, v.Id, v.Title, v.Description, modelId)
		if nil != err {
			logs.Error("Update or insert dev group failed! errmsg:%s", err.Error())
			return err
		}
	}

	//缓存分组
	err = strategy.CacheGroup(o, ctx.Model.Redis)
	if nil != err {
		logs.Error("Cache plat failed! errmsg:%s", err.Error())
		//Dont't return err
	}

	return nil
}
