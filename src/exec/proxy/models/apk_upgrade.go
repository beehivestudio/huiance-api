package models

import (
	"errors"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
)

/* 升级信息 */
type UpgradeReq struct {
	BusinessId     int64  `json:"business_id"`      // 业务ID
	PatchAlgo      int64  `json:"patch_algo"`       // 客户端支持的差分算法ID（0：下载全量包升级，1：BSDIFFPATCH差分包升级，2：HDIFFPATCH差分包升级）
	DevTypeId      int    `json:"dev_type_id"`      // 客户端设备类型种类ID
	DevId          string `json:"dev_id"`           // 客户端设备ID
	DevModel       string `json:"dev_model"`        // 客户端设备型号
	DevPlat        string `json:"dev_plat"`         // 客户端设备平台
	EuiVersion     string `json:"eui_version"`      // 客户端的EUI版本号
	ApkVersionCode int64  `json:"apk_version_code"` //客户端的应用版本
	Len            int    `json:"len"`              // 列表长度
	List           []List `json:"list"`
}

/* 批量升级请求数据 */
type UpgradeListInfo struct {
	PatchAlgo  int64  `json:"patch_algo"`  // 客户端支持的差分算法ID（0：下载全量包升级，1：BSDIFFPATCH差分包升级，2：HDIFFPATCH差分包升级）
	DevId      string `json:"dev_id"`      // 客户端设备ID
	DevModel   string `json:"dev_model"`   // 客户端设备型号
	DevPlat    string `json:"dev_plat"`    // 客户端设备平台
	EuiVersion string `json:"eui_version"` // 客户端的EUI版本号
	Len        int    `json:"len"`         // 列表长度
	List       []List `json:"list"`
}

/* 升级信息 */
type List struct {
	PackageName    string `json:"package_name"`     // 客户端应用包名
	ApkVersionCode int64  `json:"apk_version_code"` // 客户端应用版本
	ApkMd5         string `json:"apk_md5"`          // 客户端APK的MD5值
}

/******************************************************************************
 **函数名称: ApkUpgradeList
 **功    能: 应用升级（批量）
 **输入参数: o: orm.Ormer
 **			upgradeAllUrl  升级url
**			businessKey    业务key
**			uplistReq      升级请求参数
 **输出参数: NONE
 **返    回: 结构体 UpgradeInfoResp 升级返回信息
 **实现描述:
 **注意事项:
 **作    者: # yangzhao # 2020-06-04 14:10:11 #
 ******************************************************************************/
func ApkUpgradeList(upgradeAllUrl, businessKey string, list *UpgradeReq) (upr *ApkUpgradeData, err error) {

	timestamp := time.Now().Unix()

	//计算和apk接口的签名
	sign, err := ApkAllSign(list, businessKey, upgradeAllUrl, timestamp)
	if err != nil {
		return nil, err
	}

	req, err := jsoniter.Marshal(list)
	if err != nil {
		return nil, err
	}

	logs.Info("Get upgrade all url msg :%s", upgradeAllUrl)

	timeNow := time.Now()

	//请求升级接口
	header := map[string]interface{}{
		"timestamp":    time.Now().Unix(),                 // 时间戳
		"business-id":  list.BusinessId,                   // 业务方ID
		"sign":         sign,                              // 签名数据
		"Content-Type": "application/json; charset=utf-8", // BODY格式
	}

	rsp, err := comm.HttpPostMethod(upgradeAllUrl, header, req)
	if nil != err {
		logs.Error("Request apk upgrade failed! req:%s errmsg:%s", req, err.Error())
		return nil, err
	}

	var resp ApkUpgradeListResp

	//解析请求参数
	if err := jsoniter.Unmarshal(rsp, &resp); nil != err {
		logs.Error("Parsing all upgrade parameters failed! errmsg:%s", err.Error())
		return nil, err
	} else if resp.Code != comm.OK {
		return nil, errors.New(resp.Message)
	}

	logs.Info("Request apk upgrade all spend time :%v", time.Since(timeNow))

	return &resp.Data, nil
}

/* 获取和apk批量接口签名 */
func ApkAllSign(upgradeInfoReq *UpgradeReq, businessKey, url string, timestamp int64) (sign string, err error) {

	str, err := jsoniter.Marshal(upgradeInfoReq)
	if err != nil {
		logs.Info("Get sign body marshal failed! errmsg :%s", err.Error())
		return "", err
	}

	signBody := map[string]interface{}{}
	signBody["body"] = string(str)
	signBody["timestamp"] = timestamp
	signBody["url"] = comm.HttpGetUrlPath(url)

	logs.Info("Get sign body msg :%v，business key :%s", signBody, businessKey)

	sign, err = comm.SignMapByMd5(signBody, businessKey, true)
	if err != nil {
		logs.Error("Get sign failed , errmsg :%s, sign :%s", err.Error(), sign)
		return "", err
	}

	logs.Info("Get sing msg :%s", sign)

	return sign, nil
}

/* apk批量升级应答结果 */
type ApkUpgradeListResp struct {
	Code    int            `json:"code"`    // 错误码
	Message string         `json:"message"` // 错误描述
	Data    ApkUpgradeData `json:"data"`    // 升级信息
}

type ApkUpgradeData struct {
	Len  int               `json:"len"`
	List []UpgradeInfoResp `json:"list"`
}

/* apk单条升级应答结果 */
type ApkUpgradeRsp struct {
	Code    int             `json:"code"`    // 错误码
	Message string          `json:"message"` // 错误描述
	Data    UpgradeInfoResp `json:"data"`    // 升级信息
}

/******************************************************************************
 **函数名称: ApkUpgradeOne
 **功    能: 应用升级（单条）
 **输入参数: o: orm.Ormer
 **			upgradeInfo 升级请求信息
 **输出参数: NONE
 **返    回: 结构体 UpgradeInfoResp 升级返回信息
 **实现描述:
 **注意事项:
 **作    者: # yangzhao # 2020-06-04 14:10:11 #
 ******************************************************************************/
func ApkUpgradeOne(upgradeUrl, businessKey, human string, upgradeInfoReq UpgradeInfoReq) (upr *UpgradeInfoResp, err error) {

	timestamp := time.Now().Unix()

	//计算和apk接口的签名
	sign, err := ApkSign(upgradeUrl+human, businessKey, upgradeInfoReq, timestamp)
	if err != nil {
		return nil, err
	}

	req, err := jsoniter.Marshal(upgradeInfoReq)
	if err != nil {
		return nil, err
	}

	logs.Info("Apk upgrade request sign :%s, body :%s", sign, req)

	timeNow := time.Now()

	//请求升级接口
	header := map[string]interface{}{
		"timestamp":    timestamp,                         // 时间戳
		"business-id":  upgradeInfoReq.BusinessId,         // 业务方ID
		"sign":         sign,                              // 签名数据
		"Content-Type": "application/json; charset=utf-8", // BODY格式
	}

	rsp, err := comm.HttpPostMethod(upgradeUrl+human, header, req)
	if nil != err {
		logs.Error("Apk upgrade request failed! req:%s errmsg:%s", req, err.Error())
		return nil, err
	}

	logs.Info("Apk upgrade response! rsp:%v", string(rsp))

	//解析请求参数
	up := ApkUpgradeRsp{}
	if err := jsoniter.Unmarshal(rsp, &up); nil != err {
		logs.Error("Parsing upgrade rsp parameters failed! errmsg:%s", err.Error())
		return nil, err
	} else if up.Code != comm.OK {
		return nil, errors.New(up.Message)
	}

	logs.Info("Request api send time :%v", time.Since(timeNow))

	return &up.Data, nil
}

/* 获取和apk接口签名 */
func ApkSign(url string, businessKey string, upgradeInfoReq UpgradeInfoReq, timestamp int64) (sign string, err error) {

	str, err := jsoniter.Marshal(upgradeInfoReq)
	if err != nil {
		logs.Info("Get sign body marshal failed! errmsg :%s", err.Error())
		return "", err
	}

	signBody := map[string]interface{}{}
	signBody["body"] = string(str)
	signBody["timestamp"] = timestamp
	signBody["url"] = comm.HttpGetUrlPath(url)

	logs.Info("Get sign body msg :%v，business key :%s", signBody, businessKey)

	sign, err = comm.SignMapByMd5(signBody, businessKey, true)
	if err != nil {
		logs.Error("Get sign failed , errmsg :%s, sign :%s", err.Error(), sign)
		return "", err
	}

	logs.Info("Get sing msg :%s", sign)

	return sign, nil
}

/* 代理升级请求参数 */
type UpgradeInfo struct {
	PatchAlgo          int64  `json:"patch_algo"`           // 客户端支持的差分算法ID（0：下载全量包升级，1：BSDIFFPATCH差分包升级，2：HDIFFPATCH差分包升级）
	DevId              string `json:"dev_id"`               // 客户端设备ID
	DevModel           string `json:"dev_model"`            // 客户端设备型号
	DevPlat            string `json:"dev_plat"`             // 客户端设备平台
	EuiVersion         string `json:"eui_version"`          // 客户端的EUI版本号
	AppointVersionCode int    `json:"appoint_version_code"` // 升级指定版本
	PackageName        string `json:"package_name"`         // 客户端应用包名
	ApkMd5             string `json:"apk_md5"`              // 客户端APK的MD5值
	ApkVersionCode     int64  `json:"apk_version_code"`     // 客户端应用版本
}

/* 代理升级返回参数 */
type UpgradeInfoResp struct {
	PatchAlgo      int64  `json:"patch_algo"`       // 客户端支持的差分算法（0：下载全量包升级，1：BSDIFFPATCH差分包升级，2：HDIFFPATCH差分包升级）
	PackageName    string `json:"package_name"`     // 应用包名
	ApkVersionCode int64  `json:"apk_version_code"` // 应用包版本号
	PackageUrl     string `json:"package_url"`      // 全量包下载地址
	PackageMd5     string `json:"package_md5"`      // 全量包MD5值
	PackageSize    int64  `json:"package_size"`     // 全量包的大小
	PatchUrl       string `json:"patch_url"`        // 差分包下载地址
	PatchMd5       string `json:"patch_md5"`        // 差分包MD5值
	PatchSize      int64  `json:"patch_size"`       // 差分包的大小
}

//请求apk升级参数
type UpgradeInfoReq struct {
	BusinessId         int64  `json:"business_id"`          // 业务ID
	PatchAlgo          int64  `json:"patch_algo"`           // 客户端支持的差分算法ID（0：下载全量包升级，1：BSDIFFPATCH差分包升级，2：HDIFFPATCH差分包升级）
	DevTypeId          int    `json:"dev_type_id"`          // 客户端设备类型种类ID
	DevId              string `json:"dev_id"`               // 客户端设备ID
	DevModel           string `json:"dev_model"`            // 客户端设备型号
	DevPlat            string `json:"dev_plat"`             // 客户端设备平台
	EuiVersion         string `json:"eui_version"`          // 客户端的EUI版本号
	AppointVersionCode int    `json:"appoint_version_code"` // 升级指定版本
	PackageName        string `json:"package_name"`         // 客户端应用包名
	ApkMd5             string `json:"apk_md5"`              // 客户端APK的MD5值
	ApkVersionCode     int64  `json:"apk_version_code"`     //客户端的应用版本
}
