package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/strategy"
)

// ApkController operations for Apk
type ApkController struct {
	CommonController
}

// URLMapping ...
func (c *ApkController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Put", c.Put)
	c.Mapping("Get", c.Get)
	c.Mapping("GetList", c.GetList)
}

// Apk创建
// @Title Apk创建
// @Description Apk创建
// @Param	body	body	controller.Apk	true	"body for Apk content"
// @Success 201 {int} comm.AddResp
// @Failure 403 body is empty
// @router / [post]
// @Author Shuangpeng.guo 2020-06-17 13:43:16
func (this *ApkController) Post() {
	ctx := GetBackendCntx()

	// 获取请求数据&检查参数
	req, code, err := this.getApkPost()
	if nil != err {
		logs.Warn("Paramater is invalid! errmsg:%s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	} else if 0 == req.AppId {
		logs.Warn("Apk appId not allowed empty.")
		this.ErrorMessage(comm.ERR_PARAM_MISS, "AppId not allowed empty.")
		return
	} else if "" == req.Bcode {
		logs.Warn("Apk bcode not allowed empty.")
		this.ErrorMessage(comm.ERR_PARAM_MISS, "Bcode not allowed empty.")
		return
	}

	// 检查APP是否存在
	app, err := upgrade.GetAppById(ctx.Model.Mysql.O, req.AppId)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get app by id failed! Id:%d errmsg:%s",
			req.AppId, err.Error())
		code, msg := database.MysqlFormatError(err)
		this.ErrorMessage(code, msg)
		return
	} else if orm.ErrNoRows == err {
		logs.Warn("App not exist! Id:%d", req.AppId)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist")
		return
	}

	// 检查上传地址合法性
	apkUrl, err := url.Parse(req.Url)
	if nil != err {
		logs.Warn("Apk url invalid! url:%s", req.Url)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "Apk url invalid")
		return
	}

	req.Url = fmt.Sprintf("%s://%s%s?platid=%d&splatid=%d",
		apkUrl.Scheme, apkUrl.Host, apkUrl.Path, app.CdnPlatId, app.CdnSplatId)

	apk := &upgrade.Apk{
		AppId:             req.AppId,
		Bcode:             req.Bcode,
		Url:               req.Url,
		Enable:            req.Enable,
		EuiLowVersion:     req.EuiLowVersion,
		EuiLowVersionInt:  strategy.DealEui(req.EuiLowVersion),
		EuiHighVersion:    req.EuiHighVersion,
		EuiHighVersionInt: strategy.DealEui(req.EuiHighVersion),
		Status:            upgrade.APK_STATUS_PROCESSING,
		Description:       req.Description,
		Memo:              req.Memo,
		CreateTime:        time.Now(),
		CreateUser:        this.mail,
		UpdateTime:        time.Now(),
		UpdateUser:        this.mail,
	}

	id, err := ctx.Model.CreateApk(apk)
	if nil != err {
		logs.Error("Apk create failed! param:%v errmsg:%s", apk, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}

	this.FormatAddResp(http.StatusCreated, comm.OK, "Ok", id)
}

// 修改APK信息
// @Title 修改APK信息
// @Description 修改APK信息
// @Param	id		path 	string			true	"APK ID"
// @Param	body	body	controller.Apk	true	"body for Apk content"
// @Success 201 {int} comm.AddResp
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @Author Shuangpeng.guo 2020-06-17 13:43:39
func (this *ApkController) Put() {
	ctx := GetBackendCntx()

	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Get id failed! errmsg:%s", err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 获取请求数据&检查参数
	req, code, err := this.getApkPost()
	if nil != err {
		logs.Error("Parameter format is invalid! errmsg:%s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}

	// 检查APK是否存在
	apk, err := upgrade.GetApkById(ctx.Model.Mysql.O, id)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get apk by id failed! Id:%d errmsg:%s",
			id, err.Error())
		code, msg := database.MysqlFormatError(err)
		this.ErrorMessage(code, msg)
		return
	} else if orm.ErrNoRows == err {
		logs.Warn("Apk not exist! Id: %d", id)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "Apk not exist")
		return
	}

	// 检查APP是否存在
	app, err := upgrade.GetAppById(ctx.Model.Mysql.O, apk.AppId)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get app by id failed! Id:%d errmsg:%s",
			req.AppId, err.Error())
		code, msg := database.MysqlFormatError(err)
		this.ErrorMessage(code, msg)
		return
	} else if orm.ErrNoRows == err {
		logs.Warn("APK PUT: app not exist! AppId:%d", req.AppId)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist")
		return
	}

	// 检查上传地址合法性
	apkUrl, err := url.Parse(req.Url)
	if nil != err {
		logs.Warn("Apk url invalid! url:%s", req.Url)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "Apk url invalid")
		return
	}

	doCheck := false
	_url := fmt.Sprintf("%s://%s%s?platid=%d&splatid=%d",
		apkUrl.Scheme, apkUrl.Host, apkUrl.Path, app.CdnPlatId, app.CdnSplatId)
	if _url != apk.Url || upgrade.APK_STATUS_NORMAL != apk.Status {
		doCheck = true
		apk.Url = _url
		apk.Status = upgrade.APK_STATUS_PROCESSING
	}

	apk.Id = id
	apk.Enable = req.Enable
	apk.EuiLowVersion = req.EuiLowVersion
	apk.EuiLowVersionInt = strategy.DealEui(req.EuiLowVersion)
	apk.EuiHighVersion = req.EuiHighVersion
	apk.EuiHighVersionInt = strategy.DealEui(req.EuiHighVersion)
	apk.Description = req.Description
	apk.Memo = req.Memo
	apk.UpdateTime = time.Now()
	apk.UpdateUser = this.mail

	err = ctx.Model.UpdateApk(apk, doCheck)
	if nil != err {
		logs.Error("Apk update failed! param:%v errmsg:%s", apk, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}

	this.FormatResp(http.StatusCreated, comm.OK, "Ok")
}

// 查询单条APK
// @Title 查询单条APK
// @Description 查询单条APK
// @Param	id	path	string	true	"app id"
// @Success 200 {object} comm.InterfaceResp
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @Author Shuangpeng.guo 2020-06-17 13:43:59
func (this *ApkController) Get() {
	ctx := GetBackendCntx()

	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Get id failed! errmsg:%s", err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	apkModel, err := ctx.Model.GetApkById(id)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Apk query failed! id:%v errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("Apk not exist! id:%d", id)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist.")
		return
	}

	apk := &Apk{
		Id:             apkModel.Id,
		AppId:          apkModel.AppId,
		Enable:         apkModel.Enable,
		AppName:        apkModel.AppName,
		BusinessId:     apkModel.BusinessId,
		PackageName:    apkModel.PackageName,
		VersionCode:    apkModel.VersionCode,
		VersionName:    apkModel.VersionName,
		Url:            apkModel.Url,
		Md5:            apkModel.Md5,
		Size:           apkModel.Size,
		EuiLowVersion:  apkModel.EuiLowVersion,
		EuiHighVersion: apkModel.EuiHighVersion,
		Status:         apkModel.Status,
		Description:    apkModel.Description,
		Memo:           apkModel.Memo,
		CreateTime:     apkModel.CreateTime.Format("2006-01-02 15:04:05"),
		CreateUser:     apkModel.CreateUser,
		UpdateTime:     apkModel.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateUser:     apkModel.UpdateUser,
	}

	this.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", apk)
}

// 获取APK列表
// @Title 获取APK列表
// @Description 获取APK列表
// @Param	app_id			query	int64	true	"应用ID"
// @Param	id				query	int64	false	"APKID"
// @Param	version			query	string	false	"版本号"
// @Param	version_name	query	string	false	"版本名称"
// @Param	status			query	int		false	"状态"
// @Param	page			query	int		false	"页号 默认为1"
// @Param	page_size		query	int		false	"条目数 默认为20"
// @Success 200 {object} comm.InterfaceResp
// @Failure 403
// @router /list [get]
// @Author Shuangpeng.guo 2020-06-17 13:44:27
func (this *ApkController) GetList() {
	ctx := GetBackendCntx()

	fields, value, page, pageSize, err := this.getListParam()
	if nil != err {
		logs.Error("Get parameter failed! errmsg:%s", err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	total, apks, err := ctx.Model.GetApkList(fields, value, page, pageSize)
	if nil != err {
		logs.Error("Get app list failed! fields:%v errmsg:%s",
			fields, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}

	l := len(apks)

	listData := &comm.ListData{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
		Len:      l,
		List:     make([]interface{}, l),
	}

	for k, apk := range apks {
		listData.List[k] = &Apk{
			Id:             apk.Id,
			AppId:          apk.AppId,
			AppName:        apk.AppName,
			Enable:         apk.Enable,
			BusinessId:     apk.BusinessId,
			PackageName:    apk.PackageName,
			VersionCode:    apk.VersionCode,
			VersionName:    apk.VersionName,
			Url:            apk.Url,
			Md5:            apk.Md5,
			Size:           apk.Size,
			EuiLowVersion:  apk.EuiLowVersion,
			EuiHighVersion: apk.EuiHighVersion,
			Status:         apk.Status,
			Description:    apk.Description,
			Memo:           apk.Memo,
			CreateTime:     apk.CreateTime.Format("2006-01-02 15:04:05"),
			CreateUser:     apk.CreateUser,
			UpdateTime:     apk.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateUser:     apk.UpdateUser,
		}
	}

	this.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", listData)
}

////////////////////////////////////////////////////////////////////////////////

/* APK信息 */
type Apk struct {
	Id             int64  `json:"id"`               // APK ID
	AppId          int64  `json:"app_id"`           // 应用ID
	AppName        string `json:"app_name"`         // 应用名称
	BusinessId     int64  `json:"business_id"`      // 业务线ID
	Bcode          string `json:"bcode"`            // 业务方Code
	Enable         int    `json:"enable"`           // 是否启用（0：禁用 1：启用）
	VersionCode    int64  `json:"version_code"`     // 版本号
	VersionName    string `json:"version_name"`     // 版本名
	PackageName    string `json:"package_name"`     // 包名
	Url            string `json:"url"`              // 下载地址
	Md5            string `json:"md5"`              // APK MD5值
	Size           int64  `json:"size"`             // APK大小(字节)
	EuiLowVersion  string `json:"eui_low_version"`  // 依赖的EUI最低版本
	EuiHighVersion string `json:"eui_high_version"` // 依赖的EUI最高版本
	Status         int    `json:"status"`           // 状态(0:禁用 1:启用)
	Description    string `json:"description"`      // 描述信息
	Memo           string `json:"memo"`             // 备注信息
	CreateTime     string `json:"create_time"`      // 创建时间
	CreateUser     string `json:"create_user"`      // 创建者
	UpdateTime     string `json:"update_time"`      // 更新时间
	UpdateUser     string `json:"update_user"`      // 更新者
}

/******************************************************************************
 **函数名称: getApkPost
 **功    能: 内部方法 获取post参数
 **输入参数:
 **输出参数: Apk post 参数
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-17 13:40:15 #
 ******************************************************************************/
func (this *ApkController) getApkPost() (*Apk, int, error) {
	req := &Apk{}

	if err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("Unmarshal post failed! body:%s errmsg:%s",
			this.Ctx.Input.RequestBody, err.Error())
		return nil, comm.ERR_PARAM_INVALID, err
	} else if 0 == len(req.Url) {
		logs.Error("Apk url not allowed empty! body: %s", this.Ctx.Input.RequestBody)
		return nil, comm.ERR_PARAM_MISS, errors.New("Url not allowed empty.")
	}

	return req, 0, nil
}

/******************************************************************************
 **函数名称: getListParam
 **功    能: 内部方法 获取get参数
 **输入参数:
 **输出参数:
 ** 	fields: 查询参数
 ** 	page: 分页页码
 ** 	pageSize: 单页条数
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-17 13:42:44 #
 ******************************************************************************/
func (this *ApkController) getListParam() (
	fields []string, values []interface{}, page, pageSize int, err error) {

	fields = []string{}
	values = []interface{}{}
	page = comm.PAGE_START
	pageSize = comm.PAGE_SIZE

	if appIdStr := this.GetString("app_id"); "" != appIdStr {
		appId, err := strconv.ParseInt(appIdStr, 10, 64)
		if nil != err {
			logs.Error("Get appId id failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, errors.New("appId type error")
		}
		fields = append(fields, "apk.app_id = ?")
		values = append(values, appId)
	}

	if appName := this.GetString("app_name"); "" != appName {
		fields = append(fields, "app.app_name LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", appName))
	}

	if appIdStr := this.GetString("business_id"); "" != appIdStr {
		appId, err := strconv.ParseInt(appIdStr, 10, 64)
		if nil != err {
			logs.Error("Get businessId id failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, errors.New("businessId type error")
		}
		fields = append(fields, "app.business_id = ?")
		values = append(values, appId)
	}

	if idStr := this.GetString("id"); "" != idStr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if nil != err {
			logs.Error("Get id failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "apk.id = ?")
		values = append(values, id)
	}

	if enableStr := this.GetString("enable"); "" != enableStr {
		enable, err := strconv.ParseInt(enableStr, 10, 64)
		if nil != err {
			logs.Error("Get enable failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "apk.enable = ?")
		values = append(values, enable)
	}

	if pkgName := this.GetString("package_name"); "" != pkgName {
		fields = append(fields, "app.package_name LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", pkgName))
	}

	if versionCodeStr := this.GetString("version_code"); "" != versionCodeStr {
		versionCode, err := strconv.ParseInt(versionCodeStr, 10, 64)
		if nil != err {
			logs.Error("Get version code failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, errors.New("version code type error")
		}
		fields = append(fields, "apk.version_code = ?")
		values = append(values, versionCode)
	}

	if verName := this.GetString("version_name"); "" != verName {
		fields = append(fields, "apk.version_name LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", verName))
	}

	if statusStr := this.GetString("status"); "" != statusStr {
		status, err := strconv.ParseInt(statusStr, 10, 64)
		if nil != err {
			logs.Error("Get status failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, errors.New("status type error")
		}
		fields = append(fields, "apk.status = ?")
		values = append(values, status)
	}

	if pageStr := this.GetString("page"); "" != pageStr {
		page, err = strconv.Atoi(pageStr)
		if nil != err {
			logs.Error("Get page failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
	}

	if pageSizeStr := this.GetString("page_size"); "" != pageSizeStr {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if nil != err {
			logs.Error("Get page size failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
	}

	return fields, values, page, pageSize, nil
}
