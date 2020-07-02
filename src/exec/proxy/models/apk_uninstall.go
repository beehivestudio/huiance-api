package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
)

/******************************************************************************
 **函数名称: ApkUninstallList
 **功    能: 应用卸载（批量）
 **输入参数:
 **     uninstallAllUrl: 卸载url
 **	    businessKey: 业务key
 **	    uplistReq: 卸载信息结构体
 **输出参数: NONE
 **返    回: 卸载返回信息
 **实现描述:
 **注意事项:
 **作    者: # yangzhao # 2020-06-04 14:10:11 #
 ******************************************************************************/
func ApkUninstallList(uninstallAllUrl, businessKey string,
	uplistReq *UpgradeReq) (list *ApkUninstallData, err error) {

	timestamp := time.Now().Unix()

	// 计算和APK接口的签名
	sign, err := ApkAllSign(uplistReq, businessKey, uninstallAllUrl, timestamp)
	if nil != err {
		logs.Error("Compute sign failed! req:%v key:%s url:%s",
			uplistReq, businessKey, uninstallAllUrl)
		return nil, err
	}

	req, err := jsoniter.Marshal(uplistReq)
	if nil != err {
		logs.Error("Marshal request failed! req:%v key:%s url:%s",
			uplistReq, businessKey, uninstallAllUrl)
		return nil, err
	}

	//请求升级接口
	header := map[string]interface{}{
		"timestamp":    timestamp,                         // 时间戳
		"business-id":  uplistReq.BusinessId,              // 业务方ID
		"sign":         sign,                              // 签名数据
		"Content-Type": "application/json; charset=utf-8", // BODY格式
	}

	rsp, err := comm.HttpPostMethod(uninstallAllUrl, header, req)
	if nil != err {
		logs.Error("Apk uninstall request failed, req:%s errmsg:%s", req, err.Error())
		return nil, err
	}

	logs.Info("Get apk uninstall response body :%v", string(rsp))

	var resp ApkUninstallResp

	//解析请求参数
	if err := jsoniter.Unmarshal(rsp, &resp); nil != err {
		logs.Error("Parsing all uninstall parameters failed! errmsg:%s", err.Error())
		return nil, err
	} else if resp.Code != comm.OK {
		logs.Error("Apk uninstall response failed! body:%v", resp)
		return nil, errors.New(resp.Message)
	}

	return &resp.Data, nil
}

/* apk批量卸载应答结果 */
type ApkUninstallResp struct {
	Code    int              `json:"code"`    // 错误码
	Message string           `json:"message"` // 错误描述
	Data    ApkUninstallData `json:"data"`    // 升级信息
}

type ApkUninstallData struct {
	Len  int                 `json:"len"`
	List []UninstallInfoResp `json:"list"`
}

/* 代理卸载返回参数 */
type UninstallInfoResp struct {
	PackageName    string `json:"package_name"`     // 应用包名
	ApkVersionCode int64  `json:"apk_version_code"` // 应用包版本
}
