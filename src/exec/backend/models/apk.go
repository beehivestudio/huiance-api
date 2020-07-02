package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"

	"upgrade-api/src/share/cdn"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	libApk "upgrade-api/src/share/lib/apk"
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/utils"
)

type ApkCheckJob struct {
	ApkId int64
	Model *Models
}

/******************************************************************************
 **函数名称: Do
 **功    能: 异步处理方法
 **输入参数: None
 **输出参数:
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-17 13:42:44 #
 ******************************************************************************/
func (this *ApkCheckJob) Do(int) bool {

	var apkUpload upgrade.ApkUpload

	// 查询APK信息
	apk, err := upgrade.GetApkById(this.Model.Mysql.O, this.ApkId)
	if nil != err {
		logs.Error("<APK CHECK>: Query apk failed. id: %d, msg: %s", this.ApkId, err.Error())
		return false
	} else if upgrade.APK_STATUS_PROCESSING != apk.Status {
		logs.Info("<APK CHECK>: Apk worker status not PROCESSING, status: %d", apk.Status)
		return true
	}

	if err := this.DoHandler(apk, &apkUpload); nil != err {
		logs.Error("<APK CHECK>: Apk check failed. id: %d, msg: %s", this.ApkId, err.Error())
		// no return
	}

	// 删除关联patch包
	if err := upgrade.DeletePatch(this.Model.Mysql.O, apk.AppId, apk.VersionCode); nil != err {
		logs.Error("<APK CHECK>: Delete apk patch failed! param:%v msg: %s", apk, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		// no return
	}

	if _, err := this.Model.Mysql.O.Update(apk); nil != err {
		logs.Error("<APK CHECK>: Update apk status failed! param:%v msg: %s", apk, err.Error())
		return false
	}

	// 存储ApkUpload包
	if len(apkUpload.PackageName) > 0 {
		if _, _, err = this.Model.Mysql.O.ReadOrCreate(
			&apkUpload, "package_name", "version_code", "md5"); nil != err {
			logs.Error("<APK CHECK>: Read or create apk upload failed, param: %v, msg: %s", apkUpload, err.Error())
			return false
		}
	}

	return true
}

/******************************************************************************
 **函数名称: DoHandler
 **功    能: Apk 验证方法
 **输入参数:
 **输出参数:
 ** 	apk: apk信息
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-17 13:42:44 #
 ******************************************************************************/
func (this *ApkCheckJob) DoHandler(
	apk *upgrade.Apk, apkUpload *upgrade.ApkUpload) error {

	localPath, err := cdn.Download(apk.Url,
		fmt.Sprintf("%s%s_", comm.DOWNLOAD_APK_PATH, cdn.GenApkOutKey(apk.Id)), "")
	if nil != err {
		logs.Error("<APK CHECK>: Download apk failed. url: %s msg: %s", apk.Url, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	}

	// 解析APK包
	apkInfo, err := libApk.GetApkInfo(localPath, comm.APK_CN_CONF)
	if nil != err {
		logs.Error("<APK CHECK>: Analysis apk failed. apkId:%d, localPath: %s, msg: %s",
			apk.Id, localPath, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	}

	// 验证版本号
	cnt, err := this.Model.Mysql.O.QueryTable(upgrade.UPGRADE_TAB_APK).Filter(
		"app_id", apk.AppId).Filter("enable", comm.ENABLE).Filter("status",
		upgrade.APK_STATUS_NORMAL).Filter("version_code__gte", apkInfo.VersionCode).Count()
	if nil != err {
		logs.Error("<APK CHECK>: Query apk ver failed. apkId:%d, ver: %d, msg: %s",
			apk.Id, apkInfo.VersionCode, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	} else if cnt > 0 {
		logs.Error("<APK CHECK>: Version number conflict. apkId:%d, ver: %d",
			apk.Id, apkInfo.VersionCode)
		apk.Status = upgrade.APK_STATUS_VERSION_FAIL
		return err
	}

	// 验证包名
	app, err := upgrade.GetAppById(this.Model.Mysql.O, apk.AppId)
	if nil != err {
		logs.Error("<APK CHECK>: Query app failed. apkId:%d, ver: %d, msg: %s",
			apk.Id, apkInfo.VersionCode, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	}
	if app.PackageName != apkInfo.PackageName {
		logs.Error("<APK CHECK>: Pkg name not match. apkId:%d, apkPkg: %s, appPkg: %s",
			apk.Id, apkInfo.VersionCode, app.PackageName)
		apk.Status = upgrade.APK_STATUS_PKG_FAIL
		return err
	}

	// 验证签名
	if 0 == len(apkInfo.PublicKey) {
		logs.Error("<APK CHECK>: Apk public key cannot be empty apkId: %d", apk.Id)
		apk.Status = upgrade.APK_STATUS_SK_FAIL
		return err
	}

	if len(app.AppPublicKey) > 0 && app.AppPublicKey != apkInfo.PublicKey {
		logs.Error("<APK CHECK>: Apk public key not match apkId: %d, apkPublicKey: %s, appPublicKey: %s",
			apk.Id, apkInfo.PublicKey, app.AppPublicKey)
		apk.Status = upgrade.APK_STATUS_SK_FAIL
		return err
	}

	// 应用存储公钥不存在时创建
	if len(app.AppPublicKey) == 0 {
		app.AppPublicKey = apkInfo.PublicKey
		app.UpdateTime = time.Now()
		this.Model.Mysql.O.Update(app, "app_public_key", "update_time")
	}

	apk.VersionName = apkInfo.VersionName
	apk.VersionCode = apkInfo.VersionCode
	apk.Md5 = apkInfo.Md5
	apk.Size = apkInfo.Size
	apk.Status = upgrade.APK_STATUS_NORMAL
	apk.UpdateTime = time.Now()

	// 存储APK包
	apkUpload.FileName = utils.GetFileNameFromUrl(localPath)
	apkUpload.Md5 = apk.Md5
	apkUpload.Size = apk.Size
	apkUpload.PackageName = apkInfo.PackageName
	apkUpload.Label = apkInfo.Label
	apkUpload.VersionCode = apk.VersionCode
	apkUpload.VersionName = apk.VersionName
	apkUpload.PublicKey = apkInfo.PublicKey
	apkUpload.LocalPath = localPath
	apkUpload.CdnUrl = apk.Url
	apkUpload.Status = upgrade.APKFILEUPLOAD_SUCCESS
	apkUpload.CreateTime = time.Now()
	apkUpload.UpdateTime = time.Now()
	apkUpload.CreateUser = comm.SYSYTEM_USER
	apkUpload.UpdateUser = comm.SYSYTEM_USER

	return nil
}

/******************************************************************************
 **函数名称: CreateApk
 **功    能: 创建APK
 **输入参数:
 **     apk: 应用信息
 **输出参数: NONE
 **返    回: 创建Id
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-03 14:14:35 #
 ******************************************************************************/
func (ctx *Models) CreateApk(apk *upgrade.Apk) (int64, error) {

	o := mysql.GetMysqlPool(ctx.Mysql.AliasName).O

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return 0, err
	}

	id, err := upgrade.AddApk(o, apk)
	if nil != err {
		logs.Warn("Apk create failed! param:%v errmsg:%s", apk, err.Error())
		o.Rollback()
		return 0, err
	}

	apkJob := &ApkCheckJob{
		ApkId: id,
		Model: ctx,
	}

	if len(*ctx.Worker.GetChannel()) >= cap(*ctx.Worker.GetChannel()) {
		logs.Warn("apk worker channel full, len: %d", len(*ctx.Worker.GetChannel()))
		o.Rollback()
		return 0, errors.New("apk worker channel full")
	}
	*ctx.Worker.GetChannel() <- apkJob

	o.Commit()

	return id, nil
}

/******************************************************************************
 **函数名称: UpdateApk
 **功    能: 修改APK
 **输入参数:
 **     apk: 应用信息
 **     doCheck: 是否走异步check流程
 **输出参数: NONE
 **返    回: 创建Id
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-03 14:14:35 #
 ******************************************************************************/
func (ctx *Models) UpdateApk(apk *upgrade.Apk, doCheck bool) error {

	o := mysql.GetMysqlPool(ctx.Mysql.AliasName).O

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return err
	}

	_, err = o.Update(apk)
	if nil != err {
		logs.Warn("Apk Update failed! param:%v errmsg:%s", apk, err.Error())
		o.Rollback()
		return err
	}

	if doCheck {
		apkJob := &ApkCheckJob{
			ApkId: apk.Id,
			Model: ctx,
		}

		if len(*ctx.Worker.GetChannel()) >= cap(*ctx.Worker.GetChannel()) {
			logs.Warn("Apk worker channel full! len: %d", len(*ctx.Worker.GetChannel()))
			o.Rollback()
			return errors.New("Apk worker channel full")
		}
		*ctx.Worker.GetChannel() <- apkJob
	}

	o.Commit()

	return nil
}

type Apk struct {
	Id             int64
	Enable         int
	Bcode          string
	AppId          int64
	AppName        string
	BusinessId     int64
	PackageName    string
	VersionCode    int64
	VersionName    string
	Url            string
	Md5            string
	Size           int64
	EuiLowVersion  string
	EuiHighVersion string
	Status         int
	Description    string
	Memo           string
	CallbackUrl    string
	CreateTime     time.Time
	CreateUser     string
	UpdateTime     time.Time
	UpdateUser     string
}

/******************************************************************************
 **函数名称: GetApkList
 **功    能: 获取ApkList
 **输入参数:
 **     fields: 过滤条件
 **     values: 过滤值
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-03 10:14:35 #
 ******************************************************************************/
func (ctx *Models) GetApkList(fields []string,
	values []interface{}, page, pageSize int) (int64, []*Apk, error) {

	var apks []*Apk

	sql := fmt.Sprintf(`SELECT
		apk.id,
		apk.bcode,
		apk.app_id,
		apk.enable,
		apk.version_code,
		apk.version_name,
		apk.url,
		apk.md5,
		apk.size,
		apk.eui_low_version,
		apk.eui_high_version,
		apk.status,
		apk.description,
		apk.memo,
		apk.callback,
		apk.create_time,
		apk.create_user,
		apk.update_time,
		apk.update_user,
		app.name as app_name,
		app.package_name,
		app.business_id
	FROM
		%s AS apk
	LEFT JOIN %s AS app
		ON apk.app_id = app.id
	`, upgrade.UPGRADE_TAB_APK, upgrade.UPGRADE_TAB_APP)

	sqlCount := fmt.Sprintf(`SELECT
		COUNT(apk.id) as total
	FROM
		%s AS apk
	LEFT JOIN %s AS app
		ON apk.app_id = app.id
	`, upgrade.UPGRADE_TAB_APK, upgrade.UPGRADE_TAB_APP)

	// 如果where条件不为空追加 where 条件
	if 0 < len(fields) {
		where := strings.Join(fields, " AND ")
		sql = fmt.Sprintf("%s WHERE %s ", sql, where)
		sqlCount = fmt.Sprintf("%s WHERE %s ", sqlCount, where)
	}

	res := orm.ParamsList{}
	if _, err := ctx.Mysql.O.Raw(sqlCount, values).ValuesFlat(&res); nil != err {
		return 0, nil, err
	} else if len(res) == 0 {
		res = append(res, 0)
	}

	total := int64(0)
	totalStr, ok := res[0].(string)
	if ok {
		total, _ = strconv.ParseInt(totalStr, 10, 64)
	}

	// 添加排序&分页条件
	sql = fmt.Sprintf("%s ORDER BY app.id DESC LIMIT ? OFFSET ? ", sql)
	values = append(values, pageSize, (page-1)*pageSize)

	if _, err := ctx.Mysql.O.Raw(sql, values...).QueryRows(&apks); nil != err {
		return 0, nil, err
	}

	return total, apks, nil
}

/******************************************************************************
 **函数名称: GetApkById
 **功    能: 获取Apk
 **输入参数:
 **     id: apkId
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-03 10:14:35 #
 ******************************************************************************/
func (ctx *Models) GetApkById(id int64) (*Apk, error) {

	apk := new(Apk)

	sql := fmt.Sprintf(`SELECT
		apk.id,
		apk.bcode,
		apk.app_id,
		apk.enable,
		apk.version_code,
		apk.version_name,
		apk.url,
		apk.md5,
		apk.size,
		apk.eui_low_version,
		apk.eui_high_version,
		apk.status,
		apk.description,
		apk.memo,
		apk.callback,
		apk.create_time,
		apk.create_user,
		apk.update_time,
		apk.update_user,
		app.name as app_name,
		app.package_name,
		app.business_id
	FROM
		%s AS apk
	LEFT JOIN %s AS app
		ON apk.app_id = app.id
	WHERE
		apk.id = ?
	`, upgrade.UPGRADE_TAB_APK, upgrade.UPGRADE_TAB_APP)

	if err := ctx.Mysql.O.Raw(sql, id).QueryRow(apk); nil != err {
		return nil, err
	}

	return apk, nil
}
