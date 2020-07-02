package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	ApiBaseController
}

// URLMapping ...
func (c *AppController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
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

// Post ...
// @Title Post
// @Description create App
// @Param	body		body 	models.App	true		"body for App content"
// @Success 201 {int} models.App
// @Failure 403 body is empty
// @router  / [post]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *AppController) Post() {
	req := &App{}
	ctx := GetApiCntx()

	/* 解析body */
	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("Unmarshal request failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	/* 解析Header 取business_id */
	if businessId := c.Ctx.Input.Header("business-id"); businessId != "" {
		bId, err := strconv.ParseInt(businessId, 10, 64)
		if nil != err {
			logs.Error("Parse business id failed! errmsg:%s", err.Error())
			c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
			return
		}
		req.BusinessId = bId
	}

	if 0 == len(req.Name) {
		logs.Error("App name not allowed empty.")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "App name not allowed empty.")
		return
	} else if 0 == len(req.PackageName) {
		logs.Error("App package name not allowed empty.")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Package name not allowed empty.")
		return
	} else if 0 == req.DevTypeId {
		logs.Error("App devTypeId not allowed empty.")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "App devTypeId not allowed empty.")
		return
	}

	/* 验证平台信息 */
	s, err := ctx.Model.CdnSplat.CheckSplat(req.CdnSplatId)
	if nil != err {
		logs.Error("App SplatId check exception! id:%d errmsg:%s", req.CdnPlatId, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}
	req.CdnPlatId = s.PlatId

	// 验证设备平台信息
	extantApps, err := upgrade.GetAppAndPlatByPkg(ctx.Model.Mysql.O, req.PackageName)
	if nil != err {
		logs.Error("App check extant failed! req:%v errmsg:%s", req, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
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
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "app dev plat conflict")
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
		CreateUser:  comm.SYSYTEM_USER,
		UpdateTime:  time.Now(),
	}

	id, err := ctx.Model.CreateApp(app, req.DevPlatIds)
	if nil != err {
		logs.Error("Create app failed! param:%v errmsg:%s", app, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusCreated, comm.OK, "Ok", id)
}

// GetOne ...
// @Title Get One
// @Description get App by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.App
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *AppController) GetOne() {
	ctx := GetApiCntx()

	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Get id failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}
	/* 获取应用及关联平台信息 */
	appModel, err := ctx.Model.GetAppAndPlat(id)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Get app and plat failed! id:%v errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("App not exist! id:%d", id)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist.")
		return
	}

	_plats := strings.Split(appModel.DevPlatIds, ",")
	plats := make([]int64, len(_plats))
	for k, v := range _plats {
		plats[k], _ = strconv.ParseInt(v, 10, 64)
	}

	app := &App{
		Id:          appModel.Id,
		Name:        appModel.Name,
		PackageName: appModel.PackageName,
		BusinessId:  appModel.BusinessId,
		CdnPlatId:   appModel.CdnPlatId,
		CdnSplatId:  appModel.CdnSplatId,
		DevTypeId:   appModel.DevTypeId,
		Enable:      appModel.Enable,
		HasDevPlat:  appModel.HasDevPlat,
		DevPlatIds:  plats,
		Description: appModel.Description,
		CreateTime:  appModel.CreateTime.Format("2006-01-02 15:04:05"),
		CreateUser:  appModel.CreateUser,
		UpdateUser:  appModel.UpdateUser,
		UpdateTime:  appModel.UpdateTime.Format("2006-01-02 15:04:05"),
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", app)

}

// GetAll ...
// @Title Get All
// @Description get App
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.App
// @Failure 403
// @router /list [get]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *AppController) GetAll() {
	ctx := GetApiCntx()

	fields, values, page, pageSize, err := c.getListParam()
	if nil != err {
		logs.Error("Get parameter failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}
	/* 获取应用及关联平台信息 */
	total, apps, err := ctx.Model.GetAppAndPlatList(fields, values, page, pageSize)
	if nil != err {
		logs.Error("Get app list failed! fields:%v values:%v errmsg:%s",
			fields, values, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
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
		_plats := strings.Split(app.DevPlatIds, ",")
		plats := make([]int64, len(_plats))
		for k, v := range _plats {
			plats[k], _ = strconv.ParseInt(v, 10, 64)
		}

		listData.List[k] = &App{
			Id:          app.Id,
			Name:        app.Name,
			PackageName: app.PackageName,
			BusinessId:  app.BusinessId,
			CdnPlatId:   app.CdnPlatId,
			CdnSplatId:  app.CdnSplatId,
			DevTypeId:   app.DevTypeId,
			Enable:      app.Enable,
			HasDevPlat:  app.HasDevPlat,
			DevPlatIds:  plats,
			Description: app.Description,
			CreateTime:  app.CreateTime.Format("2006-01-02 15:04:05"),
			CreateUser:  app.CreateUser,
			UpdateTime:  app.CreateTime.Format("2006-01-02 15:04:05"),
			UpdateUser:  app.UpdateUser,
		}
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", listData)
}

// Put ...
// @Title Put
// @Description update the App
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.App	true		"body for App content"
// @Success 200 {object} models.App
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # jiapeng # 2020-06-08 16:36:35 #
func (c *AppController) Put() {
	ctx := GetApiCntx()

	/* 取id */
	id, err := strconv.ParseInt(c.Ctx.Input.Param(":id"), 10, 64)
	if nil != err {
		logs.Error("Get id failed! errmsg:%s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	/* 解析body */
	req := &App{}

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, req); nil != err {
		logs.Error("Unmarshal request failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	} else if 0 == len(req.Name) {
		logs.Error("App name not allowed empty.")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "Name not allowed empty.")
		return
	} else if 0 == req.DevTypeId {
		logs.Error("App devTypeId not allowed empty.")
		c.ErrorMessage(comm.ERR_PARAM_MISS, "App devTypeId not allowed empty.")
		return
	}

	app, err := upgrade.GetAppById(ctx.Model.Mysql.O, id)
	if err != nil && err != orm.ErrNoRows {
		logs.Error("app not exist ! id:%v errmsg:%s", id, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	} else if err == orm.ErrNoRows {
		logs.Error("app not exist ! id:%v errmsg:%s", id, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	app.Name = req.Name
	app.DevTypeId = req.DevTypeId
	app.Enable = req.Enable
	app.HasDevPlat = req.HasDevPlat
	app.Description = req.Description
	app.UpdateUser = comm.SYSYTEM_USER
	app.UpdateTime = time.Now()

	// 验证设备平台信息
	extantApps, err := upgrade.GetAppAndPlatByPkg(ctx.Model.Mysql.O, req.PackageName)
	if nil != err {
		logs.Error("App check extant failed! req:%v errmsg:%s", req, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
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
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "app dev plat conflict")
		return
	}

	/* 修改应用 */
	err = ctx.Model.UpdateApp(app, req.DevPlatIds)
	if nil != err && orm.ErrNoRows != err {
		logs.Error("Create app failed! param:%v errmsg:%s", app, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	} else if orm.ErrNoRows == err {
		logs.Error("App not exist! id:%d", id)
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "App not exist.")
		return
	}

	c.FormatResp(http.StatusCreated, comm.OK, "Ok")
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
func (c *AppController) getListParam() (
	fields []string, values []interface{}, page, pageSize int, err error) {

	fields = []string{}
	values = []interface{}{}
	page = comm.PAGE_START
	pageSize = comm.PAGE_SIZE

	bIdStr := c.GetString("business_id")
	if "" != bIdStr {
		bid, err := strconv.ParseInt(bIdStr, 10, 64)
		if nil != err {
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "business_id = ?")
		values = append(values, bid)
	}

	idStr := c.GetString("id")
	if "" != idStr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if nil != err {
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "id = ?")
		values = append(values, id)
	}
	name := c.GetString("name")
	if "" != name {
		fields = append(fields, "name LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", name))
	}

	packageName := c.GetString("package_name")
	if "" != packageName {
		fields = append(fields, "package_name LIKE ?")
		values = append(values, fmt.Sprintf("%%%s%%", packageName))
	}

	cdnSplatIdStr := c.GetString("cdn_splat_id")
	if "" != cdnSplatIdStr {
		cdnSplatId, err := strconv.ParseInt(cdnSplatIdStr, 10, 64)
		if nil != err {
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "cdn_splat_id = ?")
		values = append(values, cdnSplatId)
	}

	devTypeIdStr := c.GetString("dev_type_id")
	if "" != devTypeIdStr {
		devTypeId, err := strconv.ParseInt(devTypeIdStr, 10, 64)
		if nil != err {
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "dev_type_id = ?")
		values = append(values, devTypeId)
	}

	enableStr := c.GetString("enable")
	if "" != enableStr {
		enable, err := strconv.ParseInt(idStr, 10, 64)
		if nil != err {
			return nil, nil, 0, 0, err
		}
		fields = append(fields, "enable = ?")
		values = append(values, enable)
	}

	pageStr := c.GetString("page")
	if "" != pageStr {
		page, err = strconv.Atoi(idStr)
		if nil != err {
			return nil, nil, 0, 0, err
		}
	}

	pageSizeStr := c.GetString("page_size")
	if "" != pageSizeStr {
		pageSize, err = strconv.Atoi(pageSizeStr)
		if nil != err {
			return nil, nil, 0, 0, err
		}
	}

	return fields, values, page, pageSize, nil
}
