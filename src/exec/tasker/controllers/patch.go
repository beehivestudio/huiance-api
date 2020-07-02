package controllers

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/patch"
	"upgrade-api/src/share/lib/rds"
	"upgrade-api/src/share/strategy"
)

/******************************************************************************
 **函数名称: ApkPatchHandleTask
 **功    能: APK差分处理任务
 **输入参数: NONE
 **输出参数: NONE
 **返    回: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-1 13:15:42 #
 ******************************************************************************/
func ApkPatchHandleTask() (err error) {
	ctx := GetTaskerCntx()

	logs.Info("Begin apk patch handler")

	mc := ctx.Model.Mysql
	rc := ctx.Model.Redis.Get()
	defer rc.Close()

	for {
		/* 1.获取差分任务 */
		patchData := strategy.PatchData{}

		patchDateStr, err := rds.RedisRPOP(rc, comm.RDS_KEY_APK_PATCH_TASK_LIST)
		if nil != err {
			logs.Error("Pop patch task failed! errmsg:%s", err.Error())
			if redis.ErrNil == err {
				break // End of task
			}
		}

		err = jsoniter.UnmarshalFromString(patchDateStr, &patchData)
		if nil != err {
			logs.Error("Unmarshal data failed! data:%v errmsg:%s",
				patchDateStr, err.Error())
			return err
		}

		/* 2.进行差分处理 */
		err = patch.ApkPatchHandle(mc.O, patchData)
		if nil != err {
			logs.Error("Apk patch handle failed! data:%v errmsg:%s",
				patchData, err.Error())
			return err
		}
	}

	logs.Info("End apk patch handler.")

	return nil
}

/******************************************************************************
 **函数名称: GetPatchMysqlStatusTask
 **功    能: 从mysql中获取待处理状态的差分任务
 **输入参数: NONE
 **输出参数: NONE
 **返    回: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Linlin.guo # 2020-6-1 13:15:42 #
 ******************************************************************************/
func GetPatchMysqlStatusTask() (err error) {
	ctx := GetTaskerCntx()

	mc := ctx.Model.Mysql
	rc := ctx.Model.Redis.Get()
	defer rc.Close()

	apkPatchList, err := upgrade.GetApkPatchByStatus(mc.O, upgrade.APK_PATCH_STATUS_WAIT)
	if nil != err {
		logs.Error("Get apk patch by status failed! status:%d errmsg:%s",
			upgrade.APK_PATCH_STATUS_WAIT, err.Error())
		return err
	}

	for _, apkPatch := range apkPatchList {
		patchData := strategy.PatchData{
			ApkPatchId: apkPatch.Id,
		}

		err = rds.RedisLPUSHJsonData(rc, comm.RDS_KEY_APK_PATCH_TASK_LIST, &patchData)
		if nil != err {
			logs.Error("Push patch task failed! patchData:%v errmsg:%s",
				patchData, err.Error())
			continue
		}
	}

	return nil
}
