package controllers

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/json-iterator/go"

	"upgrade-api/src/share/comm"
)

type ApkBaseController struct {
	comm.BaseController
}

/******************************************************************************
 **函数名称: Prepare
 **功    能: 与客户端签名校验
 **输入参数:
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项: 相关wiki地址:http://wiki.letv.cn/pages/viewpage.action?pageId=37323874#API%E5%AE%89%E5%85%A8%E8%A7%84%E8%8C%832-%E7%AD%BE%E5%90%8D%E4%BB%A3%E7%A0%81-Golang
 **作    者: # zhao.yang # 2020-06-09 10:14:35 #
 ******************************************************************************/
func (c *ApkBaseController) Prepare() {
	ctx := GetProxyCntx()

	//加载配置文件
	my_ak := ctx.Conf.Sign.AkSign
	my_sk := ctx.Conf.Sign.SkSign

	if len(my_ak) == 0 || len(my_sk) == 0 {
		return
	}

	sign, date, err := c.CheckParameter(my_ak)
	if err != nil {
		c.FormatResp(http.StatusUnauthorized, comm.ERR_AUTH, err.Error())
		return
	}

	my_sign := ""

	if string(c.Ctx.Request.Method) == "GET" {

		forms := make(map[string][]string)
		params := c.Ctx.Input.Params()
		for k, v := range params {
			forms[k] = []string{fmt.Sprintf("%v", v)}
		}

		my_sign = Sign(my_sk, "GET", GetPath(string(c.Ctx.Request.RequestURI)), []byte{}, date, forms)

	} else if string(c.Ctx.Request.Method) != "Get" {

		forms := make(map[string][]string)
		body := map[string]interface{}{}

		err := jsoniter.Unmarshal(c.Ctx.Input.RequestBody, &body)
		if err != nil {
			c.FormatResp(http.StatusUnauthorized, comm.ERR_AUTH, "Parsing parameter failed")
			return
		}
		for k, v := range body {
			forms[k] = []string{fmt.Sprintf("%v", v)}
		}

		logs.Info("Sign parameter my_sk :%s, url :%s, body :%s", GetPath(string(c.Ctx.Request.RequestURI)), c.Ctx.Input.RequestBody)

		my_sign = Sign(my_sk, c.Ctx.Request.Method, GetPath(string(c.Ctx.Request.RequestURI)), c.Ctx.Input.RequestBody, date, forms)
	}

	logs.Info("Get my_sign msg:%s , req_sign msg:%s", my_sign, sign)

	//验证签名
	if sign != my_sign {
		c.FormatResp(http.StatusUnauthorized, comm.ERR_AUTH, "sign failed")
		return
	}

	logs.Info("Sign success !")
}

/* 参数验证 */
func (c *ApkBaseController) CheckParameter(myAk string) (sign, date string, err error) {

	date = c.Ctx.Request.Header.Get("Date")
	if len(date) == 0 {
		return "", "", errors.New("sign failed! date err")
	}

	logs.Info("Get date msg:%s", date)

	auth := c.Ctx.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		return "", "", errors.New("sign failed! auth err")
	}

	logs.Info("Get authorization msg:%s", auth)

	ss := strings.Split(auth, " ")

	hdr := ss[0]
	ak := ss[1]
	if len(ss) < 3 || hdr != "LETV" || ak != myAk {
		return "", "", errors.New("sign failed! auth len err")
	}

	sign = ss[2]
	logs.Info("Get sign msg %s", sign)

	return sign, date, nil

}

/* 与客户端的签名 */
func Sign(key string, method string, path string, body []byte, date string, forms map[string][]string) string {

	logs.Info("Get Parameter forms msg:%s", forms)

	var params []string
	for k, a := range forms {
		for _, v := range a {
			// if v is a empty string, skip
			if v == "" {
				continue
			}
			params = append(params, fmt.Sprintf("%s=%s", k, v))
		}
	}

	sort.Strings(params)
	paramString := strings.Join(params, "&")

	logs.Info("Get paramString msg:%s", paramString)

	bodyMD5 := ""
	if len(body) > 0 {
		h := md5.New()
		h.Write(body)
		bodyMD5 = hex.EncodeToString(h.Sum(nil))
	}

	logs.Info("Get bodyMd5 msg :%s", bodyMD5)

	stringToSign := method + "\n" +
		path + "\n" +
		bodyMD5 + "\n" + date

	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(stringToSign))

	logs.Info("Get String to sign msg :%s", stringToSign)

	return hex.EncodeToString(mac.Sum(nil))
}

//获取请求path
func GetPath(ul string) string {
	u, err := url.Parse(ul)
	if nil != err {
		logs.Error("Parse url failed! errmsg:%s", err.Error())
		return ""
	}
	return u.Path
}
