package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
	"upgrade-api/src/share/strategy"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
)

// ApkController operations for Apk
type ApkController struct {
	ApiBaseController
}

// URLMapping ...
func (c *ApkController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("GetOneByBcode", c.GetOneByBcode)
}

/* 创建APK请求参数 */
type Apk struct {
	Id             int64  `json:"id"`               // APK ID
	AppId          int64  `json:"app_id"`           // 应用ID
	Enable         int    `json:"enable"`           // 是否启用（0：禁用 1：启用）
	Bcode          string `json:"bcode"`            // 业务方Code
	VersionCode    int64  `json:"version_code"`     // 版本号
	VersionName    string `json:"version_name"`     // 版本名
	Url            string `json:"url"`              // 下载地址
	Md5            string `json:"md5"`              // MD5值
	Size           int64  `json:"size"`             // APK大小(字节)
	EuiLowVersion  string `json:"eui_low_version"`  // 依赖的EUI最低版本
	EuiHighVersion string `json:"eui_high_version"` // 依赖的EUI最高版本
	Status         int    `json:"status"`           // 状态
	Description    string `json:"description"`      // 描述
	Memo           string `json:"memo"`             // 备注信息
	CreateTime     string `json:"create_time"`      // 创建时间
	CreateUser     string `json:"create_user"`      // 创建者
	UpdateTime     string `json:"update_time"`      // 更新时间
	UpdateUser     string `json:"update_user"`      // 更新者
	CallBackUrl    string `json:"callback_url"`     // 回调地址
}

// Post ...
// @Title Post
// @Description create Apk
// @Param	body		body 	models.Apk	true		"body for Apk content"
// @Success 201 {int} models.Apk
// @Failure 403 body is empty
// @router  / [post]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkController) Post() {
	req := &Apk{}
	ctx := GetApiCntx()

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("App parameter is invalid! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(req.Url) {
		logs.Error("Apk url not allowed empty! body:%s", c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Url not allowed empty!")
	} else if 0 == req.AppId {
		logs.Error("Apk appId not allowed empty. body:%s", c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "App id not allowed empty!")
		return
	} else if "" == req.Bcode {
		logs.Error("Apk bcode not allowed empty.")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "bcode not allowed empty. ")
		return
	}

	//检查版本号
	lastApk, err := upgrade.GetApkMaxVersionCodeByAppId(ctx.Model.Mysql.O, req.AppId)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("get apk max version code failed, appId: %d, msg: %s",
			req.AppId, err.Error())
		code, msg := database.MysqlFormatError(err)
		c.ErrorMessage(code, msg)
		return
	} else if orm.ErrNoRows != err {
		if lastApk.VersionCode >= req.VersionCode {
			logs.Warn("apk version code illegal, versionCode: %d", req.VersionCode)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "apk version code illegal")
			return
		}
	}

	// 检查APP是否存在
	app, err := upgrade.GetAppById(ctx.Model.Mysql.O, req.AppId)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("APK POST: get app by id failed, Id: %d, msg: %s",
			req.AppId, err.Error())
		code, msg := database.MysqlFormatError(err)
		c.ErrorMessage(code, msg)
		return
	} else if orm.ErrNoRows == err {
		logs.Warn("APK POST: app not exist, Id: %d", req.AppId)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "app not exist")
		return
	}

	// 检查上传地址合法性
	apkUrl, err := url.Parse(req.Url)
	if nil != err {
		logs.Warn("APK POST: apk url invalid, url: %s", req.Url)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "apk url invalid")
		return
	}

	req.Url = fmt.Sprintf("%s://%s%s?platid=%d&splatid=%d",
		apkUrl.Scheme, apkUrl.Host, apkUrl.Path, app.CdnPlatId, app.CdnSplatId)

	apk := &upgrade.Apk{
		AppId:             req.AppId,
		Enable:            req.Enable,
		Url:               req.Url,
		EuiLowVersion:     req.EuiLowVersion,
		EuiLowVersionInt:  strategy.DealEui(req.EuiLowVersion),
		EuiHighVersion:    req.EuiHighVersion,
		EuiHighVersionInt: strategy.DealEui(req.EuiHighVersion),
		Status:            upgrade.APK_STATUS_PROCESSING,
		Description:       req.Description,
		Memo:              req.Memo,
		CallbackUrl:       req.CallBackUrl,
		Size:              req.Size,
		Md5:               req.Md5,
		VersionCode:       req.VersionCode,
		VersionName:       req.VersionName,
		CreateTime:        time.Now(),
		CreateUser:        comm.SYSYTEM_USER,
		UpdateTime:        time.Now(),
		Bcode:             req.Bcode,
	}
	/* 创建apk */
	id, err := ctx.Model.CreateApk(apk, req.CallBackUrl)
	if nil != err {
		logs.Error("Apk create failed! param:%v errmsg:%s", apk, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusCreated, comm.OK, "Ok", id)
}

// GetOne ...
// @Title Get One
// @Description get Apk by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Apk
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkController) GetOne() {
	ctx := GetApiCntx()

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Get apk id failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	apkModel, err := upgrade.GetApkById(ctx.Model.Mysql.O, id)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get apk failed! id:%d errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("Apk not exist! id:%d", id)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist!")
		return
	}

	apk := &Apk{
		Id:             apkModel.Id,
		AppId:          apkModel.AppId,
		Enable:         apkModel.Enable,
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
		Bcode:          apkModel.Bcode,
		CallBackUrl:    apkModel.CallbackUrl,
		CreateTime:     apkModel.CreateTime.Format("2006-01-02 15:04:05"),
		CreateUser:     apkModel.CreateUser,
		UpdateTime:     apkModel.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateUser:     apkModel.UpdateUser,
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", apk)
}

// GetAll ...
// @Title Get All
// @Description get Apk
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Apk
// @Failure 403
// @router /list [get]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkController) GetAll() {
	ctx := GetApiCntx()

	fields, page, pageSize, err := c.getListParam()
	if nil != err {
		logs.Error("Get parameter failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	total, apks, err := ctx.Model.GetApkList(fields, page, pageSize)
	if nil != err {
		logs.Error("Get apk list failed! fields:%v errmsg:%s",
			fields, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
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
			Enable:         apk.Enable,
			VersionCode:    apk.VersionCode,
			VersionName:    apk.VersionName,
			Url:            apk.Url,
			Md5:            apk.Md5,
			Size:           apk.Size,
			EuiLowVersion:  apk.EuiLowVersion,
			EuiHighVersion: apk.EuiHighVersion,
			Status:         apk.Status,
			CallBackUrl:    apk.CallbackUrl,
			Bcode:          apk.Bcode,
			Description:    apk.Description,
			Memo:           apk.Memo,
			CreateTime:     apk.CreateTime.Format("2006-01-02 15:04:05"),
			CreateUser:     apk.CreateUser,
			UpdateTime:     apk.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateUser:     apk.UpdateUser,
		}
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", listData)
}

// Put ...
// @Title Put
// @Description update the Apk
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Apk	true		"body for Apk content"
// @Success 200 {object} models.Apk
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkController) Put() {
	ctx := GetApiCntx()

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Get apk id failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 获取请求数据&检查参数
	req := &Apk{}
	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("Unmarshal parameter failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(req.Url) {
		logs.Error("Apk url not allowed empty. body:%s", c.Ctx.Input.RequestBody)
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Url not allowed empty!")
		return
	}

	apk, err := upgrade.GetApkById(ctx.Model.Mysql.O, id)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("APK PUT: get apk by id failed, Id: %d, msg: %s",
			id, err.Error())
		code, msg := database.MysqlFormatError(err)
		c.ErrorMessage(code, msg)
		return
	} else if orm.ErrNoRows == err {
		logs.Warn("APK PUT: apk not exist, Id: %d", id)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "apk not exist")
		return
	}

	// 检查APP是否存在
	app, err := upgrade.GetAppById(ctx.Model.Mysql.O, apk.AppId)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("APK PUT: get app by id failed, Id: %d, msg: %s",
			req.AppId, err.Error())
		code, msg := database.MysqlFormatError(err)
		c.ErrorMessage(code, msg)
		return
	} else if orm.ErrNoRows == err {
		logs.Warn("APK PUT: app not exist, Id: %d", req.AppId)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "app not exist")
		return
	}

	// 检查上传地址合法性
	apkUrl, err := url.Parse(req.Url)
	if nil != err {
		logs.Warn("APK PUT: apk url invalid, url: %s", req.Url)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "apk url invalid")
		return
	}

	doCheck := false
	if _url := fmt.Sprintf("%s://%s%s?platid=%d&splatid=%d",
		apkUrl.Scheme, apkUrl.Host, apkUrl.Path, app.CdnPlatId, app.CdnSplatId); _url != apk.Url || upgrade.APK_STATUS_NORMAL != apk.Status {
		doCheck = true
		apk.Url = _url
		apk.Status = upgrade.APK_STATUS_PROCESSING
	}

	if doCheck {
		// 检查版本号
		lastApk, err := upgrade.GetApkMaxVersionCodeByAppId(ctx.Model.Mysql.O, req.AppId)
		if nil != err {
			logs.Error("get apk max version code failed, appId: %d, msg: %s",
				req.AppId, err.Error())
			code, msg := database.MysqlFormatError(err)
			c.ErrorMessage(code, msg)
			return
		} else if lastApk.Id != id && lastApk.VersionCode >= req.VersionCode {
			logs.Warn("apk version code illegal, versionCode: %d", req.VersionCode)
			c.ErrorMessage(comm.ERR_PARAM_INVALID, "apk version code illegal")
			return
		}
	}

	apk.Enable = req.Enable
	apk.VersionCode = req.VersionCode
	apk.VersionName = req.VersionName
	apk.Md5 = req.Md5
	apk.Size = req.Size
	apk.EuiLowVersion = req.EuiLowVersion
	apk.EuiLowVersionInt = strategy.DealEui(req.EuiLowVersion)
	apk.EuiHighVersion = req.EuiHighVersion
	apk.EuiHighVersionInt = strategy.DealEui(req.EuiHighVersion)
	apk.Description = req.Description
	apk.Memo = req.Memo
	apk.CallbackUrl = req.CallBackUrl
	apk.UpdateUser = comm.SYSYTEM_USER

	err = ctx.Model.UpdateApk(apk, req.CallBackUrl, doCheck)
	if nil != err {
		logs.Error("Update apk failed! param:%v errmsg:%s", apk, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatResp(http.StatusCreated, comm.OK, "Ok")
}

// GetOneByBcode ...
// @Title Get One
// @Description get Apk by id
// @Param	bcode		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Apk
// @Failure 403 :bcode is empty
// @router /bcode/:bcode [get]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *ApkController) GetOneByBcode() {
	ctx := GetApiCntx()

	bcode := c.Ctx.Input.Param(":bcode")
	if "" == bcode {
		logs.Error("Get apk bcode failed!")
		c.ErrorMessage(comm.ERR_PARAM_INVALID, errors.New("Get apk bcode failed!").Error())
		return
	}

	apkModel, err := upgrade.GetApkByBcode(ctx.Model.Mysql.O, bcode)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get apk failed! bcode:%d errmsg:%s", bcode, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("Apk not exist! bcode:%d", bcode)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "apk not exist!")
		return
	}

	apk := &Apk{
		Id:             apkModel.Id,
		AppId:          apkModel.AppId,
		Enable:         apkModel.Enable,
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
		Bcode:          apkModel.Bcode,
		CallBackUrl:    apkModel.CallbackUrl,
		CreateTime:     apkModel.CreateTime.Format("2006-01-02 15:04:05"),
		CreateUser:     apkModel.CreateUser,
		UpdateTime:     apkModel.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateUser:     apkModel.UpdateUser,
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", apk)
}

/******************************************************************************
 **函数名称: getListParam
 **功    能: 获取参数列表
 **输入参数: NONE
 **输出参数: NONE
 **     fields:字段
 **     values:参数
 **     page:页码
 **     pageSize:条数
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # peng.jia # 2020-06-02 10:14:35 #
 ******************************************************************************/
func (c *ApkController) getListParam() (
	fields map[string]interface{}, page, pageSize int, err error) {

	fields = make(map[string]interface{})
	page = comm.PAGE_START
	pageSize = comm.PAGE_SIZE

	appId, err := c.GetInt64("app_id")
	if nil != err {
		return nil, 0, 0, errors.New("get app_id failed,can't nil")
	}
	fields["app_id"] = appId

	if idStr := c.GetString("id"); "" != idStr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if nil != err {
			return nil, 0, 0, err
		}
		fields["id"] = id
	}

	if enableStr := c.GetString("enable"); "" != enableStr {
		enable, err := strconv.ParseInt(enableStr, 10, 64)
		if nil != err {
			return nil, 0, 0, err
		}
		fields["enable"] = enable
	}

	if versionCodeStr := c.GetString("version_code"); "" != versionCodeStr {
		versionCode, err := strconv.ParseInt(versionCodeStr, 10, 64)
		if nil != err {
			return nil, 0, 0, err
		}
		fields["version_code"] = versionCode
	}

	if versionName := c.GetString("version_name"); "" != versionName {
		fields["version_name__icontains"] = versionName

	}

	if statusStr := c.GetString("status"); "" != statusStr {
		status, err := strconv.ParseInt(statusStr, 10, 64)
		if nil != err {
			return nil, 0, 0, err
		}
		fields["status"] = status
	}

	if pageStr := c.GetString("page"); "" != pageStr {
		page, err = strconv.Atoi(pageStr)
		if nil != err {
			return nil, 0, 0, err
		}
	}

	if pageSizeStr := c.GetString("page_size"); "" != pageSizeStr {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if nil != err {
			return nil, 0, 0, err
		}
	}

	return fields, page, pageSize, nil
}
