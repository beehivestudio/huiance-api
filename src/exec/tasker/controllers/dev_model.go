package controllers

import (
	"github.com/astaxie/beego/logs"
	"upgrade-api/src/share/strategy"

	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/dms"
)

/******************************************************************************
 **函数名称: UpdateDevModel
 **功    能: 定时更新设备机型表
 **输入参数:NONE
 **输出参数: NONE
 **返    回: NONE
 **实现描述: 定时从http://internal.scloud.letv.com/device/manager/api/v1/tvgroup
 **          接口更新设备机型表
 **作    者: # zhao.yang # 2019-06-01 15:44:28 #
 ******************************************************************************/
func UpdateDevModel() (err error) {
	ctx := GetTaskerCntx()

	o := ctx.Model.Mysql.O

	//获取分组信息
	list, err := dms.GetTvGroupList()
	if err != nil {
		logs.Error("Get tv group list failed! errmsg:%s", err.Error())
		return err
	}

	models := map[string]string{}

	//遍历设备机型信息
	for _, v := range list {
		model := v.Model
		models[model] = v.Description
	}

	logs.Info("Get models len :%v", len(models))

	for model, description := range models {
		//查询平台id
		platId, err := upgrade.GetPlatIdByModel(o, model)
		if err != nil {
			logs.Error("Get plat id by model failed! errmsg:%s", err.Error())
			return err
		}

		//更新设备机型表信息
		err = upgrade.UpdateModel(o, model, description, platId)
		if err != nil {
			logs.Error("Update or insert dev model failed! errmsg:%s", err.Error())
			return err
		}
	}

	//缓存机型
	err = strategy.CacheModel(o, ctx.Model.Redis)
	if nil != err {
		logs.Error("Cache plat failed! errmsg:%s", err.Error())
		//Don't return err
	}

	return nil
}
