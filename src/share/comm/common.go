package comm

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"

	"upgrade-api/src/share/lib/crypt"
	"upgrade-api/src/share/lib/utils"
)

// 基本类型
const (
	STRING = "string"
	INT    = "int"
	INT8   = "int8"
	UINT8  = "uint8"
	INT16  = "int16"
	UINT16 = "uint16"
	INT32  = "int32"
	UINT32 = "uint32"
	INT64  = "int64"
	UINT64 = "uint64"
	BOOL   = "bool"
	FLOAT  = "float64"
)

/******************************************************************************
 **函数名称: ParamTypeConversion
 **功    能: 参数类型转换，用于进行反射处理时转换为对应的类型
 **输入参数:
 **   	param: 参数值
 **     typ: 目标类型
 **输出参数:
 **返    回:
 **实现描述:
 **作    者: # guoshuangpeng@le.com # 2019-04-16 19:44:28 #
 ******************************************************************************/
func ParamTypeConversion(param string, typ string) (out interface{}, err error) {
	typ = strings.TrimPrefix(typ, "*")
	switch typ {
	case STRING:
		return param, nil
	case INT:
		return strconv.Atoi(param)
	case INT8:
		_out, err := strconv.ParseInt(param, 10, 8)
		return int8(_out), err
	case UINT8:
		_out, err := strconv.ParseUint(param, 10, 8)
		return uint8(_out), err
	case INT16:
		_out, err := strconv.ParseInt(param, 10, 16)
		return int16(_out), err
	case UINT16:
		_out, err := strconv.ParseUint(param, 10, 16)
		return uint16(_out), err
	case INT32:
		_out, err := strconv.ParseInt(param, 10, 32)
		return int32(_out), err
	case UINT32:
		_out, err := strconv.ParseUint(param, 10, 32)
		return uint32(_out), err
	case INT64:
		_out, err := strconv.ParseInt(param, 10, 64)
		return int64(_out), err
	case UINT64:
		_out, err := strconv.ParseUint(param, 10, 64)
		return uint64(_out), err
	case BOOL:
		return strconv.ParseBool(param)
	case FLOAT:
		return strconv.ParseFloat(param, 64)
	}
	return nil, errors.New("type undefined")
}

// 签名key相关
const (
	SIGN_KEY      = "Sign"
	SIGN_TYPE_KEY = "SignType"
)

// 签名类型相关
const (
	SIGN_TYPE_VALUE_MD5 = "MD5"
)

/******************************************************************************
 **函数名称: SignMapByMd5
 **功    能: MD5方式签名
 **输入参数:
 **   	v: 原参数
 **     key: 秘钥
 **		ignore: 是否忽略参数为空数据
 **输出参数:
 **返    回:
 **实现描述:
 **作    者: # guoshuangpeng@le.com # 2020-06-13 10:03:47 #
 ******************************************************************************/
func SignMapByMd5(v map[string]interface{},
	key string, ignore bool) (sign string, err error) {
	var arr []string
	for i, j := range v {
		// 非签名相关参数略过
		if SIGN_KEY == i || SIGN_TYPE_KEY == i {
			continue
		}

		value, err := InterfaceToString(j)
		if nil != err {
			continue
		} else if ignore && len(value) == 0 {
			continue
		}
		arr = append(arr, fmt.Sprintf("%s=%s", i, value))
	}

	// 排序
	sort.Strings(arr)

	sign = strings.Join(arr, "&")

	sign = fmt.Sprintf("%s&key=%s", sign, key)
	logs.Debug("sign str: ", sign)
	sign = crypt.Md5Sum(sign)
	return sign, nil
}

/******************************************************************************
 **函数名称: SignByMd5
 **功    能: MD5方式签名
 **输入参数:
 **   	v: 原参数
 **     businessKey: 机构秘钥
 **输出参数:
 **返    回: 连接后query参数, 秘钥
 **实现描述:
 **作    者: # guoshuangpeng@le.com # 2019-04-18 10:03:47 #
 ******************************************************************************/
func SignByMd5(v interface{},
	key string) (query string, sign string, code int, err error) {

	sortParam, code, err := SortParam(v)
	if nil != err {
		return "", "", code, err
	}

	// 签名
	for _, p := range sortParam {

		// 拼接签名
		if p.Sign {
			if "" == sign {
				sign = p.Key + "=" + p.Value
			} else {
				sign = strings.Join([]string{sign, p.Key + "=" + p.Value}, "&")
			}
		}

		// urlEncode处理
		_value := url.QueryEscape(p.Value)

		// 拼接query字符串
		if "" == query {
			query = p.Key + "=" + _value
		} else {
			query = strings.Join([]string{query, p.Key + "=" + _value}, "&")
		}
	}

	sign = strings.Join([]string{sign, "key=" + key}, "&")
	logs.Debug("sign str: ", sign)
	sign = crypt.Md5Sum(sign)

	return query, sign, OK, nil
}

/******************************************************************************
 **函数名称: SignByMd5ForSso
 **功    能: MD5方式签名
 **输入参数:
 **   	v: 原参数
 **     businessKey: 机构秘钥
 **输出参数:
 **返    回: 连接后query参数, 秘钥
 **实现描述:
 **作    者: # guoshuangpeng@le.com # 2019-04-18 10:03:47 #
 ******************************************************************************/
func SignByMd5ForSso(v interface{},
	key string) (query string, sign string, code int, err error) {

	sortParam, code, err := SortParam(v)
	if nil != err {
		return "", "", code, err
	}
	// 签名
	for _, p := range sortParam {
		// 拼接签名
		if p.Sign {
			if "" == sign {
				sign = p.Key + "=" + p.Value
			} else {
				sign = strings.Join([]string{sign, p.Key + "=" + p.Value}, "&")
			}
		}
		// urlEncode处理
		_value := url.QueryEscape(p.Value)

		// 拼接query字符串
		if "" == query {
			query = p.Key + "=" + _value
		} else {
			query = strings.Join([]string{query, p.Key + "=" + _value}, "&")
		}
	}

	sign = sign + key
	sign = crypt.Md5Sum(sign)

	return query, sign, OK, nil
}

type SortParamStr struct {
	Key   string
	Value string
	Sign  bool
}

// 参数排序
func SortParam(v interface{}) (result []SortParamStr, code int, err error) {

	defer func() {
		// 将panic转换为error 返回
		if r := recover(); r != nil {
			result = nil
			code = ERR_INTERNAL_SERVER_ERROR
			err = errors.New("sign unmarshal failed")
		}
	}()

	var keys []string

	result = make([]SortParamStr, 0)

	o := reflect.ValueOf(v)

	// 判断如果传入类型不是指针或者是空值或者值不可改则返回验证失败
	if o.Kind() != reflect.Ptr || o.IsNil() || !o.Elem().CanSet() {
		return nil, ERR_INTERNAL_SERVER_ERROR, errors.New("sign param type error")
	}

	vt := reflect.TypeOf(v).Elem()
	vr := o.Elem()

	// 取出签名key值
	for i := 0; i < vt.NumField(); i++ {

		// 参数为空，略过
		if vr.Field(i).Kind() == reflect.Ptr && vr.Field(i).IsNil() {
			continue
		} else if _t, _ := InterfaceToString(vr.Field(i).Interface()); vr.Field(i).Kind() == reflect.String && "" == _t {
			continue
		}

		_name := vt.Field(i).Name

		// 非签名相关参数略过
		if SIGN_KEY == _name || SIGN_TYPE_KEY == _name {
			continue
		}

		keys = append(keys, _name)
	}

	// 排序
	sort.Strings(keys)

	// 签名
	for _, k := range keys {
		var _value string
		sf, isExist := vt.FieldByName(k)
		if false == isExist {
			return nil, ERR_INTERNAL_SERVER_ERROR, errors.New("key not exist")
		}

		// 判断如果json字段未设置 选用Name
		_key := sf.Tag.Get("json")
		if "" == _key {
			_key = utils.ToSnake(k)
		} else {
			_key = strings.Replace(_key, " ", "", -1)
			_key = strings.Split(_key, ",")[0]
		}

		// 取值
		if vr.FieldByName(k).Kind() == reflect.Ptr {
			_value, _ = InterfaceToString(vr.FieldByName(k).Elem().Interface())
		} else {
			_value, _ = InterfaceToString(vr.FieldByName(k).Interface())
		}

		_sign := false
		if doSign := sf.Tag.Get("sign"); "" != doSign {
			_sign = true
		}

		item := SortParamStr{
			Key:   _key,
			Value: _value,
			Sign:  _sign,
		}

		// 在result 中插值
		result = append(result, item)
	}

	return result, OK, nil
}

// Interface(基本类型 转 string)
func InterfaceToString(i interface{}) (string, error) {
	switch i.(type) {
	case string:
		return i.(string), nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", i), nil
	case float32:
		return strconv.FormatFloat(float64(i.(float32)), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(i.(float64), 'f', -1, 64), nil
	case bool:
		return strconv.FormatBool(i.(bool)), nil
	case []string:
		return strings.Join(i.([]string), ","), nil
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64:
		return strings.Replace(strings.Trim(fmt.Sprint(i), "[]"), " ", ",", -1), nil
	}

	return "", errors.New("NO MATCH")
}

/******************************************************************************
 **函数名称: Struct2Map
 **功    能: 结构体转MAP
 **输入参数:
 **     s: 结构体对象
 **输出参数: NONE
 **返    回: MAP对象
 **实现描述:
 **注意事项:
 **作    者: # Zhao.yang # 2020-06-19 11:32:09 #
 ******************************************************************************/
func Struct2Map(s interface{}) map[string]interface{} {
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
