package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/utils"
)

// AppController operations for App
type AppController struct {
	CommonController
}

// URLMapping ...
func (this *AppController) URLMapping() {
	this.Mapping("Post", this.Post)
	this.Mapping("Put", this.Put)
	this.Mapping("Get", this.Get)
	this.Mapping("GetList", this.GetList)
}

/* 创建应用请求参数 */
type App struct {
	Id          int64   `json:"id" description:"应用ID"`
	Name        string  `json:"name" description:"应用名称"`
	PackageName string  `json:"package_name" description:"应用包名（不唯一）；注：存在包名相同，但分属不同的平台"`
	BusinessId  int64   `json:"business_id" description:"业务ID；所属业务线"`
	CdnPlatId   int64   `json:"cdn_plat_id" description:"CDN 平台ID"`
	CdnSplatId  int64   `json:"cdn_splat_id" description:"CDN 子平台ID；SPLATID：子平台ID。用于业务计费，由CDN统一分配"`
	DevTypeId   int64   `json:"dev_type_id" description:"设备种类ID；"`
	Enable      int     `json:"enable" description:"是否启用；可通过此开关控制应用上下线；0：禁用（下线）；1：启用（上线）"`
	HasDevPlat  int     `json:"has_dev_plat" description:"是否关联设备平台；0：不关联。即：适用于所有设备平台；1：关联。即：只适用于指定设备平台。如果关联表中无设备平台，表示所有设备平台均不支持。"`
	DevPlatIds  []int64 `json:"dev_plat_ids" description:"关联设备平台"`
	Description string  `json:"description"  description:"描述信息"`
	CreateTime  string  `json:"create_time" description:"创建时间"`
	CreateUser  string  `json:"create_user" description:"页面创建者"`
	UpdateTime  string  `json:"update_time" description:"修改时间"`
	UpdateUser  string  `json:"update_user" description:"页面修改者"`
}

/* 创建应用请求参数 */
type AppResp struct {
	Id           int64                 `json:"id" description:"应用ID"`
	Name         string                `json:"name" description:"应用名称"`
	PackageName  string                `json:"package_name" description:"应用包名（不唯一）；注：存在包名相同，但分属不同的平台"`
	BusinessId   int64                 `json:"business_id" description:"业务ID；所属业务线"`
	BusinessName string                `json:"business_name" description:"业务名称；所属业务线"`
	CdnPlatId    int64                 `json:"cdn_plat_id" description:"CDN 平台ID"`
	CdnSplatId   int64                 `json:"cdn_splat_id" description:"CDN 子平台ID；SPLATID：子平台ID。用于业务计费，由CDN统一分配"`
	DevTypeId    int64                 `json:"dev_type_id" description:"设备种类ID；"`
	DevTypeName  string                `json:"dev_type_name" description:"设备种类名称"`
	Enable       int                   `json:"enable" description:"是否启用；可通过此开关控制应用上下线；0：禁用（下线）；1：启用（上线）"`
	HasDevPlat   int                   `json:"has_dev_plat" description:"是否关联设备平台；0：不关联。即：适用于所有设备平台；1：关联。即：只适用于指定设备平台。如果关联表中无设备平台，表示所有设备平台均不支持。"`
	DevPlats     []*upgrade.DevPlatOut `json:"dev_plats" description:"关联设备平台"`
	Description  string                `json:"description"  description:"描述信息"`
	CreateTime   string                `json:"create_time" description:"创建时间"`
	CreateUser   string                `json:"create_user" description:"页面创建者"`
	UpdateTime   string                `json:"update_time" description:"修改时间"`
	UpdateUser   string                `json:"update_user" description:"页面修改者"`
}

// 新增应用
// @Title 新增应用
// @Description 新增应用
// @Param    body    body    controllers.AppReq    true    "body for App content"
// @Success 201 {int} comm.AddResp
// @Failure 403 body is empty
// @router / [post]
// @Author Shuangpeng.guo 2020-06-17 13:37:02
func (this *AppController) Post() {
	ctx := GetBackendCntx()

	// 获取请求数据&检查参数
	req, code, err := this.getAppPost()
	if nil != err {
		logs.Error("App parameter format is invalid! errmsg:%s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	} else if 0 == len(req.PackageName) {
		logs.Error("App package name not allowed empty.")
		this.ErrorMessage(comm.ERR_PARAM_MISS, "package name not allowed empty. ")
		return
	}

	// 验证平台信息
	s, err := ctx.Model.CdnSplat.CheckSplat(req.CdnPlatId)
	if nil != err {
		logs.Error("App SplatId check exception! id:%d errmsg:%s",
			req.CdnPlatId, err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "splat id check exception. ")
		return
	}
	req.CdnPlatId = s.PlatId

	// 验证设备平台信息
	extantApps, err := upgrade.GetAppAndPlatByPkg(ctx.Model.Mysql.O, req.PackageName)
	if nil != err {
		logs.Error("App check extant failed! req:%v errmsg:%s", req, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}

	platCheck := true
	for _, extantApp := range extantApps {
		if extantApp.HasDevPlat == 0 || req.HasDevPlat == 0 {
			platCheck = false
			break
		}
		_platIds := make([]int64, len(extantApp.DevPlats))
		for k, v := range extantApp.DevPlats {
			_platIds[k] = v.Id
		}
		if len(utils.Intersect(_platIds, req.DevPlatIds)) > 0 {
			platCheck = false
			break
		}
	}

	if !platCheck {
		logs.Error("App check extant conflict! req:%v extantApps:%v", req, extantApps)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "app dev plat conflict")
		return
	}

	app := &upgrade.App{
		Name:        req.Name,
		PackageName: req.PackageName,
		BusinessId:  req.BusinessId,
		CdnPlatId:   req.CdnPlatId,
		CdnSplatId:  req.CdnSplatId,
		DevTypeId:   req.DevTypeId,
		Enable:      req.Enable,
		HasDevPlat:  req.HasDevPlat,
		Description: req.Description,
		CreateTime:  time.Now(),
		CreateUser:  this.mail,
		UpdateTime:  time.Now(),
		UpdateUser:  this.mail,
	}

	id, err := ctx.Model.CreateApp(app, req.DevPlatIds)
	if nil != err {
		logs.Error("App create failed! param:%v errmsg:%s", app, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}

	this.FormatAddResp(http.StatusCreated, comm.OK, "Ok", id)
}

// 修改应用
// @Title 修改应用
// @Description 修改应用
// @Param    id        path     string                true    "The id you want to update"
// @Param    body    body    controllers.AppReq    true    "body for App content"
// @Success 201 {int} comm.BaseResp
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @Author Shuangpeng.guo 2020-06-17 13:37:22
func (this *AppController) Put() {
	ctx := GetBackendCntx()

	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Get id failed! errmsg:%s", err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	// 获取请求数据&检查参数
	req, code, err := this.getAppPost()
	if nil != err {
		logs.Error("Parameter format is invalid! errmsg:%s", err.Error())
		this.ErrorMessage(code, err.Error())
		return
	}

	// 检查APP是否存在
	app := &upgrade.App{
		Id: id,
	}
	if err = ctx.Model.Mysql.O.Read(app); nil != err {
		if orm.ErrNoRows == err {
			logs.Error("Update app: app not exist, id: %d", id)
			this.ErrorMessage(comm.ERR_PARAM_INVALID, "app not exist")
			return
		}
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}

	// 验证设备平台信息
	extantApps, err := upgrade.GetAppAndPlatByPkg(ctx.Model.Mysql.O, app.PackageName)
	if nil != err {
		logs.Error("App check extant failed! req:%v errmsg:%s", req, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}

	platCheck := true
	for _, extantApp := range extantApps {
		if extantApp.Id == app.Id {
			continue
		}
		if extantApp.HasDevPlat == 0 || req.HasDevPlat == 0 {
			platCheck = false
			break
		}
		_platIds := make([]int64, len(extantApp.DevPlats))
		for k, v := range extantApp.DevPlats {
			_platIds[k] = v.Id
		}
		if len(utils.Intersect(_platIds, req.DevPlatIds)) > 0 {
			platCheck = false
			break
		}
	}

	if !platCheck {
		logs.Error("App check extant conflict! req:%v extantApps:%v", req, extantApps)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "app dev plat conflict")
		return
	}

	app.Name = req.Name
	app.DevTypeId = req.DevTypeId
	app.Enable = req.Enable
	app.HasDevPlat = req.HasDevPlat
	app.Description = req.Description
	app.UpdateTime = time.Now()
	app.UpdateUser = this.mail

	err = ctx.Model.UpdateApp(app, req.DevPlatIds)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("App update failed! param:%v errmsg:%s", app, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("App not exist! id:%d", id)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist!")
		return
	}

	this.FormatResp(http.StatusCreated, comm.OK, "Ok")
}

// 查询单条应用
// @Title 查询单条应用
// @Description 查询单条应用
// @Param    id    path    string    true    "app id"
// @Success 200 {object} comm.InterfaceResp
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @Author Shuangpeng.guo 2020-06-17 13:37:44
func (this *AppController) Get() {
	ctx := GetBackendCntx()
	id, err := strconv.ParseInt(this.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Parase id failed! errmsg:%s", err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	appModel, err := upgrade.GetAppAndPlatById(ctx.Model.Mysql.O, id)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Query app and plat failed! id:%v errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("App not exist! id: %d", id)
		this.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist!")
		return
	}

	app := &AppResp{
		Id:           appModel.Id,
		Name:         appModel.Name,
		PackageName:  appModel.PackageName,
		BusinessId:   appModel.BusinessId,
		BusinessName: appModel.BusinessName,
		CdnPlatId:    appModel.CdnPlatId,
		CdnSplatId:   appModel.CdnSplatId,
		DevTypeId:    appModel.DevTypeId,
		DevTypeName:  appModel.DevTypeName,
		Enable:       appModel.Enable,
		HasDevPlat:   appModel.HasDevPlat,
		DevPlats:     appModel.DevPlats,
		Description:  appModel.Description,
		CreateTime:   appModel.CreateTime.Format("2006-01-02 15:04:05"),
		CreateUser:   appModel.CreateUser,
		UpdateTime:   appModel.CreateTime.Format("2006-01-02 15:04:05"),
		UpdateUser:   appModel.UpdateUser,
	}

	this.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", app)
}

// 获取应用列表
// @Title 获取应用列表
// @Description 获取应用列表
// @Param    business_id	query    int64		false    "业务线ID     非管理员用户必传（不做逻辑限制）"
// @Param    id				query    int64		false    "应用ID"
// @Param    name			query    string		false    "应用名称"
// @Param    package_name	query    string		false    "应用包名"
// @Param    cdn_splat_id	query    int64		false    "子平台ID    "
// @Param    dev_type_id	query    int64		false    "设备种类ID"
// @Param    enable			query    int64		false    "是否启用"
// @Param    page			query    int        false    "页号 默认为1"
// @Param    page_size		query    int        false    "条目数 默认为20"
// @Success 200 {object} comm.InterfaceResp
// @Failure 403
// @router /list [get]
// @Author Shuangpeng.guo 2020-06-17 13:38:02
func (this *AppController) GetList() {
	ctx := GetBackendCntx()

	fields, values, page, pageSize, err := this.getListParam()
	if nil != err {
		logs.Error("Get query param failed! errmsg:%s", err.Error())
		this.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	total, apps, err := upgrade.GetAppAndPlatList(
		ctx.Model.Mysql.O, fields, values, page, pageSize)
	if nil != err {
		logs.Error("App list query failed! fields:%v values:%v errmsg:%s",
			fields, values, err.Error())
		code, e := database.MysqlFormatError(err)
		this.ErrorMessage(code, e)
		return
	}
	l := len(apps)

	listData := &comm.ListData{
		Page:     page,
		PageSize: pageSize,
		Total:    int(total),
		Len:      l,
		List:     make([]interface{}, l),
	}

	for k, app := range apps {
		listData.List[k] = &AppResp{
			Id:           app.Id,
			Name:         app.Name,
			PackageName:  app.PackageName,
			BusinessId:   app.BusinessId,
			BusinessName: app.BusinessName,
			CdnPlatId:    app.CdnPlatId,
			CdnSplatId:   app.CdnSplatId,
			DevTypeId:    app.DevTypeId,
			DevTypeName:  app.DevTypeName,
			Enable:       app.Enable,
			HasDevPlat:   app.HasDevPlat,
			DevPlats:     app.DevPlats,
			Description:  app.Description,
			CreateTime:   app.CreateTime.Format("2006-01-02 15:04:05"),
			CreateUser:   app.CreateUser,
			UpdateTime:   app.UpdateTime.Format("2006-01-02 15:04:05"),
			UpdateUser:   app.UpdateUser,
		}
	}

	this.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", listData)
}

////////////////////////////////////////////////////////////////////////////////

/******************************************************************************
 **函数名称: getAppPost
 **功    能: 内部方法 获取post参数
 **输入参数:
 **输出参数: App post 参数
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-17 13:40:15 #
 ******************************************************************************/
func (this *AppController) getAppPost() (*App, int, error) {
	req := &App{}

	if err := jsoniter.Unmarshal(this.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("App parameter format is invalid! body:%s errmsg:%s",
			this.Ctx.Input.RequestBody, err.Error())
		return nil, comm.ERR_PARAM_INVALID, err
	} else if 0 == len(req.Name) {
		logs.Error("App name not allowed empty.")
		return nil, comm.ERR_PARAM_MISS, errors.New("appName not allowed empty")
	} else if 0 == req.DevTypeId {
		logs.Error("App devTypeId not allowed empty.")
		return nil, comm.ERR_PARAM_MISS, errors.New("app devTypeId not allowed empty")
	}

	/*
	   // 检查业务线
	   _, err := upgrade.GetBusinessById(ctx.Model.Mysql.O, req.BusinessId)
	   if nil != err {
	       logs.Error("App business check failed! businessId:%d errmsg:%s",
	           req.BusinessId, err.Error())
	       code, e := database.MysqlFormatError(err)
	       return nil, code, errors.New(e)
	   }

	   // 检查设备种类
	   _, err = upgrade.GetDevTypeById(ctx.Model.Mysql.O, req.DevTypeId)
	   if nil != err {
	       logs.Error("App dev type check failed! devType:%d errmsg:%s",
	           req.DevTypeId, err.Error())
	       code, e := database.MysqlFormatError(err)
	       return nil, code, errors.New(e)
	   }
	*/

	return req, 0, nil
}

/******************************************************************************
 **函数名称: getListParam
 **功    能: 内部方法 获取get参数
 **输入参数:
 **输出参数:
 ** 	fields: 查询参数
 ** 	values: 对应查询条件
 ** 	page: 分页页码
 ** 	pageSize: 单页条数
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2020-06-17 13:42:44 #
 ******************************************************************************/
func (this *AppController) getListParam() (
	fields []string, values []interface{}, page, pageSize int, err error) {

	fields = []string{}
	values = []interface{}{}
	page = comm.PAGE_START
	pageSize = comm.PAGE_SIZE

	if bIdStr := this.GetString("business_id"); "" != bIdStr {
		bid, err := strconv.ParseInt(bIdStr, 10, 64)
		if nil != err {
			logs.Error("Get business id failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "app.business_id = ?")
		values = append(values, bid)
	}

	if idStr := this.GetString("id"); "" != idStr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if nil != err {
			logs.Error("Get id failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "app.id = ?")
		values = append(values, id)
	}

	if name := this.GetString("name"); "" != name {
		fields = append(fields, "app.name LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", name))
	}

	if packageName := this.GetString("app.package_name"); "" != packageName {
		fields = append(fields, "package_name LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", packageName))
	}

	if cdnSplatIdStr := this.GetString("cdn_splat_id"); "" != cdnSplatIdStr {
		cdnSplatId, err := strconv.ParseInt(cdnSplatIdStr, 10, 64)
		if nil != err {
			logs.Error("Get cdn splat id failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "app.cdn_splat_id = ?")
		values = append(values, cdnSplatId)
	}

	if devTypeIdStr := this.GetString("dev_type_id"); "" != devTypeIdStr {
		devTypeId, err := strconv.ParseInt(devTypeIdStr, 10, 64)
		if nil != err {
			logs.Error("Get dev type id failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "app.dev_type_id = ?")
		values = append(values, devTypeId)
	}

	if enableStr := this.GetString("enable"); "" != enableStr {
		enable, err := strconv.ParseInt(enableStr, 10, 64)
		if nil != err {
			logs.Error("Get enable failed! errmsg:%s", err.Error())
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "app.enable = ?")
		values = append(values, enable)
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
