package models

import (
	"fmt"
	"time"

	"github.com/astaxie/beego/logs"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/mysql"
)

/******************************************************************************
 **函数名称: CreateApp
 **功    能: 创建应用
 **输入参数:
 **     app: 应用信息
 **		devPlats: 关联平台信息
 **输出参数: NONE
 **返    回: 创建Id
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-01 14:14:35 #
 ******************************************************************************/
func (ctx *Models) CreateApp(
	app *upgrade.App, devPlatIds []int64) (int64, error) {

	o := mysql.GetMysqlPool(ctx.Mysql.AliasName).O

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return 0, err
	}

	id, err := upgrade.AddApp(o, app)
	if nil != err {
		logs.Warn("App create failed! param:%v errmsg:%s", app, err.Error())
		o.Rollback()
		return 0, err
	}

	if app.Enable == 0 || 0 == len(devPlatIds) {
		o.Commit()
		return id, nil
	}

	appPlatRels := make([]*upgrade.AppPlatRel, len(devPlatIds))

	for k, v := range devPlatIds {
		appPlatRels[k] = &upgrade.AppPlatRel{
			AppId:      id,
			DevPlatId:  v,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
	}

	_, err = o.InsertMulti(len(devPlatIds), appPlatRels)
	if nil != err {
		logs.Warn("Insert app plats failed! param:%v errmsg:%s",
			appPlatRels, err.Error())
		o.Rollback()
		return 0, err
	}

	o.Commit()

	return id, nil
}

/******************************************************************************
 **函数名称: UpdateApp
 **功    能: 修改应用
 **输入参数:
 **     app: 应用信息
 **		devPlatIds: 关联平台信息
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (ctx *Models) UpdateApp(
	app *upgrade.App, devPlatIds []int64) error {

	o := mysql.GetMysqlPool(ctx.Mysql.AliasName).O

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return err
	}

	// 更新app信息
	_, err = o.Update(app, "name", "dev_type_id", "enable",
		"has_dev_plat", "description", "update_time", "update_user")
	if nil != err {
		logs.Warn("App update failed! param:%v errmsg:%s", app, err.Error())
		o.Rollback()
		return err
	}

	// 删除app平台关联信息
	deleteSql := fmt.Sprintf(
		`DELETE
		 FROM %s
         WHERE app_id=?`, upgrade.UPGRADE_TAB_APP_PLAT_REL)
	_, err = o.Raw(deleteSql, app.Id).Exec()
	if nil != err {
		logs.Warn("App plat rel delete failed! errmsg:%s", err.Error())
		o.Rollback()
		return err
	}

	if app.HasDevPlat == 0 || 0 == len(devPlatIds) {
		o.Commit()
		return nil
	}

	appPlatRels := make([]*upgrade.AppPlatRel, len(devPlatIds))

	for k, v := range devPlatIds {
		appPlatRels[k] = &upgrade.AppPlatRel{
			AppId:      app.Id,
			DevPlatId:  v,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
	}

	_, err = o.InsertMulti(len(devPlatIds), appPlatRels)
	if nil != err {
		logs.Warn("Insert app plats failed! param: %v, errmsg:%s",
			appPlatRels, err.Error())
		o.Rollback()
		return err
	}

	o.Commit()

	return nil
}
