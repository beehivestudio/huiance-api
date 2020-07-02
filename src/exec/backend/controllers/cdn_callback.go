package controllers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/database"
	"upgrade-api/src/share/database/upgrade"
	"upgrade-api/src/share/lib/patch"
)

type CdnCallBackController struct {
	comm.BaseController
}

// Post CDN回调
// @Title CDN回调
// @Description CDN回调
// @Param	outkey	    query	string	true	"文件唯一标识"
// @Param	storeurl	query	string	true	"cdn存储短路径"
// @Param	status	    query	string	true	"状态"
// @Param	md5	        query	string	true	"文件md5"
// @Success 201 {int} models.Business
// @Failure 403 body is empty
// @router / [post]
// @author # yangzhao # 2020-06-10 18:00:05 #
func (c *CdnCallBackController) Post() {

	ctx := GetBackendCntx()

	//校验参数
	outKey, storeUrl, code, err := c.checkParameter()
	if err != nil {
		logs.Error("Parameter is invalid! code:%d errmsg:%s", code, err.Error())
		c.ErrorMessage(code, err.Error())
		return
	}

	//拼接cdn回调短路径
	cdnUrl := comm.CDN_UPLOAD_CDN_HOST + storeUrl

	logs.Info("Cdn callback short url msg :%s", storeUrl)

	if strings.Contains(outKey, comm.CDN_CALLBACK_PATCH) { // 差分包上传CDN回调
		err = patch.PatchCdnCallBack(ctx.Model.Mysql.O, ctx.Model.Redis, cdnUrl, outKey)
		if err != nil {
			logs.Error("Callback patch failed! errmsg:%s", err.Error())
			code, e := database.MysqlFormatError(err)
			c.ErrorMessage(code, e)
			return
		}
	} else if strings.Contains(outKey, comm.CDN_CALLBACK_APK) { // APK上传CDN回调
		err = upgrade.UpdateApkUpload(ctx.Model.Mysql.O, cdnUrl, outKey)
		if err != nil {
			logs.Error("Callback apk failed! errmsg:%s", err.Error())
			code, e := database.MysqlFormatError(err)
			c.ErrorMessage(code, e)
			return
		}
	}

	c.FormatResp(http.StatusOK, comm.OK, "Ok")

	return
}

/* CDN回调参数 */
type CdnCallBack struct {
	Status   int    `json:"status"`   // 状态 (200:成功 400:重复)
	Storeurl string `json:"storeurl"` // 存储路径
	Md5      string `json:"md5"`      // MD5值
	Outkey   string `json:"outkey"`   // OUT KEY
}

/* 检查参数合法性 */
func (c *CdnCallBackController) checkParameter() (outKey, storeUrl string, code int, err error) {

	logs.Info("Cdn callback request body :%s", c.Ctx.Request.RequestURI)

	storeurl := c.GetString("storeurl")
	if 0 == len(storeurl) {
		return "", "", comm.ERR_PARAM_MISS, errors.New("Callback store url is invalid")
	}

	outkey := c.GetString("outkey")
	if 0 == len(outkey) {
		return "", "", comm.ERR_PARAM_MISS, errors.New("Callback out key is invalid")
	}

	status, err := c.GetInt("status")
	if err != nil {
		logs.Error("Get status failed. errmsg:%s", err.Error())
		return "", "", comm.ERR_PARAM_MISS, err
	} else if status != comm.CDN_CODE_SUCC {
		logs.Error("Callback status is not success!")
		return "", "", comm.ERR_PARAM_MISS, errors.New("Callback status not success")
	}

	return outkey, storeurl, 0, nil
}
