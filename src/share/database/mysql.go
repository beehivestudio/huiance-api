package database

import (
	"regexp"
	"strconv"

	"upgrade-api/src/share/comm"
	"upgrade-api/src/share/lib/mysql"
)

/******************************************************************************
 **函数名称: MysqlFormatError
 **功    能: 格式化MySQL异常信息
 **输入参数:
 **     err: 错误信息
 **输出参数: NONE
 **返    回:
 **     code: 错误码
 **     message: 错误信息
 **实现描述:
 **注意事项:
 **作    者: # Shuangpeng.guo # 2018-11-30 10:55:24 #
 ******************************************************************************/
func MysqlFormatError(err error) (code int, message string) {
	if nil == err {
		return comm.OK, "Ok"
	}

	reg := regexp.MustCompile(`^Error ([0-9]+): (.*)$`)

	a := reg.FindStringSubmatch(err.Error())
	if len(a) < 3 {
		return comm.ERR_INTERNAL_SERVER_ERROR, err.Error()
	}

	code, _ = strconv.Atoi(a[1])
	switch code {
	case mysql.ERR_SYNTAX:
		return comm.ERR_SYS_MYSQL, "Syntax error!"
	case mysql.ERR_SAFE_UPDATE:
		return comm.ERR_SYS_MYSQL, "Unsafe update!"
	case mysql.ERR_DATA_TYPE:
		return comm.ERR_SYS_MYSQL, "Data type error!"
	case mysql.ERR_ACCESS_DENIED:
		return comm.ERR_SYS_MYSQL, "Access denied!"
	case mysql.ERR_MASTER_NODE:
		return comm.ERR_SYS_MYSQL, "Master node crash!"
	case mysql.ERR_TOO_MANY_OPEN_FILE:
		return comm.ERR_SYS_MYSQL, "Too many open file!"
	case mysql.ERR_DUPLICATE_ENTRY:
		return comm.ERR_DUPLICATE_ENTRY, "Pk conflict!"
	case mysql.ERR_CHARSET:
		return comm.ERR_SYS_MYSQL, "Charset error!"
	case mysql.ERR_LIMIT:
		return comm.ERR_PARAM_EXCEED_LIMIT, "Exceed limit!"
	case mysql.ERR_FOREIGN_KEY:
		return comm.ERR_DUPLICATE_ENTRY, "Fk conflict!"
	}

	return comm.ERR_SYS_MYSQL, err.Error()
}
