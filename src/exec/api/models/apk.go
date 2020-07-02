package models

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/cdn"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database/upgrade"
	libApk "upgrade-api/src/share/lib/apk"
	"upgrade-api/src/share/lib/mysql"
	"upgrade-api/src/share/lib/utils"
)

/* apk下载任务*/
type ApkCheckJob struct {
	ApkId       int64
	Model       *Models
	CallbackUrl string
}

/******************************************************************************
 **函数名称: Do
 **功    能: 执行任务
 **输入参数: NONE
 **输出参数:
 **返    回:
 **     bool: 执行成功
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (job *ApkCheckJob) Do(int) bool {

	var apkUpload upgrade.ApkUpload

	apk, err := upgrade.GetApkById(job.Model.Mysql.O, job.ApkId)
	if nil != err {
		logs.Error("<APK CHECK>: Query apk failed! id:%d errmsg:%s", job.ApkId, err.Error())
		return false
	} else if upgrade.APK_STATUS_PROCESSING != apk.Status {
		logs.Info("<APK CHECK>: Apk worker status not processing! status:%d", apk.Status)
		return true
	}

	if err := job.DoHandler(apk, &apkUpload); nil != err {
		logs.Error("<APK CHECK>: Apk check failed! id:%d errmsg:%s", job.ApkId, err.Error())
		// Don't return
	}

	// 删除关联patch包
	if err := upgrade.DeletePatch(job.Model.Mysql.O, apk.AppId, apk.VersionCode); nil != err {
		logs.Error("<APK CHECK>: Delete apk patch failed! param:%v msg: %s", apk, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		// no return
	}

	if _, err := job.Model.Mysql.O.Update(apk); nil != err {
		logs.Error("<APK CHECK>: Update apk status failed! param:%v errmsg:%s", apk, err.Error())
		return false
	}

	if err := job.Callback(); nil != err {
		logs.Error("<APK CHECK>: Call back failed! apkId:%v errmsg:%s ", apk.Id, err.Error())
		return true
	}

	// 存储ApkUpload包
	if len(apkUpload.PackageName) > 0 {
		if _, _, err := job.Model.Mysql.O.ReadOrCreate(
			&apkUpload, "package_name", "version_code", "md5"); nil != err {
			logs.Error("<APK CHECK>: Read or create apk upload failed, param: %v, msg: %s", apkUpload, err.Error())
			return false
		}
	}

	return true
}

/******************************************************************************
 **函数名称: DoHandler
 **功    能: 下载apk, 校验参数
 **输入参数:
 ** 	apk: 应用信息
 **输出参数:
 **     error: 异常返回
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (job *ApkCheckJob) DoHandler(
	apk *upgrade.Apk, apkUpload *upgrade.ApkUpload) error {

	/*  */
	app, err := upgrade.GetAppById(job.Model.Mysql.O, apk.AppId)
	if nil != err {
		logs.Error("<APK CHECK>: Query app failed! apkId:%d  errmsg:%s",
			apk.Id, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	}
	/* 下载 */
	localPath, err := cdn.Download(apk.Url,
		fmt.Sprintf("%s%s_", comm.DOWNLOAD_APK_PATH, cdn.GenApkOutKey(apk.Id)), "")
	if nil != err {
		logs.Error("<APK CHECK>: Download apk failed! url:%s errmsg:%s",
			apk.Url, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	}

	/* 解析APK包 */
	apkInfo, err := libApk.GetApkInfo(localPath, comm.APK_CN_CONF)
	if nil != err {
		logs.Error("<APK CHECK>: Analysis apk failed. apkId:%d, localPath: %s, msg: %s",
			apk.Id, localPath, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	}

	/* 验证版本号 */
	cnt, err := job.Model.Mysql.O.QueryTable(upgrade.UPGRADE_TAB_APK).Filter(
		"app_id", apk.AppId).Filter("enable", comm.ENABLE).Filter("status",
		upgrade.APK_STATUS_NORMAL).Filter("version_code__gte", apkInfo.VersionCode).Count()
	if nil != err {
		logs.Error("<APK CHECK>: Query apk version code failed! apkId:%d versionCode:%d errmsg:%s",
			apk.Id, apkInfo.VersionCode, err.Error())
		apk.Status = upgrade.APK_STATUS_OTHER_ERR
		return err
	} else if cnt > 0 {
		logs.Error("<APK CHECK>: Version code conflict! apkId:%d, versionCode:%d",
			apk.Id, apkInfo.VersionCode)
		apk.Status = upgrade.APK_STATUS_VERSION_FAIL
		return err
	}

	/* 验证包名 */
	if app.PackageName != apkInfo.PackageName {
		logs.Error("<APK CHECK>: Pkg name not match. apkId:%d apkPkgName:%s appPkgName:%s",
			apk.Id, apkInfo.PackageName, app.PackageName)
		apk.Status = upgrade.APK_STATUS_PKG_FAIL
		return err
	}

	if 0 == len(apkInfo.PublicKey) { /* 验证签名 */
		logs.Error("<APK CHECK>: Apk public key is invalid! apkId:%d", apk.Id)
		apk.Status = upgrade.APK_STATUS_SK_FAIL
		return err
	}

	if len(app.AppPublicKey) > 0 && app.AppPublicKey != apkInfo.PublicKey {
		logs.Error("<APK CHECK>: Apk public key not match! apkId:%d apkPublicKey:%s appPublicKey:%s",
			apk.Id, apkInfo.PublicKey, app.AppPublicKey)
		apk.Status = upgrade.APK_STATUS_SK_FAIL
		return err
	}

	if 0 == len(app.AppPublicKey) {
		app.AppPublicKey = apkInfo.PublicKey
		app.UpdateTime = time.Now()
		job.Model.Mysql.O.Update(app, "app_public_key", "update_time")
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
 **函数名称: Callback
 **功    能: 下载完成回调
 **输入参数: NONE
 **输出参数:
 **     error: 异常返回
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (job *ApkCheckJob) Callback() error {
	/* 检查url是否有效 */
	if job.CallbackUrl == "" {
		logs.Info("Apk dont't need call back ")
		return nil
	}
	/* 查询apk */
	apk, err := upgrade.GetApkById(job.Model.Mysql.O, job.ApkId)
	if nil != err {
		logs.Error("Apk call back  for query apk failed, apkId:%v,errmsg:%s",
			apk.Id, err.Error())
		return err
	}

	/*查询app */
	app, err := upgrade.GetAppById(job.Model.Mysql.O, apk.AppId)
	if nil != err {
		logs.Error("Apk call back for query app failed, appId:%v,errmsg:%s",
			app.Id, err.Error())
		return err
	}

	/* 请求体 */
	body := map[string]interface{}{
		"bcode":            apk.Bcode,          // 业务方APK唯一码（注：该字段由业务方自行定义，但需保证该值的唯一性）
		"app_id":           apk.AppId,          // 增量升级平台APP ID
		"apk_id":           apk.Id,             // 增量升级平台APK ID
		"package_name":     app.PackageName,    // 应用包名
		"version_code":     apk.VersionCode,    // 版本号
		"version_name":     apk.VersionName,    // 版本名
		"url":              apk.Url,            // APK下载路径
		"md5":              apk.Md5,            // APK包MD5值
		"size":             apk.Size,           // APK包大小（字节）
		"eui_low_version":  apk.EuiLowVersion,  // 依赖的EUI最低版本
		"eui_high_version": apk.EuiHighVersion, // 依赖的EUI最高版本
		"status":           apk.Status,         // 当前状态
		"description":      apk.Description,    // 描述信息
	}

	jsonBody, err := jsoniter.Marshal(body)
	if nil != err {
		logs.Error("Call back failed! errmsg:%v", err.Error())
		return err
	}

	req, err := http.NewRequest(http.MethodPut, job.CallbackUrl, bytes.NewBuffer(jsonBody))
	if nil != err {
		logs.Error("Call back failed! errmsg:%v", err.Error())
		return err
	}
	/* 签名 */
	url, err := url.Parse(job.CallbackUrl)
	if nil != err {
		logs.Error("call back url invalid, url: %s", job.CallbackUrl)
		return err
	}

	signBody := map[string]interface{}{}
	signBody["body"] = string(jsonBody)
	signBody["url"] = url.RequestURI()

	business, err := upgrade.GetBusinessById(job.Model.Mysql.O, app.BusinessId)
	if nil != err {
		logs.Error("Get business by id failed! businessId:%d errmsg:%s",
			app.BusinessId, err.Error())
		return err
	}

	sign, err := comm.SignMapByMd5(signBody, business.Key, true)
	if nil != err {
		logs.Error("Sign md5 failed! body:%s errmsg:%s", jsonBody, err.Error())
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("timestamp", fmt.Sprintf("%v", time.Now().Unix()))
	req.Header.Set("sign", sign)
	req.Header.Set("business-id", business.Key)

	client := new(http.Client)
	resp, err := client.Do(req)
	if nil != err {
		logs.Error("Call back failed! errmsg:%v", err.Error())
		return err
	}
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		logs.Error("Call back failed! errmsg:%v", err.Error())
		return err
	}

	if resp.Status != "200 OK" {
		logs.Error("Call back failed!")
		return errors.New("Request status not 200")
	}

	logs.Info("Call back success! %s", string(bodyByte))

	return nil
}

/******************************************************************************
 **函数名称: CreateApk
 **功    能: 新建APK
 **输入参数:
 ** 	apk: 应用信息
 ** 	callbackUrl: 回调路径
 **输出参数:
 ** 	int64: id
 **     error: 异常返回
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (mod *Models) CreateApk(apk *upgrade.Apk, callbackUrl string) (int64, error) {
	o := orm.NewOrm()
	o.Using(mysql.DEFAULT_ALIAS_NAME)

	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return 0, err
	}

	id, err := upgrade.AddApk(o, apk)
	if nil != err {
		o.Rollback()
		logs.Error("Apk create failed! param:%v errmsg:%s", apk, err.Error())
		return 0, err
	}

	apkJob := &ApkCheckJob{
		ApkId:       id,
		Model:       mod,
		CallbackUrl: callbackUrl,
	}

	/* 执行apk任务 */
	if len(*mod.Worker.GetChannel()) >= cap(*mod.Worker.GetChannel()) {
		o.Rollback()
		logs.Error("Apk worker channel full! len:%d", len(*mod.Worker.GetChannel()))
		return 0, errors.New("apk worker channel full")
	}
	logs.Info("CreateApk -> job")
	*mod.Worker.GetChannel() <- apkJob

	o.Commit()

	return id, nil
}

/******************************************************************************
 **函数名称: UpdateApk
 **功    能: 修改APK
 **输入参数:
 **     apk: 应用信息
 **     callbackUrl: 回调地址
 *8     doCheck: 是否走异步处理
 **输出参数: NONE
 **返    回: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-03 14:14:35 #
 ******************************************************************************/
func (mod *Models) UpdateApk(
	apk *upgrade.Apk, callbackUrl string, doCheck bool) error {

	o := orm.NewOrm()
	o.Using(mysql.DEFAULT_ALIAS_NAME)

	// 开始事物
	err := o.Begin()
	if nil != err {
		logs.Error("Begin transaction failed! errmsg:%s", err.Error())
		return err
	}

	err = upgrade.UpdateApkById(o, apk)
	if nil != err {
		o.Rollback()
		logs.Error("Apk Update failed! param:%v errmsg:%s", apk, err.Error())
		return err
	}

	if doCheck {
		apkJob := &ApkCheckJob{
			ApkId:       apk.Id,
			Model:       mod,
			CallbackUrl: callbackUrl,
		}
		/* 执行apk任务 */
		if len(*mod.Worker.GetChannel()) >= cap(*mod.Worker.GetChannel()) {
			o.Rollback()
			logs.Warn("apk worker channel full! len:%d", len(*mod.Worker.GetChannel()))
			return errors.New("apk worker channel full")
		}
		*mod.Worker.GetChannel() <- apkJob
	}

	o.Commit()

	return nil
}

/******************************************************************************
 **函数名称: GetApkList
 **功    能: 获取apk分页信息
 **输入参数:
 **     fields:{字段:值}
 **     page:页码
 **     pageSize:条数
 **输出参数: NONE
 **返    回:
 ** 	int64: 总长度
 ** 	[]*upgrade.Apk: APK列表
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-03 14:14:35 #
 ******************************************************************************/
func (mod *Models) GetApkList(
	fields map[string]interface{}, page, pageSize int) (int64, []*upgrade.Apk, error) {

	apks := new([]*upgrade.Apk)

	cond := orm.NewCondition()

	for k, v := range fields {
		cond = cond.And(k, v)
	}
	apkTable := upgrade.Apk{}
	qs := mod.Mysql.O.QueryTable(apkTable.TableName()).SetCond(cond)
	total, err := qs.Count()
	if nil != err {
		return 0, nil, err
	}

	_, err = qs.OrderBy("-id").Limit(pageSize, (page-1)*pageSize).All(apks)
	if nil != err {
		return 0, nil, err
	}

	return total, *apks, err
}
