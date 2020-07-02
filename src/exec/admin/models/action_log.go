package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

/* 日志行为数据 */
type ActionLog struct {
	Id         int64     `json:"id" description:"日志id"`
	Name       string    `json:"name" description:"操作人名称"`
	Email      string    `json:"email" description:"操作人邮箱"`
	ActionId   int64     `json:"action_id" description:"用户行为id"`
	Action     string    `json:"action" description:"用户行为描述"`
	Detail     string    `json:"detail" description:"日志详细信息"`
	CreateTime time.Time `json:"create_time" description:"创建时间"`
}

/******************************************************************************
 **函数名称: GetActionLogList
 **功    能: 获取用户列表
 **输入参数:
 **     fields: 过滤条件
 **输出参数: NONE
 **返    回:
 **实现描述:
 **注意事项:
 **作    者: # Zhao.yang # 2020-06-17 14:14:35 #
 ******************************************************************************/
func GetActionLogList(o orm.Ormer, fields []string,
	values []interface{}, page, pageSize int) (int64, []ActionLog, error) {

	logs.Info("Get action log fields msg:%v, values :%v", fields, values)

	var actionLog []ActionLog

	sql := fmt.Sprintf(`SELECT
		actionlog.id,
		userinfor.name,
		userinfor.email,
		actionlog.action_id,
		act.action,
		actionlog.detail,
		actionlog.create_time 
	FROM
		action_log actionlog
	LEFT JOIN user userinfor ON userinfor.id = actionlog.user_id
	LEFT JOIN action act ON actionlog.action_id = act.id`)

	// 添加where条件
	if 0 < len(fields) {
		where := strings.Join(fields, " AND ")
		sql = fmt.Sprintf("%s WHERE %s ", sql, where)
	}

	// 添加排序&分页条件
	sql = fmt.Sprintf("%s LIMIT ? OFFSET ? ", sql)
	values = append(values, pageSize, (page-1)*pageSize)

	logs.Info("Get action log list value msg:%v", values)

	_, err := o.Raw(sql, values...).QueryRows(&actionLog)
	if nil != err {
		logs.Error("%s() Get action log list failed! errmsg:%s",
			utils.GetRunFuncName(), err.Error())
		return 0, nil, err
	}

	return int64(len(actionLog)), actionLog, err
}
