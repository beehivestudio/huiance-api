package utils

import (
	"strings"
)

// 驼峰转蛇形
func ToSnake(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

// 蛇形转驼峰
func ToCamel(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

// 去除空格和换行
func RemoveEmpty(s string) string {
	// 去除空格
	s = strings.Replace(s, " ", "", -1)
	// 去除换行符
	s = strings.Replace(s, "\n", "", -1)
	// 去除制表符
	s = strings.Replace(s, "\t", "", -1)
	return s
}

// 去除邮箱后缀
func RemoveMailSuffix(s string) string {
	s = strings.Replace(s, "@le.com", "", -1)
	s = strings.Replace(s, "@letv.com", "", -1)
	return s
}

/******************************************************************************
 **函数名称: InterfaceSliceToStringSlice
 **功    能: InterfaceSlice 转 StringSlice
 **输入参数:
 **      i: InterfaceSlice
 **输出参数: NONE
 **返    回: s: StringSlice
 **实现描述:
 **注意事项:
 **作    者: # taoshengbo # 2020-05-29 17:46:06 #
 ******************************************************************************/
func InterfaceSliceToStringSlice(inter []interface{}) (str []string) {
	str = make([]string, len(inter))
	for i := range inter {
		str[i] = inter[i].(string)
	}
	return str
}
