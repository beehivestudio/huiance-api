package controllers

import (
	"net/http"
	"strconv"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/strategy"
)

// ApkStrategyController operations for ApkUpgrade
type ApkUpgradeController struct {
	ApiBaseController
}

// URLMapping ...
func (c *ApkUpgradeController) URLMapping() {
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
}

// GetOne ...ApkUpgradeController
// @Title Get One
// @Description get Apk by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Apk
// @Failure 403 :id is empty
// @router / [post]
// @author # taoshengbo # 2020-06-08 16:36:35 #
func (c *ApkUpgradeController) GetOne() {
	ctx := GetApiCntx()

	var upgradeInfo strategy.UpgradeInfo

	businessId := c.Ctx.Request.Header.Get("business-id")
	upgradeInfo.BusinessId, _ = strconv.ParseInt(businessId, 10, 64)
	human, _ := strconv.Atoi(c.GetString("human"))
	logs.Info("human=", human)

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &upgradeInfo); nil != err {
		logs.Error("Unmarshal request body failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	var upgradeInfoResp *strategy.UpgradeInfoResp
	var cerr *comm.Error
	if upgradeInfoResp, cerr = strategy.ApkUpgradeOne(
		ctx.Model.Mysql.O, ctx.Model.Redis, ctx.Model.Quota,
		comm.UPGRADE_SILENCE, human, &upgradeInfo); cerr != nil && nil != cerr.Err {
		logs.Error("Get upgrade apk version failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, cerr.Err.Error())
		code, e := database.MysqlFormatError(cerr.Err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", upgradeInfoResp)
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
// @router /list [post]
// @author # taoshengbo # 2020-06-08 16:36:35 #
func (c *ApkUpgradeController) GetAll() {
	ctx := GetApiCntx()

	var upgradeListInfo strategy.UpgradeListInfo

	businessId := c.Ctx.Request.Header.Get("business-id")
	upgradeListInfo.BusinessId, _ = strconv.ParseInt(businessId, 10, 64)
	action, _ := strconv.Atoi(c.GetString("action"))
	human, _ := strconv.Atoi(c.GetString("human"))

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &upgradeListInfo); nil != err {
		logs.Error("Unmarshal request body failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	var upgradeInfoResp []*strategy.UpgradeInfoResp
	var cerr *comm.Error
	if upgradeInfoResp, cerr = strategy.ApkUpgradeBatch(
		ctx.Model.Mysql.O, ctx.Model.Redis, ctx.Model.Quota,
		action, comm.UPGRADE_SILENCE, human, &upgradeListInfo); cerr != nil && nil != cerr.Err {
		logs.Error("Batch get upgrade apk list failed! errmsg:%s", cerr.Err.Error())
		code, e := database.MysqlFormatError(cerr.Err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceListResp(http.StatusOK,
		comm.OK, len(upgradeInfoResp), "Ok", upgradeInfoResp)
}
