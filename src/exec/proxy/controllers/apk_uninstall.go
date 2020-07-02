package controllers

import (
	"errors"
	"net/http"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"

	"upgrade-api/src/exec/proxy/models"
	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
)

// ApkController operations for Apk
type ApkUninstallController struct {
	ApkBaseController
}

// URLMapping ...
func (c *ApkUninstallController) URLMapping() {
	c.Mapping("GetAll", c.GetAll)
}

// GetAll ...
// @Title Get All
// @Description 批量卸载apk
// @Param	dev_id	 query	string	true	"设备ID"
// @Param	model	     query	string	true	"设备型号"
// @Param	platform     query	string	true	"设备平台"
// @Param	len     	 query	string	true	"列表长度"
// @Param	list	     query	string	true	"客户端已安装程序列表"
// @Success 200 {object} models.Apk
// @Failure 403
// @router /list [post]
// @author # yangzhao # 2020-06-08 18:00:05 #
func (c *ApkUninstallController) GetAll() {

	//解析数据
	req, code, err := c.checkUninstallList()
	if nil != err {
		logs.Warn("Apk upgrade parameter format is invalid! msg: %s", err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	//批量卸载处理
	rsp, err := models.ApkUninstallList(ApkUninstallListUrl(), BusinessKey(), req)
	if err != nil {
		logs.Error("Apk uninstall all is failed errmsg :%s", err.Error())
		code, e := database.MysqlFormatError(err)
		c.ErrorMessage(code, e)
		return
	}

	c.FormatInterfaceListResp(http.StatusOK, comm.OK, len(rsp.List), "Ok", rsp.List)
}

//解析参数和校验参数
func (c *ApkUninstallController) checkUninstallList() (ur *models.UpgradeReq, code int, err error) {

	var ul models.UpgradeListInfo

	//解析参数
	err = jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &ul)
	if err != nil {
		c.ErrorMessage(comm.ERR_PARAM_INVALID, err.Error())
		return
	}

	//校验参数
	if len(ul.DevPlat) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("Device plat is empty")
	} else if len(ul.DevId) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("Device id is empty")
	} else if len(ul.DevModel) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("Device model is empty")
	} else if len(ul.List) == 0 {
		return nil, comm.ERR_PARAM_MISS, errors.New("Apk information list is empty")
	}

	return UninstallReq(&ul), comm.OK, nil
}

//apk卸载接口请求参数
func UninstallReq(ul *models.UpgradeListInfo) *models.UpgradeReq {

	uninstallInfoReq := &models.UpgradeReq{
		BusinessId: int64(BusinessId()),
		DevTypeId:  DevTypeId(),
		DevId:      ul.DevId,
		DevModel:   ul.DevModel,
		DevPlat:    ul.DevPlat,
		Len:        len(ul.List),
		List:       ul.List,
	}

	logs.Info("Get request apk body :%s", uninstallInfoReq)

	return uninstallInfoReq
}

//批量卸载url
func ApkUninstallListUrl() string {
	ctx := GetProxyCntx()
	return ctx.Conf.AkpUpgrade.ApkUninstall
}
