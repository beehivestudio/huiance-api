package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/strategy"
)

// ApkStrategyController operations for ApkUpgrade
type ApkStrategyController struct {
	CommonController
}

// URLMapping ...
func (c *ApkStrategyController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
}

// Post ...
// @Title Post
// @Description create Apk
// @Param	body		body 	models.Apk	true		"body for Apk content"
// @Success 201 {int} models.Apk
// @Failure 403 body is empty
// @router / [post]
// @author # taoshengbo # 2020-06-08 16:36:35 #
func (c *ApkStrategyController) Post() {
	ctx := GetBackendCntx()

	var v upgrade.ApkUpgradeStrategy
	var err error

	v.CreateUser = c.mail
	v.UpdateUser = c.mail

	if err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); nil != err {
		logs.Error("Unmarshal paramater failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v.BeginDatetime = time.Unix(v.ReqBeginDatetime, 0)
	v.EndDatetime = time.Unix(v.ReqEndDatetime, 0)

	var id int64
	if id, err = strategy.Create(ctx.Model.Mysql, ctx.Model.Redis, ctx.Model.Quota, &v); nil != err {
		logs.Error("Create strategy failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatAddResp(http.StatusOK, comm.OK, "Ok", id)
}

// Put ...
// @Title Put
// @Description update the Apk
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Apk	true		"body for Apk content"
// @Success 200 {object} models.Apk
// @Failure 403 :id is not int
// @router /:id([0-9]+) [put]
// @author # taoshengbo # 2020-06-08 16:36:35 #
func (c *ApkStrategyController) Put() {
	ctx := GetBackendCntx()

	var v upgrade.ApkUpgradeStrategy

	v.CreateUser = c.mail
	v.UpdateUser = c.mail

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &v); nil != err {
		logs.Error("Unmarshal parameter failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	v.BeginDatetime = time.Unix(v.ReqBeginDatetime, 0)
	v.EndDatetime = time.Unix(v.ReqEndDatetime, 0)

	_id := c.Ctx.Input.Param(":id")
	v.Id, _ = strconv.ParseInt(_id, 10, 64)

	if err := strategy.Update(ctx.Model.Mysql, ctx.Model.Redis, ctx.Model.Quota, &v); nil != err {
		logs.Error("Update strategy failed! id:%d body:%s errmsg:%s",
			_id, c.Ctx.Input.RequestBody, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.ErrorMessage(comm.OK, "Ok")
}

// GetOne ...
// @Title Get One
// @Description get Apk by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Apk
// @Failure 403 :id is empty
// @router /:id([0-9]+) [get]
// @author # taoshengbo # 2020-06-08 16:36:35 #
func (c *ApkStrategyController) GetOne() {
	var err error
	var v *upgrade.ApkUpgradeStrategy

	ctx := GetBackendCntx()

	_id := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(_id, 10, 64)

	if v, err = strategy.GetStrategyById(ctx.Model.Mysql.O, id); nil != err {
		logs.Error("Get strategy failed! id:%d body:%s errmsg:%s",
			id, c.Ctx.Input.RequestBody, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", v)
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
// @author # taoshengbo # 2020-06-08 16:36:35 #
func (c *ApkStrategyController) GetAll() {
	var err error
	var v []*upgrade.ApkUpgradeStrategy

	ctx := GetBackendCntx()

	_apk_id := c.GetString("apk_id")
	apkId, _ := strconv.ParseInt(_apk_id, 10, 64)
	if 0 == apkId {
		logs.Error("Apk id invalid!")
		c.ErrorMessage(comm.ERR_PARAM_INVALID, "Apk id invlaid!")
		return
	}

	if v, err = strategy.GetStrategyListByApkId(ctx.Model.Mysql.O, apkId); nil != err {
		logs.Error("Get strategy list by apk id failed! apkId:%d errmsg:%s!", apkId, err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(v), "Ok", v)
}
