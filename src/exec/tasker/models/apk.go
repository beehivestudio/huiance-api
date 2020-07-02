package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/mysql"
)

/******************************************************************************
 **函数名称: GetExpireApk
 **功    能: 扫描APK表
 **输入参数:
 **    len: 单次请求长度
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-09 14:14:35 #
 ******************************************************************************/
func (mod *Models) GetExpireApk(len int) ([]*upgrade.Apk, error) {

	list := new([]*upgrade.Apk)

	m, _ := time.ParseDuration("-" + comm.DEFAULT_APK_EXPIRE)
	tmStr := time.Now().Add(m).Format("2006-01-02 15:04:05")

	cond := orm.NewCondition()

	cond = cond.And("status", upgrade.APK_STATUS_PROCESSING)
	cond = cond.And("update_time__lte", tmStr)

	qs := mod.Mysql.O.QueryTable(upgrade.UPGRADE_TAB_APK).SetCond(cond)

	_, err := qs.Limit(len, 0).All(list)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get apk timeout list failed! errmsg:%s", err.Error())
		return nil, err
	}

	return *list, nil

}

/******************************************************************************
 **函数名称: ApkExpireHandler
 **功    能: 超时处理
 **输入参数:
 **    apk: APK参数
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-09 14:14:35 #
 ******************************************************************************/
func (this *Models) ApkExpireHandler(apk *upgrade.Apk) error {
	m, _ := time.ParseDuration("-" + comm.DEFAULT_APK_EXPIRE)
	expTm := time.Now().Add(m)

	logs.Debug("Apk expire handler! apk:%v", apk)

	// 开始事物
	o := mysql.GetMysqlPool(this.Mysql.AliasName).O

	err := o.Begin()
	if nil != err {
		logs.Error("Begin tx failed! apk:%v errmsg:%s", apk, err.Error())
		return err
	}

	// 获取APK信息
	apk = &upgrade.Apk{
		Id: apk.Id,
	}

	err = o.Read(apk)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get apk failed! apkId:%d errmsg:%s",
			apk.Id, err.Error())
		o.Rollback()
		return err
	} else if orm.ErrNoRows == err {
		logs.Error("Apk isn't exists! apkId: %d", apk.Id)
		o.Rollback()
		return errors.New("apk not exist")
	} else if expTm.Unix() < apk.UpdateTime.Unix() ||
		upgrade.APK_STATUS_PROCESSING != apk.Status {
		o.Rollback()
		logs.Error("Apk task not expire! apkId:%d", apk.Id)
		return errors.New("Apk task not expire! ")
	}

	// 更新APK状态
	apk.Status = upgrade.APK_STATUS_TIMEOUT
	apk.UpdateTime = time.Now()
	// apk.UpdateUser = comm.SYSYTEM_USER

	_, err = o.Update(apk)
	if nil != err {
		o.Rollback()
		logs.Error("Update apk failed! apk:%v errmsg:%s", apk, err.Error())
		return err
	}

	o.Commit()

	// TODO 异步通知业务方

	return nil

}
