package comm

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego/logs"
)

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("XRealIP"); "" != ip {
		remoteAddr = ip
	} else if ip = req.Header.Get("XForwardedFor"); "" != ip {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if "::1" == remoteAddr {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

/******************************************************************************
 **函数名称: HttpGetMethod
 **功    能: 发送http请求
 **输入参数:
 **     url: 请求地址
 **输出参数: NONE
 **返    回:
 **     respBody：服务器返回信息
 **     err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2018-12-03 15:42:19 #
 ******************************************************************************/
func HttpGetMethod(url string) (respBody []byte, err error) {
	response, err := http.Get(url)
	if nil != err {
		return nil, err
	}

	defer response.Body.Close()

	/* 读取返回值信息 */
	respBody, err = ioutil.ReadAll(response.Body)
	if nil != err {
		return nil, err
	}

	return respBody, nil
}

/******************************************************************************
 **函数名称: HttpPostMethod
 **功    能: 发送http请求
 **输入参数:
 **     url: 请求地址
 **     header: HTTP请求HEADER参数
 **     data: JSON数据
 **输出参数: NONE
 **返    回:
 **     respBody：服务器返回信息
 **     err: 错误信息
 **实现描述:
 **注意事项:
 ** 	1.如果BODY为JSON格式, 则HEADER参数中需要设置"Content-Type"的值为application/json; charset=utf-8".
 **作    者: # Shuangpeng.guo # 2018-12-03 15:42:19 #
 ******************************************************************************/
func HttpPostMethod(url string,
	header map[string]interface{}, data []byte) (respBody []byte, err error) {
	client := &http.Client{}
	if nil != err {
		return nil, err
	}

	body := bytes.NewBuffer([]byte(data))

	req, err := http.NewRequest("POST", url, body)
	if nil != err {
		return nil, err
	}

	/* 设置HEADER参数 */
	for k, v := range header {
		req.Header.Set(k, fmt.Sprint(v))
	}

	/* 发送消息 */
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	/* 读取返回值信息 */
	respBody, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return respBody, nil
}

/******************************************************************************
 **函数名称: HttpPostByXwww
 **功    能: 发送http请求
 **输入参数:
 **     url: 请求地址
 **    data: json数据
 **输出参数: NONE
 **返    回:
 **     respBody：服务器返回信息
 **     err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2018-12-03 15:42:19 #
 ******************************************************************************/
func HttpPostByXwww(url string, data string) (respBody []byte, err error) {
	resp, err := http.Post(url, "application/x-www-form-urlencoded",
		strings.NewReader(data))
	if nil != err {
		return nil, err
	}

	defer resp.Body.Close()

	/* 读取返回值信息 */
	respBody, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return respBody, nil
}

/******************************************************************************
 **函数名称: HttpPutMethod
 **功    能: 发送http请求[PUT]
 **输入参数:
 **     url: 请求地址
 **     data: json数据
 **输出参数: NONE
 **返    回:
 **     respBody：服务器返回信息
 **     err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2019-05-07 16:54:26 #
 ******************************************************************************/
func HttpPutMethod(url string, data []byte) (respBody []byte, err error) {
	client := &http.Client{}
	if nil != err {
		return nil, err
	}

	body := bytes.NewBuffer([]byte(data))

	req, err := http.NewRequest("PUT", url, body)
	if nil != err {
		return nil, err
	}

	/* 设置请求json格式 */
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	/* 发送消息 */
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	/* 读取返回值信息 */
	respBody, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return respBody, nil
}

/******************************************************************************
 **函数名称: HttpPutByXwww
 **功    能: 发送http请求[PUT]
 **输入参数:
 **     url: 请求地址
 **     data: json数据
 **输出参数: NONE
 **返    回:
 **     respBody：服务器返回信息
 **     err: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Qifeng.zou # 2019-05-07 16:54:26 #
 ******************************************************************************/
func HttpPutByXwww(url string, data string) (respBody []byte, err error) {
	client := &http.Client{}
	if nil != err {
		return nil, err
	}

	body := bytes.NewBuffer([]byte(data))

	req, err := http.NewRequest("PUT", url, body)
	if nil != err {
		return nil, err
	}

	/* 设置请求json格式 */
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	/* 发送消息 */
	resp, err := client.Do(req)
	if nil != err {
		return nil, err
	}
	defer resp.Body.Close()

	/* 读取返回值信息 */
	respBody, err = ioutil.ReadAll(resp.Body)
	if nil != err {
		return nil, err
	}

	return respBody, nil
}

/******************************************************************************
 **函数名称: HttpGetUrlPath
 **功    能: 获取请求路径后缀
 **输入参数:
 **     _url: HTTP请求地址
 **输出参数: NONE
 **返    回: PATH及QUERY参数
 **实现描述:
 **     假如url为"wwww.api.upgrade.le.com/upgrade/api/v1/action?uid=1", 经此函数
 **     处理后返回"/upgrade/v1/action?uid=1"
 **注意事项:
 **作    者: # Zhao.yang # 2020-06-19 11:32:09 #
 ******************************************************************************/
func HttpGetUrlPath(_url string) string {
	u, err := url.Parse(_url)
	if nil != err {
		logs.Error("Parse url failed! errmsg:%s", err.Error())
		return ""
	} else if "" == u.RawQuery {
		return u.Path
	}

	return u.Path + "?" + u.RawQuery
}
