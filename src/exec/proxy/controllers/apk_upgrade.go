package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"

	"upgrade-api/src/exec/proxy/models"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
)

// ApkController operations for Apk
type ApkUpgradeController struct {
	ApkBaseController
}

// URLMapping ...
func (c *ApkUpgradeController) URLMapping() {
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)

}

// Get One
// @Title Get One
// @Description   代理层：获取一条指定包最新版本或指定版本APP安装包和差分包的地址
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Apk
// @Failure 403 :id is empty
// @router / [post]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *ApkUpgradeController) GetOne() {

	logs.Info("Get upgrade one request body :%v", string(c.Ctx.Input.RequestBody))

	var err error

	human := c.GetString("human")
	if human == "" {
		c.ErrorMessage(comm.ERR_PARAM_MISS, fmt.Sprint(errors.New("human is empty")))
		return
	}

	// 校验参数和获取请求数据
	req, code, err := c.checkUpgradeInfo()
	if nil != err {
		logs.Warn("Apk upgrade parameter format is invalid! msg: %s", err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	//升级处理
	rsp, err := models.ApkUpgradeOne(ApkUpgradeOneUrl(), BusinessKey(), human, req)
	if err != nil {
		logs.Error("Apk upgrade one is failed errmsg :%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	logs.Info("Apk upgrade response body :%s", rsp)

	c.FormatInterfaceResp(http.StatusOK, comm.OK, "Ok", rsp)
}

//校验参数
func (c *ApkUpgradeController) checkUpgradeInfo() (u models.UpgradeInfoReq, code int, err error) {

	var up models.UpgradeInfo

	//解析参数
	err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &up)
	if err != nil {
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return u, comm.ERR_PARAM_INVALID, errors.New("apk upgrade parameter format err")
	}

	logs.Info("Get request body msg :%v", up)

	//参数校验
	if len(up.DevPlat) == 0 {
		return u, comm.ERR_PARAM_MISS, errors.New("device plat is empty")
	}

	if len(up.DevId) == 0 {
		return u, comm.ERR_PARAM_MISS, errors.New("device id is empty")
	}

	if len(up.DevModel) == 0 {
		return u, comm.ERR_PARAM_MISS, errors.New("device model is empty")
	}

	if len(up.ApkMd5) == 0 {
		return u, comm.ERR_PARAM_MISS, errors.New("appMd5 is empty")
	}

	if len(up.EuiVersion) == 0 {
		return u, comm.ERR_PARAM_MISS, errors.New("euiVersion is empty")
	}

	if len(up.PackageName) == 0 {
		return u, comm.ERR_PARAM_MISS, errors.New("packageName is empty")
	}

	if up.ApkVersionCode == 0 {
		return u, comm.ERR_PARAM_MISS, errors.New("apkVersionCode is empty")
	}

	u = models.UpgradeInfoReq{
		BusinessId:         int64(BusinessId()),
		PatchAlgo:          up.PatchAlgo,
		DevTypeId:          DevTypeId(),
		DevId:              up.DevId,
		DevModel:           up.DevModel,
		DevPlat:            up.DevPlat,
		EuiVersion:         up.EuiVersion,
		AppointVersionCode: up.AppointVersionCode,
		PackageName:        up.PackageName,
		ApkMd5:             up.ApkMd5,
		ApkVersionCode:     up.ApkVersionCode,
	}

	logs.Info("Get request one upgrade body mes :%s", u)

	return u, comm.OK, nil
}

// GetAll ...
// @Title Get All
// @Description 批量升级apk
// @Param	patch_algo	query	int	false	"差分算法ID"
// @Param	dev_id		query	string	true	"客户端设备ID"
// @Param	dev_plat	query	string	true	"客户端设备平台"
// @Param	dev_model	query	string	true	"客户端设备型号"
// @Param	eui_version	query	string	true	"客户端的EUI版本号"
// @Param	len			query	string	true	"列表长度"
// @Param	list		query	string	true	"应用信息"
// @Success 200 {object} models.Apk
// @Failure 403
// @router /list [post]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *ApkUpgradeController) GetAll() {

	logs.Info("Get upgrade all request body :%v", string(c.Ctx.Input.RequestBody))

	action := c.GetString("action")
	human := c.GetString("human")

	if action == "" || human == "" {
		logs.Error("Get action msg:%s, human msg:%s", action, human)
		c.ErrorMessage(comm.ERR_PARAM_INVALID,
			fmt.Sprint(errors.New("human or action is empty")))
		return
	}

	//解析数据
	req, code, err := c.checkUpgradeListInfo()
	if nil != err {
		logs.Warn("Apk Upgrade parameter format is invalid! msg: %s", err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	//批量升级处理
	rep, err := models.ApkUpgradeList(ApkUpgradeListUrl(action, human), BusinessKey(), req)
	if err != nil {
		logs.Error("Apk upgrade all is failed errmsg :%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(rep.List), "Ok", rep.List)
}

//解析参数和校验参数
func (c *ApkUpgradeController) checkUpgradeListInfo() (ur *models.UpgradeReq, code int, err error) {

	var ul models.UpgradeListInfo

	//解析参数
	err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &ul)
	if err != nil {
		logs.Warn("Parsing parameters failed, msg: %s", err.Error())
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	//校验参数
	if len(ul.DevPlat) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("devPlat is empty")
	}

	if len(ul.DevId) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("devId is empty")
	}

	if len(ul.DevModel) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("devModel is empty")
	}

	if len(ul.EuiVersion) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("euiVersion is empty")
	}

	if len(ul.List) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("apk infor list is empty")
	}

	for _, v := range ul.List {
		if v.ApkVersionCode == 0 {
			return nil, comm.ERR_PARAM_MISS, errors.New("apk version code is empty")
		}
	}

	return UpgradeReq(&ul), 0, nil
}

//apk升级接口请求参数
func UpgradeReq(ul *models.UpgradeListInfo) *models.UpgradeReq {

	upgradeInfoReq := &models.UpgradeReq{
		BusinessId: int64(BusinessId()),
		PatchAlgo:  ul.PatchAlgo,
		DevTypeId:  DevTypeId(),
		DevId:      ul.DevId,
		DevModel:   ul.DevModel,
		DevPlat:    ul.DevPlat,
		EuiVersion: ul.EuiVersion,
		Len:        len(ul.List),
		List:       ul.List,
	}

	logs.Info("Get request apk body :%s", upgradeInfoReq)

	return upgradeInfoReq
}

//单个应用升级url
func ApkUpgradeOneUrl() string {
	ctx := GetProxyCntx()
	return ctx.Conf.AkpUpgrade.ApkUpgrade
}

//批量应用升级url
func ApkUpgradeListUrl(action, human string) string {
	ctx := GetProxyCntx()
	return ctx.Conf.AkpUpgrade.ApkUpAradeAll + action + "&human=" + human
}

//业务id
func BusinessId() int {
	ctx := GetProxyCntx()
	return ctx.Conf.BusinessInfor.BusinessId
}

//设备类型Id
func DevTypeId() int {
	ctx := GetProxyCntx()
	return ctx.Conf.BusinessInfor.DevTypeId
}

//业务key
func BusinessKey() string {
	ctx := GetProxyCntx()
	return ctx.Conf.BusinessKey
}
