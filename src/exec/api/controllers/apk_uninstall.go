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

// ApkStrategyController operations for ApkUninstall
type ApkUninstallController struct {
	ApiBaseController
}

// URLMapping ...
func (c *ApkUninstallController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
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
func (c *ApkUninstallController) GetAll() {
	ctx := GetApiCntx()

	var list strategy.UpgradeListInfo

	businessId := c.Ctx.Request.Header.Get("business-id")
	list.BusinessId, _ = strconv.ParseInt(businessId, 10, 64)

	if err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &list); nil != err {
		logs.Error("Unmarshal request failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	var upgradeInfoResp []*strategy.UpgradeInfoResp
	var cerr *comm.Error
	if upgradeInfoResp, cerr = strategy.ApkUpgradeBatch(
		ctx.Model.Mysql.O, ctx.Model.Redis, ctx.Model.Quota, comm.ACTION_UPGRADE, comm.UPGRADE_UNINSTALL,
		comm.HUMAN_NOT, &list); cerr != nil && nil != cerr.Err {

		logs.Error("Batch get apk list failed! body:%s errmsg:%s",
			c.Ctx.Input.RequestBody, cerr.Err.Error())

		code, e := database.MysqlFormatError(cerr.Err)
		c.ErrorMessage(code, e)
		return
	}

	uninstallInfoResp := upgradeInfoRespToUninstallInfoResp(upgradeInfoResp)

	c.FormatInterfaceListResp(http.StatusOK,
		comm.OK, len(uninstallInfoResp), "Ok", uninstallInfoResp)
}

/******************************************************************************
 **函数名称: upgradeInfoRespToUninstallInfoResp
 **功    能: 卸载信息格式化
 **输入参数: upgradeInfoResp: 卸载信息列表
 **输出参数: NONE
 **返    回: 结构体 uninstallInfoResp 格式化后卸载信息列表
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-06-04 15:10:34 #
 ******************************************************************************/
func upgradeInfoRespToUninstallInfoResp(upgradeInfoResp []*strategy.UpgradeInfoResp) (
	uninstallInfoResp []*strategy.UninstallInfoResp) {

	for _, _v := range upgradeInfoResp {
		v := strategy.UninstallInfoResp{}
		v.PackageName = _v.PackageName
		v.ApkVersionCode = _v.ApkVersionCode
		uninstallInfoResp = append(uninstallInfoResp, &v)
	}

	return uninstallInfoResp
}
